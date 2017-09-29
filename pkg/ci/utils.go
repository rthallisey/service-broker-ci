package ci

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// TODO: Consider renaming to getRepoAddr
func getScriptAddr(repoScriptAndArgs string, dir string) (string, string) {
	var script, args string
	// Check for a valid git repo, otherwise look locally
	repo := resolveGitRepo(repoScriptAndArgs)

	// Split rthallisey/service-broker-ci/wait-for-resource.sh create mediawiki
	//
	// addr="rthallisey/service-broker-ci/wait-for-resource.sh"
	// args="create mediawiki"
	//
	if repo == "" {
		items := strings.Split(repoScriptAndArgs, " ")
		script = items[0]
		if len(items) > 1 {
			args = strings.Join(items[1:len(items)], " ")
		}
		if dir == "template" {
			return fmt.Sprintf("template/%s.yaml", script), args
		}
		return script, args
	} else {
		script, args = getScriptAndArgs(repo, repoScriptAndArgs)
		if dir == "template" {
			fmt.Println(script)
			return fmt.Sprintf("%s/%s/%s/template/%s.yaml", BaseURL, repo, Branch, script), args
		} else if dir == "script" {
			return fmt.Sprintf("%s/%s/%s/%s", BaseURL, repo, Branch, script), args
		}
		return "", ""
	}
}

func resolveGitRepo(repo string) string {
	var validRepo string

	// Loop through each string in a git repo and combine them to test
	// for a valid git repo. If there is no valid repo found, look locally
	// for the file.
	//
	// `curl https://github.com/fake-git-user/fake-git-repo`  - FAIL
	// `curl https://github.com/rthallisey/service-broker-ci` - PASS
	//
	addr := strings.Split(repo, "/")
	if len(addr) >= 2 {
		// A git repo's address is always the first two items
		//     rthallisey/service-broker-ci/...
		baseRepo := addr[0:2]
		gitURL := []string{"https://github.com"}
		for count, _ := range addr {
			// Combine 0...N items to form the url for testing
			validRepo = strings.Join(addr[0:count], "/")
			if validRepo == "" {
				validRepo = strings.Join(baseRepo, "/")
			}

			// Combine: 'https://github.com' + '/' + 'rthallisey/service-broker-ci'
			gitURL = append(gitURL, validRepo)
			validRepo = strings.Join(gitURL, "/")

			req, _ := http.Get(validRepo)
			defer req.Body.Close()

			if req.StatusCode == http.StatusOK {
				validRepo = strings.Split(validRepo, "https://github.com/")[1]
				fmt.Printf("REPO: %s\n", validRepo)
				break
			}
		}
	}
	return validRepo
}

func getScriptAndArgs(repo string, repoScriptAndArgs string) (string, string) {
	// Split 'openshift/ansible-service-broker' and
	// '/scripts/broker-ci/wait-for-resource.sh create mediawiki'
	scriptAndArgs := strings.Split(repoScriptAndArgs, repo)[1]

	// Verify we have ARGS
	s := strings.Split(scriptAndArgs, " ")
	if len(s) < 2 {
		return s[0], ""
	}

	script := s[0]
	fmt.Printf("SCRIPT: %s\n", script)

	listArgs := s[1:len(s)]
	args := strings.Join(listArgs, " ")
	fmt.Printf("ARGS: %s\n", args)

	return script, args
}

func findBindTarget(repo string, provisioned []string) (string, []string, error) {
	var usedTargets []int
	foundTarget := false
	foundBind := false
	var bindTarget string

	// The config in imperative so order matters
	for count, r := range provisioned {
		// Remove the first Provisioned app that matches the Bind repo
		// and the first Provisioned app that doesn't.

		// The first Provisioned app that doesn't match the bindApp is
		// the bindTarget.
		if r != repo && !foundTarget {
			bindTarget = r
			foundTarget = true
			usedTargets = append(usedTargets, count)
		}

		// The first Provisioned app that matches the Bind repo is the
		// bind app.
		if r == repo && !foundBind {
			foundBind = true
			usedTargets = append(usedTargets, count)
		}

		if foundBind && foundTarget {
			cleanupUsedTargets(usedTargets, provisioned)
			return bindTarget, provisioned, nil
		}
	}

	return "", provisioned, errors.New("Failed to find a provisioned bind target and bind app")
}

func cleanupUsedTargets(usedTargets []int, provisioned []string) []string {
	for _, d := range usedTargets {
		// Cleanup the bindTarget and the bind app
		if len(provisioned) == 1 {
			provisioned = provisioned[:0]
		} else {
			provisioned = append(provisioned[:d], provisioned[d+1:]...)
		}
	}
	return provisioned
}

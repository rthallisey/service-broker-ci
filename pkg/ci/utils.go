package ci

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rthallisey/service-broker-ci/pkg/action"
)

func getTemplateAddr(repo string) string {
	fmt.Printf("REPO: %s\n", repo)

	r := strings.Split(repo, "/")

	// TODO: combine with actions.resourceName
	// [ ansibleplaybookbundle, mediawiki123 ]
	apb := r[len(r)-1]
	gitOrg := strings.Join(r[0:len(r)-1], "/")

	// APB template will be in the template directory
	return fmt.Sprintf("%s/%s/%s/template/%s.yaml", BaseURL, gitOrg, Branch, apb)
}

func getScriptAddr(repoScriptAndArgs string) (string, string) {
	// Check for a valid git repo, otherwise look locally
	repo := resolveGitRepo(repoScriptAndArgs)

	// Split rthallisey/service-broker-ci/wait-for-resource.sh create mediawiki
	//
	// addr="rthallisey/service-broker-ci/wait-for-resource.sh"
	// args="create mediawiki"
	//
	script, args := getScriptAndArgs(repo, repoScriptAndArgs)
	return fmt.Sprintf("%s/%s/%s/%s", BaseURL, repo, Branch, script), args
}

func resolveGitRepo(repo string) string {
	var validRepo string

	// Loop through each string in a git repo and combine them to test
	// for a valid git repo. If there is no valid repo found, look locally
	// for the file. The local file check hasn't been implemented yet.
	//
	// `git ls-remote https://github.com/rthallisey`                   - FAIL
	// `git ls-remote https://github.com/rthallisey/service-broker-ci` - PASS
	//
	addr := strings.Split(repo, "/")
	for count, _ := range addr {
		gitURL := []string{"https://github.com"}

		// Combine 0...N items to form the url for testing
		validRepo = strings.Join(addr[0:count], "/")

		// Combine: 'https://github.com' + '/' + 'rthallisey'
		gitURL = append(gitURL, validRepo)
		validRepo = strings.Join(gitURL, "/")

		// Test if the repo is a valid git repo
		_, err := action.RunCommand("git", "ls-remote", validRepo)
		if err == nil {
			// Return without 'https://github.com/'
			validRepo = strings.Split(validRepo, "https://github.com/")[1]
			fmt.Printf("REPO: %s\n", validRepo)
			break
		}
	}
	//TODO: If there's no valid github repo found, look locally for the file

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

		// The first Provisioned app that doesn't match Bind is the
		// bindTarget.
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

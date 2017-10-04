# service-broker-ci
GO package that provides a CI framework for testing Service Catalog Instances.

### Syntax

```yaml
<API KEYS>: <FILE>
```


##### API Keys
- provision
- bind
- unbind
- deprovision
- verify


##### File
The file field accepts a valid git repo ```rthallisey/service-broker-ci/postgresql```
of the apb or a local file. The API Keys provision, bind, unbind, and
deprovision will add a .yaml extention to the path.


##### Verify
Verify is used to check if an action is successful.  Verify accepts a script
from git repo ```rthallisey/service-broker-ci/wait-for-resource.sh``` or a local
script ```wait-for-resource.sh```.

The Verify API Key is also a shell. It can run any shell command and return the
output.
```yaml
verify: oc get pods
```


### Config file format
Templates used by provision, bind, unbind, and deprovision are expected to be in
the template directory. Everything else uses the full path provided.
```bash
.
|── template
│   ├── mediawiki123.yaml
│   ├── postgresql-mediawiki123-bind.yaml
│   └── postgresql.yaml
```


##### Using Local Paths
The config file accepts local paths to scripts and templates.

Every template will be searched for in the ```templates``` directory locally
while other scripts will use the top level directory.
```yaml
provision: mediawiki123
verify: wait-for-resource.sh create pod mediawiki
```


##### Matching Paths
When describing a path to a template, that path will be the key used to identify
which app is being acted upon.

For example, to provision and deprovision the same postgresql app, use matching
paths.
```yaml
provision: postgresql
verify: wait-for-resource.sh create pod postgresql

deprovision: postgresql
verify: wait-for-resource.sh delete pod postgresql
```


### Bind Ordering
There are two applications that are used in a bind, the **bindApp** and the
**bindTarget**. The bindApp is the application that will be binded to another
application. If I say: "I want to bind postgresql to mediawiki". Then, the
bindApp will be postgresql. The application that is being binded to, or is
receiving the bind credentials, is the bindTarget. Mediaiwiki is the bindTarget
from the example.

In config.yaml, the bindTarget is determined by the first application
provisioned that's not the same as the bindApp.

This will bind postgresql to mediawiki123.
```yaml
provision: rthallisey/service-broker-ci/mediawiki123
provision: rthallisey/service-broker-ci/postgresql

bind: rthallisey/service-broker-ci/postgresql
```

You can do multiple bindings with different applications where the next bind
call will bind to the next availble provisioned application.

This will bind postgresql to mediawiki123 and mariadb to elasticsearch.
```yaml
provision: rthallisey/service-broker-ci/mediawiki123
provision: rthallisey/service-broker-ci/elasticsearch

provision: rthallisey/service-broker-ci/postgresql
provision: rthallisey/service-broker-ci/mariadb

bind: rthallisey/service-broker-ci/postgresql
bind: rthallisey/service-broker-ci/mariadb
```

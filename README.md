# service-broker-ci
GO package that provides a CI framework for testing Service Catalog Instances

### Syntax

```yaml
<API KEYS>: <FILE>
```


#### API Keys
provision
bind
unbind
deprovision
verify


#### File
The file field accepts a valid git repo ```rthallisey/service-broker-ci/postgresql```
of the apb or a local file ```postgresql```.


#### Verify
Verify is used to check if an action is successful.  Verify accepts a script
from git repo ```rthallisey/service-broker-ci/wait-for-resource.sh``` or a local
script ```wait-for-resource.sh```.


### Directory Structure
Templates are expected to be in the template directory. Everything else uses the
full path provided.
.
|── template
│   ├── mediawiki123.yaml
│   ├── postgresql-mediawiki123-bind.yaml
│   └── postgresql.yaml


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


### Local Paths

The config file accepts local paths to scripts and templates.

Any template will be search for in the templates directory locally
and other scripts will use the top level directory.
```yaml
provision: mediawiki123
verify: wait-for-resource.sh create pod mediawiki
```


### Matching Paths

When describing a path to a template, that path will the key used to identify
which app is being acted upon.

To provision and deprovision the same postgresql app, use matching paths.
```yaml
provision: postgresql
verify: wait-for-resource.sh create pod postgresql

deprovision: postgresql
verify: wait-for-resource.sh delete pod postgresql
```

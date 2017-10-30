#!/bin/bash

oc delete -f templates/postgresql-mediawiki123-apb-bind.yaml
./wait-for-resource.sh delete servicebinding binding

oc delete serviceinstance mediawiki postgresql
./wait-for-resource.sh delete serviceinstance postgresql
./wait-for-resource.sh delete serviceinstance mediawiki

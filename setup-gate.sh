#!/bin/bash

wget -O /tmp/glide.tar.gz https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
tar xfv /tmp/glide.tar.gz -C /tmp
sudo mv $(find /tmp -name "glide") /usr/bin
sudo wget -O /bin/oc https://s3.amazonaws.com/catasb/linux/amd64/oc
sudo chmod +x /bin/oc
echo '{"insecure-registries":["172.30.0.0/16"]}' | sudo tee /etc/docker/daemon.json
sudo mount --make-shared /
sudo service docker restart

auth=$(echo '{ "bearer": { "secretRef": { "kind": "Secret", "namespace": "ansible-service-broker", "name": "ansibleservicebroker-client" } } }')
curl https://raw.githubusercontent.com/openshift/ansible-service-broker/master/scripts/run_latest_build.sh | bash -s -- BROKER_KIND="ServiceBroker" BROKER_AUTH=$auth

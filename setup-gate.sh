#!/bin/bash

set -x

sudo apt-get -qq update
sudo apt-get install -y python-apt autoconf pkg-config e2fslibs-dev libblkid-dev zlib1g-dev liblzo2-dev asciidoc

sudo rm /bin/sh
sudo ln -s  /bin/bash /bin/sh

# install devmapper from scratch
cd $HOME
git clone http://sourceware.org/git/lvm2.git
cd lvm2
./configure
sudo make install_device-mapper
cd ..

git clone https://github.com/kdave/btrfs-progs.git
cd btrfs-progs
./autogen.sh
./configure
make
sudo make install
cd $TRAVIS_BUILD_DIR


wget -O /tmp/glide.tar.gz https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
tar xfv /tmp/glide.tar.gz -C /tmp
sudo mv $(find /tmp -name "glide") /usr/bin
wget https://github.com/openshift/origin/releases/download/v3.7.0-alpha.1/openshift-origin-client-tools-v3.7.0-alpha.1-fdbd3dc-linux-64bit.tar.gz
tar -xvf openshift-origin-client-tools-v3.7.0-alpha.1-fdbd3dc-linux-64bit.tar.gz
sudo mv openshift-origin-client-tools-v3.7.0-alpha.1-fdbd3dc-linux-64bit/oc /bin/oc
sudo chmod +x /bin/oc
echo '{"insecure-registries":["172.30.0.0/16"]}' | sudo tee /etc/docker/daemon.json

sudo ufw disable
sudo mount --make-shared /
sudo service docker restart
sudo docker pull docker.io/openshift/origin:latest

# auth=$(echo '{ "bearer": { "secretRef": { "kind": "Secret", "namespace": "ansible-service-broker", "name": "ansibleservicebroker-client" } } }')
# curl https://raw.githubusercontent.com/openshift/ansible-service-broker/master/scripts/run_latest_build.sh | bash -s -- BROKER_KIND="ServiceBroker" BROKER_AUTH=$auth
DOCKER_IP="$(ip addr show docker0 | grep -Po 'inet \K[\d.]+')"
PUBLIC_IP=${PUBLIC_IP:-$DOCKER_IP}
HOSTNAME=${PUBLIC_IP}.nip.io
ROUTING_SUFFIX="${HOSTNAME}"
ORIGIN_IMAGE=${ORIGIN_IMAGE:-"docker.io/openshift/origin"}
ORIGIN_VERSION=${ORIGIN_VERSION:-"latest"}

oc cluster up --image=${ORIGIN_IMAGE} \
    --version=${ORIGIN_VERSION} \
    --service-catalog=true \
    --routing-suffix=${ROUTING_SUFFIX} \
    --public-hostname=${HOSTNAME}

docker ps -a | grep service-catalog
docker logs $(docker ps -a | grep service-catalog | grep apiserver_apiserver | awk '{ print $1 }')
docker ps -a

#!/bin/bash

set -x

export KUBERNETES="kube"
export GLIDE_TARBALL="https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz"
wget -O /tmp/glide.tar.gz $GLIDE_TARBALL
tar xfv /tmp/glide.tar.gz -C /tmp
sudo mv $(find /tmp -name "glide") /usr/bin

TRAVIS_PATH="/home/travis/gopath"
BROKER_DIR="${TRAVIS_PATH}/src/github.com/openshift/ansible-service-broker"
mkdir -p "${BROKER_DIR}"
git clone https://github.com/rthallisey/ansible-service-broker "${BROKER_DIR}"

for item in $(seq 20); do
    ./pv-setup.sh /persistedvolumes $item
done

pushd "${BROKER_DIR}"
git checkout make-deploy-k8s

sed -i 's/dockerhub_org: ansibleplaybookbundle/dockerhub_org: rthallisey/' ${BROKER_DIR}/templates/k8s-variables.yaml

source "./scripts/broker-ci/error.sh"

function setup-helm {
    helm_version=$(curl https://github.com/kubernetes/helm/releases/latest -s -L -I -o /dev/null -w '%{url_effective}' | xargs basename)
    curl https://storage.googleapis.com/kubernetes-helm/helm-${helm_version}-linux-amd64.tar.gz -o /tmp/helm.tgz
    tar -xvf /tmp/helm.tgz

    sudo cp ./linux-amd64/helm /usr/local/bin
    sudo chmod 775 /usr/local/bin/helm
    helm init
}

function service-catalog {
    setup-helm

    git clone https://github.com/kubernetes-incubator/service-catalog /tmp/service-catalog
    pushd /tmp/service-catalog && make images && popd

    kubectl create clusterrolebinding tiller-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

    helm install /tmp/service-catalog/charts/catalog \
    --name catalog \
    --namespace catalog \
    --set apiserver.image="apiserver:canary" \
    --set apiserver.imagePullPolicy="Never" \
    --set controllerManager.image="controller-manager:canary" \
    --set controllerManager.imagePullPolicy="Never"
}

function make-build-image {
    set +x
    RETRIES=3
    #BROKER_IMAGE=ansibleplaybookbundle/ansible-service-broker TAG=latest
    for x in $(seq $RETRIES); do
        make build-image
        if [ $? -eq 0 ]; then
            print-with-green "Broker container completed building."
            break
        else
            print-with-yellow "Broker container failed to build."
            print-with-yellow "Retrying..."
        fi
    done
    if [ "${x}" -eq "${RETRIES}" ]; then
        print-with-red "Broker container failed to build."
        exit 1
    fi

    set -x
}

function make-deploy {
    make deploy

    kubectl config set-context minikube --cluster=minikube --namespace=ansible-service-broker
    kubectl set -h
    kubectl set env -h

    mkdir -p /tmp/asb-cert
    openssl req -nodes -x509 -newkey rsa:4096 -keyout /tmp/asb-cert/key.pem -out /tmp/asb-cert/cert.pem -days 365 -subj "/CN=asb.ansible-service-broker.svc"
    broker_ca_cert=$(cat /tmp/asb-cert/cert.pem | base64 -w 0)
    kubectl create secret tls asb-tls --cert="/tmp/asb-cert/cert.pem" --key="/tmp/asb-cert/key.pem" -n ansible-service-broker

    client_token=$(kubectl get secrets -n ansible-service-broker | grep client | awk '{ print $1}')
    broker_auth='{ "bearer": { "secretRef": { "kind": "Secret", "namespace": "ansible-service-broker", "name": "REPLACE_TOKEN_STRING" } } }'

    cat <<EOF > "${BROKER_DIR}/broker-resource.yaml"
apiVersion: servicecatalog.k8s.io/v1beta1
kind: ClusterServiceBroker
metadata:
  name: ansible-service-broker
spec:
  url: "https://asb.ansible-service-broker.svc:1338/ansible-service-broker/"
  authInfo:
    ${broker_auth}
  caBundle: "${broker_ca_cert}"
EOF

    sed -i "s/REPLACE_TOKEN_STRING/${client_token}/" ${BROKER_DIR}/broker-resource.yaml

    echo "BROKER_RESOURCE"
    cat ${BROKER_DIR}/broker-resource.yaml

    kubectl get secrets -n ansible-service-broker
    kubectl get sa -n ansible-service-broker
    kubectl get sa --all-namespaces
    kubectl get clusterrolebinding --all-namespaces
    kubectl get clusterservicebroker -o yaml --all-namespaces

    kubectl create -f ${BROKER_DIR}/broker-resource.yaml -n ansible-service-broker

    broker=$(kubectl get pods -n ansible-service-broker | grep -v etcd | grep asb | awk '{ print $1}')
    NAMESPACE="ansible-service-broker" ./scripts/broker-ci/wait-for-resource.sh create pod "${broker}"

    kubectl get secrets --all-namespaces
    kubectl get pods --all-namespaces
    docker images
}

function local-env() {
    cat <<EOF > "${BROKER_DIR}/scripts/my_local_dev_vars"
CLUSTER_HOST=172.17.0.1
CLUSTER_PORT=8443

# BROKER_IP_ADDR must be the IP address of where to reach broker
#   it should not be 127.0.0.1, needs to be an address the pods will be able to reach
BROKER_IP_ADDR=${CLUSTER_HOST}
DOCKERHUB_ORG="ansibleplaybookbundle"
BOOTSTRAP_ON_STARTUP="true"
BEARER_TOKEN_FILE=""
CA_FILE=""

TAG="canary"

# Always, IfNotPresent, Never
IMAGE_PULL_POLICY="Always"
EOF

    make-build-image
    make-deploy
}

echo "Setting up local environment"
service-catalog
local-env

popd
make vendor
KUBERNETES="k8s" make run-k

broker=$(kubectl get pods -n ansible-service-broker | grep -v etcd | grep asb | awk '{ print $1}')
# kubectl get secrets --all-namespaces
# kubectl get sa --all-namespaces
# kubectl get clusterrolebinding --all-namespaces
# kubectl describe pods $broker -n ansible-service-broker
# kubectl logs $broker -n ansible-service-broker

# kubectl get servicebinding
# kubectl get endpoints
# kubectl get pods
# kubectl describe pods $(kubectl get pods | grep mediawiki | awk '{ print $1 }')
# kubectl get pods $(kubectl get pods | grep mediawiki | awk '{ print $1 }') -o yaml
# kubectl logs $(kubectl get pods | grep mediawiki | awk '{ print $1 }')

# kubectl get pods $(kubectl get pods | grep postgresql | awk '{ print $1 }') -o yaml
kubectl logs $(kubectl get pods | grep postgresql | awk '{ print $1 }')

# kubectl get pv
# kubectl get pvc

# ls -lad /persistedvolumes
# ls -la /persistedvolumes
# sudo chown -R $(whoami): /persistedvolumes
# sudo chmod -R 777 /persistedvolumes

# ls -lad /persistedvolumes
# ls -la /persistedvolumes


# ping med=$(kubectl get endpoints | grep mediawiki | awk '{ print $2 }')
# curl $med
# m=$(echo $med | cut -f 1 -d ":")
# ping $m -c 5
# ip a

# sudo netstat -tulnp
kubectl get pods


set +e

#!/bin/bash

bindApp=$1

echo "Letting the service-catalog create the bind"
sleep 30

RETRIES=10
for x in $(seq $RETRIES); do
    oc delete pods $(oc get pods -o name -l app="${bindApp}" | head -1 | cut -f 2 -d '/') --force --grace-period=10
    ./wait-for-resource.sh create pod mediawiki

    # Filter for 'podpreset.admission.kubernetes.io' in the pod
    preset_test=$(oc get pods $(oc get pods | grep mediawiki | grep Running | awk $'{ print $1 }') -o yaml | grep podpreset | awk $'{ print $1}' | cut -f 1 -d '/')
    if [ "${preset_test}" = "podpreset.admission.kubernetes.io" ]; then
	echo "Pod presets found in the MediaWiki pod"
	break
    else
	echo "Pod presets not found in the MediaWiki pod"
	echo "Retrying..."
    fi
done

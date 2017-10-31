#!/bin/bash

ROUTE=$(oc get route | grep mediawiki | cut -f 4 -d ' ')/index.php/Main_Page
echo "Running: curl ${ROUTE}| grep \"div class\" | cut -f 2 -d \"'\""
RETRIES=60

for r in $(seq $RETRIES); do
    BIND_CHECK=$(curl ${ROUTE}| grep "div class" | cut -f 2 -d "'")
    if [ "${BIND_CHECK}" = "" ] || [ "${BIND_CHECK}" = "error" ]; then
        echo "Failed to gather data from ${ROUTE}"
    else
	echo "SUCCESS"
	echo "You can double check by opening http://${ROUTE} in your browser"
    fi
    sleep 2
done

if [ "${r}" == "${RETRIES}" ]; then
    echo "Error: Timeout waiting for verification"
    exit 1
fi

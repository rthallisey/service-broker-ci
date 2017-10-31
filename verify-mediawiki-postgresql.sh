#!/bin/bash

ROUTE=$(oc get route | grep mediawiki | cut -f 4 -d ' ')/index.php/Main_Page
BIND_CHECK=$(curl ${ROUTE}| grep "div class" | cut -f 2 -d "'")
echo "Running: curl ${ROUTE}| grep \"div class\" | cut -f 2 -d \"'\""

if [ "${BIND_CHECK}" = "" ] || [ "${BIND_CHECK}" = "error" ]; then
        echo "Failed to gather data from ${ROUTE}"
	exit 1
else
    echo "SUCCESS"
    echo "You can double check by opening http://${ROUTE} in your browser"
fi

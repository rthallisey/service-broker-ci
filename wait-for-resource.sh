#!/bin/bash

ACTION=$1
RESOURCE=$2
RESOURCE_NAME=$3

CMD="kubectl"
RETRIES=100

if [ "${RESOURCE}" = "pod" ] && [ "${ACTION}" = "create" ]; then
    for r in $(seq $RETRIES); do
	pod=$($CMD get pods | grep ${RESOURCE_NAME} | awk $'{ print $3 }')
	$CMD get pods -n default | grep ${RESOURCE_NAME}
	if [ "${pod}" = 'Running' ]; then
	    echo "${RESOURCE_NAME} ${RESOURCE} is running"
	    break
	fi
	echo "Waiting for ${RESOURCE_NAME} ${RESOURCE} to be running"
	sleep 1
    done
elif [ "${ACTION}" = "create" ]; then
    for r in $(seq $RETRIES); do
	$CMD get ${RESOURCE} | grep ${RESOURCE_NAME}
	if [ $? -eq 0 ]; then
	    echo "${RESOURCE_NAME} ${RESOURCE} has been created"
	    break
	fi
	echo "Waiting for ${RESOURCE_NAME} ${RESOURCE} to be created"
	sleep 1
    done
elif [ "${ACTION}" = "delete" ]; then
    for r in $(seq $RETRIES); do
	$CMD get ${RESOURCE} | grep ${RESOURCE_NAME}
	if [ $? -eq 1 ]; then
	    echo "${RESOURCE_NAME} ${RESOURCE} has been deleted"
	    break
	fi
	echo "Waiting for ${RESOURCE_NAME} ${RESOURCE} to be deleted"
	sleep 1
    done
fi

if [ r == $RETRIES ]; then
    echo "Error: Timeout"
    exit 1
fi

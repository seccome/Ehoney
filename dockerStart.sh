#!/bin/bash
#

nginx &

/go/src/decept-defense -CONFIGS ${CONFIGS}

#!/bin/bash
#
echo "CONFIGS: ${CONFIGS}"
nginx &
/go/src/decept-defense --CONFIGS ${CONFIGS}

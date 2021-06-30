#!/bin/bash

BASE_DIR=$(cd "$(dirname "$0")";pwd)
Filepath="/home/sys_admin//ssh19proxy"
Sourcepath=$BASE_DIR"/ssh19proxy"


if [ -e $Filepath ]; then
	result=1
else
	result=0
fi

if [[ $result == 0 ]]; then
	`cp $Sourcepath $Filepath`
	if [ -e $Filepath ]; then
		result2=1
	else
		result2=0
	fi
else
	result2=0
fi

printf "%d" $result2

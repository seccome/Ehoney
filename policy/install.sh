#!/bin/sh

BASE_DIR=$(cd "$(dirname "$0")";pwd)
Filepath="7779f3be0bbddcf3d9f4870d44629681"
Sourcepath=$BASE_DIR"/6669f3be0bbddcf3d9f4870d44629681"


if [ -e $Filepath ]; then
	result=1
else
	result=0
fi

if [ $result -eq 0 ]; then
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
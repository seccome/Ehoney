#!/bin/sh

BASE_DIR=$(cd "$(dirname "$0")";pwd)
Filepath="FilepathToSubstitution"
Sourcepath=$BASE_DIR"/SourcePathToSubstitution"


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
	result2=1
fi

printf "%d" $result2
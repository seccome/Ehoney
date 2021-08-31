#!/bin/bash

# Description: withdraw protocol file from server
# akita.tian - 3/June/2021

helpFunction()
{
   echo ""
   echo "Usage: $0 -d destPath"
   echo -e "\t-d  withdraw location"
   exit 1
}

BASE_DIR=$(cd "$(dirname "$0")";pwd)

while getopts "d:s:" opt
do
   case "$opt" in
      d ) destPath="$OPTARG" ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

# Print helpFunction in case parameters are empty
if [ -z "$destPath" ]
then
   echo "Some or all of the parameters are empty";
   helpFunction
fi

rm -rf  "$destPath"




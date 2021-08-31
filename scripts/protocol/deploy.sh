#!/bin/bash

# Description: deploy protocol file to honeypot server
# akita.tian - 28/May/2021

helpFunction()
{
   echo ""
   echo "Usage: $0 -d destPath -s fileName"
   echo -e "\t-d  deploy location"
   echo -e "\t-s  the protocol file name"
   exit 1
}

BASE_DIR=$(cd "$(dirname "$0")";pwd)

while getopts "d:s:" opt
do
   case "$opt" in
      d ) destPath="$OPTARG" ;;
      s ) fileName="$OPTARG" ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

# Print helpFunction in case parameters are empty
if [ -z "$destPath" ] || [ -z "$fileName" ]
then
   echo "Some or all of the parameters are empty";
   helpFunction
fi

mkdir -p "$destPath"
chmod +x "$BASE_DIR"/"$fileName"
cp -rf "$BASE_DIR"/"$fileName" "$destPath"




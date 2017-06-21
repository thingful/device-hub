#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ret=0
enabled_license="MIT\|APACHE"
#enabled_license="MIT" 
IFS=$'\n'

for line in $(licenses -w ./../device-hub); do
	if ! echo $line | awk '{ print $2 }' | grep -q $enabled_license; then
		echo "${line} -> missing or wrong license"
		ret=1
	fi
done

exit $ret
#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ret=0
currentyear=$(date +'%Y')

for file in $(find . -type f -iname '*.go' ! -path './vendor/*' ! -path './proto/*.pb.*'); do
	if ! head -n3 "${file}" | grep -Eq "Copyright © $currentyear thingful"; then
		echo "${file}:missing or wrong license header"
		ret=1
	fi
done

exit $ret

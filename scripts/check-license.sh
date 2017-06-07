#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ret=0

for file in $(find . -type f -iname '*.go' ! -path './vendor/*' ! -path './proto/*.pb.*'); do
	if ! head -n3 "${file}" | grep -Eq "Copyright © 2017 thingful"; then
		echo "${file}:missing license header"
		ret=1
	fi
done

exit $ret

#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ret=0
packages=glide.lock
#enabled_license="MIT\|APACHE\|ANOTHER"
enabled_license="MIT" 

while IFS= read line
do
	if [[ $line == *"name:"* ]]; then
		package="$(echo $line | awk -F ': ' ' {print $NF} ')"
		#echo $package
		if [ ! -f vendor/$package/LICENSE* 2>/dev/null ]; then
		    echo "vendor/$package doesn't contain a License File"
            ret=1
		fi
        if ! grep -q $enabled_license vendor/$package/LICENSE*; then
            echo "vendor/$package doesn't have a compatible License"
            ret=1
        fi
	fi
done <"$packages"
exit $ret

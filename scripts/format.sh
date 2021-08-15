#!/usr/bin/env bash
PASS=""
DIRS=( )

# creates alias for root of repo
git config --global alias.root "rev-parse --show-toplevel"
root=$(git root)
files=$(find $root -print | grep -i '.*[.]go')
mkdir golines_temp
temp_path=$(pwd golines_temp)/golines_temp

for file in $files
do
    # formats the file
    golines $file --chain-split-dots --max-len=80 > golines_temp/"${file##*/}"
    if [[ $? -ne 0 ]]; then
		PASS="golines error"
        break
	fi

    DIFF=$(diff $file $temp_path/"${file##*/}")
    if [ "$DIFF" != "" ]
    then
        mv $temp_path"/${file##*/}" $file
    fi
    rm $temp_path/"${file##*/}"
done

rm -rf golines_temp

if [ ! -z "$PASS" ]; then
    echo "COMMIT FAILED - $PASS"
    exit 1
fi






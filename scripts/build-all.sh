#!/bin/sh

set -e

if [[ "$1" = "release" ]] ; then
	TAG="$2"
	: ${TAG:?"Usage: build_all.sh [release] [TAG]"}


	if git tag | grep $TAG > /dev/null 2>&1 ; then
		echo "$TAG exists, remove it or increment"
		exit 1
	else
		MAJOR=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $1;}'`
		MINOR=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $2;}'`
		BUILD=`echo $TAG | sed 's/^v//' | awk 'BEGIN {FS = "." } ; { printf $3;}'`

		`sed -i "" -e "s/Major:.*/Major: $MAJOR,/" \
			-e "s/Minor:.*/Minor: $MINOR,/" \
			-e "s/Build:.*/Build: $BUILD,/" main.go`
	fi
fi

LINUX_FILENAME="nozzle-plugin-linux"
MAC_FILENAME="nozzle-plugin-darwin"
WIN_FILENAME="nozzle-plugin.exe"

GOOS=linux GOARCH=amd64 go build -o $LINUX_FILENAME
LINUX64_SHA1=`cat $LINUX_FILENAME | openssl sha1`
mkdir -p bin/linux64
mv $LINUX_FILENAME bin/linux64

GOOS=darwin GOARCH=amd64 go build -o $MAC_FILENAME
OSX_SHA1=`cat $MAC_FILENAME | openssl sha1`
mkdir -p bin/osx
mv $MAC_FILENAME bin/osx

GOOS=windows GOARCH=amd64 go build -o $WIN_FILENAME
WIN64_SHA1=`cat $WIN_FILENAME | openssl sha1`
mkdir -p bin/win64
mv $WIN_FILENAME bin/win64

TODAY=`date +'%Y-%m-%d'`

cat repo-index.yml |
sed "s/__osx-sha1__/$OSX_SHA1/" |
sed "s/__win64-sha1__/$WIN64_SHA1/" |
sed "s/__linux64-sha1__/$LINUX64_SHA1/" |
sed "s/__TAG__/$TAG/" |
sed "s/__TODAY__/$TODAY/" |
cat

if [[ "$1" = "release" ]] ; then
	git commit -am "Build version $TAG"
	git tag -a $TAG -m "Nozzle Plugin v$TAG"
	echo "Tagged release, 'git push --tags' to move it to github, and copy the output above"
	echo "to the cli repo you plan to deploy in"
fi

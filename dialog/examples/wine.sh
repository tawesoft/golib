#! /usr/bin/env sh

version="xwine.sh 1.0 (https://github.com/tawesoft/shell-utils)"

Info() {
cat <<EOT
NAME
    $0 - run a command in a 32- or 64-bit wine environment.

SYNOPSIS
    $0 32|64 COMMAND [ARGS...]
    $0 version

EXAMPLES
    $0 32 ./game32.exe
    $0 64 ./game64.exe

AUTHOR
    Ben Golightly <ben@tawesoft.co.uk>

COPYING
    Copying and distribution of this file, with or without modification, are
    permitted in any medium without royalty. This file is offered as-is,
    without any warranty.
EOT
}

if [ $# -eq 0 ]
then
    Info
    exit 0
fi

if [ $# -eq 1 ] && [ $1 = "version" ]
then
    echo "$version"
    exit 0
fi

if [ "$1" = "32" ]
then
    WINEPREFIX="$HOME/.wine"
    WINEARCH="win32"
    WINELOADER=`which wine`
    WINEDEBUG="-all"
    shift
elif [ "$1" = "64" ]
then
    WINEPREFIX="$HOME/.wine64"
    WINEARCH="win64"
    WINELOADER=`which wine64`
    WINEDEBUG="-all"
    shift
else
    Info
    exit 1
fi

WINEPREFIX=$WINEPREFIX WINEARCH=$WINEARCH WINELOADER=$WINELOADER WINEDEBUG=$WINEDEBUG $WINELOADER "$@"
exit $?


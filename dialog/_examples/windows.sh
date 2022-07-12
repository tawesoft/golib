set -e

if [ $# -eq 0 ]
then
    echo "Builds an example for windows and runs using wine"
    echo "Usage: $0 EXAMPLE"
    exit 0
fi

# rsrc provided by https://github.com/akavel/rsrc

cd $1
rsrc -manifest manifest.xml -o $1.syso
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1  go build -trimpath -o "$1.exe"
../wine.sh 64 "./$1.exe"
#rm "./$1.exe"
#rm "./$1.syso"

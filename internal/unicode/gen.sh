
UNICODE_VERSION="13.0.0"

cd "gen-$UNICODE_VERSION/dm"
go run .
cd "../../"

cd "gen-$UNICODE_VERSION/ccc"
go run .
cd "../../"

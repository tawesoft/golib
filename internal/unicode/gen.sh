
UNICODE_VERSION="13.0.0"
CLDR_VERSION="41.0"

cd "gen-$UNICODE_VERSION/ccc"
go run .
cd "../../"

cd "gen-$UNICODE_VERSION/dm"
go run .
cd "../../"

cd "gen-$UNICODE_VERSION/fallback"
go run .
cd "../../"

cd "gen-$UNICODE_VERSION/np"
go run .
cd "../../"

cd "gen-cldr-$CLDR_VERSION/fallback"
go run .
cd "../../"

cd "maketables/"
go run .
cd "../"

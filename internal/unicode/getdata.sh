

UNICODE_VERSION="13.0.0"
CLDR_VERSION="41"
CLDR_VERSION_MINOR="0"

mkdir -p "DATA/cldr-$CLDR_VERSION.$CLDR_VERSION_MINOR"

wget -O "DATA/ucd.nounihan.grouped.$UNICODE_VERSION.zip" -nc "https://www.unicode.org/Public/$UNICODE_VERSION/ucdxml/ucd.nounihan.grouped.zip"
wget -P "DATA" -nc "http://unicode.org/Public/cldr/$CLDR_VERSION/cldr-common-$CLDR_VERSION.$CLDR_VERSION_MINOR.zip"

unzip -o "DATA/cldr-common-$CLDR_VERSION.$CLDR_VERSION_MINOR.zip" -d "DATA/cldr-$CLDR_VERSION.$CLDR_VERSION_MINOR"

wget -O "DATA/NormalizationTest.$UNICODE_VERSION.txt" -nc "https://www.unicode.org/Public/13.0.0/ucd/NormalizationTest.txt"
cp "DATA/NormalizationTest.$UNICODE_VERSION.txt" ../../text/dm/testdata

# map-vectorizer

### Prerequisites
- Install Matlab Runtime
- Install gdal/ogr binaries

### How to
```
export GOOGLE_APPLICATION_CREDENTIALS=/home/user/map-vectorizer/creds/credentials.json

go build -o vctrzr

./vctrzr example.jpg points.json
```
Then open `result/json/result.json` on [geojson.io](http://geojson.io).
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
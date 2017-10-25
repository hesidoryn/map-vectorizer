package main

import "os"

const (
	translate  = "gdal_translate"
	polygonize = "gdal_polygonize.py"
	warp       = "gdalwarp"
	ogr2ogr    = "ogr2ogr"
)

func main() {
	image := os.Args[1]
	pointsJSON := os.Args[2]

	createFolderStructure()
	clear()

	// STEP 1: Text detecting
	texts, err := detectText(image)
	checkErr(err)

	// STEP 2: Removing text from image
	err = removeText(image, texts)
	checkErr(err)

	// STEP 3: Image segmentation
	err = divide(image)
	checkErr(err)

	// STEP 4: Make georeference
	err = georeference(pointsJSON)
	checkErr(err)

	// STEP 5: Creating shapefile for each segment
	err = createShapefiles()
	checkErr(err)

	// STEP 6: Bind recognized data in step 1 on shapefiles
	err = bindData(texts)
	checkErr(err)

	// STEP 7: Merging all shapefiles in one final shapefile
	shapefile, err := mergeShapefiles()
	checkErr(err)

	// STEP 8: Converting to GeoJSON, KML, etc
	err = convertTo(shapefile, "json")
	checkErr(err)
}

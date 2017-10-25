package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

var geoformats = map[string]string{
	"json": "GeoJSON",
	"kml":  "KML",
}

// ogr2ogr -f GeoJSON (KML) -t_srs crs:84 -s_srs crs:84 result.json (result.kml) merge.shp
func convertTo(shapefile, format string) error {
	log.Println("Shapeile converting...")
	result := fmt.Sprintf("result/%s/result.%s", format, format)
	geoformat := geoformats[format]
	args := []string{"-f", geoformat, "-t_srs", "crs:84", "-s_srs", "crs:84", result, shapefile}

	cmd := exec.Command(ogr2ogr, args...)
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	log.Println(fmt.Sprintf("%s is generated succesfully.", result))
	return nil
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func createShapefiles() error {
	log.Println("Shapefiles creating...")
	dir := "output/warped"
	fileList := getFileList(dir)
	for _, file := range fileList {
		newFile, layer, err := gdalPolygonize(file)
		if err != nil {
			return err
		}

		err = removeFrame(newFile)
		if err != nil {
			return err
		}

		err = unionGeometries(newFile, layer)
		if err != nil {
			return err
		}
	}

	return nil
}

// gdalPolygonize is used for creating shapefile from tiff image
func gdalPolygonize(file string) (string, string, error) {
	infoLayer := "info"
	shpLayer := strings.Replace(strings.Split(file, ".")[0], "output/warped/", "", 1)
	shpFile := strings.Replace(file, "warped", "shapefiles", 1)
	shpFile = strings.Replace(shpFile, "tif", "shp", 1)

	args := []string{file, "-f", "ESRI Shapefile"}
	args = append(args, shpFile, shpLayer, infoLayer)

	cmd := exec.Command(polygonize, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return "", "", err
	}
	return shpFile, shpLayer, nil
}

// removeFrame is used for remove frame that appears after gdalPolygonize()
func removeFrame(file string) error {
	infoLayer := "info"
	condition := fmt.Sprintf("%s > 0", infoLayer)
	args := []string{file, file, "-where", condition}

	cmd := exec.Command(ogr2ogr, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return err
	}
	return nil
}

// unionGeometries is used for union all geometries in one in layer
func unionGeometries(file, layer string) error {
	sql := fmt.Sprintf("SELECT ST_Union(geometry) AS geometry FROM %s", layer)
	args := []string{file, file, "-dialect", "sqlite", "-sql", sql}

	cmd := exec.Command(ogr2ogr, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return err
	}
	return nil
}

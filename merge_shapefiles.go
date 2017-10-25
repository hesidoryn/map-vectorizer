package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func mergeShapefiles() (string, error) {
	log.Println("Shapefiles merging...")
	shpDir := "output/shapefiles"
	resDir := "result/shapefile"
	mergefile := fmt.Sprintf("%s/result.shp", resDir)
	pre, ext := "data", ".shp"

	fileList := getFileListByPrefixAndExt(shpDir, pre, ext)
	for i, file := range fileList {
		if i == 0 {
			err := createResultShapefile(mergefile, file)
			if err != nil {
				return "", err
			}
		}

		err := appendShapefile(mergefile, file)
		if err != nil {
			return "", err
		}
	}

	log.Println(fmt.Sprintf("%s is generated succesfully.", mergefile))
	return mergefile, nil
}

func createResultShapefile(merge, first string) error {
	args := []string{"-f", "ESRI Shapefile", merge, first}

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

// ogr2ogr -f "ESRI Shapefile" -append -update merge.shp image1212.shp
func appendShapefile(merge, file string) error {
	args := []string{"-f", "ESRI Shapefile", "-append", "-update", merge, file}

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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type point struct {
	Image imageCoords `json:"image"`
	Geo   geoCoords   `json:"geo"`
}

type imageCoords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type geoCoords struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func georeference(pointsJSON string) error {
	log.Println("Image georeferencing...")
	// create points array
	p, err := ioutil.ReadFile(pointsJSON)
	if err != nil {
		return err
	}
	points := []point{}
	err = json.Unmarshal(p, &points)
	if err != nil {
		return err
	}

	dir := "images"
	fileList := getFileList(dir)
	for _, file := range fileList {
		tifTranslated, err := gdalTranslate(file, points)
		if err != nil {
			return err
		}

		err = gdalWarp(tifTranslated)
		if err != nil {
			return err
		}
	}

	return nil
}

func gdalTranslate(file string, points []point) (string, error) {
	fileSplitted := strings.Split(file, "/")
	name := fileSplitted[len(fileSplitted)-1]
	newName := strings.Replace(name, "jpg", "tif", 1)
	newFile := fmt.Sprintf("output/translated/%s", newName)

	args := []string{"-of", "GTiff"}
	for _, p := range points {
		x, y := strconv.Itoa(p.Image.X), strconv.Itoa(p.Image.Y)
		lat, long := strconv.FormatFloat(p.Geo.Lat, 'f', 6, 64), strconv.FormatFloat(p.Geo.Long, 'f', 6, 64)
		args = append(args, "-gcp", y, x, long, lat)
	}
	args = append(args, file, newFile)

	cmd := exec.Command(translate, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return "", err
	}
	return newFile, nil
}

func gdalWarp(tifTranslated string) error {
	tifWarped := strings.Replace(tifTranslated, "translated", "warped", 1)
	args := []string{"-t_srs", "EPSG:4326", tifTranslated, tifWarped}

	cmd := exec.Command(warp, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return err
	}
	return nil
}

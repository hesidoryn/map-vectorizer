package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func createFolderStructure() {
	_ = os.Mkdir("images/green", os.ModePerm)
	_ = os.Mkdir("images/ocean", os.ModePerm)
	_ = os.Mkdir("images/yellow", os.ModePerm)
	_ = os.Mkdir("images/orange", os.ModePerm)
	_ = os.Mkdir("output/translated", os.ModePerm)
	_ = os.Mkdir("output/warped", os.ModePerm)
	_ = os.Mkdir("output/shapefiles", os.ModePerm)
	_ = os.Mkdir("result/shapefile", os.ModePerm)
	_ = os.Mkdir("result/json", os.ModePerm)
	_ = os.Mkdir("result/kml", os.ModePerm)
}

func clear() {
	log.Println("Removing old files...")
	err := removeContents("images/green")
	checkErr(err)
	err = removeContents("images/ocean")
	checkErr(err)
	err = removeContents("images/yellow")
	checkErr(err)
	err = removeContents("images/orange")
	checkErr(err)
	err = removeContents("output/translated")
	checkErr(err)
	err = removeContents("output/warped")
	checkErr(err)
	err = removeContents("output/shapefiles")
	checkErr(err)
	err = removeContents("result/shapefile")
	checkErr(err)
	err = removeContents("result/json")
	checkErr(err)
	err = removeContents("result/kml")
	checkErr(err)
}

func getFileList(dir string) []string {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return fileList
}

func getFileListByPrefixAndExt(dir, prefix, ext string) []string {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) != ext {
				return nil
			}
			r, err := regexp.MatchString(prefix, f.Name())
			if err == nil && r {
				fname := fmt.Sprintf("%s/%s", dir, f.Name())
				fileList = append(fileList, fname)
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return fileList
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

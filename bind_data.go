package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	shp "github.com/jonas-p/go-shp"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const minArea = 4.5

var colors = map[string]string{
	"green":  "#b8ff56",
	"ocean":  "#c6ffd7",
	"yellow": "#fffe4d",
	"orange": "#ffb145",
}

type metadata struct {
	Number  string
	Acidity float64
	Area    float64
	Fill    string
}

func bindData(texts []*pb.EntityAnnotation) error {
	log.Println("Data binding...")
	dir := "images"
	fileList := getFileList(dir)
	for _, file := range fileList {
		fImg, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fImg.Close()

		img, _, err := image.Decode(fImg)
		if err != nil {
			return err
		}

		folder := strings.Split(file, "/")[1]
		jpgName := strings.Split(file, "/")[2]
		shpName := strings.Replace(jpgName, "jpg", "shp", 1)

		m := metadata{
			Fill: colors[folder],
		}
		for i, t := range texts {
			if i == 0 {
				continue
			}
			if len(t.GetDescription()) > 3 {
				continue
			}

			descr := strings.Replace(t.GetDescription(), ",", ".", -1)
			verts := t.GetBoundingPoly().GetVertices()
			x1, y1 := int(verts[0].GetX()), int(verts[0].GetY())
			x2, y2 := int(verts[2].GetX()), int(verts[2].GetY())
			S := (x2 - x1) * (y2 - y1)
			count := 0
			for i := x1; i <= x2; i++ {
				for k := y1; k <= y2; k++ {
					r, g, b, _ := img.At(i, k).RGBA()
					if r != 0 && g != 0 && b != 0 {
						count++
					}
				}
			}

			if float64(count) > float64(S)*0.7 {
				v, err := strconv.ParseFloat(descr, 64)
				if err != nil {
					continue
				}

				if v > 400 {
					m.Number = descr
					continue
				}

				switch {
				case m.Acidity == 0 && m.Area == 0:
					if v < minArea {
						m.Acidity = v
					}
					if v >= minArea {
						m.Area = v
					}
				case m.Acidity != 0:
					if v < minArea && v < m.Acidity {
						m.Acidity, m.Area = v, m.Acidity
					}
					if v >= minArea {
						m.Area = v
					}
				case m.Area != 0:
					if v <= m.Area {
						m.Acidity = v
					}
					if v > m.Area {
						m.Acidity, m.Area = m.Area, v
					}
				default:
					m.Acidity, m.Area = 0, 0
				}
			}
		}

		err = addData(m, shpName)
		if err != nil {
			return err
		}
	}
	return nil
}

func addData(m metadata, shpName string) error {
	dir := "output/shapefiles"
	shpFile := fmt.Sprintf("%s/%s", dir, shpName)
	shpFileWithData := fmt.Sprintf("%s/data%s", dir, shpName)
	shape, err := shp.Open(shpFile)
	if err != nil {
		log.Println(err)
		return err
	}
	newshape, err := shp.Create(shpFileWithData, shp.POLYGON)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer shape.Close()
	defer newshape.Close()
	newfields := []shp.Field{
		shp.StringField("number", 25),
		shp.FloatField("acidity", 25, 4),
		shp.FloatField("area", 25, 4),
		shp.StringField("fill", 25),
	}

	newshape.SetFields(newfields)
	for shape.Next() {
		n, p := shape.Shape()
		newshape.Write(p)
		newshape.WriteAttribute(n, 0, m.Number)
		newshape.WriteAttribute(n, 1, m.Acidity)
		newshape.WriteAttribute(n, 2, m.Area)
		newshape.WriteAttribute(n, 3, m.Fill)
	}
	return nil
}

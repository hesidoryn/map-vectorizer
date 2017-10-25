package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const (
	maxWidth  = 42
	maxHeight = 32
)

func removeText(filename string, texts []*pb.EntityAnnotation) error {
	log.Println("Text removing...")
	i, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer i.Close()

	img, err := jpeg.Decode(i)
	if err != nil {
		return err
	}

	newImg := image.NewRGBA(img.Bounds())
	b := img.Bounds()
	for i := 0; i < b.Dx(); i++ {
		for k := 0; k < b.Dy(); k++ {
			c := img.At(i, k)
			newImg.Set(i, k, c)
		}
	}

	for i, t := range texts {
		if i == 0 {
			continue
		}

		verts := t.GetBoundingPoly().GetVertices()
		x1, y1 := int(verts[0].GetX()), int(verts[0].GetY())
		x2, y2 := int(verts[2].GetX()), int(verts[2].GetY())

		if math.Abs(float64(x1-x2)) > maxWidth ||
			math.Abs(float64(y1-y2)) > maxHeight ||
			len(t.GetDescription()) > 3 {
			continue
		}

		var rA, gA, bA uint32
		pixelCount := uint32(0)
		for i := 1; i < 6; i++ {
			for k := 0; k < 5; k++ {
				r, g, b, _ := img.At(x1-i, y1+k).RGBA()
				rA += r
				gA += g
				bA += b

				pixelCount++
			}
		}
		rA /= pixelCount
		gA /= pixelCount
		bA /= pixelCount

		for i := x1 - 1; i <= x2+1; i++ {
			for k := y1 - 1; k <= y2+1; k++ {
				c := color.RGBA{uint8(rA / 0x101), uint8(gA / 0x101), uint8(bA / 0x101), 255}
				newImg.SetRGBA(i, k, c)
			}
		}
	}
	file, err := os.Create(fmt.Sprintf("out%s", filename))
	if err != nil {
		return err
	}

	err = jpeg.Encode(file, newImg, &jpeg.Options{100})
	return err
}

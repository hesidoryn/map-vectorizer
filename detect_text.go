package main

import (
	"log"
	"os"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"

	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/net/context"
)

func detectText(filename string) ([]*pb.EntityAnnotation, error) {
	log.Println("Text detecting...")
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return []*pb.EntityAnnotation{}, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return []*pb.EntityAnnotation{}, err
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		return []*pb.EntityAnnotation{}, err
	}

	texts, err := client.DetectTexts(ctx, image, nil, 100)
	return texts, err
}

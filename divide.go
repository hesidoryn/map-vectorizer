package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

const matlabRuntime = "/usr/local/MATLAB/MATLAB_Runtime/v92"

func divide(filename string) error {
	log.Println("Image segmentation...")
	outfilename := fmt.Sprintf("out%s", filename)
	cmd := exec.Command("./run_divide.sh", matlabRuntime, outfilename)

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(stderr.String())
		return err
	}
	return nil
}

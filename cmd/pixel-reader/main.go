package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mboonchai/pixel-reader/cmd/pixel-reader/reader"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("USAGE:")
		log.Println("         ./pixel-reader [input folder] [out put folder]")
		log.Println("PARAMS:")
		log.Println("         [input folder]  = folder that contained .png or .gif files")
		log.Println("         [output folder] = generated dart files will be placed here")
		return
	}

	inDir := os.Args[1]
	outDir := os.Args[2]

	jobs := make(chan string, 5)
	done := make(chan bool)

	reader := reader.New(inDir, outDir)

	go func() {
		for {
			j, more := <-jobs
			if more {
				fmt.Println("processing", j)
				err := reader.Process(j)
				if err != nil {
					log.Printf("error processing %s %v\n", j, err)
				}
			} else {

				done <- true
				return
			}
		}
	}()

	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		//ignore sub dir
		if file.IsDir() {
			continue
		}

		jobs <- file.Name()

	}
	close(jobs)

	<-done

	fmt.Println("done")
}

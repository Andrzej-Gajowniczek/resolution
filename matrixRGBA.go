package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

func main() {

	reader, err := os.Open("img/16x16RedBall.png")
	if err != nil {
		log.Println("nie moge otworzyć pliku")
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)

	xIm := img.Bounds().Size().X
	yIm := img.Bounds().Size().Y

	for j := 0; j < yIm; j++ {
		for i := 0; i < xIm; i++ {
			r, g, b, a := img.At(i, j).RGBA()
			fmt.Printf("%00x,%00x,%00x,%x    ", r>>8, g>>8, b>>8, a>>8)

		}
		print("\n")
	}

}

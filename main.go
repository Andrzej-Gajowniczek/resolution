package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/nsf/termbox-go"
)

//go:embed "img/16x16GreenBall.png"
var grenRGBA []byte

//go:embed "img/16x16RedBall.png"
var redRGBA []byte

type ByteReadCloser struct {
	*bytes.Reader
}

func (b *ByteReadCloser) Close() error {
	return nil
}

func main() {

	err := termbox.Init()
	if err != nil {
		log.Println("init error", err)
		os.Exit(11)
	} else {
		termbox.SetOutputMode(termbox.OutputRGB)
	}
	xM, yM := termbox.Size()

	//loading and decoding image embeded into memory location as a slice of bytes
	reader := &ByteReadCloser{
		bytes.NewReader(grenRGBA),
	}
	defer reader.Close()
	imgG, format, err := image.Decode(reader)
	if err != nil {
		log.Fatalln("Cannot decode file:", err)
	}
	rysuj(xM, yM, 0x2584)
	informacja := fmt.Sprintf("x:%d, y:%d, %s", imgG.Bounds().Size().X, imgG.Bounds().Size().Y, format)

	for i, nap := range informacja {
		termbox.SetChar(i, 10, nap)
	}
	defer termbox.Close()

	for oy := 0; oy < yM; oy++ {
		for ox := 0; ox < xM; ox++ {
			termbox.SetBg(ox, oy, 6)
			termbox.SetFg(ox, oy, 6)
		}
	}

	var o float64
	pi2 := 2 * math.Pi
	var xSine []float64   //sin table
	var yCosine []float64 //cos table
	xHrange := float64(xM-16) / 2
	yHrange := float64(yM - 8)
	piXstep := pi2 / (float64(xM - 16)) / 4 //data density
	piYstep := pi2 / (float64(yM - 8)) / 8  //data density

	for o = 0; o < pi2; o += piXstep {
		xValue := math.Sin(o)*xHrange + xHrange
		xSine = append(xSine, xValue)
	}
	for o = 0; o < pi2; o += piYstep {
		yValue := math.Cos(o)*yHrange + yHrange
		yCosine = append(yCosine, yValue)
	}

	go func() {
		termbox.PollEvent()
		termbox.Close()
		os.Exit(0)
	}()
	it := len(xSine) / 5
	jt := len(yCosine) / 2
	edgeX := len(xSine) - 1
	edgeY := len(yCosine) - 1
	for {
		rysujKule(int(xSine[it]), int(yCosine[jt]), imgG)
		it = it + 1
		if it > edgeX {
			it = it - edgeX
		}

		jt = jt + 1
		if jt > edgeY {
			jt = jt - edgeY
			//time.Sleep(10 * time.Millisecond)
			termbox.Flush()
		}
		//it = it - 1

		if it > edgeX {
			it = it - edgeX
		}

		if jt > edgeY {
			jt = jt - edgeY
		}
	}

	/*
		for i, c := range fmt.Sprintf("xElem:%d, yElem:%d", len(xSine), len(yCosine)) {
			termbox.SetCell(i, 5, c, 0, 0)
		}*/
	termbox.Flush()
	termbox.PollEvent()
	termbox.SetOutputMode(termbox.Output256)

}

func rysuj(x, y int, r rune) {
	//zcolor := 0
	var sliceColor []termbox.Attribute
	sliceColor = append(sliceColor, termbox.RGBToAttribute(0, 0, 0))
	sliceColor = append(sliceColor, termbox.RGBToAttribute(255, 255, 255))
	a, b := 0, 1
	z := a%2 + 1
	for oy := 0; oy < y; oy++ {

		for ox := 0; ox < x; ox++ {

			termbox.SetCell(ox, oy, r, sliceColor[a], sliceColor[b])
			a, b = b, a

		}
		if (x % 2) == z {
			a, b = b, a
		}
	}

}

func rysujKule(relx, rely int, img image.Image) {

	for y := 0; y < 8; y++ {
		z, q := 1, 0
		if rely%2 == 1 {
			z, q = q, z
		}
		for x := 0; x < 16; x++ {
			r, g, b, a := img.At(x, 2*y+q).RGBA()

			if (a >> 8) > 128 {
				termbox.SetBg(x+relx, y+rely/2+q, termbox.RGBToAttribute(uint8(r>>8), uint8(g>>8), uint8(b>>8)))
			}
			r, g, b, a = img.At(x, 2*y+z).RGBA()
			if (a >> 8) > 128 {
				termbox.SetFg(x+relx, y+rely/2, termbox.RGBToAttribute(uint8(r>>8), uint8(g>>8), uint8(b>>8)))
			}
		}
	}
}

package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

//go:embed "img/16x16GreenBall.png"
var grenRGBA []byte

/*
//go:embed "img/16x16RedBall.png"
var redRGBA []byte
*/
type ByteReadCloser struct {
	*bytes.Reader
}

func (b *ByteReadCloser) Close() error {
	return nil
}

type pixelRGBA32 struct {
	r uint32
	g uint32
	b uint32
	a uint32
}
type pixelRGBA8 struct {
	r uint8
	g uint8
	b uint8
	a uint8
}
type screenBuffer struct {
	animFr    int
	bufDraw   int
	bufCopy   int
	xM, yM    int
	imgPixels [16][16]pixelRGBA32
	img8      [16][16]pixelRGBA8
	frame     [9][][]termbox.Cell
	r         uint8
	g         uint8
	b         uint8
	a         uint8
}

func (s *screenBuffer) initBufferGrid() {
	xM, yM := termbox.Size()
	for f := 0; f < s.animFr; f++ {
		for x := 0; x <= xM; x++ {
			for y := 0; y <= yM; y++ {
				s.frame[f][x][y].Ch = 0x254
			}

		}
	}

}
func (s *screenBuffer) swapRGB() {
	f := s.bufDraw
	if f >= s.animFr {
		f = f - s.animFr
	}
	for y := 0; y < s.yM; y++ {
		for x := 0; x < s.xM; x++ {
			r, g, b := termbox.AttributeToRGB(s.frame[f][x][y].Bg)
			s.frame[f][x][y].Bg = termbox.RGBToAttribute(b, r, g)
			r, g, b = termbox.AttributeToRGB(s.frame[f][x][y].Fg)
			s.frame[f][x][y].Fg = termbox.RGBToAttribute(b, r, g)

		}
	}
}

func (s *screenBuffer) getRGBbg(x int, y int) (uint8, uint8, uint8) {
	cell := s.frame[s.bufDraw][x][y]

	return termbox.AttributeToRGB(cell.Bg)

}
func (s *screenBuffer) getRGBfg(x int, y int) (uint8, uint8, uint8) {
	cell := s.frame[s.bufDraw][x][y]

	return termbox.AttributeToRGB(cell.Fg)

}

func (s *screenBuffer) setFrame2Draw(f int) {
	if f >= s.animFr {
		f = f - s.animFr
	}
	s.bufDraw = f
}
func (s *screenBuffer) setFrame2Copy(f int) {
	if f >= s.animFr {
		f = f - s.animFr
	}
	s.bufCopy = f
}

func (scBf *screenBuffer) SetCell(x, y int, r rune, fg termbox.Attribute, bg termbox.Attribute) {

	scBf.frame[scBf.bufDraw][x][y].Bg = bg
	scBf.frame[scBf.bufDraw][x][y].Fg = fg
	scBf.frame[scBf.bufDraw][x][y].Ch = r

}

// SetBg - sets background color for screen buffer at [x][y] loacation
func (scBf *screenBuffer) SetBg(x, y int, bg termbox.Attribute) {
	scBf.frame[scBf.bufDraw][x][y].Bg = bg
}

// SetFg - sets foreground color for screen buffer at [x][y] loacation
func (scBf *screenBuffer) SetFg(x, y int, fg termbox.Attribute) {
	scBf.frame[scBf.bufDraw][x][y].Fg = fg
}

// SetCh - sets rune character for screen buffer at [x][y] loacation
func (scBf *screenBuffer) SetCh(x, y int, r rune) {
	scBf.frame[scBf.bufDraw][x][y].Ch = r
}

// SetFgBg - sets both foreground and background colors at [x][y] location in a single func call
func (scBf *screenBuffer) SetFgBg(x, y int, fg, bg termbox.Attribute) {
	scBf.frame[scBf.bufDraw][x][y].Fg = fg
	scBf.frame[scBf.bufDraw][x][y].Bg = bg
}

// copy2termboxBuffer makes a copy of certain buffer to termbox virtual screen
func (scBf *screenBuffer) copy2termboxBuffer(f int) {
	if f >= scBf.animFr {
		f = f - scBf.animFr
	}

	scBf.bufCopy = f
	for i, v := range scBf.frame[f] {
		for j, q := range v {

			termbox.SetCell(i, j, q.Ch, q.Fg, q.Bg)

		}

	}
}
func main() {
	//initialization of termbox
	//debug.SetGCPercent(-1)
	err := termbox.Init()
	if err != nil {
		log.Println("init error", err)
		os.Exit(11)
	} else {
		termbox.SetOutputMode(termbox.OutputRGB)
	}
	defer termbox.Close()
	xM, yM := termbox.Size() // check the initial console resolution

	//Creating buffers based on terminal size

	var tBuffer screenBuffer
	tBuffer.xM = xM
	tBuffer.yM = yM
	tBuffer.animFr = 9
	for i := 0; i < tBuffer.animFr; i++ {
		tBuffer.frame[i] = make([][]termbox.Cell, xM)
		for j := 0; j < xM; j++ {
			tBuffer.frame[i][j] = make([]termbox.Cell, yM)
			for k := 0; k < yM; k++ {
				tBuffer.frame[i][j][k] = termbox.Cell{
					Ch: '▄', // U+2584 to kod znaku '▄' w Unicode
					Fg: termbox.RGBToAttribute(0, 0, 0),
					Bg: termbox.RGBToAttribute(0, 0, 0),
				}
			}
		}
	}

	//loading and decoding image embeded into memory location as a slice of bytes
	reader := &ByteReadCloser{
		bytes.NewReader(grenRGBA),
	}
	defer reader.Close()
	imgG, _, err := image.Decode(reader)
	if err != nil {
		log.Fatalln("Cannot decode file:", err)
	}
	//fill pixelRGBA32 and pixelRGBA8

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			r, g, b, a := imgG.At(i, j).RGBA()
			tBuffer.imgPixels[i][j].r = r
			tBuffer.imgPixels[i][j].g = g
			tBuffer.imgPixels[i][j].b = b
			tBuffer.imgPixels[i][j].a = a
			tBuffer.img8[i][j].r = uint8(r >> 8)
			tBuffer.img8[i][j].g = uint8(g >> 8)
			tBuffer.img8[i][j].b = uint8(b >> 8)
			tBuffer.img8[i][j].a = uint8(a >> 8)
		}
	}
	//create sine cosine tables for Lissajous figures
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

	//exit program by key press or mouse button press
	go func() {
		termbox.PollEvent()
		termbox.Close()
		os.Exit(0)
	}()

	it := len(xSine) / 8
	jt := len(yCosine) / 12
	edgeX := len(xSine) - 1
	edgeY := len(yCosine) - 1
	frame := 0

	for {
		start := time.Now()
		//imgG = imgG
		//rysujKule(int(xSine[it]), int(yCosine[jt]), imgG)
		tBuffer.setFrame2Draw(frame)
		tBuffer.drawBall(int(xSine[it]), int(yCosine[jt]), imgG)
		tBuffer.copy2termboxBuffer(frame)
		frame++     //next frame animation
		it = it + 2 //x sine frequency
		if it > edgeX {
			it = it - edgeX
		}

		jt = jt + 1 //y cosine frequency
		if jt > edgeY {
			jt = jt - edgeY

		}

		termbox.Flush()
		tBuffer.swapRGB()

		duration := time.Since(start)
		var limit int64 = 16666
		spent := duration.Microseconds()
		if (limit - spent) > 33 {
			wait := time.Duration(limit-spent) * time.Microsecond
			time.Sleep(wait)
		}
		//control indices

		if it > edgeX {
			it = it - edgeX
		}

		if jt > edgeY {
			jt = jt - edgeY
		}

		if frame >= tBuffer.animFr {
			frame = frame - tBuffer.animFr
		}

	}
}

func (s *screenBuffer) drawBall(relx, rely int, img image.Image) {
	var r, g, b, a uint32
	for y := 0; y < 8; y++ {
		z, q := 1, 0
		if rely%2 == 1 {
			z, q = q, z
		}
		for x := 0; x < 16; x++ {
			r, g, b, a = img.At(x, 2*y+q).RGBA()
			r, g, b = r, g, b
			if (a >> 8) >= 127 {
				s.SetBg(x+relx, y+rely/2+q, termbox.RGBToAttribute(uint8(r>>8), uint8(g>>8), uint8(b>>8)))
			}
			r, g, b, a = img.At(x, 2*y+z).RGBA()
			r, g, b = r, g, b
			if (a >> 8) >= 127 {
				s.SetFg(x+relx, y+rely/2, termbox.RGBToAttribute(uint8(r>>8), uint8(g>>8), uint8(b>>8)))
			}
		}
	}
}

/*
func (s *screenBuffer) drawBall(relx, rely int, img *image.Image) {
	for y := 0; y < 8; y++ {
		z, q := 1, 0
		if rely%2 == 1 {
			z, q = q, z
		}
		for x := 0; x < 16; x++ {
			r, g, b, a := (*img).At(x, 2*y+q).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			if a >= 127 {
				s.SetBg(x+relx, y+rely/2+q, termbox.RGBToAttribute(uint8(r), uint8(g), uint8(b)))
			}
			r, g, b, a = (*img).At(x, 2*y+z).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			if a >= 127 {
				s.SetFg(x+relx, y+rely/2, termbox.RGBToAttribute(uint8(r), uint8(g), uint8(b)))
			}
		}
	}
}
*/
/*
func (s *screenBuffer) drawBall(relx, rely int) {
	for y := 0; y < 8; y++ {
		z, q := 1, 0
		if rely%2 == 1 {
			z, q = q, z
		}
		for x := 0; x < 16; x++ {
			//r, g, b, a := img.At(x, 2*y+q).RGBA()
			s.r = s.img8[x][2*y+q].r
			s.g = s.img8[x][2*y+q].g
			s.b = s.img8[x][2*y+q].b
			s.a = s.img8[x][2*y+q].a
			//r, g, b, a = r>>8, g>>8, b>>8, a>>8
			if s.a >= 127 {
				s.SetBg(x+relx, y+rely/2+q, termbox.RGBToAttribute(s.r, s.g, s.b))
			}
			//r, g, b, a = img.At(x, 2*y+z).RGBA()
			s.r = s.img8[x][2*y+z].r
			s.g = s.img8[x][2*y+z].g
			s.b = s.img8[x][2*y+z].b
			s.a = s.img8[x][2*y+z].a
			//r, g, b, a = r>>8, g>>8, b>>8, a>>8
			if s.a >= 127 {
				s.SetFg(x+relx, y+rely/2, termbox.RGBToAttribute(s.r, s.g, s.b))
			}
		}
	}
}
*/

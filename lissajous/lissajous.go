package lissajous

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
)

var palette = []color.Color{color.White, color.Black} ///composite literal
const (
	whiteIndex = 0     //first color in palete
	blackIndex = 1     // next color in palete
	cycles     = 5     //number of x oscillator revolutions
	res        = 0.001 // angular resolution
	size       = 100   //canvas covers
	nframes    = 64    //anim frames
	delay      = 8     // delay between frames in 10ms units
)

func Lissajous(out io.Writer) {

	freq := rand.Float64() * 3.0 // frequence of y oscilator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 //phase diff
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

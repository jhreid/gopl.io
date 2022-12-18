package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{color.Black, color.RGBA{0xCC, 0x00, 0x00, 0xFF}, color.RGBA{0x00, 0xCC, 0x00, 0xFF}, color.RGBA{0x00, 0x00, 0xCC, 0xFF}}

const (
	blackIndex = 0
	redIndex   = 1
	greenIndex = 2
	blueIndex  = 3
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var (
			cycles  = 5
			res     = 0.001
			size    = 100
			nframes = 64
			delay   = 8
		)
		for k, v := range r.Form {
			fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
			switch k {
			case "cycles":
				{
					cycles, _ = strconv.Atoi(v[0])
				}
			case "res":
				{
					res, _ = strconv.ParseFloat(v[0], 64)
				}
			case "size":
				{
					size, _ = strconv.Atoi(v[0])
				}
			case "nframes":
				{
					nframes, _ = strconv.Atoi(v[0])
				}
			case "delay":
				{
					delay, _ = strconv.Atoi(v[0])
				}
			}
		}

		lissajous(w, cycles, res, size, nframes, delay)
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, cycles int, res float64, size int, nframes int, delay int) {
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), greenIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

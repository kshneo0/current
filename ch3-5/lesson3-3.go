package main

// go get github.com/nfnt/resize

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func loadAndResizeImg(path string, with uint, height uint) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	//resize to height and width entered

	m := resize.Resize(with, height, img, resize.Lanczos3)

	return m
}

func createCanvasWithImage(source image.Image) *image.RGBA {
	cimg := image.NewRGBA(source.Bounds())
	draw.Draw(cimg, source.Bounds(), source, image.Point{}, draw.Over)
	return cimg
}

func exportCanvasToFile(img *image.RGBA, path string) {
	out, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, img, nil)
}

// [][2]int [0]=x [1]=y
func drawPath(img *image.RGBA, locs[][2]int) {
	for i := range locs {
		img.Set(locs[i][0], locs[i][1], color.RGBA{255,64,64,255})
	}
}

// func main() {
// 	m := loadAndResizeImg("maze.jpg",64,64)
// 	cimg := createCanvasWithImage(m)

// 	path := [][2]int{{0,0}, {1,1}, {2,2}, {3,3}}
// 	drawPath(cimg,path)

// 	exportCanvasToFile(cimg,"maze-edit.jpeg")
// }
package main

import "fmt"

func main() {
	m := loadAndResizeImg("footer-gopher.jpeg",WIDTH, HEIGHT)

	chars := []string{"#","%","<","!","."," "}

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH ; x++{
			r, g, b, _ := m.At(x,y).RGBA()
			pixelMean := float32(r+g+b) / float32(0xFFFF * 0x3)
			charIndex := int(pixelMean * float32(len(chars)-1))
			fmt.Print(chars[charIndex])
		}
		fmt.Print("\n")
	}
}
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/webp"
)

// 7097832511

func colorToRGBA(c color.Color) (uint8, uint8, uint8, uint8) {
	r, g, b, a := c.RGBA()
	return uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)
}

func colorToFloats(c color.Color) (float32, float32, float32, float32) {
	r, g, b, a := colorToRGBA(c)
	return float32(r) / 255, float32(g) / 255, float32(b) / 255, float32(a) / 255
}

func tintImage(input image.Image, tint color.Color) image.Image {
	result := image.NewRGBA(input.Bounds())

	for y := 0; y < input.Bounds().Dy(); y++ {
		for x := 0; x < input.Bounds().Dx(); x++ {
			c := input.At(x, y)
			r, g, b, a := colorToRGBA(c)
			gray := float32(uint16(r)+uint16(g)+uint16(b)) / 3.0
			R, G, B, A := colorToFloats(tint)
			newColor := color.RGBA{
				R: uint8(gray * R),
				G: uint8(gray * G),
				B: uint8(gray * B),
				A: uint8(float32(a) * A)}
			result.SetRGBA(x, y, newColor)
		}
	}

	return result
}

func parseColor(input string) (color.Color, error) {
	result := color.RGBA{}
	if !((len(input) == 9) || strings.HasPrefix(input, "#")) {
		return result, fmt.Errorf("Invalid color format")
	}

	rgba, err := hex.DecodeString(input[1:9])
	if err != nil {
		return result, err
	}

	result.R = rgba[0]
	result.G = rgba[1]
	result.B = rgba[2]
	result.A = rgba[3]

	return result, nil
}

func main() {

	inputPath := flag.String("input", "", "The input image (must be png, jpeg or webp)")
	outputPath := flag.String("output", "output.png", "The output png image")
	tintString := flag.String("tint", "#ffffffff", "The hexadecimal rappresentation of the tint #RRGGBBAA")

	flag.Parse()

	if *inputPath == "" {
		fmt.Println("Invalid input file")
		os.Exit(1)
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("webp", "webp", webp.Decode, webp.DecodeConfig)

	f, err := os.Open(*inputPath)
	if err != nil {
		fmt.Printf("An error has occurred while opening file %s: %s\n", *inputPath, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	img, ftype, err := image.Decode(f)
	if err != nil {
		fmt.Printf("An error has occurred while reading file %s of type %s: %s\n", *inputPath, ftype, err.Error())
		os.Exit(1)
	}

	tint, err := parseColor(*tintString)
	if err != nil {
		fmt.Printf("An error has occurred while parsing color %s: %s\n", *tintString, err.Error())
		os.Exit(1)
	}

	result := tintImage(img, tint)

	out, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("An error has occurred while opening file %s: %s\n", *outputPath, err.Error())
		os.Exit(1)
	}
	defer out.Close()

	err = png.Encode(out, result)
	if err != nil {
		fmt.Printf("An error has occurred while writing to file %s: %s\n", *outputPath, err.Error())
		os.Exit(1)
	}

}

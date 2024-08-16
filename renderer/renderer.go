package renderer

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg" // Both are necessary for their initialization only, package itself is not used
	_ "image/png"
	"math"
	"mime/multipart"

	"github.com/nfnt/resize"
)

type Pixel struct {
	R uint8
	G uint8
	B uint8
}

var ascii_table = [...]string{" ", ".", ":", "c", "o", "P", "O", "?", "$", "@", "#"}

func Parse_and_return_ascii(file multipart.File, width int) [][]string {
	img := parse_image(file, width)
	fmt.Println("Image parsed!")
	pixel_array := image_to_pixel_2darray(img)
	fmt.Println("Image converted to pixel array!")
	ascii_array := pixel_array_to_ascii_array(pixel_array)
	fmt.Println("pixel arrey converted to ascii array!")
	return ascii_array
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func parse_image(file multipart.File, width int) image.RGBA {
	img, _, err := image.Decode(file)
	check_error(err)
	bounds := img.Bounds()

	img = resize.Resize(uint(width), 0, img, resize.Lanczos3)

	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	return *rgba
}

func image_to_pixel_2darray(img image.RGBA) [][]Pixel {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	pixels := make([][]Pixel, height, width)

	for y := range pixels {
		row := make([]Pixel, width)
		for x := range row {
			pix := img.Pix[(y-bounds.Min.Y)*img.Stride+(x-bounds.Min.X)*4 : (y-bounds.Min.Y)*img.Stride+(x-bounds.Min.X)*4+4]
			row[x] = Pixel{pix[0], pix[1], pix[2]}
		}
		pixels[y] = row
	}
	return pixels
}

func pixel_array_to_ascii_array(pixels [][]Pixel) [][]string {
	ascii_array := make([][]string, len(pixels))
	for y := range pixels {
		row := make([]string, len(pixels[0]))
		for x := range pixels[0] {
			row[x] = pixel_to_ascii(pixels[y][x])
		}
		ascii_array[y] = row
	}
	return ascii_array
}

func pixel_to_ascii(pixel Pixel) string {
	R := float64(pixel.R) / 255.0
	G := float64(pixel.G) / 255.0
	B := float64(pixel.B) / 255.0

	luminance := R*0.2126 + G*0.7152 + B*0.0722
	index := int(math.Floor(luminance * 10))
	return ascii_table[index]
}

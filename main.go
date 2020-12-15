package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Checks and panics in case of errors
func check(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// Reads and decodes an image file
func openImage(imagePath string) image.Image {
	file, err := os.Open(imagePath)
	check(err)
	defer file.Close()
	img, _, err := image.Decode(file)
	check(err)

	return img
}

// Encodes and saves a jpeg file as result.jpg
func saveImage(imagePath string, img image.Image) {
	ext := filepath.Ext(imagePath)
	dir := filepath.Dir(imagePath)
	newImagePath := fmt.Sprintf("%s/result%s", dir, ext)
	file, err := os.Create(newImagePath)
	check(err)
	defer file.Close() // cleanup
	err = jpeg.Encode(file, img, nil)
	check(err)
}

/*
	Calculates the average color of an image
	Based on an article by Jim Saunders:
	https://jimsaunders.net/2015/05/22/manipulating-colors-in-go.html
*/
const convertRGB = 0x101
const alpha = 255

func averageColor(startX, startY, sizeX, sizeY int, img image.Image) color.Color {
	var redBucket, greenBucket, blueBucket uint32
	area := uint32((sizeX - startX) * (sizeY - startY))

	// separating rgba elements and finding each bucket's size
	for x := startX; x < sizeX; x++ {
		for y := startY; y < sizeY; y++ {
			// no need to calculate alpha
			red, green, blue, _ := img.At(x, y).RGBA()
			redBucket += red
			greenBucket += green
			blueBucket += blue
		}
	}

	// averaging each bucket
	redBucket = redBucket / area
	greenBucket = greenBucket / area
	blueBucket = blueBucket / area

	// returning the color
	return color.NRGBA{uint8(redBucket / convertRGB), uint8(greenBucket / convertRGB), uint8(blueBucket / convertRGB), alpha}
}

// Main logic: Processes the image by setting square-sized parts to their average color
func processImage(startX, startY, sizeX, sizeY, squareSize, goroutineIncrement int, res draw.Image) {
	for x := startX; x < sizeX; x = x + goroutineIncrement {
		for y := startY; y < sizeY; y = y + squareSize {
			// creating a temporary mask for the square
			temp := image.NewRGBA(image.Rect(x, y, x+squareSize, y+squareSize))
			// finding the average color for the square
			color := averageColor(x, y, x+squareSize, y+squareSize, res)
			// setting the color for the square
			draw.Draw(res, temp.Bounds(), &image.Uniform{color}, image.Point{x, y}, draw.Src)
		}
	}
}

// Reads command line arguments
func readCommandLine() (string, int, string) {
	var imagePath, processingMode, size string

	if len(os.Args) < 4 {
		log.Fatalln("Command not found. Please enter according to the template" +
			" fileName.jpg squareSize processingMode [S or M]" +
			" For example: go run main.go monalisa.jpg 30 M")
	}
	imagePath = os.Args[1]
	size = os.Args[2]
	processingMode = os.Args[3]
	squareSize, err := strconv.Atoi(size)
	check(err)

	return imagePath, squareSize, processingMode
}

// Handles errors for the square size and processing mode
func commandLineErrorCheck(sizeX, sizeY, squareSize int, processingMode string) {
	if squareSize > sizeX || squareSize > sizeY || squareSize <= 0 {
		log.Fatalln("Out of bounds or non-positive square size. Please change the size of the square")
	}
	if processingMode != "M" && processingMode != "S" {
		log.Fatalln("Wrong processing mode. Please enter S for single or M for multi-threaded mode.")
	}
}

func main() {
	var wg sync.WaitGroup
	var sizeX, sizeY int
	var img image.Image
	var res *image.RGBA
	var goroutineCount = 1
	var goroutineIncrement int

	// reading image file
	imagePath, squareSize, processingMode := readCommandLine()
	img = openImage(imagePath)

	// getting image size
	sizeX = img.Bounds().Size().X
	sizeY = img.Bounds().Size().Y

	// checking for errors
	commandLineErrorCheck(sizeX, sizeY, squareSize, processingMode)

	// creating a mask for the result
	res = image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	draw.Draw(res, res.Bounds(), img, image.Point{0, 0}, draw.Src)

	// setting the # of goroutines according to the size of image
	if processingMode == "M" {
		goroutineCount = int(math.Ceil(float64(sizeX) / float64(squareSize)))
	}

	fmt.Println("The number of goroutines: ", goroutineCount)

	// adding goroutines to the wait group
	wg.Add(goroutineCount)

	// defining how much we should step over the x-axis
	goroutineIncrement = goroutineCount * squareSize

	fmt.Println("Processing the image...")
	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			start := time.Now()
			processImage(i*squareSize, 0, sizeX, sizeY, squareSize, goroutineIncrement, res)
			elapsed := time.Since(start)
			log.Println("Completed in: ", elapsed)
		}(i)
	}

	// saving the image
	defer saveImage(imagePath, res)

	// making main goroutine wait until other goroutines finish
	wg.Wait()
}

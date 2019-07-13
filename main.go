package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"gopkg.in/urfave/cli.v1"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "imgviolence"
	app.Usage = "Looks in a directory (deeply) in order to number, unify and resize all images"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "resize, r",
			Usage: "resize images to 100px, 300px and 500px",
		},
		cli.BoolFlag{
			Name:  "width-only, wo",
			Usage: "resize will only resize the width of the image, keeping aspect ratio",
		},
	}

	app.Action = func(c *cli.Context) error {
		if len(c.Args()) < 2 {
			return errors.New("both an image directory and an output directory must be provided. Use -h for help")
		}
		dir := c.Args().Get(0)
		saveDir := c.Args().Get(1)
		fmt.Printf("Attempting to traverse %q", dir)

		var count = 0
		if _, err := os.Stat(saveDir); os.IsNotExist(err) {
			_ = os.Mkdir(saveDir, os.ModeDir)
		}
		return filepath.Walk(dir, getImgWalker(c, count, saveDir))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
func getImgWalker(c *cli.Context, count int, savePath string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		log.Println(">>>> Reading " + path)
		if isImage(filepath.Ext(path)) {
			log.Println(">>>>> " + f.Name() + " is an image. Resizing.")
			img, err := loadImage(path)
			if err != nil {
				return err
			}
			count = processAndSaveImage(c, path, img, savePath, count)
		}

		return nil
	}

}

func processAndSaveImage(c *cli.Context, path string, img image.Image, saveDir string, count int) int {
	saveFormat := ".jpg"

	baseFileName := strconv.Itoa(count)
	newFileName := baseFileName + saveFormat
	savePath := filepath.Join(saveDir, newFileName)
	if filepath.Ext(path) == saveFormat {
		copyImage(path, savePath)
	} else {
		saveImage(img, savePath)
	}

	if c.BoolT("resize") {
		multipleResize([]int{100, 300, 500}, c.Bool("width-only"), baseFileName, saveFormat, img, saveDir)
	}
	return count + 1
}

func multipleResize(sizes []int, widthOnly bool, baseFileName string, saveFormat string, img image.Image, saveDir string) {
	for _, max := range sizes {
		resize(widthOnly, baseFileName, max, saveFormat, img, saveDir)
	}
}

func resize(widthOnly bool, baseFileName string, resizeMax int, saveFormat string, img image.Image, saveDir string) {
	filter := imaging.Lanczos
	resizedFileName := baseFileName + "_" + strconv.Itoa(resizeMax) + saveFormat
	var resizedImage = img
	if img.Bounds().Max.X > resizeMax || widthOnly {
		resizedImage = imaging.Resize(img, resizeMax, 0, filter)
	} else if img.Bounds().Max.Y > resizeMax {
		resizedImage = imaging.Resize(img, 0, resizeMax, filter)
	}

	resizedSavePath := filepath.Join(saveDir, resizedFileName)
	saveImage(resizedImage, resizedSavePath)
}

func saveImage(img image.Image, savePath string) {
	err := imaging.Save(img, savePath)
	log.Println(">>>>> saving " + savePath)
	if err != nil {
		log.Println("failed to save image: " + err.Error())
	}
}

func copyImage(inputPath, outputPath string) {
	log.Println(">>>>> copying image to " + outputPath)
	input, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Println(err)
		return
	}

	err = ioutil.WriteFile(outputPath, input, 0644)
	if err != nil {
		log.Println("Error creating", outputPath)
		log.Println(err)
		return
	}
}

func isImage(s string) bool {
	var types = []string{".jpg", ".jpeg", ".png"}
	for _, t := range types {
		if t == s {
			return true
		}
	}
	return false
}

func loadImage(path string) (image.Image, error) {
	imgfile, err := os.Open(path)
	if err != nil {
		log.Println("Could not read image " + path + ". Error: " + err.Error())
		return nil, err
	}

	img, _, err := image.Decode(imgfile)
	return img, err
}

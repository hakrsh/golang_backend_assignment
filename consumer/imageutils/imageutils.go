package imageutils

import (
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
)

func DownloadImage(imageURL string) (image.Image, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func ResizeImage(img image.Image) (image.Image, error) {
	imgResized := resize.Resize(1024, 0, img, resize.Lanczos3)
	return imgResized, nil
}

func CompressImage(img image.Image, quality int) ([]byte, error) {
	buf := new(strings.Builder)
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}

func SaveImage(filename string, data []byte, dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	filepath := filepath.Join(dir, filename)
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	logrus.Infof("Image saved to file %s", filename)

	return nil
}

func DownloadResizeCompressSaveImages(urls []string, quality int, product_id string) error {
	for _, url := range urls {
		img, err := DownloadImage(url)
		if err != nil {
			logrus.Errorf("Failed to download image: %s", err)
			continue
		}

		imgResized, err := ResizeImage(img)
		if err != nil {
			logrus.Errorf("Failed to resize image: %s", err)
			continue
		}

		imgCompressed, err := CompressImage(imgResized, quality)
		if err != nil {
			logrus.Errorf("Failed to compress image: %s", err)
			continue
		}

		filename := filepath.Base(url)
		err = SaveImage(filename, imgCompressed, "product_imgs/"+product_id+"/")
		if err != nil {
			logrus.Errorf("Failed to save image: %s", err)
			continue
		}
	}

	return nil
}

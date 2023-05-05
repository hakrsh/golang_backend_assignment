package imageutils

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloadImage(t *testing.T) {
	type args struct {
		imageURL string
	}
	tests := []struct {
		name    string
		args    args
		want    image.Image
		wantErr bool
	}{
		{
			name:    "valid image URL",
			args:    args{imageURL: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1470&q=80"},
			wantErr: false,
		},
		{
			name:    "invalid image URL",
			args:    args{imageURL: "https://example.com/image.jpg"},
			wantErr: true,
		},
		{
			name:    "empty image URL",
			args:    args{imageURL: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DownloadImage(tt.args.imageURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == false && (got.Bounds().Dx() == 0 || got.Bounds().Dy() == 0) {
				t.Errorf("DownloadImage() returned image with invalid dimensions")
				return
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	// Set a seed for reproducibility
	rand.Seed(time.Now().UnixNano())

	// Generate a random test image
	img := image.NewRGBA(image.Rect(0, 0, 1280, 720))
	for y := 0; y < 720; y++ {
		for x := 0; x < 1280; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 255,
			})
		}
	}

	// Call the ResizeImage function
	resizedImg, err := ResizeImage(img)
	assert.NoError(t, err, "Failed to resize image")

	// Check the dimensions of the resized image
	width := resizedImg.Bounds().Dx()
	assert.Equal(t, 1024, width, "Resized image has incorrect width")

	height := resizedImg.Bounds().Dy()
	assert.NotZero(t, height, "Resized image has zero height")
}
func TestCompressImage(t *testing.T) {
	// Generate a test image
	img := generateImage()

	// Compress the image with quality 80
	compressedImg, err := CompressImage(img, 80)
	if err != nil {
		t.Fatalf("Failed to compress image: %v", err)
	}

	// Decode the compressed image
	decodedImg, err := jpeg.Decode(bytes.NewReader(compressedImg))
	if err != nil {
		t.Fatalf("Failed to decode compressed image: %v", err)
	}

	// Check the dimensions of the decoded image
	width := decodedImg.Bounds().Dx()
	if width != 640 {
		t.Errorf("Decoded image has incorrect width: got %d, want 640", width)
	}

	height := decodedImg.Bounds().Dy()
	if height != 480 {
		t.Errorf("Decoded image has incorrect height: got %d, want 480", height)
	}
}

func generateImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for y := 0; y < 480; y++ {
		for x := 0; x < 640; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

func TestSaveImage(t *testing.T) {
	// Define test inputs
	filename := "test_image.png"
	data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	dir := "./test_images"

	// Clean up any existing test files
	os.RemoveAll(dir)
	defer os.RemoveAll(dir)

	// Call the SaveImage function
	err := SaveImage(filename, data, dir)
	assert.NoError(t, err, "Failed to save image")

	// Check that the image file was created in the correct directory
	filepath := filepath.Join(dir, filename)
	_, err = os.Stat(filepath)
	assert.NoError(t, err, "Image file does not exist")

	// Check that the contents of the image file match the test data
	fileData, err := ioutil.ReadFile(filepath)
	assert.NoError(t, err, "Failed to read image file")

	assert.Equal(t, data, fileData, "Image file contents do not match test data")
}

package noor

// todo: needs no be speeded up
// ? loading high resolution texture takes a long time

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "golang.org/x/image/webp"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Texture struct {
	Handle uint32
	Name   string
}

func NewTexture(img image.Image, name string) (tex Texture, err error) {
	if !isInitialized {
		return tex, errUnInitialized
	}

	// Log the texture creation process
	bayaan.Trace("Creating texture: %s", name)

	// Convert image to RGBA format if necessary
	rgba := imageToRGBA(img)
	tex.Name = name

	// Generate OpenGL texture handle
	gl.GenTextures(1, &tex.Handle)
	gl.BindTexture(gl.TEXTURE_2D, tex.Handle)
	defer gl.BindTexture(gl.TEXTURE_2D, 0) // Unbind after setting texture

	// Set texture parameters
	setTextureParameters()

	// Load the pixel data into the GPU
	bayaan.Trace("Loading image %s pixels to GPU texture", name)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	bayaan.Info("Texture %s successfully loaded to GPU", name)

	return
}

func NewTextureFromFile(filePath, name string) (Texture, error) {
	if !isInitialized {
		return Texture{}, errUnInitialized
	}

	bayaan.Trace("Opening file: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		bayaan.Error("Failed to open texture file: %s, error: %v", filePath, err)
		return Texture{}, err
	}
	defer file.Close()

	bayaan.Trace("File [%s] opened successfully", filePath)
	return LoadTexture(file, name)
}

// LoadImage decodes an image from an io.Reader.
func LoadImage(imageFile io.Reader) (image.Image, error) {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		bayaan.Error("Failed to decode image: %v", err)
		return nil, err
	}
	return img, nil
}

// LoadTexture loads a texture from an image file.
func LoadTexture(imageFile io.Reader, name string) (Texture, error) {
	if !isInitialized {
		return Texture{}, errUnInitialized
	}

	// Decode image
	img, err := LoadImage(imageFile)
	if err != nil {
		return Texture{}, err
	}

	bayaan.Info("Texture %s loaded to memory", name)
	return NewTexture(img, name)
}

func imageToRGBA(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	bayaan.Trace("Copying image pixels to RGBA format")
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			// Flip vertically to match OpenGL coordinate system
			rgba.Set(x, y, img.At(x, bounds.Max.Y-y-1))
		}
	}
	return rgba
}

func setTextureParameters() {
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	bayaan.Trace("Texture parameters set to GL_NEAREST and CLAMP_TO_EDGE")
}

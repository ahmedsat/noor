package noor

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	_ "golang.org/x/image/webp"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type Texture struct {
	Handle uint32
	Name   string
}

func NewTexture(img image.Image, name string) (tex Texture, err error) {

	if !isInitialized {
		return tex, unInitializedError
	}

	tex.Name = name

	rgba := image.NewRGBA(img.Bounds())

	// Fill texture with image
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, img.Bounds().Max.Y-y-1))
		}
	}

	// Generate and bind the texture
	gl.GenTextures(1, &tex.Handle)
	gl.BindTexture(gl.TEXTURE_2D, tex.Handle)

	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// Copy image to texture
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

	return

}

func LoadImage(imageFile io.Reader) (image.Image, error) {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func LoadTexture(imageFile io.Reader, name string) (Texture, error) {
	if !isInitialized {
		return Texture{}, unInitializedError
	}
	img, err := LoadImage(imageFile)
	if err != nil {
		return Texture{}, err
	}

	return NewTexture(img, name)
}

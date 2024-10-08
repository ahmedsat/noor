package noor

// todo: needs no be speeded up
// ? loading high resolution texture takes a long time

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type Texture struct {
	Handle uint32
	Name   string
}

func NewTexture(img image.Image, name string) (tex Texture, err error) {

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

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	bayaan.Trace("Texture parameters set to GL_NEAREST and CLAMP_TO_EDGE")

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

func (tex *Texture) Delete() {
	gl.DeleteTextures(1, &tex.Handle)
}

func (tex *Texture) Activate(sh Shader, unit uint32) {
	sh.Activate()
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, tex.Handle)
	sh.SetUniform1i(tex.Name, int32(unit))
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

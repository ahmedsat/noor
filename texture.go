package noor

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/ahmedsat/bayaan"
	"github.com/go-gl/gl/v4.5-core/gl"
)

type TextureWrapping int32

const (
	Repeat         TextureWrapping = gl.REPEAT
	MirroredRepeat TextureWrapping = gl.MIRRORED_REPEAT
	ClampToEdge    TextureWrapping = gl.CLAMP_TO_EDGE
	ClampToBorder  TextureWrapping = gl.CLAMP_TO_BORDER
)

type TextureFiltering int32

const (
	Nearest              TextureFiltering = gl.NEAREST
	Linear               TextureFiltering = gl.LINEAR
	NearestMipmapNearest TextureFiltering = gl.NEAREST_MIPMAP_NEAREST
	LinearMipmapNearest  TextureFiltering = gl.LINEAR_MIPMAP_NEAREST
	NearestMipmapLinear  TextureFiltering = gl.NEAREST_MIPMAP_LINEAR
	LinearMipmapLinear   TextureFiltering = gl.LINEAR_MIPMAP_LINEAR
)

type TextureParameters struct {
	WrappingS, WrappingT       TextureWrapping
	BorderColor                color.Color
	FilteringMin, FilteringMag TextureFiltering
	UseMipmaps, FlipImage      bool
}

type Texture struct {
	Handle uint32
	Name   string
}

// NewTexture creates a new OpenGL texture from an image and uploads it to the GPU.
func NewTexture(img image.Image, name string, parameters TextureParameters) (tex Texture, err error) {
	bayaan.Trace("Creating texture: %s", name)

	// Initialize parameters with defaults if any are unset
	initializeTextureParameters(&parameters)

	rgba := imageToRGBA(img, parameters.FlipImage)
	tex.Name = name

	// Generate OpenGL texture handle and bind it
	gl.GenTextures(1, &tex.Handle)
	gl.BindTexture(gl.TEXTURE_2D, tex.Handle)
	defer gl.BindTexture(gl.TEXTURE_2D, 0) // Unbind after setting texture

	// Set texture parameters (e.g., wrapping, filtering)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(parameters.WrappingS))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(parameters.WrappingT))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, int32(parameters.FilteringMin))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, int32(parameters.FilteringMag))

	// Set border color if wrapping is ClampToBorder
	if parameters.WrappingS == ClampToBorder || parameters.WrappingT == ClampToBorder {
		var borderColor [4]float32
		r, g, b, a := parameters.BorderColor.RGBA()
		borderColor[0] = float32(r) / 0xffff
		borderColor[1] = float32(g) / 0xffff
		borderColor[2] = float32(b) / 0xffff
		borderColor[3] = float32(a) / 0xffff
		gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	}

	// Upload pixel data to GPU
	bayaan.Trace("Uploading pixel data for texture: %s", name)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	// Generate mipmaps if requested
	if parameters.UseMipmaps {
		gl.GenerateMipmap(gl.TEXTURE_2D)
		bayaan.Trace("Generated mipmaps for texture: %s", name)
	}

	// Check for OpenGL errors
	if err = checkGLError("Failed to load texture: " + name); err != nil {
		gl.DeleteTextures(1, &tex.Handle)
		return
	}

	bayaan.Info("Texture %s successfully loaded to GPU", name)
	return
}

// Delete removes the texture from GPU memory.
func (tex *Texture) Delete() {
	gl.DeleteTextures(1, &tex.Handle)
	bayaan.Trace("Deleted texture: %s", tex.Name)
}

// Activate binds the texture to a specific texture unit and sets it in the shader.
func (tex *Texture) Activate(sh Shader, unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, tex.Handle)
	sh.Activate()
	sh.SetUniform1i(tex.Name, int32(unit))
	bayaan.Trace("Activated texture %s on unit %d", tex.Name, unit)
}

// initializeTextureParameters sets default values for any unset parameters.
func initializeTextureParameters(params *TextureParameters) {
	if params.WrappingS == 0 {
		params.WrappingS = ClampToEdge
	}
	if params.WrappingT == 0 {
		params.WrappingT = ClampToEdge
	}
	if params.BorderColor == nil {
		params.BorderColor = color.RGBA{255, 255, 255, 255}
	}
	if params.FilteringMin == 0 {
		params.FilteringMin = Nearest
	}
	if params.FilteringMag == 0 {
		params.FilteringMag = Linear
	}
}

// imageToRGBA converts an image.Image to RGBA format,allow flipping it vertically to match OpenGL's coordinate system if needed.
func imageToRGBA(img image.Image, flipImage bool) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	bayaan.Trace("Converting image to RGBA format")
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			if flipImage {
				// Flip vertically to match OpenGL coordinate system
				rgba.Set(x, bounds.Dy()-y-1, img.At(x, y))
			} else {
				// Keep original orientation
				rgba.Set(x, y, img.At(x, y))
			}
		}
	}
	return rgba
}

// checkGLError checks for any OpenGL errors and logs them if found.
func checkGLError(msg string) error {
	if errCode := gl.GetError(); errCode != gl.NO_ERROR {
		errMsg := fmt.Sprintf("%s: OpenGL error: 0x%x", msg, errCode)
		bayaan.Warn(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

package noor

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"os"

	_ "github.com/chai2010/webp" // Register WEBP format

	"github.com/go-gl/gl/v4.6-core/gl"
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

type TextureType uint32

const (
	Texture2D      TextureType = gl.TEXTURE_2D
	TextureArray2D TextureType = gl.TEXTURE_2D_ARRAY
	TextureCubemap TextureType = gl.TEXTURE_CUBE_MAP
)

type TextureFormat uint32

const (
	FormatRGBA8   TextureFormat = gl.RGBA8
	FormatRGB8    TextureFormat = gl.RGB8
	FormatRG8     TextureFormat = gl.RG8
	FormatR8      TextureFormat = gl.R8
	FormatSRGBA8  TextureFormat = gl.SRGB8_ALPHA8
	FormatRGBA16F TextureFormat = gl.RGBA16F
	FormatRGBA32F TextureFormat = gl.RGBA32F
)

type TextureParameters struct {
	WrappingS, WrappingT       TextureWrapping
	BorderColor                color.Color
	FilteringMin, FilteringMag TextureFiltering
	UseMipmaps, FlipImage      bool
	Format                     TextureFormat
	Type                       TextureType
	GenerateMipmaps            bool
	AnisotropyLevel            float32
}

type Texture struct {
	Handle     uint32
	Name       string
	Type       TextureType
	Format     TextureFormat
	Width      int32
	Height     int32
	Depth      int32
	Parameters TextureParameters
}

// NewTextureFromFile creates a new texture from a file path
func NewTextureFromFile(filepath string, parameters TextureParameters) (*Texture, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open texture file %s: %w", filepath, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode texture image %s (format: %s): %w", filepath, format, err)
	}

	tex, err := NewTexture(img, filepath, parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to create texture from image %s: %w", filepath, err)
	}

	return &tex, nil
}

// NewTexture creates a new OpenGL texture from an image and uploads it to the GPU.
func NewTexture(img image.Image, name string, parameters TextureParameters) (tex Texture, err error) {

	// Initialize parameters with defaults if any are unset
	initializeTextureParameters(&parameters)

	rgba := imageToRGBA(img, parameters.FlipImage)
	tex = Texture{
		Name:       name,
		Type:       parameters.Type,
		Format:     parameters.Format,
		Width:      int32(rgba.Rect.Size().X),
		Height:     int32(rgba.Rect.Size().Y),
		Parameters: parameters,
	}

	if err := tex.createAndSetup(rgba); err != nil {
		return tex, fmt.Errorf("failed to create texture: %w", err)
	}

	return tex, nil
}

// createAndSetup handles the OpenGL texture creation and setup
func (tex *Texture) createAndSetup(rgba *image.RGBA) error {
	gl.GenTextures(1, &tex.Handle)
	gl.BindTexture(uint32(tex.Type), tex.Handle)
	defer gl.BindTexture(uint32(tex.Type), 0)

	// Set texture parameters
	gl.TexParameteri(uint32(tex.Type), gl.TEXTURE_WRAP_S, int32(tex.Parameters.WrappingS))
	gl.TexParameteri(uint32(tex.Type), gl.TEXTURE_WRAP_T, int32(tex.Parameters.WrappingT))
	gl.TexParameteri(uint32(tex.Type), gl.TEXTURE_MIN_FILTER, int32(tex.Parameters.FilteringMin))
	gl.TexParameteri(uint32(tex.Type), gl.TEXTURE_MAG_FILTER, int32(tex.Parameters.FilteringMag))

	// Set anisotropic filtering if supported
	if tex.Parameters.AnisotropyLevel > 0 {
		gl.TexParameterf(uint32(tex.Type), gl.TEXTURE_MAX_ANISOTROPY, tex.Parameters.AnisotropyLevel)
	}

	// Set border color if using ClampToBorder
	if tex.Parameters.WrappingS == ClampToBorder || tex.Parameters.WrappingT == ClampToBorder {
		var borderColor [4]float32
		r, g, b, a := tex.Parameters.BorderColor.RGBA()
		borderColor[0] = float32(r) / 0xffff
		borderColor[1] = float32(g) / 0xffff
		borderColor[2] = float32(b) / 0xffff
		borderColor[3] = float32(a) / 0xffff
		gl.TexParameterfv(uint32(tex.Type), gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	}

	// Upload pixel data
	gl.TexImage2D(
		uint32(tex.Type),
		0,
		int32(tex.Format),
		tex.Width,
		tex.Height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	if err := checkGLError("uploading texture data"); err != nil {
		return err
	}

	// Generate mipmaps if requested
	if tex.Parameters.GenerateMipmaps {
		gl.GenerateMipmap(uint32(tex.Type))
		if err := checkGLError("generating mipmaps"); err != nil {
			return err
		}
	}

	return nil
}

// UpdateData updates the texture data for a region of the texture
func (tex *Texture) UpdateData(xOffset, yOffset int32, width, height int32, data []byte) error {
	gl.BindTexture(uint32(tex.Type), tex.Handle)
	defer gl.BindTexture(uint32(tex.Type), 0)

	gl.TexSubImage2D(
		uint32(tex.Type),
		0,
		xOffset,
		yOffset,
		width,
		height,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data),
	)

	return checkGLError("updating texture data")
}

// Resize resizes the texture to the specified dimensions
func (tex *Texture) Resize(width, height int32) error {
	gl.BindTexture(uint32(tex.Type), tex.Handle)
	defer gl.BindTexture(uint32(tex.Type), 0)

	gl.TexImage2D(
		uint32(tex.Type),
		0,
		int32(tex.Format),
		width,
		height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		nil,
	)

	tex.Width = width
	tex.Height = height

	if tex.Parameters.GenerateMipmaps {
		gl.GenerateMipmap(uint32(tex.Type))
	}

	return checkGLError("resizing texture")
}

// Delete removes the texture from GPU memory.
func (tex *Texture) Delete() {
	if tex.Handle != 0 {
		gl.DeleteTextures(1, &tex.Handle)
		tex.Handle = 0
	}
}

// Activate binds the texture to a specific texture unit and sets it in the shader.
func (tex *Texture) Activate(sh Shader, unit uint32, uniformName string) error {
	sh.Activate()
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(uint32(tex.Type), tex.Handle)
	sh.SetUniformInt32(uniformName, int32(unit))
	return checkGLError("activating texture")
}

// initializeTextureParameters sets default values for any unset parameters.
func initializeTextureParameters(params *TextureParameters) {
	if params.WrappingS == 0 {
		params.WrappingS = Repeat
	}
	if params.WrappingT == 0 {
		params.WrappingT = Repeat
	}
	if params.FilteringMin == 0 {
		if params.UseMipmaps {
			params.FilteringMin = LinearMipmapLinear
		} else {
			params.FilteringMin = Linear
		}
	}
	if params.FilteringMag == 0 {
		params.FilteringMag = Linear
	}
	if params.Format == 0 {
		params.Format = FormatRGBA8
	}
	if params.Type == 0 {
		params.Type = Texture2D
	}
	if params.GenerateMipmaps {
		params.UseMipmaps = true
	}
	if params.UseMipmaps {
		params.GenerateMipmaps = true
	}
}

// imageToRGBA converts an image.Image to RGBA format,allow flipping it vertically to match OpenGL's coordinate system if needed.
func imageToRGBA(img image.Image, flipImage bool) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

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
		return fmt.Errorf("%s: OpenGL error: 0x%x", msg, errCode)
	}
	return nil
}

func DefaultTextureParameters() TextureParameters {
	return TextureParameters{
		WrappingS:       Repeat,
		WrappingT:       Repeat,
		FilteringMin:    Linear,
		FilteringMag:    Linear,
		Format:          FormatRGBA8,
		Type:            Texture2D,
		GenerateMipmaps: true,
	}
}

func DefaultDiffuseTextureMap() *Texture {

	// create a new texture Image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// create a new texture
	tex, _ := NewTexture(img, "diffuseMap", DefaultTextureParameters())

	return &tex
}

func DefaultSpecularTextureMap() *Texture {

	// create a new texture Image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// create a new texture
	tex, _ := NewTexture(img, "specularMap", DefaultTextureParameters())

	return &tex
}

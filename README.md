# Noor Rendering Engine - Examples

This document lists example implementations demonstrating how to use the `Noor` rendering engine.
Each example serves as a reference to help understand the core concepts and advanced features of `Noor`.

## 📌 Basics
These examples cover fundamental concepts essential for working with `Noor`:

- [ ] [**Window Creation**](examples/window/main.go) - Setting up a basic window.
- [ ] [**Rendering a Triangle**](examples/triangle/main.go) - Drawing a simple triangle on the screen.
- [ ] [**Shaders**](examples/shader/main.go) - Using vertex and fragment shaders.
- [ ] [**Textures**](examples/texture/main.go) - Loading and applying textures to objects.
- [ ] [**Meshes**](examples/mesh/main.go) - Creating and rendering mesh-based objects.
- [ ] [**Camera System**](examples/camera/main.go) - Implementing a camera for navigation.
- [ ] [**Lighting**](examples/lighting/main.go) - Basic lighting and shading techniques.
- [ ] [**Transformations**](examples/transformations/main.go) - Applying translation, rotation, and scaling.

## 🎮 2D Game Examples
These examples demonstrate how to use `Noor` for developing 2D games:

- [ ] [**Conway's Game of Life**](examples/cgol/main.go) - Cellular automaton simulation.
- [ ] [**Pong**](examples/pong/main.go) - Classic paddle and ball game.
- [ ] [**Space Invaders**](examples/space_invaders/main.go) - Simple arcade shooter.
- [ ] [**Asteroids**](examples/asteroids/main.go) - Classic asteroid destruction game.
- [ ] [**Tetris**](examples/tetris/main.go) - Block-stacking puzzle game.
- [ ] [**Platformer**](examples/platformer/main.go) - Basic side-scrolling platformer.

## 🕶️ 3D Game Examples
These examples introduce 3D rendering concepts in `Noor`:

- [ ] [**3D Triangle**](examples/3d/triangle/main.go) - Basic 3D rendering with a single triangle.
- [ ] [**3D Mesh**](examples/3d/mesh/main.go) - Loading and rendering 3D models.
- [ ] [**3D Lighting**](examples/3d/lighting/main.go) - Implementing light sources in 3D.
- [ ] [**3D First-Person Camera**](examples/3d/fps_camera/main.go) - Navigating a 3D environment with a first-person view.

## 🔧 Noor API Design
Below is a high-level design of Noor's API:

### Initialization & Context Management
- `noor.Init() error` – Initializes the Noor engine.
- `noor.Terminate()` – Cleans up resources before exiting.
- `noor.CreateWindow(title string, width, height int) (*Window, error)` – Creates a rendering window.

### Rendering Loop & Context
- `Window.RenderLoop(func())` – Runs the main rendering loop.
- `Window.SwapBuffers()` – Swaps the frame buffers.

### Meshes & Objects
- `noor.NewMesh(vertices []Vertex, indices []uint32) *Mesh`
- `Mesh.Draw(shader *Shader)`

### Shaders
- `noor.NewShader(vertexSrc, fragmentSrc string) (*Shader, error)`
- `Shader.Use()`
- `Shader.SetUniform(name string, value interface{})`

### Textures
- `noor.NewTexture(filepath string) (*Texture, error)`
- `Texture.Bind(slot int)`

### Camera
- `noor.NewCamera(position, target Vector3) *Camera`
- `Camera.SetPerspective(fov, aspect, near, far float32)`

### Lighting
- `noor.NewLight(position Vector3, color Vector3, intensity float32) *Light`

### Transformations
- `Object.Translate(x, y, z float32)`
- `Object.Rotate(angle float32, axis Vector3)`
- `Object.Scale(x, y, z float32)`

### Utilities
- `noor.GetTime() float32`
- `noor.GetDeltaTime() float32`
- `noor.GetWindowWidth() int`
- `noor.GetWindowHeight() int`

For more information, visit [Noor's website](https://noor-engine.github.io/noor/).

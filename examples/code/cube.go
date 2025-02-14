package main

import "github.com/ahmedsat/noor"

var vertices = []noor.Vertex{
	// Front face
	{Position: [3]float32{-0.5, -0.5, 0.5}, Color: [3]float32{1, 0, 0}, UV: [2]float32{0, 0}, Normal: [3]float32{0, 0, 1}}, // Bottom-left
	{Position: [3]float32{0.5, -0.5, 0.5}, Color: [3]float32{0, 1, 0}, UV: [2]float32{1, 0}, Normal: [3]float32{0, 0, 1}},  // Bottom-right
	{Position: [3]float32{0.5, 0.5, 0.5}, Color: [3]float32{0, 0, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{0, 0, 1}},   // Top-right
	{Position: [3]float32{-0.5, 0.5, 0.5}, Color: [3]float32{1, 1, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{0, 0, 1}},  // Top-left

	// Back face
	{Position: [3]float32{-0.5, -0.5, -0.5}, Color: [3]float32{1, 0, 1}, UV: [2]float32{0, 0}, Normal: [3]float32{0, 0, -1}}, // Bottom-left
	{Position: [3]float32{0.5, -0.5, -0.5}, Color: [3]float32{0, 1, 1}, UV: [2]float32{1, 0}, Normal: [3]float32{0, 0, -1}},  // Bottom-right
	{Position: [3]float32{0.5, 0.5, -0.5}, Color: [3]float32{1, 1, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{0, 0, -1}},   // Top-right
	{Position: [3]float32{-0.5, 0.5, -0.5}, Color: [3]float32{0, 0, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{0, 0, -1}},  // Top-left

	// Left face
	{Position: [3]float32{-0.5, -0.5, -0.5}, Color: [3]float32{1, 0, 0}, UV: [2]float32{0, 0}, Normal: [3]float32{-1, 0, 0}}, // Bottom-back
	{Position: [3]float32{-0.5, -0.5, 0.5}, Color: [3]float32{0, 1, 0}, UV: [2]float32{1, 0}, Normal: [3]float32{-1, 0, 0}},  // Bottom-front
	{Position: [3]float32{-0.5, 0.5, 0.5}, Color: [3]float32{0, 0, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{-1, 0, 0}},   // Top-front
	{Position: [3]float32{-0.5, 0.5, -0.5}, Color: [3]float32{1, 1, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{-1, 0, 0}},  // Top-back

	// Right face
	{Position: [3]float32{0.5, -0.5, -0.5}, Color: [3]float32{1, 0, 1}, UV: [2]float32{0, 0}, Normal: [3]float32{1, 0, 0}}, // Bottom-back
	{Position: [3]float32{0.5, -0.5, 0.5}, Color: [3]float32{0, 1, 1}, UV: [2]float32{1, 0}, Normal: [3]float32{1, 0, 0}},  // Bottom-front
	{Position: [3]float32{0.5, 0.5, 0.5}, Color: [3]float32{1, 1, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{1, 0, 0}},   // Top-front
	{Position: [3]float32{0.5, 0.5, -0.5}, Color: [3]float32{0, 0, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{1, 0, 0}},  // Top-back

	// Top face
	{Position: [3]float32{-0.5, 0.5, -0.5}, Color: [3]float32{1, 0, 0}, UV: [2]float32{0, 0}, Normal: [3]float32{0, 1, 0}}, // Back-left
	{Position: [3]float32{0.5, 0.5, -0.5}, Color: [3]float32{0, 1, 0}, UV: [2]float32{1, 0}, Normal: [3]float32{0, 1, 0}},  // Back-right
	{Position: [3]float32{0.5, 0.5, 0.5}, Color: [3]float32{0, 0, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{0, 1, 0}},   // Front-right
	{Position: [3]float32{-0.5, 0.5, 0.5}, Color: [3]float32{1, 1, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{0, 1, 0}},  // Front-left

	// Bottom face
	{Position: [3]float32{-0.5, -0.5, -0.5}, Color: [3]float32{1, 0, 1}, UV: [2]float32{0, 0}, Normal: [3]float32{0, -1, 0}}, // Back-left
	{Position: [3]float32{0.5, -0.5, -0.5}, Color: [3]float32{0, 1, 1}, UV: [2]float32{1, 0}, Normal: [3]float32{0, -1, 0}},  // Back-right
	{Position: [3]float32{0.5, -0.5, 0.5}, Color: [3]float32{1, 1, 1}, UV: [2]float32{1, 1}, Normal: [3]float32{0, -1, 0}},   // Front-right
	{Position: [3]float32{-0.5, -0.5, 0.5}, Color: [3]float32{0, 0, 0}, UV: [2]float32{0, 1}, Normal: [3]float32{0, -1, 0}},  // Front-left
}

var indices = []uint32{
	// Front face
	0, 1, 2, 2, 3, 0,
	// Back face
	4, 5, 6, 6, 7, 4,
	// Left face
	8, 9, 10, 10, 11, 8,
	// Right face
	12, 13, 14, 14, 15, 12,
	// Top face
	16, 17, 18, 18, 19, 16,
	// Bottom face
	20, 21, 22, 22, 23, 20,
}

#version 460 core

layout(location = 0) in vec3 aPosition;
layout(location = 1) in vec3 aColor;
layout(location = 2) in vec2 aUv;
layout(location = 3) in vec3 aNormal;

out vec3 vColor;
out vec2 vUv;
out vec3 vNormal;

void main() {
  gl_Position = vec4(aPosition, 1.0);
  vColor = aColor;
  vUv = aUv;
  vNormal = aNormal;
}
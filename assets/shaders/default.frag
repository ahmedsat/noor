#version 460
out vec4 fragColor;

in vec3 vColor;
in vec2 vUv;
in vec3 vNormal;

void main() {
  fragColor = vec4(vColor, 1.0);
}
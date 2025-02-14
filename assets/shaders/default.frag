#version 460
out vec4 fragColor;

in vec3 vColor;
in vec2 vUv;
in vec3 vNormal;

uniform sampler2D uTexture;

void main() {
  vec4 texColor = texture(uTexture, vUv);
  float r = texColor.r * int(texColor.r != 0) + vColor.r * int(texColor.r == 0);
  float g = texColor.g * int(texColor.g != 0) + vColor.g * int(texColor.g == 0);
  float b = texColor.b * int(texColor.b != 0) + vColor.b * int(texColor.b == 0);
  fragColor = vec4(r, g, b, 1.0);
}
#version 460 core

out vec4 fragColor;

in vec3 vColor;
in vec2 vUv;
in vec3 vNormal;

uniform sampler2D wallTexture;

void main() {
  fragColor = mix(vec4(vColor, 1.0), texture(wallTexture, vUv),.5);
  // fragColor = vec4(vColor, 1.0);
}
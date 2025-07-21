#version 410 core

out vec4 FragColor;

in vec2 uv;
in vec3 fg;
in vec3 bg;

uniform sampler2D fontAtlas;

void main() {
    float alpha = texture(fontAtlas, uv).r;
    FragColor = vec4(mix(bg, fg, alpha), 1.0);
}

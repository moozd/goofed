#version 330 core
in vec2 TexCoords;
in vec3 FGColor;
in vec3 BGColor;

out vec4 FragColor;

uniform sampler2D uFontAtlas;

void main() {
    float alpha = texture(uFontAtlas, TexCoords).r; // swizzled from red
    vec3 color = mix(BGColor, FGColor, alpha);
    FragColor = vec4(TexCoords, 0.0, 1.0);
}

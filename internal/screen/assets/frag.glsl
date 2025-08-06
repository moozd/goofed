#version 410 core

in vec2 uv;
in vec3 fg;
in vec3 bg;

out vec4 FragColor;

uniform sampler2D fontAtlas;
uniform float pixelRange; // SDF pixel range (typically 4.0-8.0)
uniform vec2 atlasSize; // Atlas texture dimensions

void main() {
    float sdf = texture(fontAtlas, uv).r;
    // float alpha = pow(sdf, 5.0);
    float distance = sdf - 0.5;
    distance *= pixelRange;

    vec2 unitRange = pixelRange / atlasSize;
    vec2 screenTexSize = vec2(1.0) / fwidth(uv);
    float screenPxRange = max(0.5 * dot(unitRange, screenTexSize), 1.0);

    float alpha = clamp(distance * screenPxRange + 0.5, 0.0, 1.0);
    vec3 color = mix(bg, fg, alpha);
    FragColor = vec4(color, 1.0);
}

#version 330 core
layout(location = 0) in vec2 aPos;
layout(location = 1) in vec2 aUV;
layout(location = 2) in vec3 aFGColor;
layout(location = 3) in vec3 aBGColor;

out vec2 TexCoords;
out vec3 FGColor;
out vec3 BGColor;

uniform mat4 uProjection;

void main() {
    gl_Position = uProjection * vec4(aPos, 0.0, 1.0);
    TexCoords = aUV;
    FGColor = aFGColor;
    BGColor = aBGColor;
}

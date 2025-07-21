#version 410 core

layout(location = 0) in vec3 aPos;
layout(location = 1) in vec2 aUV;
layout(location = 2) in vec3 aFg;
layout(location = 3) in vec3 aBg;

out vec3 fg;
out vec3 bg;
out vec2 uv;

 

uniform mat4 projection;

void main() {
    gl_Position = projection * vec4(aPos, 1.0);
    fg = aFg;
    uv = aUV;
    bg = aBg;
}

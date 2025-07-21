#version 410 core

out vec4 FragColor;

in vec2 uv;
in vec3 fg;
in vec3 bg;

void main() {
    float borderSize = 0.1;

    bool isBorder =
        uv.x < borderSize || uv.x > 1.0 - borderSize ||
            uv.y < borderSize || uv.y > 1.0 - borderSize;

    FragColor = isBorder ? vec4(fg, 1.0) : vec4(bg, 1.0);
}

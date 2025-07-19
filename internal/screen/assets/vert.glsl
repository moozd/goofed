#version 410 core

layout(location = 0) in vec3 aPos; // Vertex position input (from VBO)
layout(location = 1) in vec3 aColor; // Vertex color input (from VBO)

out vec3 vertexColor; // Pass color to fragment shader

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
    // Apply transformations
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    vertexColor = aColor; // Forward color attribute
}

#version 410 core

out vec4 FragColor;

uniform float squareSize; // size of each square in pixels

void main() {
    // Get the pixel coordinates
    int x = int(floor(gl_FragCoord.x / squareSize));
    int y = int(floor(gl_FragCoord.y / squareSize));

    // Compute sum of x and y
    int sum = x + y;

    // If sum is even, color white, else black
    if (sum % 2 == 0) {
        FragColor = vec4(1.0); // white
    } else {
        FragColor = vec4(0.0, 0.0, 0.0, 1.0); // black
    }
}

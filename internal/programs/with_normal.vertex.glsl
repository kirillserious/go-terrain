#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec3 normal;

out vec3 Normal;
out vec3 FragPos;

void main() {
    Normal = normalize(mat3(transpose(inverse(model))) * normal);
    FragPos = vec3(model * vec4(vert, 1));
    gl_Position = projection * camera * model * vec4(vert, 1);
}
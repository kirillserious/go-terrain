package programs

import gl2 "terrain/internal/gl"

const terrainVertex = `
#version 330
uniform int max_i;
uniform int max_j;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec3 normal;

out vec3 Normal;
out vec3 FragPos;
out vec2 FragTexCoord;

void main() {
    Normal = normalize(mat3(transpose(inverse(model))) * normal);
    FragPos = vec3(model * vec4(vert, 1));
	FragTexCoord = vec2(vert.x / max_i, vert.z / max_j);
    gl_Position = projection * camera * model * vec4(vert, 1);
}
`

const terrainFragment = `
#version 330
uniform sampler2D tex;
uniform vec3 viewPos;

in vec3 Normal;
in vec3 FragPos;
in vec2 FragTexCoord;

out vec4 outputColor;

const vec3 LightPos = vec3(50., 50., 50.);
const vec4 LightColor = vec4(1., 1., 1., 1.);

float specularStrength = 0.5f;

void main() {
    vec3 lightDir = normalize(LightPos - FragPos);
    vec3 viewDir = normalize(viewPos - FragPos);
    float diff = max(dot(Normal, lightDir), 0.0);
    vec3 diffuse = diff * vec3(LightColor);
    vec3 reflectDir = reflect(lightDir, -Normal);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32);
    vec3 specular = specularStrength * spec * vec3(LightColor);

    vec3 ambient = vec3(0.5, 0.5, 0.5) * vec3(LightColor);
    vec3 result = ambient + diffuse + specular;

    outputColor = vec4(result, 1.0f) * texture(tex, FragTexCoord);
}
`

var Terrain = LasyProgram(map[gl2.ShaderKind]string{
	gl2.VertexShaderKind:   terrainVertex,
	gl2.FragmentShaderKind: terrainFragment,
})

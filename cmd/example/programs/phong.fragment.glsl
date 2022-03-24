#version 330
uniform vec4 Color;
uniform vec3 viewPos;

in vec3 Normal;
in vec3 FragPos;

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

    outputColor = vec4(result, 1.0f) * Color;
}
uniform sampler2D tDiffuse;
uniform float uIntensity;
uniform float uTime;

varying vec2 vUv;

float Random(vec2 seed) {
    return fract(sin(dot(seed, vec2(12.9898, 78.233))) * 43758.5453);
}

void main() {
    vec2 uv = vUv;

    if (uIntensity > 0.001) {
        float blockX = floor(uv.x * 500.0);
        float lineNoise = Random(vec2(blockX, floor(uTime * 18.0)));

        if (lineNoise < uIntensity * 0.4) {
            uv.x += (lineNoise - 0.5) * uIntensity * 0.15;
        }

        float shift = uIntensity * 0.01;
        float r = texture2D(tDiffuse, uv + vec2(shift, 0.0)).r;
        float g = texture2D(tDiffuse, uv).g;
        float b = texture2D(tDiffuse, uv - vec2(shift, 0.0)).b;

        float staticNoise = (Random(uv * uTime) - 0.5) * uIntensity * 0.25;

        vec3 color = vec3(r, g, b) + staticNoise;
        color = mix(color, vec3(0.3, 0.7, 0.5), uIntensity * 0.12);

        gl_FragColor = vec4(color, 0.5);
    } else {
        gl_FragColor = texture2D(tDiffuse, uv);
    }
}
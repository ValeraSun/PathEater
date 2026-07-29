uniform sampler2D tDiffuse;

uniform float uIntensity;
uniform float uTime;
uniform float uAspect;

uniform vec2 uDropletPositions[12];
uniform float uDropletSizes[12];
uniform float uDropletSpeeds[12];
uniform float uDropletDarkness[12];

varying vec2 vUv;

const int DROPLETS = 12;

float Drop(vec2 uv, vec2 center, float radius) {
    vec2 p = (uv - center) / radius;
    p.x *= uAspect;

    float width = mix(1.05, 0.3, clamp(p.y * 0.5 + 0.5, 0.0, 1.0));
    float d = length(vec2(p.x / max(width, 0.18), p.y * 0.88));

    return 1.0 - smoothstep(0.76, 1.0, d);
}

float Trail(vec2 uv, vec2 center, float radius) {
    vec2 p = uv - center;
    p.x *= uAspect;

    float start = center.y + radius * 0.15;
    float end = center.y + radius * 3.8;

    if (uv.y < start || uv.y > end) {
        return 0.0;
    }

    float t = (uv.y - start) / (end - start);
    float w = mix(radius * 0.3, radius * 0.06, t);

    float x = 1.0 - smoothstep(w * 0.45, w, abs(p.x));
    float fade = pow(1.0 - t, 0.65);

    return x * fade;
}

void main() {
    vec4 color = texture2D(tDiffuse, vUv);

    if (uIntensity <= 0.001) {
        gl_FragColor = color;
        return;
    }

    vec3 result = color.rgb;

    for (int i = 0; i < DROPLETS; i++) {
        float radius = uDropletSizes[i];
        vec2 pos = uDropletPositions[i];

        pos.y -= uTime * uDropletSpeeds[i];

        if (pos.y < -radius * 2.0) {
            continue;
        }

        float mask = Drop(vUv, pos, radius);
        float center = Drop(vUv, pos, radius * 0.58);
        float streak = Trail(vUv, pos, radius);

        vec3 edge = mix(
            vec3(0.19, 0.015, 0.27),
            vec3(0.09, 0.004, 0.13),
            uDropletDarkness[i]
        );

        result = mix(result, vec3(0.045, 0.002, 0.065), streak * uIntensity * 0.42);
        result = mix(result, edge, mask * uIntensity * 0.88);
        result = mix(result, vec3(0.004, 0.0, 0.007), center * uIntensity * 0.94);
    }

    gl_FragColor = vec4(result, color.a);
}
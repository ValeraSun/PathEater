import * as THREE from "three";
import { EffectComposer } from "three/examples/jsm/postprocessing/EffectComposer.js";
import { OutputPass } from "three/examples/jsm/postprocessing/OutputPass.js";
import { RenderPass } from "three/examples/jsm/postprocessing/RenderPass.js";
import { ShaderPass } from "three/examples/jsm/postprocessing/ShaderPass.js";

import fragmentShader from "../Shaders/DamageGlitch.frag.glsl";
import vertexShader from "../Shaders/DamageGlitch.vert.glsl";

const DAMAGE_DECAY_PER_SECOND = 2.2;
const DEBUG_GLITCH_ENABLED = false;
const DEBUG_GLITCH_INTENSITY = 0.8;

const DamageGlitchShader = {
    uniforms: {
        tDiffuse: { value: null },
        uIntensity: { value: 0 },
        uTime: { value: 0 }
    },
    vertexShader,
    fragmentShader
};

export class DamageGlitchEffect {
    private readonly composer: EffectComposer;
    private readonly glitchPass: ShaderPass;

    private intensity = 0;
    private elapsedTime = 0;

    public constructor(renderer: THREE.WebGLRenderer, scene: THREE.Scene, camera: THREE.Camera) {
        this.composer = new EffectComposer(renderer);
        this.composer.setPixelRatio(window.devicePixelRatio);

        const renderPass = new RenderPass(scene, camera);
        this.glitchPass = new ShaderPass(DamageGlitchShader);
        const outputPass = new OutputPass();

        this.composer.addPass(renderPass);
        this.composer.addPass(this.glitchPass);
        this.composer.addPass(outputPass);
    }

    public TriggerDamage(strength: number): void {
        this.intensity = Math.min(1, this.intensity + strength);
    }

    public Render(deltaTime: number): void {
        this.elapsedTime += deltaTime;

        if (DEBUG_GLITCH_ENABLED) {
            this.glitchPass.uniforms.uIntensity.value = DEBUG_GLITCH_INTENSITY;
        } else {
            this.intensity = Math.max(0, this.intensity - DAMAGE_DECAY_PER_SECOND * deltaTime);
            this.glitchPass.uniforms.uIntensity.value = this.intensity;
        }

        this.glitchPass.uniforms.uTime.value = this.elapsedTime;
        this.composer.render(deltaTime);
    }

    public SetSize(width: number, height: number): void {
        this.composer.setSize(width, height);
    }

    public Dispose(): void {
        this.glitchPass.dispose();
        this.composer.dispose();
    }
}
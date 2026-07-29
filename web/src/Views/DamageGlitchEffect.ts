import * as THREE from "three";
import { EffectComposer } from "three/examples/jsm/postprocessing/EffectComposer.js";
import { OutputPass } from "three/examples/jsm/postprocessing/OutputPass.js";
import { RenderPass } from "three/examples/jsm/postprocessing/RenderPass.js";
import { ShaderPass } from "three/examples/jsm/postprocessing/ShaderPass.js";

import vertexShader from "../Shaders/DamageGlitch.vert.glsl";
import fragmentShader from "../Shaders/DamageGlitch.frag.glsl";
import bloodFragmentShader from "../Shaders/MonsterBloodSplatter.frag.glsl";

const DamageDecayPerSecond = 2.2;
const BloodDecayPerSecond = 0.4;
const DropletCount = 12;

const DebugAlwaysGlitch = false;
const DebugGlitchIntensity = 0.35;

const DamageShader = {
    uniforms: {
        tDiffuse: { value: null },
        uIntensity: { value: 0 },
        uTime: { value: 0 }
    },
    vertexShader,
    fragmentShader
};

const BloodShader = {
    uniforms: {
        tDiffuse: { value: null },
        uIntensity: { value: 0 },
        uTime: { value: 0 },
        uAspect: { value: window.innerWidth / window.innerHeight },
        uDropletPositions: {
            value: Array.from(
                { length: DropletCount },
                () => new THREE.Vector2()
            )
        },
        uDropletSizes: {
            value: new Array(DropletCount).fill(0)
        },
        uDropletSpeeds: {
            value: new Array(DropletCount).fill(0)
        },
        uDropletDarkness: {
            value: new Array(DropletCount).fill(0)
        }
    },
    vertexShader,
    fragmentShader: bloodFragmentShader
};

export class DamageGlitchEffect {
    private Composer: EffectComposer;
    private GlitchPass: ShaderPass;
    private BloodPass: ShaderPass;

    private GlitchIntensity = 0;
    private BloodIntensity = 0;
    private ElapsedTime = 0;

    public constructor(
        renderer: THREE.WebGLRenderer,
        scene: THREE.Scene,
        camera: THREE.Camera
    ) {
        this.Composer = new EffectComposer(renderer);
        this.Composer.addPass(new RenderPass(scene, camera));

        this.GlitchPass = new ShaderPass(DamageShader);
        this.Composer.addPass(this.GlitchPass);

        this.BloodPass = new ShaderPass(BloodShader);
        this.Composer.addPass(this.BloodPass);

        this.Composer.addPass(new OutputPass());
    }

    public TriggerDamage(strength: number): void {
        this.GlitchIntensity = Math.min(1, this.GlitchIntensity + strength);
    }

    public TriggerMonsterBloodSplatter(): void {
        this.BloodIntensity = 1;

        const uniforms = this.BloodPass.uniforms;
        const positions = uniforms.uDropletPositions.value as THREE.Vector2[];
        const sizes = uniforms.uDropletSizes.value as number[];
        const speeds = uniforms.uDropletSpeeds.value as number[];
        const darkness = uniforms.uDropletDarkness.value as number[];

        uniforms.uTime.value = 0;

        for (let index = 0; index < DropletCount; index++) {
            positions[index].set(Math.random() * 0.9 + 0.05, Math.random() * 0.84 + 0.08);
            sizes[index] = Math.random() * 0.011 + 0.009;
            speeds[index] = Math.random() * 0.028 + 0.018;
            darkness[index] = Math.random() * 0.8 + 0.2;
        }
    }

    public Render(deltaTime: number): void {
        this.ElapsedTime += deltaTime;
        this.UpdateGlitch(deltaTime);
        this.UpdateBlood(deltaTime);
        this.UpdateUniforms();
        this.Composer.render(deltaTime);
    }

    public SetSize(width: number, height: number): void {
        this.Composer.setSize(width, height);
        this.BloodPass.uniforms.uAspect.value = width / height;
    }

    public Dispose(): void {
        this.Composer.dispose();
        this.GlitchPass.dispose();
        this.BloodPass.dispose();
    }

    private UpdateGlitch(deltaTime: number): void {
        if (DebugAlwaysGlitch) {
            this.GlitchIntensity = DebugGlitchIntensity;
            return;
        }

        this.GlitchIntensity = Math.max(0, this.GlitchIntensity - DamageDecayPerSecond * deltaTime);
    }

    private UpdateBlood(deltaTime: number): void {
        if (this.BloodIntensity <= 0) {
            return;
        }

        this.BloodIntensity = Math.max(0, this.BloodIntensity - BloodDecayPerSecond * deltaTime);
        this.BloodPass.uniforms.uTime.value += deltaTime;
    }

    private UpdateUniforms(): void {
        this.GlitchPass.uniforms.uIntensity.value = this.GlitchIntensity;
        this.GlitchPass.uniforms.uTime.value = this.ElapsedTime;
        this.BloodPass.uniforms.uIntensity.value = this.BloodIntensity;
    }
}
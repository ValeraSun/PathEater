import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { clone as cloneSkeleton } from "three/examples/jsm/utils/SkeletonUtils.js";

const MODEL_URL = "/models/walking_astronaut_webp.glb";
const FADE_DURATION = 0.2;

interface LoadedModel {
    scene: THREE.Group;
    animations: THREE.AnimationClip[];
}

let modelPromise: Promise<LoadedModel> | null = null;

function loadModel(): Promise<LoadedModel> {
    if (!modelPromise) {
        const loader = new GLTFLoader();
        modelPromise = loader.loadAsync(MODEL_URL).then(gltf => ({
            scene: gltf.scene,
            animations: gltf.animations
        }));
    }
    return modelPromise;
}

export class PlayerView 
{
    public mesh = new THREE.Group();
    private mixer: THREE.AnimationMixer | null = null;
    private actions: Record<string, THREE.AnimationAction> = {};
    private currentAction: THREE.AnimationAction | null = null;
    
    public constructor() {
        loadModel()
            .then(({ scene, animations }) => {
                const model = cloneSkeleton(scene) as THREE.Group;
                this.mesh.add(model);

                this.mixer = new THREE.AnimationMixer(model);

                for (const clip of animations) {
                    this.actions[clip.name] = this.mixer.clipAction(clip);
                }

                this.playAnimation("idle");
            })
            .catch(error => {
                console.error("Не удалось загрузить модель игрока:", error);
            });
    }

    public AdvanceAnimation(dt: number): void {
        this.mixer?.update(dt);
    }

    public SetMoving(isMoving: boolean): void {
        this.playAnimation(
            isMoving ? "moon_walk" : "idle",
            isMoving ? 2.3 : 1
        );
    }

    private playAnimation(name: string, speed = 1): void {
        const next = this.actions[name];

        if (!next) return;

        next.timeScale = speed;

        if (next === this.currentAction) return;

        next.reset().fadeIn(FADE_DURATION).play();

        this.currentAction?.fadeOut(FADE_DURATION);
        this.currentAction = next;
    }
}
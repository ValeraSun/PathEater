import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { clone as cloneSkeleton } from "three/examples/jsm/utils/SkeletonUtils.js";

const MODEL_URL = "/models/alien_walk_attack.glb";
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

export class AlienView {
    public mesh = new THREE.Group();
    private mixer: THREE.AnimationMixer | null = null;
    private actions: Record<string, THREE.AnimationAction> = {};
    private currentAction: THREE.AnimationAction | null = null;
    private attacking = false;

    public constructor() {
        loadModel()
            .then(({ scene, animations }) => {
                const model = cloneSkeleton(scene) as THREE.Group;
                this.mesh.add(model);
                this.mixer = new THREE.AnimationMixer(model);
                for (const clip of animations) {
                    console.log("Alien animation:", clip.name);

                    this.actions[clip.name] = this.mixer.clipAction(clip);
                }

                this.ApplyAnimationState();
            })
            .catch(error => {
                console.error("Не удалось загрузить модель пришельца:", error);
            });
    }

    public AdvanceAnimation(dt: number): void {
        this.mixer?.update(dt);
    }

    public SetAttacking(attacking: boolean): void {
        if (this.attacking === attacking) {
            return;
        }

        this.attacking = attacking;
        this.ApplyAnimationState();
    }

    private ApplyAnimationState(): void {
        const actionName = this.attacking
            ? this.FindActionName("attack")
            : this.FindActionName("walk");

        if (!actionName) {
            console.warn('Не найдена анимация');
            return;
        }

        this.PlayAnimation(actionName,this.attacking);
    }

    private FindActionName(part: "walk" | "attack"): string | null {
        const names = Object.keys(this.actions);
        return (
            names.find(name =>
                name
                    .toLowerCase()
                    .includes(part)
            ) ?? null
        );
    }

    private PlayAnimation(name: string,playOnce: boolean): void {
        const next = this.actions[name];
        if (!next || next === this.currentAction) {
            return;
        }

        next.reset();
        next.enabled = true;
        next.setEffectiveTimeScale(1);
        next.setEffectiveWeight(1);

        if (playOnce) {
            next.setLoop(THREE.LoopOnce, 1);
            next.clampWhenFinished = true;
        } else {
            next.setLoop(
                THREE.LoopRepeat,
                Infinity
            );

            next.clampWhenFinished = false;
        }

        next.fadeIn(FADE_DURATION).play();
        this.currentAction?.fadeOut(FADE_DURATION);
        this.currentAction = next;
    }
}
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { clone as cloneSkeleton } from "three/examples/jsm/utils/SkeletonUtils.js";

const MODEL_URL = "/models/metal_fnaf_door.glb";

const ANIMATION_SPEED = 1;



let modelPromise: Promise<THREE.Group> | null = null;

function loadModel(): Promise<THREE.Group> {
    if (!modelPromise) {
        const loader = new GLTFLoader();
        modelPromise = loader.loadAsync(MODEL_URL).then(glft => glft.scene);
    }
    return modelPromise;
}

export class DoorView {
    public mesh = new THREE.Group();
    public isOpen: boolean = false
    private doorModel: THREE.Object3D | null = null;
    private progress: number = 0;
    private targetProgress: number = 0;
    private state: string = "close";
    private animationSpeed: number = ANIMATION_SPEED;
    private doorParts: { wholeDoor?: THREE.Object3D } = {};
    private doorHeight: number = 3.0;

   constructor() {
    console.log("[DoorView] Constructor called");
        loadModel()      
            .then(( scene ) => {
                const model = cloneSkeleton(scene) as THREE.Group;
                model.position.y = 0;

                this.mesh.add(model);
                this.doorParts.wholeDoor = model;

            })
            .catch(error => {
                console.error("Не удалось загрузить модель двери:", error);
            });;
}

    // Реализация AnimatedEntityView
    public AdvanceAnimation(deltaTime: number): void {
        if (this.isOpen && this.state === "close") {
            this.Open()
        }

        if (!this.isOpen && this.state === "open") {
            this.Close()
        }

        this.Update(deltaTime);
    }

    public SetState(isOpen: boolean, progress?: number): void {
        if (progress !== undefined) {
            this.progress = Math.max(0, Math.min(1, progress));
            this.targetProgress = this.progress;
            this.state = isOpen ? "open" : "close";
            this.UpdateDoorPosition();
            return;
        }
        if (isOpen) this.Open();
        else this.Close();
    }

    public Open(): void {
        console.log("открыли")
        if (this.state === "close") {
            this.state = "open";
            this.targetProgress = 1;
        }
    }

    public Close(): void {
          console.log("закрыли")
        if (this.state === "open") {
            this.state = "close";
            this.targetProgress = 0;
        }
    }

    private Update(deltaTime: number): void {
        if (this.state === "open") {
            this.progress += deltaTime * this.animationSpeed;
            if (this.progress >= 1) { this.progress = 1; }
            this.UpdateDoorPosition();
        } else if (this.state === "close") {
            this.progress -= deltaTime * this.animationSpeed;
            if (this.progress <= 0) { this.progress = 0; }
            this.UpdateDoorPosition();
        }
    }

    private UpdateDoorPosition(): void {
        if (!this.doorParts.wholeDoor) return;
        const offset = this.progress * this.doorHeight;
        this.doorParts.wholeDoor.position.y = -offset;
    }
}
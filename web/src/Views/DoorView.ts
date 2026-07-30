import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { AnimatedEntityView } from "../Views/EntityViewFactory"; // путь может отличаться

const MODEL_URL = "/models/space_door.glb";
const ANIMATION_SPEED = 2.0;

export enum DoorState { CLOSED, OPENING, OPEN, CLOSING }

export class DoorView implements AnimatedEntityView {
    public mesh = new THREE.Group();
    private doorModel: THREE.Object3D | null = null;
    private state: DoorState = DoorState.CLOSED;
    private progress: number = 0;
    private targetProgress: number = 0;
    private animationSpeed: number = ANIMATION_SPEED;
    private doorParts: { wholeDoor?: THREE.Object3D } = {};
    private doorHeight: number = 3.0;
    private doorWidth: number = 2.0;

   constructor() {
    console.log("[DoorView] Constructor called");
    this.LoadModel();
}

    // Реализация AnimatedEntityView
    public AdvanceAnimation(deltaTime: number): void {
        this.Update(deltaTime);
    }

    public SetState(isOpen: boolean, progress?: number): void {
        if (progress !== undefined) {
            this.progress = Math.max(0, Math.min(1, progress));
            this.targetProgress = this.progress;
            this.state = isOpen ? DoorState.OPEN : DoorState.CLOSED;
            this.UpdateDoorPosition();
            return;
        }
        if (isOpen) this.Open();
        else this.Close();
    }

    public Open(): void {
        if (this.state === DoorState.CLOSED || this.state === DoorState.CLOSING) {
            this.state = DoorState.OPENING;
            this.targetProgress = 1;
        }
    }

    public Close(): void {
        if (this.state === DoorState.OPEN || this.state === DoorState.OPENING) {
            this.state = DoorState.CLOSING;
            this.targetProgress = 0;
        }
    }

    public Toggle(): void {
        if (this.state === DoorState.OPEN || this.state === DoorState.OPENING) this.Close();
        else this.Open();
    }

    private Update(deltaTime: number): void {
        if (this.state === DoorState.OPENING) {
            this.progress += deltaTime * this.animationSpeed;
            if (this.progress >= 1) { this.progress = 1; this.state = DoorState.OPEN; }
            this.UpdateDoorPosition();
        } else if (this.state === DoorState.CLOSING) {
            this.progress -= deltaTime * this.animationSpeed;
            if (this.progress <= 0) { this.progress = 0; this.state = DoorState.CLOSED; }
            this.UpdateDoorPosition();
        }
    }

    private LoadModel(): void {
    console.log("[DoorView] Loading model from", MODEL_URL);
    const loader = new GLTFLoader();
    loader.load(MODEL_URL, (gltf) => {
        console.log("[DoorView] Model loaded successfully");
        this.doorModel = gltf.scene;
        this.doorParts.wholeDoor = this.doorModel;
        this.mesh.add(this.doorModel);
        this.UpdateDoorPosition();
    }, undefined, (error) => {
        console.error("[DoorView] Error loading model, using fallback", error);
        this.CreateFallbackDoor();
    });
}

    private CreateFallbackDoor(): void {
        const geometry = new THREE.BoxGeometry(this.doorWidth, this.doorHeight, 0.2);
        const material = new THREE.MeshStandardMaterial({ color: 0x446688, metalness: 0.8, roughness: 0.3 });
        const door = new THREE.Mesh(geometry, material);
        this.mesh.add(door);
        this.doorParts.wholeDoor = door;
        this.UpdateDoorPosition();
    }

    private UpdateDoorPosition(): void {
        if (!this.doorParts.wholeDoor) return;
        const offset = this.progress * this.doorHeight;
        this.doorParts.wholeDoor.position.y = -offset;
    }

    // Дополнительные методы для получения состояния (опционально)
    public IsOpen(): boolean { return this.state === DoorState.OPEN; }
    public GetProgress(): number { return this.progress; }
}
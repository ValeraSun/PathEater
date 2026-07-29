import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

const MODEL_URL = "/models/space_door.glb"; // Путь к вашей модели
const ANIMATION_SPEED = 2.0; // Скорость открытия/закрытия

export enum DoorState {
    CLOSED = 0,
    OPENING = 1,
    OPEN = 2,
    CLOSING = 3
}

export class DoorView {
    public mesh = new THREE.Group();
    private doorModel: THREE.Object3D | null = null;
    private state: DoorState = DoorState.CLOSED;
    private progress: number = 0; // 0 - закрыто, 1 - открыто
    private targetProgress: number = 0;
    private animationSpeed: number = ANIMATION_SPEED;
    
    // Части двери для анимации
    private doorParts: {
        leftDoor?: THREE.Object3D;
        rightDoor?: THREE.Object3D;
        topDoor?: THREE.Object3D;
        bottomDoor?: THREE.Object3D;
        wholeDoor?: THREE.Object3D;
    } = {};

    // Параметры анимации
    private doorHeight: number = 3.0;
    private doorWidth: number = 2.0;

    public constructor() {
        this.LoadModel();
    }

    public SetState(isOpen: boolean, progress?: number): void {
        if (progress !== undefined) {
            this.progress = Math.max(0, Math.min(1, progress));
            this.targetProgress = this.progress;
            this.state = isOpen ? DoorState.OPEN : DoorState.CLOSED;
            this.UpdateDoorPosition();
            return;
        }

        if (isOpen) {
            this.Open();
        } else {
            this.Close();
        }
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
        if (this.state === DoorState.OPEN || this.state === DoorState.OPENING) {
            this.Close();
        } else {
            this.Open();
        }
    }

    public Update(deltaTime: number): void {
        if (this.state === DoorState.OPENING) {
            this.progress += deltaTime * this.animationSpeed;
            if (this.progress >= 1) {
                this.progress = 1;
                this.state = DoorState.OPEN;
            }
            this.UpdateDoorPosition();
        } else if (this.state === DoorState.CLOSING) {
            this.progress -= deltaTime * this.animationSpeed;
            if (this.progress <= 0) {
                this.progress = 0;
                this.state = DoorState.CLOSED;
            }
            this.UpdateDoorPosition();
        }
    }

    public IsOpen(): boolean {
        return this.state === DoorState.OPEN;
    }

    public IsClosed(): boolean {
        return this.state === DoorState.CLOSED;
    }

    public GetProgress(): number {
        return this.progress;
    }

    public GetState(): DoorState {
        return this.state;
    }

    private LoadModel(): void {
        const loader = new GLTFLoader();

        loader.load(
            MODEL_URL,
            gltf => {
                this.doorModel = gltf.scene;
                
                // Ищем части модели по именам
                this.doorModel.traverse((child) => {
                    if (child.isObject3D) {
                        const name = child.name.toLowerCase();
                        if (name.includes("left") || name.includes("левая")) {
                            this.doorParts.leftDoor = child;
                        } else if (name.includes("right") || name.includes("правая")) {
                            this.doorParts.rightDoor = child;
                        } else if (name.includes("top") || name.includes("верх")) {
                            this.doorParts.topDoor = child;
                        } else if (name.includes("bottom") || name.includes("низ")) {
                            this.doorParts.bottomDoor = child;
                        } else if (name.includes("door") || name.includes("дверь")) {
                            this.doorParts.wholeDoor = child;
                        }
                    }
                });

                this.mesh.add(this.doorModel);
                this.UpdateDoorPosition();
            },
            undefined,
            error => {
                console.error("Ошибка загрузки space_door.glb:", error);
                this.CreateFallbackDoor();
            }
        );
    }

    // Запасной вариант, если модель не загрузилась
    private CreateFallbackDoor(): void {
        const geometry = new THREE.BoxGeometry(this.doorWidth, this.doorHeight, 0.2);
        const material = new THREE.MeshStandardMaterial({ 
            color: 0x446688,
            metalness: 0.8,
            roughness: 0.3
        });
        const door = new THREE.Mesh(geometry, material);
        this.mesh.add(door);
        this.doorParts.wholeDoor = door;
        this.UpdateDoorPosition();
    }

    private UpdateDoorPosition(): void {
        if (!this.doorModel && !this.doorParts.wholeDoor) return;

        if (this.doorParts.wholeDoor) {
            const offset = this.progress * this.doorHeight;
            this.doorParts.wholeDoor.position.y = -offset;
        }
    }
}
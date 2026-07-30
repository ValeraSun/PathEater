import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { clone as cloneSkeleton } from "three/examples/jsm/utils/SkeletonUtils.js";

const MODEL_URL = "/models/door.glb";

export class DoorView {
    public mesh = new THREE.Group();
    public isOpen: boolean = false;
    private basePosition: THREE.Vector3 = new THREE.Vector3(0, 0, 0);
    private openOffset: number = 3.0; // высота подъёма при открытии
    private model: THREE.Object3D | null = null;

    constructor() {
        console.log("DoorView: конструктор");
        // Создаём заглушку
        const geo = new THREE.BoxGeometry(2, 3, 0.2);
        const mat = new THREE.MeshStandardMaterial({ color: 0xff8800 });
        const cube = new THREE.Mesh(geo, mat);
        this.mesh.add(cube);
        this.model = cube;

        // Пытаемся загрузить модель
        this.LoadModel();
    }

    public SetBasePosition(x: number, y: number, z: number): void {
        this.basePosition.set(x, y, z);
        this.mesh.position.set(x, y, z);
    }

    public SetState(isOpen: boolean): void {
        this.isOpen = isOpen;
        
        // 1. Управление видимостью
        this.mesh.visible = !isOpen;
        
        // 2. Мгновенное изменение позиции (поднимаем/опускаем)
        if (isOpen) {
            this.mesh.position.y = this.basePosition.y + this.openOffset;
        } else {
            this.mesh.position.y = this.basePosition.y;
        }
        
        console.log(`DoorView: SetState(${isOpen}), позиция y = ${this.mesh.position.y}`);
    }

    private LoadModel(): void {
        const loader = new GLTFLoader();
        loader.load(MODEL_URL, (gltf) => {
            const model = cloneSkeleton(gltf.scene) as THREE.Group;
            // Удаляем заглушку
            if (this.model) {
                this.mesh.remove(this.model);
            }
            this.mesh.add(model);
            this.model = model;
            console.log("DoorView: модель загружена");
        }, undefined, (error) => {
            console.error("DoorView: ошибка загрузки модели", error);
        });
    }
}
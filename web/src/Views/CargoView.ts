import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export class CargoView 
{
    private group: THREE.Group;

    constructor() 
    {
        this.group = new THREE.Group();
        this.LoadModel();
    }

    private LoadModel(): void 
    {

        const loader = new GLTFLoader();

        loader.load("/models/cargo.glb", (gltf) => {

            gltf.scene.scale.set(1, 1, 1);

            this.group.add(gltf.scene);

        });

    }

    public GetObject(): THREE.Group 
    {
        return this.group;
    }
}
import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export class ButtonView 
{
    private group = new THREE.Group();

    public GetObject(): THREE.Group 
    {
        return this.group;
    }

    public constructor(x: number, y: number, z: number) 
    {
        this.LoadComputerModel(x, y, z);
    }

    private LoadComputerModel(x: number, y: number, z: number): void 
    {
        const loader = new GLTFLoader();

        loader.load("/models/industrial_button.glb",
            (gltf) => {
                const model = gltf.scene;
                model.position.set(x, y, z);
                model.scale.set(4, 4, 4);
                model.rotateX(Math.PI / 2)
                this.group.add(model);
                console.log("Button loaded");
            },
            undefined,
            (error) => {
                console.error(
                    "Ошибка загрузки bridge_console.glb:",
                    error
                );
            }
        );
    }
}
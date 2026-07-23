import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export class CargoView
{
    private group = new THREE.Group();
    private cargoModel: THREE.Object3D | null = null;
    private baggageStatus = 0;
    private maxVisibleCargo = 100;

    public constructor()
    {
        this.LoadModel();
    }

    public SetBaggageStatus(value: number): void
    {
        const nextStatus = THREE.MathUtils.clamp(Math.round(value), 0, 100);

        if (nextStatus === this.baggageStatus) return;

        this.baggageStatus = nextStatus;
        this.UpdateCargoObjects();
    }

    public GetObject(): THREE.Group
    {
        return this.group;
    }

    private LoadModel(): void
    {
        const loader = new GLTFLoader();

        loader.load("/models/cargo.glb",
            gltf =>
            {
                this.cargoModel = gltf.scene;
                this.cargoModel.scale.set(1.5, 1.5, 1.5);
                this.UpdateCargoObjects();
            },
            undefined,
            error =>
            {
                console.error(
                    "Ошибка загрузки cargo.glb:",
                    error
                );
            }
        );
    }

    private UpdateCargoObjects(): void
    {
        if (!this.cargoModel) {
            return;
        }

        this.group.clear();

        const visibleCount = this.baggageStatus === 0 ? 0 : Math.ceil(this.baggageStatus / 100 * this.maxVisibleCargo);
        const columns = 5;

        for (let index = 0; index < visibleCount; index++)
        {
            const cargo = this.cargoModel.clone(true);
            cargo.rotation.y = Math.PI;
            const x = index % columns + 1.5;
            const z = Math.floor(index / columns) - 3;
            cargo.position.set(x * 1.1,  1.3,  z * 1.1);

            this.group.add(cargo);
        }
    }
}
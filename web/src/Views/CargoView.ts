import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export class CargoView
{
    private readonly group = new THREE.Group();

    private cargoModel: THREE.Object3D | null = null;
    private baggageStatus = 0;

    private readonly maxVisibleCargo = 20;

    private readonly columns = 5;
    private readonly rows = 2;

    public constructor()
    {
        this.LoadModel();
    }

    public SetBaggageStatus(value: number): void
    {
        const nextStatus = THREE.MathUtils.clamp(Math.round(value), 0, 100);

        if (nextStatus === this.baggageStatus)
        {
            return;
        }

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
                this.cargoModel.scale.set(2, 2, 2);
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
        if (!this.cargoModel)
        {
            return;
        }

        this.group.clear();

        const visibleCount = this.baggageStatus === 0
                ? 0
                : Math.ceil(this.baggageStatus / 100 * this.maxVisibleCargo);

        for (let index = 0; index < visibleCount; index++)
        {
            const cargo = this.cargoModel.clone(true);
            const column = index % this.columns;
            const row = Math.floor(index / this.columns) % this.rows;

            const layer = Math.floor(index /(this.columns * this.rows));

            cargo.position.set(column * 2.2, layer * 2.1, row * 2.2 );

            this.group.add(cargo);
        }
    }
}
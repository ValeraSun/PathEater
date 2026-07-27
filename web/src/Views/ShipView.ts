import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { CargoView } from "./CargoView";
import type { ShipWireData } from "../Services/ShipWireData";

export interface ShipStateData
{
    baggage_status?: number;
    health?: number;
}

export class ShipView
{
    private group = new THREE.Group();
    private cargoView = new CargoView();
    private health = 100;

    public constructor()
    {
        this.LoadModel();
        const cargoObject = this.cargoView.GetObject();
        cargoObject.position.set(0, 0, 0);
        this.group.add(cargoObject);
    }

    public UpdateState(data: ShipWireData): void {
        if (typeof data.baggage_status === "number") {
            this.cargoView.SetBaggageStatus(data.baggage_status);
        }
        if (typeof data.health === "number") {
            this.health = data.health;
        }
    }

    public GetHealth(): number
    {
        return this.health;
    }

    public GetObject(): THREE.Group
    {
        return this.group;
    }

    private LoadModel(): void
    {
        const loader = new GLTFLoader();

        loader.load("/models/ship_hull_detailed_v8.glb",
            gltf =>
            {
                gltf.scene.scale.set(1, 1, 1);
                this.group.add(gltf.scene);
            },
            undefined,
            error =>
            {
                console.error("Ошибка загрузки ship_hull_detailed_v8.glb:", error);
            }
        );
    }
}
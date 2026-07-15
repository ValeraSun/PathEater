import * as THREE from "three";
import { MeshFactory } from "../Factories/MeshFactory";

export class SpaceView 
{
    public GetObject(): THREE.Group 
    {
        return this.group;
    }

    private group: THREE.Group;

    constructor() 
    {
        this.group = new THREE.Group();

        this.CreateFloor();
    }

    private CreateFloor(): void 
    {
        const floor = MeshFactory.CreatePlane(100, 100);
        floor.rotation.x = -Math.PI / 2;
        floor.position.y = -2;

        this.group.add(floor);
    }
}
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
        this.CreateStars();
    }

    private CreateFloor(): void 
    {
        const floor = MeshFactory.CreatePlane(100, 100);
        floor.rotation.x = -Math.PI / 2;
        floor.position.y = -2;

        this.group.add(floor);
    }

    private CreateStars(): void 
    {
        for (let i = 0; i < 100; i++) 
        {
            const star = MeshFactory.CreateStar();

            star.position.set(
                Math.random() * 100 - 50,
                Math.random() * 50,
                Math.random() * 100 - 50
            );

            this.group.add(star);
        }
    }
}
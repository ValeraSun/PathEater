import * as THREE from "three";

export class BreachView 
{
    Update(x: number, y: number, z: number) //1, 0, -2
    {
        this.mesh.position.set(x, y, z);
    }

    constructor() 
    {
        this.mesh = new THREE.Mesh(
            new THREE.PlaneGeometry(2, 2),
            new THREE.MeshBasicMaterial({
                color: 0x000000,
                side: THREE.DoubleSide
            })
        );
    }

    mesh: THREE.Mesh;

    public GetObject(): THREE.Mesh {
        return this.mesh;
    }
}
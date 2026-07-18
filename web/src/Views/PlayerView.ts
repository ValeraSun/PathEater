import * as THREE from "three";

export class PlayerView 
{
    constructor() 
    {
        this.mesh = new THREE.Mesh(
            new THREE.CapsuleGeometry(0.7, 1),
            new THREE.MeshStandardMaterial({color: 0x00ff00})
        );
    }

    mesh: THREE.Mesh;

    Update(position: {x: number, y: number, z: number})
    {
        this.mesh.position.set(position.x, position.y, position.z);
    }
}
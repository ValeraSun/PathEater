import * as THREE from "three";

export class BreakdownView {
    public mesh = new THREE.Group();

    constructor() {
        const geometry = new THREE.CircleGeometry(1.2, 24);
        const material = new THREE.MeshStandardMaterial({
            color: 0x050505,
            side: THREE.DoubleSide,
            roughness: 1,
            metalness: 0
        });

        const holeMesh = new THREE.Mesh(geometry, material);
        holeMesh.rotation.x = -Math.PI / 2;

        this.mesh.add(holeMesh);
    }

}
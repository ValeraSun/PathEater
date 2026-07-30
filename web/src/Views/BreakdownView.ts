import * as THREE from "three";

export class BreakdownView {
    public mesh = new THREE.Group();
    private sphereMesh: THREE.Mesh;  
    private radius = 1.2;

    constructor() {
        const geometry = new THREE.SphereGeometry(this.radius, 32, 16);
        const material = new THREE.MeshStandardMaterial({
            color: 0x050505,
            side: THREE.DoubleSide,  
            roughness: 1,
            metalness: 0,
            transparent: true,
            opacity: 0.8
        });

        this.sphereMesh = new THREE.Mesh(geometry, material);

        this.mesh.add(this.sphereMesh);
    }

    public SetRadius(radius: number): void {
        this.radius = radius;
        this.sphereMesh.scale.set(
            radius / 1.2,
            radius / 1.2,
            radius / 1.2
        );
    }

    public SetPosition(x: number, y: number, z: number): void {
        this.mesh.position.set(x, y, z);
    }
}
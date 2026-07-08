import * as THREE from "three";

export class MeshFactory 
{
    public static CreateBox(width: number, height: number, depth: number): THREE.Mesh 
    {
        const geometry = new THREE.BoxGeometry(width, height, depth);
        const material = new THREE.MeshStandardMaterial({color: 0x555555,});
        return new THREE.Mesh(geometry, material);
    }

    public static CreatePlane(width: number, height: number): THREE.Mesh 
    {
        const geometry = new THREE.PlaneGeometry(width, height);

        const material = new THREE.MeshStandardMaterial({
        color: 0x111111,
        side: THREE.DoubleSide,
        });

        return new THREE.Mesh(geometry, material);
    }

    public static CreateStar(): THREE.Mesh 
    {
        const geometry = new THREE.SphereGeometry(0.08, 8, 8);
        const material = new THREE.MeshBasicMaterial({
        color: 0xffffff,
        });

        return new THREE.Mesh(geometry, material);
    }
}
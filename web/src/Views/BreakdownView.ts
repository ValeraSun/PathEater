import * as THREE from "three";

export class BreakdownView {
    public mesh = new THREE.Group();
    private circleMesh: THREE.Mesh;
    private radius = 1.2;

    constructor() {
        const geometry = new THREE.CircleGeometry(this.radius, 24);
        const material = new THREE.MeshStandardMaterial({
            color: 0x050505,
            side: THREE.DoubleSide,
            roughness: 1,
            metalness: 0,
            transparent: true,
            opacity: 0.8
        });

        this.circleMesh = new THREE.Mesh(geometry, material);
        // НЕ вращаем - оставляем в плоскости XY (вертикально)
        // this.circleMesh.rotation.x = -Math.PI / 2; // Убираем эту строку

        this.mesh.add(this.circleMesh);
    }

    // Метод для обновления радиуса
    public SetRadius(radius: number): void {
        this.radius = radius;
        this.circleMesh.scale.set(
            radius / 1.2,
            radius / 1.2,
            radius / 1.2
        );
    }

    // Метод для установки позиции на стене
    public SetPosition(x: number, y: number, z: number): void {
        this.mesh.position.set(x, y, z);
    }

    // Метод для ориентации на стене (нормаль стены)
    public SetWallOrientation(normal: THREE.Vector3): void {
        // По умолчанию круг в плоскости XY (смотрит по Z)
        // Нам нужно повернуть его так, чтобы он смотрел по нормали стены
        const up = new THREE.Vector3(0, 0, 1);
        const quaternion = new THREE.Quaternion().setFromUnitVectors(up, normal.clone().normalize());
        this.mesh.quaternion.copy(quaternion);
    }
}
import * as THREE from "three";

export interface EntityModel {
    id: string;
    type: string;
    position: THREE.Vector3;
    targetPosition: THREE.Vector3;
    rotationY: number;
    targetRotationY: number;
    health?: number;
    dead?: boolean;
    attacking?: boolean;
    isOpen?: boolean;
}
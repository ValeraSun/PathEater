import * as THREE from "three";

export type Collider = 
{
    id: string;
    type: "solid" | "interactable";
    box: THREE.Box3;
};
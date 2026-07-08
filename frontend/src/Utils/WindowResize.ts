import * as THREE from "three";

export class WindowResize 
{
    public static Handle( camera: THREE.PerspectiveCamera, renderer: THREE.WebGLRenderer): void 
    {
        window.addEventListener("resize", () => {
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();

        renderer.setSize(window.innerWidth, window.innerHeight);
        });
    }
}
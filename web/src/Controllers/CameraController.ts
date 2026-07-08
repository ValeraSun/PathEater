import * as THREE from "three";
import { PointerLockControls } from "three/examples/jsm/controls/PointerLockControls.js";
import { PlayerModel } from "../Models/PlayerModel";

export class CameraController 
{
    private controls: PointerLockControls;
    private camera: THREE.PerspectiveCamera;
    private playerModel: PlayerModel;

    constructor(camera: THREE.PerspectiveCamera, playerModel: PlayerModel, domElement: HTMLElement) 
    {
        this.camera = camera;
        this.playerModel = playerModel;

        this.controls = new PointerLockControls(this.camera, domElement);

        document.addEventListener("click", () => {
            this.controls.lock();
        });
    }

    public Update(): void 
    {
        this.camera.position.set(
            this.playerModel.position.x,
            this.playerModel.position.y + 1,
            this.playerModel.position.z
        );
    }
}
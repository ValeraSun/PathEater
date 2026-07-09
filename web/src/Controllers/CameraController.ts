import * as THREE from "three";
import { PointerLockControls } from "three/examples/jsm/controls/PointerLockControls.js";
import { PlayerModel } from "../Models/PlayerModel";
import { CAMERA_SHIFT_Y, MOUSE_SENSITIVITY } from "../Config/CameraConfig";

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
        this.controls.pointerSpeed = MOUSE_SENSITIVITY;
        document.addEventListener("click", () => {
            this.controls.lock();
        });
    }

    public Update(): void 
    {
        this.camera.position.copy(this.playerModel.position);
        this.camera.position.y += CAMERA_SHIFT_Y;
    }
}
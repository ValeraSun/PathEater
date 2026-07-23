import * as THREE from "three";
import { PointerLockControls } from "three/examples/jsm/controls/PointerLockControls.js";
import { PlayerModel } from "../Models/PlayerModel";
import { CAMERA_SHIFT_Y, MOUSE_SENSITIVITY } from "../Config/CameraConfig";

const CAMERA_DISTANCE = 2;   
const CAMERA_HEIGHT = 1.5;  

export class CameraController
{
    private controls: PointerLockControls;
    private camera: THREE.PerspectiveCamera;
    private playerModel: PlayerModel;
    private forward = new THREE.Vector3();
    private desiredPosition = new THREE.Vector3();

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
        this.camera.getWorldDirection(this.forward);
        const pivot = this.playerModel.position;
        this.desiredPosition
            .copy(pivot)
            .addScaledVector(this.forward, -CAMERA_DISTANCE)
            .setY(pivot.y + CAMERA_SHIFT_Y + CAMERA_HEIGHT);

        this.camera.position.copy(this.desiredPosition);
    }
}
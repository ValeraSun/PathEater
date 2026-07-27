import * as THREE from "three";
import { PointerLockControls } from "three/examples/jsm/controls/PointerLockControls.js";
import { PlayerModel } from "../Models/PlayerModel";
import { MOUSE_SENSITIVITY } from "../Config/CameraConfig";

const CAMERA_EYE_HEIGHT = 1.9;
const CAMERA_FORWARD_OFFSET = 0.5;

export class CameraController
{
    private controls: PointerLockControls;
    private camera: THREE.PerspectiveCamera;
    private playerModel: PlayerModel;

    private forward = new THREE.Vector3();
    private desiredPosition = new THREE.Vector3();

    public constructor(camera: THREE.PerspectiveCamera, playerModel: PlayerModel, domElement: HTMLElement)
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

        this.desiredPosition
            .copy(this.playerModel.position)
            .addScaledVector(
                this.forward,
                CAMERA_FORWARD_OFFSET
            );

        this.desiredPosition.y = this.playerModel.position.y + CAMERA_EYE_HEIGHT;
        this.camera.position.copy(this.desiredPosition);
    }
}
import * as THREE from "three";
import { PointerLockControls } from "three/examples/jsm/controls/PointerLockControls.js";
import { MOUSE_SENSITIVITY } from "../Config/CameraConfig";
import { PlayerModel } from "../Models/PlayerModel";

const CAMERA_EYE_HEIGHT = 1.9;
const CAMERA_FORWARD_OFFSET = 0.5;

export class CameraController {
    private camera: THREE.PerspectiveCamera;
    private playerModel: PlayerModel;
    private controls: PointerLockControls;
    private interactionElement: HTMLElement;
    private horizontalDirection = new THREE.Vector3();
    private targetCameraPosition = new THREE.Vector3();

    public constructor(camera: THREE.PerspectiveCamera, playerModel: PlayerModel, interactionElement: HTMLElement) {
        this.camera = camera;
        this.playerModel = playerModel;
        this.interactionElement = interactionElement;
        this.controls = new PointerLockControls(this.camera, this.interactionElement);
        this.controls.pointerSpeed = MOUSE_SENSITIVITY;
        this.interactionElement.addEventListener("click", this.HandleInteractionElementClick);
    }

    public Update(): void {
        this.camera.getWorldDirection(this.horizontalDirection);
        this.horizontalDirection.y = 0;

        if (this.horizontalDirection.lengthSq() > 0) {
            this.horizontalDirection.normalize();
        }

        this.targetCameraPosition.copy(this.playerModel.position);
        this.targetCameraPosition.addScaledVector(this.horizontalDirection, CAMERA_FORWARD_OFFSET);
        this.targetCameraPosition.y = this.playerModel.position.y + CAMERA_EYE_HEIGHT;
        this.camera.position.copy(this.targetCameraPosition);
    }

    public Dispose(): void {
        this.interactionElement.removeEventListener("click", this.HandleInteractionElementClick);

        if (this.controls.isLocked) {
            this.controls.unlock();
        }

        this.controls.disconnect();
    }

    private HandleInteractionElementClick = (): void => {
        if (!this.controls.isLocked) {
            this.controls.lock();
        }
    };
}
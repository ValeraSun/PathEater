import * as THREE from "three";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { NetworkManager } from "../Services/NetworkManager";

export class PlayerController
{
    private model: PlayerModel;
    private view: PlayerView;
    private input: InputController;
    private camera: THREE.PerspectiveCamera;
    private networkManager: NetworkManager;
    private lastSendTime = 0;

    constructor(model: PlayerModel, view: PlayerView, input: InputController, camera: THREE.PerspectiveCamera, networkManager: NetworkManager) 
    {
        this.model = model;
        this.view = view;
        this.input = input;
        this.camera = camera;
        this.networkManager = networkManager;
    }

    Update()
    {
        const direction = new THREE.Vector3();

        this.camera.getWorldDirection(direction);

        direction.y = 0;
        let moved = false;
        direction.normalize();

        const right = new THREE.Vector3();
        right.crossVectors(direction, new THREE.Vector3(0, 1, 0)).normalize();

        if (this.input.keys.has("KeyW"))
        {
            this.model.position.x += direction.x * this.model.speed;
            this.model.position.z += direction.z * this.model.speed;
        }

        if (this.input.keys.has("KeyS"))
        {
            this.model.position.x -= direction.x * this.model.speed;
            this.model.position.z -= direction.z * this.model.speed;
        }

        if (this.input.keys.has("KeyA"))
        {
            this.model.position.x -= right.x * this.model.speed;
            this.model.position.z -= right.z * this.model.speed;
        }

        if (this.input.keys.has("KeyD"))
        {
            this.model.position.x += right.x * this.model.speed;
            this.model.position.z += right.z * this.model.speed;
        }

        const now = performance.now();

        if (moved && now - this.lastSendTime > 100) 
        {
            this.networkManager.SendMove(
                {
                    x: this.model.position.x,
                    y: this.model.position.y,
                    z: this.model.position.z
                },
                this.camera.rotation.y
            );

            this.lastSendTime = now;
        }

        this.view.Update(this.model.position);
    }
}
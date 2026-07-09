import * as THREE from "three";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { NetworkManager } from "../Services/NetworkManager";
import { CollisionManager } from "../Physics/CollisionManager";

export class PlayerController
{
    private model: PlayerModel;
    private view: PlayerView;
    private input: InputController;
    private camera: THREE.PerspectiveCamera;
    private networkManager: NetworkManager;
    private lastSendTime = 0;
    private collisionManager: CollisionManager;

    constructor(model: PlayerModel, view: PlayerView, input: InputController, camera: THREE.PerspectiveCamera, networkManager: NetworkManager, collisionManager: CollisionManager) 
    {
        this.model = model;
        this.view = view;
        this.input = input;
        this.camera = camera;
        this.networkManager = networkManager;
        this.collisionManager = collisionManager
    }

    public Update(): void 
    {
        const direction = new THREE.Vector3();

        this.camera.getWorldDirection(direction);
        direction.y = 0;
        direction.normalize();

        const right = new THREE.Vector3();
        right.crossVectors(direction, new THREE.Vector3(0, 1, 0)).normalize();

        const nextPosition = {
            x: this.model.position.x,
            y: this.model.position.y,
            z: this.model.position.z
        };

        let moved = false;

        if (this.input.IsKeyDown("KeyW")) {
            nextPosition.x += direction.x * this.model.speed;
            nextPosition.z += direction.z * this.model.speed;
            moved = true;
        }

        if (this.input.IsKeyDown("KeyS")) {
            nextPosition.x -= direction.x * this.model.speed;
            nextPosition.z -= direction.z * this.model.speed;
            moved = true;
        }

        if (this.input.IsKeyDown("KeyA")) {
            nextPosition.x -= right.x * this.model.speed;
            nextPosition.z -= right.z * this.model.speed;
            moved = true;
        }

        if (this.input.IsKeyDown("KeyD")) {
            nextPosition.x += right.x * this.model.speed;
            nextPosition.z += right.z * this.model.speed;
            moved = true;
        }

        const playerBox = new THREE.Box3().setFromCenterAndSize(
            new THREE.Vector3(
                nextPosition.x,
                nextPosition.y,
                nextPosition.z
            ),
            new THREE.Vector3(0.8, 1.8, 0.8)
        );

        if (moved && this.collisionManager.CanMove(playerBox)) {
            this.model.position.x = nextPosition.x;
            this.model.position.y = nextPosition.y;
            this.model.position.z = nextPosition.z;
        }

        this.view.Update(this.model.position);

        const now = performance.now();

        if (moved && now - this.lastSendTime > 100) {
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
    }
}
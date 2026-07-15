import * as THREE from "three";

import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { CollisionManager } from "../Physics/CollisionManager";
import { CollisionLayer } from "../Physics/OBB";
import { GameServerGateway } from "../Services/GameServerGateway";

export interface PlayerStatePayload {
    move_front: boolean;
    move_left: boolean;
    move_right: boolean;
    move_back: boolean;
    interact: boolean;
    attack: boolean;

    direction: {
        x: number;
        y: number;
        z: number;
    };
}

export class PlayerController {
    private readonly playerModel: PlayerModel;
    private readonly playerView: PlayerView;
    private readonly input: InputController;
    private readonly camera: THREE.Camera;
    private readonly gameServerGateway: GameServerGateway;
    private readonly collisionManager: CollisionManager;

    private inputLocked = false;

    private networkAccumulator = 0;
    private readonly networkInterval = 1 / 20;

    private readonly moveDirection = new THREE.Vector3();
    private readonly cameraDirection = new THREE.Vector3();
    private readonly cameraRight = new THREE.Vector3();
    private readonly displacement = new THREE.Vector3();

    public constructor(
        playerModel: PlayerModel,
        playerView: PlayerView,
        input: InputController,
        camera: THREE.Camera,
        gameServerGateway: GameServerGateway,
        collisionManager: CollisionManager
    ) {
        this.playerModel = playerModel;
        this.playerView = playerView;
        this.input = input;
        this.camera = camera;
        this.gameServerGateway = gameServerGateway;
        this.collisionManager = collisionManager;

        this.syncViewWithModel();
    }

    public Update(dt: number): void {
        if (!this.inputLocked) {
            this.updateMovement(dt);
        }

        this.updateNetwork(dt);
    }

    public LockInput(): void {
        this.inputLocked = true;
    }

    public UnlockInput(): void {
        this.inputLocked = false;
    }

    public IsInputLocked(): boolean {
        return this.inputLocked;
    }

    private updateMovement(dt: number): void {
        let inputX = 0;
        let inputZ = 0;

        if (this.input.IsKeyDown("KeyA")) {
            inputX -= 1;
        }

        if (this.input.IsKeyDown("KeyD")) {
            inputX += 1;
        }

        if (this.input.IsKeyDown("KeyW")) {
            inputZ += 1;
        }

        if (this.input.IsKeyDown("KeyS")) {
            inputZ -= 1;
        }

        if (inputX === 0 && inputZ === 0) {
            this.syncViewWithModel();
            return;
        }

        const inputLength = Math.hypot(inputX, inputZ);

        if (inputLength > 1) {
            inputX /= inputLength;
            inputZ /= inputLength;
        }

        this.calculateCameraDirections();

        this.moveDirection
            .set(0, 0, 0)
            .addScaledVector(
                this.cameraDirection,
                inputZ
            )
            .addScaledVector(
                this.cameraRight,
                inputX
            );

        if (this.moveDirection.lengthSq() === 0) {
            return;
        }

        this.moveDirection.normalize();

        this.displacement
            .copy(this.moveDirection)
            .multiplyScalar(
                this.playerModel.speed * dt
            );

        this.applyMovement();
        this.updatePlayerRotation();
    }

    private calculateCameraDirections(): void {
        this.camera.getWorldDirection(
            this.cameraDirection
        );

        this.cameraDirection.y = 0;

        if (this.cameraDirection.lengthSq() === 0) {
            this.cameraDirection.set(0, 0, -1);
        } else {
            this.cameraDirection.normalize();
        }

        this.cameraRight
            .crossVectors(
                this.cameraDirection,
                THREE.Object3D.DEFAULT_UP
            )
            .normalize();
    }

    private applyMovement(): void {
        this.collisionManager.MoveAndSlide(
            this.playerModel.body,
            this.displacement,
            CollisionLayer.Static
        );

        this.syncViewWithModel();
    }

    private syncViewWithModel(): void {
        this.playerView.mesh.position.copy(
            this.playerModel.position
        );
    }

    private updatePlayerRotation(): void {
        this.playerView.mesh.rotation.y = Math.atan2(
            this.moveDirection.x,
            this.moveDirection.z
        );
    }

    private updateNetwork(dt: number): void {
        this.networkAccumulator += dt;

        if (
            this.networkAccumulator <
            this.networkInterval
        ) {
            return;
        }

        this.networkAccumulator %= this.networkInterval;

        this.sendPlayerState();
    }

    private sendPlayerState(): void {
        this.camera.getWorldDirection(
            this.cameraDirection
        );

        if (this.cameraDirection.lengthSq() === 0) {
            this.cameraDirection.set(0, 0, -1);
        } else {
            this.cameraDirection.normalize();
        }

        const canSendInput = !this.inputLocked;

        const state: PlayerStatePayload = {
            move_front:
                canSendInput &&
                this.input.IsKeyDown("KeyW"),

            move_left:
                canSendInput &&
                this.input.IsKeyDown("KeyA"),

            move_right:
                canSendInput &&
                this.input.IsKeyDown("KeyD"),

            move_back:
                canSendInput &&
                this.input.IsKeyDown("KeyS"),

            interact:
                canSendInput &&
                this.input.WasPressedOnce("KeyE"),

            attack:
                canSendInput &&
                this.input.WasPressedOnce("Space"),

            direction: {
                x: this.cameraDirection.x,
                y: this.cameraDirection.y,
                z: this.cameraDirection.z
            }
        };

        this.gameServerGateway.SendPlayerState(state);
    }
}

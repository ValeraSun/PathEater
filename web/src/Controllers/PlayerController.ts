import * as THREE from "three";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { GameServerGateway } from "../Services/GameServerGateway";

export interface PlayerStatePayload {
    move_front: boolean;
    move_left: boolean;
    move_right: boolean;
    move_back: boolean;
    interact: boolean;
    attack: boolean;
    direction: { x: number; y: number; z: number };
}

export class PlayerController {
    private readonly playerModel: PlayerModel;
    private readonly playerView: PlayerView;
    private readonly input: InputController;
    private readonly camera: THREE.Camera;
    private readonly gameServerGateway: GameServerGateway;

    private inputLocked = false;

    private networkAccumulator = 0;
    private readonly networkInterval = 1 / 20;

    private readonly cameraDirection = new THREE.Vector3();

    public constructor(
        playerModel: PlayerModel,
        playerView: PlayerView,
        input: InputController,
        camera: THREE.Camera,
        gameServerGateway: GameServerGateway
    ) {
        this.playerModel = playerModel;
        this.playerView = playerView;
        this.input = input;
        this.camera = camera;
        this.gameServerGateway = gameServerGateway;

        this.syncViewWithModel();
    }

    public Update(dt: number): void {
        this.syncViewWithModel();
        this.updateAnimation();
        this.playerView.AdvanceAnimation(dt);
        this.updateNetwork(dt);
    }

    public LockInput(): void { this.inputLocked = true; }
    public UnlockInput(): void { this.inputLocked = false; }
    public IsInputLocked(): boolean { return this.inputLocked; }

    private syncViewWithModel(): void {
        this.playerView.mesh.position.copy(this.playerModel.position);
    }

    private updateNetwork(dt: number): void {
        this.networkAccumulator += dt;
        if (this.networkAccumulator < this.networkInterval) return;
        this.networkAccumulator %= this.networkInterval;
        this.sendPlayerState();
    }

    private updateAnimation(): void {
        const isMoving = !this.inputLocked && (
            this.input.IsKeyDown("KeyW") ||
            this.input.IsKeyDown("KeyA") ||
            this.input.IsKeyDown("KeyS") ||
            this.input.IsKeyDown("KeyD")
        );

        this.playerView.SetMoving(isMoving);
    }

    private sendPlayerState(): void {
        this.camera.getWorldDirection(this.cameraDirection);
        if (this.cameraDirection.lengthSq() === 0) {
            this.cameraDirection.set(0, 0, -1);
        } else {
            this.cameraDirection.normalize();
        }

        const canSendInput = !this.inputLocked;

        const state: PlayerStatePayload = {
            move_front: canSendInput && this.input.IsKeyDown("KeyW"),
            move_left: canSendInput && this.input.IsKeyDown("KeyA"),
            move_right: canSendInput && this.input.IsKeyDown("KeyD"),
            move_back: canSendInput && this.input.IsKeyDown("KeyS"),
            interact: canSendInput && this.input.WasPressedOnce("KeyE"),
            attack: canSendInput && this.input.WasPressedOnce("Space"),
            direction: {
                x: this.cameraDirection.x,
                y: this.cameraDirection.y,
                z: this.cameraDirection.z
            }
        };

        this.gameServerGateway.SendPlayerState(state);
    }
}
import * as THREE from "three";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { GameServerGateway } from "../Network/GameServerGateway";
import { type PlayerStatePayload } from "../Network/ServerContracts";

export class PlayerController {
    public constructor(
        playerModel: PlayerModel,
        playerView: PlayerView,
        input: InputController,
        camera: THREE.Camera,
        gameServerGateway: GameServerGateway
    ) {
        this.playerModel = playerModel;
        this.playerView = playerView;
        this.playerView.SetVisible(false);
        this.input = input;
        this.camera = camera;
        this.gameServerGateway = gameServerGateway;

        this.SyncViewWithModel();
    }

    public Update(dt: number): void {
        this.SyncViewWithModel();
        this.UpdateAnimation();
        this.playerView.AdvanceAnimation(dt);
        this.UpdateNetwork(dt);
    }

    public LockInput(): void { this.inputLocked = true; }
    public UnlockInput(): void { this.inputLocked = false; }
    public IsInputLocked(): boolean { return this.inputLocked; }

    private SyncViewWithModel(): void {
        this.playerView.mesh.position.copy(this.playerModel.position);
    }

    private UpdateNetwork(dt: number): void {
        this.networkAccumulator += dt;
        if (this.networkAccumulator < this.networkInterval) return;
        this.networkAccumulator %= this.networkInterval;
        this.SendPlayerState();
    }

    private UpdateAnimation(): void {
        const isMoving = !this.inputLocked && (
            this.input.IsKeyDown("KeyW") ||
            this.input.IsKeyDown("KeyA") ||
            this.input.IsKeyDown("KeyS") ||
            this.input.IsKeyDown("KeyD")
        );

        this.playerView.SetMoving(isMoving);
    }

    private SendPlayerState(): void {
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
            interact: canSendInput && this.input.IsKeyDown("KeyE"),
            attack: canSendInput && this.input.IsKeyDown("Mouse0"),
            direction: {
                x: this.cameraDirection.x,
                y: this.cameraDirection.y,
                z: this.cameraDirection.z
            }
        };

        this.gameServerGateway.SendPlayerState(state);
    }

    private playerModel: PlayerModel;
    private playerView: PlayerView;
    private input: InputController;
    private camera: THREE.Camera;
    private gameServerGateway: GameServerGateway;
    private inputLocked = false;
    private networkAccumulator = 0;
    private networkInterval = 1 / 20;
    private cameraDirection = new THREE.Vector3();
}
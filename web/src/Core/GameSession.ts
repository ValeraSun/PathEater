import * as THREE from "three";
import { CameraController } from "../Controllers/CameraController";
import { ComputerController } from "../Controllers/ComputerController";
import { InputController } from "../Controllers/InputController";
import { InteractionController } from "../Controllers/InteractionController";
import { PlayerController } from "../Controllers/PlayerController";
import { ComputerModel } from "../Models/ComputerModel";
import { PlayerModel } from "../Models/PlayerModel";
import { GameServerGateway } from "../Network/GameServerGateway";
import type { EntityTransformData, GameStartedPayload } from "../Network/ServerContracts";
import { GameView } from "../Views/GameView";
import { InteractionView } from "../Views/InteractionView";
import { PlayerView } from "../Views/PlayerView";
import { EntityMessageHandler } from "./EntityMessageHandler";
import { GameLoop } from "./GameLoop";

const POSITION_INTERPOLATION_SPEED = 12;
const ROTATION_INTERPOLATION_SPEED = 12;
const DAMAGE_GLITCH_BASE_STRENGTH = 0.35;
const DAMAGE_GLITCH_HEALTH_SCALE = 100;

export class GameSession {
    private readonly gameView: GameView;
    private readonly inputController: InputController;
    private readonly gameServerGateway: GameServerGateway;
    private readonly entityMessageHandler: EntityMessageHandler;
    private readonly computerModel: ComputerModel;
    private readonly interactionView: InteractionView;

    private gameLoop: GameLoop | null = null;
    private playerModel: PlayerModel | null = null;
    private playerView: PlayerView | null = null;
    private playerController: PlayerController | null = null;
    private cameraController: CameraController | null = null;
    private interactionController: InteractionController | null = null;
    private computerController: ComputerController | null = null;

    private readonly targetPlayerPosition = new THREE.Vector3();
    private targetPlayerRotationY = 0;
    private lastKnownHealth: number | null = null;

    public constructor(
        gameView: GameView,
        inputController: InputController,
        gameServerGateway: GameServerGateway,
        entityMessageHandler: EntityMessageHandler,
        computerModel: ComputerModel,
        interactionView: InteractionView
    ) {
        this.gameView = gameView;
        this.inputController = inputController;
        this.gameServerGateway = gameServerGateway;
        this.entityMessageHandler = entityMessageHandler;
        this.computerModel = computerModel;
        this.interactionView = interactionView;
    }

    public Start(payload: GameStartedPayload): void {
        this.Stop();

        const spawnPosition = new THREE.Vector3(payload.spawn.x, payload.spawn.y, payload.spawn.z);

        this.playerModel = new PlayerModel(spawnPosition);
        this.playerView = new PlayerView();
        this.targetPlayerPosition.copy(spawnPosition);
        this.targetPlayerRotationY = 0;
        this.lastKnownHealth = null;
        this.gameView.GetScene().add(this.playerView.mesh);

        this.cameraController = new CameraController(
            this.gameView.GetCamera(),
            this.playerModel,
            this.gameView.GetRendererDomElement(),
        );

        this.playerController = new PlayerController(
            this.playerModel,
            this.playerView,
            this.inputController,
            this.gameView.GetCamera(),
            this.gameServerGateway
        );

        this.computerController = new ComputerController(this.inputController);

        this.interactionController = new InteractionController(
            this.playerModel,
            this.computerModel,
            this.interactionView
        );

        this.gameLoop = new GameLoop(
            [
                {
                    Update: (deltaTime: number): void => {
                        this.UpdateLocalPlayer(deltaTime);
                    }
                },
                {
                    Update: (deltaTime: number): void => {
                        this.playerController?.Update(deltaTime);
                    }
                },
                {
                    Update: (): void => {
                        this.cameraController?.Update();
                    }
                },
                {
                    Update: (): void => {
                        this.interactionController?.Update();
                    }
                },
                {
                    Update: (): void => {
                        this.computerController?.Update();
                    }
                },
                {
                    Update: (deltaTime: number): void => {
                        this.entityMessageHandler.UpdateEntityViews(deltaTime);
                    }
                },
                {
                    Update: (): void => {
                        this.HandleAttackInput();
                    }
                }
            ],
            (deltaTime: number): void => {
                this.gameView.Render(deltaTime);
            },
            (): void => {
                this.inputController.FinishFrame();
            }
        );

        this.gameView.ShowPlayerHealth();
        this.gameLoop.Start();
    }

    public Stop(): void {
        this.gameLoop?.Stop();
        this.gameLoop = null;

        this.cameraController?.Dispose();
        this.computerController?.Dispose();
        this.interactionController?.Dispose();

        if (this.playerView) {
            this.gameView.GetScene().remove(this.playerView.mesh);
        }

        this.entityMessageHandler.Clear();
        this.gameView.HidePlayerHealth();

        this.playerModel = null;
        this.playerView = null;
        this.playerController = null;
        this.cameraController = null;
        this.interactionController = null;
        this.computerController = null;
    }

    public ApplyLocalPlayerState(playerState: EntityTransformData): void {
        if (playerState.position) {
            this.targetPlayerPosition.set(
                playerState.position.x,
                playerState.position.y,
                playerState.position.z
            );
        }

        if (playerState.rotation) {
            this.targetPlayerRotationY = Math.atan2(
                playerState.rotation.x,
                playerState.rotation.z
            );
        }

        if (typeof playerState.health === "number") {
            this.HandleHealthChanged(playerState.health);
        }
    }

    private HandleHealthChanged(newHealth: number): void {
        if (this.lastKnownHealth !== null && newHealth < this.lastKnownHealth) {
            const healthLost = this.lastKnownHealth - newHealth;
            const strength = Math.min(1, DAMAGE_GLITCH_BASE_STRENGTH + healthLost / DAMAGE_GLITCH_HEALTH_SCALE);
            this.gameView.TriggerPlayerDamage(strength);
        }

        this.lastKnownHealth = newHealth;
        this.gameView.SetPlayerHealth(newHealth);
    }

    private HandleAttackInput(): void {
        if (this.inputController.WasPressedOnce("KeyE")) {
            this.gameView.TriggerMonsterBloodSplatter();
        }
    }

    private UpdateLocalPlayer(deltaTime: number): void {
        if (!this.playerModel) {
            return;
        }

        const positionInterpolation = 1 - Math.exp(-POSITION_INTERPOLATION_SPEED * deltaTime);

        this.playerModel.position.lerp(
            this.targetPlayerPosition,
            positionInterpolation
        );

        if (!this.playerView) {
            return;
        }

        const rotationInterpolation = 1 - Math.exp(-ROTATION_INTERPOLATION_SPEED * deltaTime);

        this.playerView.mesh.rotation.y = this.InterpolateAngle(
            this.playerView.mesh.rotation.y,
            this.targetPlayerRotationY,
            rotationInterpolation
        );
    }

    private InterpolateAngle(currentAngle: number, targetAngle: number, interpolation: number): number {
        const angleDifference = Math.atan2(Math.sin(targetAngle - currentAngle),Math.cos(targetAngle - currentAngle));
        return currentAngle + angleDifference * interpolation;
    }
}
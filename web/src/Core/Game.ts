import * as THREE from "three";
import { GameView } from "../Views/GameView";
import { InputController } from "../Controllers/InputController";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { PlayerController } from "../Controllers/PlayerController";
import { EntityManager } from "../Services/EntityManager";
import { CollisionManager } from "../Physics/CollisionManager";
import { CameraController } from "../Controllers/CameraController";
import { MAX_DELTA_TIME, MILLISECONDS_IN_SECOND } from "../Config/GameConfig";
import { InteractionView } from "../Views/InteractionView";
import { InteractionController } from "../Controllers/InteractionController";
import { ComputerController } from "../Controllers/ComputerController";
import { GameServerGateway, type GameStartedPayload } from "../Services/GameServerGateway";

export class Game {
    private static instance: Game;

    private readonly input = new InputController();
    private readonly collisionManager = new CollisionManager();
    private readonly interactionView = new InteractionView();
    private readonly gameView: GameView;
    private readonly entityManager: EntityManager;
    private readonly gameServerGateway: GameServerGateway;

    private playerModel: PlayerModel | null = null;
    private playerView: PlayerView | null = null;
    private playerController: PlayerController | null = null;
    private cameraController: CameraController | null = null;
    private interactionController: InteractionController | null = null;
    private computerController: ComputerController | null = null;

    private lastTime = performance.now();
    private isRunning = false;

    public static GetInstance(): Game {
        if (!Game.instance) {
            Game.instance = new Game();
        }

        return Game.instance;
    }

    private constructor() {
        this.gameView = new GameView();

        this.entityManager = new EntityManager(
            this.gameView.GetScene()
        );

        this.gameServerGateway = new GameServerGateway(
            this.entityManager
        );
    }

    public GetGateway(): GameServerGateway {
        return this.gameServerGateway;
    }

    public MountTo(container: HTMLElement): void {
        container.appendChild(
            this.gameView.GetRendererDomElement()
        );
    }

    public async Connect(): Promise<void> {
        this.gameServerGateway.InitListeners();
        await this.gameServerGateway.ConnectToServer();
    }

    public async StartMatch(payload: GameStartedPayload): Promise<void> {
        if (this.isRunning) {
            return;
        }

        await this.collisionManager.LoadShipColliders(
            "/data/ship_wall_colliders_v3.json"
        );

        const spawn = new THREE.Vector3(
            payload.spawn.x,
            payload.spawn.y,
            payload.spawn.z
        );

        this.playerModel = new PlayerModel(spawn);
        this.playerView = new PlayerView();

        this.gameView.AttachPlayerView(this.playerView);

        this.cameraController = new CameraController(
            this.gameView.GetCamera(),
            this.playerModel,
            this.gameView.GetRendererDomElement()
        );

        this.playerController = new PlayerController(
            this.playerModel,
            this.playerView,
            this.input,
            this.gameView.GetCamera(),
            this.gameServerGateway,
            this.collisionManager
        );

        const computerView = this.gameView.GetComputerView();

        this.computerController = new ComputerController(
            this.input,
            this.playerController,
            computerView.GetId()
        );

        this.interactionController = new InteractionController(
            this.input,
            this.playerModel,
            computerView,
            this.interactionView,
            () => {
                this.computerController!.Enter();
            }
        );

        this.collisionManager.AddDynamic(
            this.playerModel.body
        );

        this.playerView.mesh.position.copy(
            this.playerModel.position
        );

        this.gameView.Init();

        this.isRunning = true;
        this.lastTime = performance.now();
        this.GameLoop();
    }

    private GameLoop = (): void => {
        if (!this.isRunning) {
            return;
        }

        requestAnimationFrame(this.GameLoop);

        const now = performance.now();

        const dt = Math.min(
            (now - this.lastTime) / MILLISECONDS_IN_SECOND,
            MAX_DELTA_TIME
        );

        this.lastTime = now;

        this.playerController?.Update(dt);
        this.cameraController?.Update();
        this.interactionController?.Update();
        this.computerController?.Update();
        this.entityManager.Update(dt);

        this.gameView.Render();
    };
}
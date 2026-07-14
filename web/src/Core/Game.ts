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
import { WebSocketClient } from "../Services/WebSocketClient";
import { GameServerGateway } from "../Services/GameServerGateway";

export class Timer {
    static wait(seconds: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, seconds * 1000));
    }
    
    static waitMs(milliseconds: number): Promise<void> {
        return new Promise(resolve => setTimeout(resolve, milliseconds));
    }
}

export class Game 
{
    private static instance: Game;
    private readonly input = new InputController();
    private readonly playerModel = new PlayerModel();
    private readonly playerView = new PlayerView();
    private readonly collisionManager = new CollisionManager();
    private readonly interactionView = new InteractionView();
    private readonly gameView: GameView;
    private readonly entityManager: EntityManager;
    private readonly webSocketClient: WebSocketClient;
    private readonly gameServerGateway: GameServerGateway;
    private readonly cameraController: CameraController;
    private readonly playerController: PlayerController;
    private readonly interactionController: InteractionController;
    private readonly computerController: ComputerController;

    private lastTime = performance.now();

    public static GetInstance(): Game {
        if (!Game.instance) {
            Game.instance = new Game();
        }

        return Game.instance;
    }

    private constructor() {
        this.gameView = new GameView(this.playerView);

        this.entityManager = new EntityManager(
            this.gameView.GetScene()
        );

        this.webSocketClient = new WebSocketClient();

        this.gameServerGateway = new GameServerGateway(
            this.webSocketClient,
            this.entityManager
        );


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
            this.webSocketClient,
            this.collisionManager
        );

        const computerView = this.gameView.GetComputerView();

        this.computerController =
            new ComputerController(
                this.input,
                this.playerController,
                computerView.GetId()
            );

        this.interactionController =
            new InteractionController(
                this.input,
                this.playerModel,
                computerView,
                this.interactionView,
                () => {
                    this.computerController.Enter();
                }
            );

        this.collisionManager.AddDynamic(
            this.playerModel.body
        );

        this.playerView.mesh.position.copy(
            this.playerModel.position
        );
    }

    public async Start(): Promise<void> {
        await this.collisionManager.LoadShipColliders(
            "/data/ship_wall_colliders_v3.json"
        );

        this.connectToServer();

        this.gameView.Init();

        this.lastTime = performance.now();
        this.GameLoop();
    }

    private connectToServer(): void {
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const url = `${protocol}//${window.location.host}/ws`;

        this.webSocketClient.connect(url);
    }

    private GameLoop = (): void => {
        requestAnimationFrame(this.GameLoop);

        const now = performance.now();

        const dt = Math.min(
            (now - this.lastTime) / MILLISECONDS_IN_SECOND,
            MAX_DELTA_TIME
        );

        this.lastTime = now;

        this.playerController.Update(dt);
        this.cameraController.Update();
        this.interactionController.Update();
        this.computerController.Update();

        this.gameView.Render();
    };
}

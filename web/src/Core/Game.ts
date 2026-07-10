import { GameView } from "../Views/GameView";
import { InputController } from "../Controllers/InputController";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { PlayerController } from "../Controllers/PlayerController";
import { NetworkManager } from "../Services/NetworkManager";
import { EntityManager } from "../Services/EntityManager";
import { CollisionManager } from "../Physics/CollisionManager";
import { CameraController } from "../Controllers/CameraController";
import { MAX_DELTA_TIME, MILLISECONDS_IN_SECOND } from "../Config/GameConfig";

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
    private input = new InputController();
    private playerModel = new PlayerModel();
    private playerView = new PlayerView();
    private cameraController: CameraController;
    private gameView: GameView;
    private playerController: PlayerController;
    private entityManager: EntityManager;
    private networkManager: NetworkManager;
    private collisionManager = new CollisionManager();

    private lastTime = performance.now();

    public static GetInstance(): Game 
    {
        if (!Game.instance) 
        {
            Game.instance = new Game();
        }
        return Game.instance;
    }

    public async Start(): Promise<void> {
        await this.collisionManager.LoadShipColliders("/data/ship_wall_colliders_v1.json");
        this.gameView.Init();
        //убрать перед защитой
        this.collisionManager.AddDebugHelpers(this.gameView.GetScene());
        this.networkManager.Connect();
        while (!this.networkManager.Connected) {await Timer.waitMs(100);}
        this.networkManager.CreateRoom();
        this.networkManager.CreateGameSession();
        this.GameLoop();
    }

    private constructor() 
    {
        this.gameView = new GameView(this.playerView);
        this.entityManager = new EntityManager(this.gameView.GetScene());
        this.networkManager = new NetworkManager(this.entityManager);
        this.cameraController = new CameraController(this.gameView.GetCamera(), 
                                                     this.playerModel,     
                                                     this.gameView.GetRendererDomElement()
                                                    );
        this.playerController = new PlayerController(
            this.playerModel,
            this.playerView,
            this.input,
            this.gameView.GetCamera(),
            this.networkManager,
            this.collisionManager
        );
    }

    private GameLoop = (): void => 
    {
        requestAnimationFrame(this.GameLoop);
        const now = performance.now();
        const dt = Math.min((now - this.lastTime) / MILLISECONDS_IN_SECOND, MAX_DELTA_TIME);
        this.lastTime = now;
        this.playerController.Update(dt);
        this.cameraController.Update();
        this.gameView.Render();
    };
}
import { GameView } from "../Views/GameView";
import { InputController } from "../Controllers/InputController";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { PlayerController } from "../Controllers/PlayerController";
import { CameraController } from "../Controllers/CameraController";
import { NetworkManager } from "../Services/NetworkManager";
import { EntityManager } from "../Services/EntityManager";
import { CollisionManager } from "../Physics/CollisionManager";

export class Game 
{
    private input = new InputController();
    private playerModel = new PlayerModel();
    private playerView = new PlayerView();
    private cameraController: CameraController;
    private static instance: Game;
    private gameView: GameView;
    private playerController: PlayerController;
    private entityManager: EntityManager;
    private networkManager: NetworkManager;
    private collisionManager = new CollisionManager();

    public static GetInstance(): Game 
    {
        if (!Game.instance) 
        {
            Game.instance = new Game();
        }
        return Game.instance;
    }

    public async Start(): Promise<void> {
        await this.collisionManager.LoadShipColliders();

        this.gameView.Init();

        this.collisionManager.AddDebugHelpers(this.gameView.GetScene());

        this.Animate();
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

    private Animate = (): void => 
    {
        requestAnimationFrame(this.Animate);
        this.playerController.Update();
        this.cameraController.Update();
        this.gameView.Render();
    };
}
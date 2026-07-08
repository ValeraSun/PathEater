import { GameView } from "../Views/GameView";
import { InputController } from "../Controllers/InputController";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { PlayerController } from "../Controllers/PlayerController";
import { CameraController } from "../Controllers/CameraController";

export class Game 
{
    private input = new InputController();
    private playerModel = new PlayerModel();
    private playerView = new PlayerView();
    private cameraController: CameraController;
    private static instance: Game;
    private gameView: GameView;
    private playerController: PlayerController;

    public static GetInstance(): Game 
    {
        if (!Game.instance) 
        {
            Game.instance = new Game();
        }
        return Game.instance;
    }

    public Start(): void 
    {
        this.gameView.Init();
        this.Animate();
    }

    private constructor() 
    {
        this.gameView = new GameView(this.playerView);
        this.cameraController = new CameraController(this.gameView.GetCamera(), 
                                                     this.playerModel,     
                                                     this.gameView.GetRendererDomElement()
                                                    );
        this.playerController = new PlayerController(
            this.playerModel,
            this.playerView,
            this.input,
            this.gameView.GetCamera()
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
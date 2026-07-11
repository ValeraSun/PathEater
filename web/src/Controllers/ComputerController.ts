import { InputController } from "./InputController";
import { NetworkManager } from "../Services/NetworkManager";
import { PlayerController } from "./PlayerController";

export class ComputerController 
{
    private active = false;
    private input: InputController;
    private networkManager: NetworkManager;
    private playerController: PlayerController;
    private computerId: string;
        
    public constructor(input: InputController,
        networkManager: NetworkManager,
        playerController: PlayerController,
        computerId: string
    ) 
    {
        this.input = input;
        this.networkManager = networkManager;
        this.playerController = playerController;
        this.computerId = computerId;
    }

    public Enter(): void 
    {
        this.active = true;
        this.playerController.LockInput();

        console.log("Управление кораблём включено");
    }

    public Exit(): void 
    {
        if (!this.active) 
        {
            return;
        }

        this.active = false;
        this.playerController.UnlockInput();

        this.networkManager.SendReleaseComputer(
            this.computerId
        );

        console.log("Управление кораблём выключено");
    }

    public Update(): void 
    {
        if (!this.active) 
        {
            return;
        }

        if (this.input.WasPressedOnce("Escape")) 
        {
            this.Exit();
            return;
        }

        let inputX = 0;
        let inputZ = 0;

        if (this.input.IsKeyDown("KeyA")) 
        {
            inputX -= 1;
        }

        if (this.input.IsKeyDown("KeyD")) 
        {
            inputX += 1;
        }

        if (this.input.IsKeyDown("KeyW")) 
        {
            inputZ -= 1;
        }

        if (this.input.IsKeyDown("KeyS")) 
        {
            inputZ += 1;
        }

        const length = Math.hypot(inputX, inputZ);

        if (length > 1) 
        {
            inputX /= length;
            inputZ /= length;
        }

        this.networkManager.SendShipInput(
            this.computerId,
            inputX,
            inputZ
        );
    }

    public IsActive(): boolean 
    {
        return this.active;
    }
}
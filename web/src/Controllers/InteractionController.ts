import { InputController } from "./InputController";
import { PlayerModel } from "../Models/PlayerModel";
import { ComputerView } from "../Views/ComputerView";
import { InteractionView } from "../Views/InteractionView";
import { NetworkManager } from "../Services/NetworkManager";

export class InteractionController 
{
    private interactionDistance = 2.5;
    private input: InputController;
    private playerModel: PlayerModel;
    private computerView: ComputerView;
    private interactionView: InteractionView;
    private networkManager: NetworkManager;

    public constructor(input: InputController, 
        playerModel: PlayerModel, 
        computerView: ComputerView, 
        interactionView: InteractionView, 
        networkManager:NetworkManager
    ) 
    {
        this.input = input;
        this.playerModel = playerModel;
        this.computerView = computerView;
        this.interactionView = interactionView;
        this.networkManager = networkManager;
    }

    public Update(): void 
    {
        const distance = this.playerModel.position.distanceTo(
            this.computerView.GetInteractionPosition()
        );

        const isNear = distance <= this.interactionDistance;
        const isFree = this.computerView.IsFree();

        if (!isNear || !isFree) 
        {
            this.interactionView.Hide();
            return;
        }

        this.interactionView.Show();

        if (this.input.WasPressedOnce("KeyE")) 
        {
            this.networkManager.SendUseComputer(
                this.computerView.GetId()
            );
        }
    }
}
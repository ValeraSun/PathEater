import { PlayerModel } from "../Models/PlayerModel";
import { ComputerView } from "../Views/ComputerView";
import { InteractionView } from "../Views/InteractionView";

export class InteractionController {
    private readonly interactionDistance = 2.5;

    private playerModel: PlayerModel;
    private computerView: ComputerView;
    private interactionView: InteractionView;

    public constructor(playerModel: PlayerModel,computerView: ComputerView,interactionView: InteractionView) {
        this.playerModel = playerModel;
        this.computerView = computerView;
        this.interactionView = interactionView;
    }

    public Update(): void {
        const distance = this.playerModel.position.distanceTo(
            this.computerView.GetInteractionPosition()
        );

        const isNear = distance <= this.interactionDistance;

        const isFree = this.computerView.IsFree();
        
        if (!isNear || !isFree) {
            this.interactionView.Hide();
            return;
        }

        this.interactionView.Show();

        // if (this.input.WasPressedOnce("KeyE")) {
        //     this.interactionView.Hide();
        //     this.onInteract();
        // }
    }
}
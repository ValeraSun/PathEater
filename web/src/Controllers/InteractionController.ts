import { ComputerModel } from "../Models/ComputerModel";
import { PlayerModel } from "../Models/PlayerModel";
import { InteractionView } from "../Views/InteractionView";

const COMPUTER_INTERACTION_DISTANCE = 2.5;

export class InteractionController {
    private playerModel: PlayerModel;
    private computerModel: ComputerModel;
    private interactionView: InteractionView;

    public constructor(playerModel: PlayerModel, computerModel: ComputerModel, interactionView: InteractionView) {
        this.playerModel = playerModel;
        this.computerModel = computerModel;
        this.interactionView = interactionView;
    }

    public Update(): void {
        const distanceToComputer = this.playerModel.position.distanceTo(this.computerModel.interactionPosition);
        const playerIsCloseEnough = distanceToComputer <= COMPUTER_INTERACTION_DISTANCE;
        const computerIsAvailable = this.computerModel.IsAvailable();

        if (!playerIsCloseEnough || !computerIsAvailable) {
            this.interactionView.Hide();
            return;
        }

        this.interactionView.Show("Наведитесь на терминал и нажмите E, чтобы управлять кораблём");
    }

    public Dispose(): void {
        this.interactionView.Hide();
    }
}
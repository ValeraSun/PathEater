import { PlayerModel } from "../Models/PlayerModel";
import { ComputerView } from "../Views/ComputerView";
import { InteractionView } from "../Views/InteractionView";
import { InputController } from "./InputController";

export class InteractionController {
    private readonly interactionDistance = 2.5;

    private readonly input: InputController;
    private readonly playerModel: PlayerModel;
    private readonly computerView: ComputerView;
    private readonly interactionView: InteractionView;
    private readonly onInteract: () => void;

    public constructor(
        input: InputController,
        playerModel: PlayerModel,
        computerView: ComputerView,
        interactionView: InteractionView,
        onInteract: () => void
    ) {
        this.input = input;
        this.playerModel = playerModel;
        this.computerView = computerView;
        this.interactionView = interactionView;
        this.onInteract = onInteract;
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

        if (this.input.WasPressedOnce("KeyE")) {
            this.interactionView.Hide();
            this.onInteract();
        }
    }
}
import { InputController } from "./InputController";
import { PlayerController } from "./PlayerController";

export class ComputerController {
    private active = false;

    private input: InputController;
    private playerController: PlayerController;
    private computerId: string;

    public constructor(
        input: InputController,
        playerController: PlayerController,
        computerId: string
    ) {
        this.input = input;
        this.playerController = playerController;
        this.computerId = computerId;
    }

    public Enter(): void {
        if (this.active) {
            return;
        }

        this.active = true;
        this.playerController.LockInput();

        console.log(
            `Управление компьютером ${this.computerId} включено`
        );
    }

    public Exit(): void {
        if (!this.active) {
            return;
        }

        this.active = false;
        this.playerController.UnlockInput();

        console.log(
            `Управление компьютером ${this.computerId} выключено`
        );
    }

    public Update(): void {
        if (!this.active) {
            return;
        }

        if (this.input.WasPressedOnce("Escape")) {
            this.Exit();
            return;
        }
    }

    public IsActive(): boolean {
        return this.active;
    }
}
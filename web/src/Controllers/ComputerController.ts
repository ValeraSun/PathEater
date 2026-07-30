import { InputController } from "./InputController";

export class ComputerController {
    public constructor(inputController: InputController) {
        this.inputController = inputController;
    }

    public EnterControlMode(): void {
        if (this.active) {
            return;
        }
        this.active = true;
    }

    public ExitControlMode(): void {
        if (!this.active) {
            return;
        }
        this.active = false;
    }

    public Update(): void {
        if (!this.active) {
            return;
        }

        if (this.inputController.WasPressedOnce("Escape")) {
            this.ExitControlMode();
        }
    }

    public IsControlModeActive(): boolean {
        return this.active;
    }

    public Dispose(): void {
        this.active = false;
    }

    private inputController: InputController;
    private active = false;
}
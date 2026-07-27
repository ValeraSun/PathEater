import { InputController } from "./InputController";

export class ComputerController {
    private active = false;
    private input: InputController;

    public constructor(input: InputController) {
        this.input = input;
    }

    public Enter(): void {
        if (this.active) {
            return;
        }
        this.active = true;
    }

    public Exit(): void {
        if (!this.active) {
            return;
        }
        this.active = false;
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
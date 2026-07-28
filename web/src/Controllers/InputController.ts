export class InputController {
    private heldKeys = new Set<string>();
    private pressedKeys = new Set<string>();

    public constructor() {
        window.addEventListener("keydown", this.HandleKeyDown);
        window.addEventListener("keyup", this.HandleKeyUp);
        window.addEventListener("blur", this.HandleWindowBlur);
    }

    public IsKeyDown(keyCode: string): boolean {
        return this.heldKeys.has(keyCode);
    }

    public WasPressedOnce(keyCode: string): boolean {
        return this.pressedKeys.has(keyCode);
    }

    public FinishFrame(): void {
        this.pressedKeys.clear();
    }

    public Dispose(): void {
        window.removeEventListener("keydown", this.HandleKeyDown);
        window.removeEventListener("keyup", this.HandleKeyUp);
        window.removeEventListener("blur", this.HandleWindowBlur);

        this.heldKeys.clear();
        this.pressedKeys.clear();
    }

    private HandleKeyDown = (event: KeyboardEvent): void => {
        if (!this.heldKeys.has(event.code)) {
            this.pressedKeys.add(event.code);
        }

        this.heldKeys.add(event.code);
    };

    private HandleKeyUp = (event: KeyboardEvent): void => {
        this.heldKeys.delete(event.code);
    };

    private HandleWindowBlur = (): void => {
        this.heldKeys.clear();
        this.pressedKeys.clear();
    };
}
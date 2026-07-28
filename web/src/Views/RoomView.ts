export class RoomView {
    private createRoomButton: HTMLButtonElement;
    private joinRoomButton: HTMLButtonElement;
    private backButton: HTMLButtonElement;
    private roomCodeInput: HTMLInputElement;
    private statusElement: HTMLElement;

    public constructor() {
        this.createRoomButton = this.GetRequiredButton("create-room-button");
        this.joinRoomButton = this.GetRequiredButton("join-room-button");
        this.backButton = this.GetRequiredButton("room-back-button");
        this.roomCodeInput = this.GetRequiredInput("room-code-input");
        this.statusElement = this.GetRequiredElement("room-status");
        this.roomCodeInput.addEventListener("input", this.HandleRoomCodeChanged);
    }

    public SetCreateRoomHandler(handler: () => void): void {
        this.createRoomButton.addEventListener("click", handler);
    }

    public SetJoinRoomHandler(handler: () => void): void {
        this.joinRoomButton.addEventListener("click",handler );
    }

    public SetBackHandler(handler: () => void): void {
        this.backButton.addEventListener("click",handler);
    }

    public GetRoomCode(): string {
        return this.roomCodeInput.value.trim();
    }

    public FocusRoomCodeInput(): void {
        this.roomCodeInput.focus();
    }

    public SetStatus(message: string): void {
        this.statusElement.textContent = message;
    }

    public SetControlsEnabled(enabled: boolean): void {
        this.createRoomButton.disabled = !enabled;
        this.joinRoomButton.disabled = !enabled;
        this.roomCodeInput.disabled = !enabled;
    }

    private HandleRoomCodeChanged = (): void => {
        this.roomCodeInput.value = this.roomCodeInput.value.replace(/\s/g, "")
    };

    private GetRequiredButton(elementId: string): HTMLButtonElement {
        const element = document.getElementById(elementId);

        if (!(element instanceof HTMLButtonElement)) {
            throw new Error(`Не найдена кнопка #${elementId}`);
        }

        return element;
    }

    private GetRequiredInput(elementId: string): HTMLInputElement {
        const element = document.getElementById(elementId);

        if (!(element instanceof HTMLInputElement)) {
            throw new Error(
                `Не найден input #${elementId}`
            );
        }

        return element;
    }

    private GetRequiredElement(elementId: string): HTMLElement {
        const element = document.getElementById(elementId);

        if (!element) {
            throw new Error(
                `Не найден элемент #${elementId}`
            );
        }

        return element;
    }
}
export class MenuView {
    private playButton: HTMLButtonElement;
    private statusElement: HTMLElement;

    public constructor() {
        this.playButton = this.GetRequiredButton("play-button");
        this.statusElement = this.GetRequiredElement("menu-status");
    }

    public SetPlayHandler(handler: () => void): void {
        this.playButton.addEventListener("click", handler);
    }

    public SetConnectionPending(connectionPending: boolean): void {
        this.playButton.disabled = connectionPending;
        this.statusElement.textContent = connectionPending ? "Идёт подключение к серверу..." : "";
    }

    public ShowConnectionError(errorMessage: string): void {
        this.playButton.disabled = false;
        this.statusElement.textContent = errorMessage;
    }

    private GetRequiredButton(elementId: string): HTMLButtonElement {
        const element = document.getElementById(elementId);

        if (!(element instanceof HTMLButtonElement)) {
            throw new Error(`Не найдена кнопка #${elementId}`);
        }

        return element;
    }

    private GetRequiredElement(elementId: string): HTMLElement {
        const element = document.getElementById(elementId);

        if (!element) {
            throw new Error(`Не найден элемент #${elementId}`);
        }
        return element;
    }
}
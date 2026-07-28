const MAX_PLAYERS = 4;

export class LobbyView {
    private playerListElement: HTMLElement;
    private statusElement: HTMLElement;
    private leaveButton: HTMLButtonElement;
    private startGameButton: HTMLButtonElement;
    private hostMode = false;

    public constructor() {
        this.playerListElement = this.GetRequiredElement("player-list");
        this.statusElement = this.GetRequiredElement("lobby-status");
        this.leaveButton = this.GetRequiredButton("leave-lobby-button");
        this.startGameButton = this.GetRequiredButton("start-game-button");
    }

    public SetLeaveHandler(handler: () => void): void {
        this.leaveButton.addEventListener("click", handler);
    }

    public SetStartGameHandler(handler: () => void): void {
        this.startGameButton.addEventListener("click", handler);
    }

    public SetHostMode(hostMode: boolean): void {
        this.hostMode = hostMode;
        this.startGameButton.hidden = !hostMode;
    }

    public IsHostMode(): boolean {
        return this.hostMode;
    }

    public SetStatus(message: string): void {
        this.statusElement.textContent = message;
    }

    public SetStartButtonEnabled(enabled: boolean): void {
        this.startGameButton.disabled = !enabled;
    }

    public RenderPlayers(playerIds: string[], localPlayerId: string | null): void {
        this.playerListElement.replaceChildren();

        playerIds.forEach((playerId: string, playerIndex: number): void => {
                const playerElement = this.CreatePlayerElement(playerId, playerIndex, localPlayerId);
                this.playerListElement.appendChild(playerElement);
            }
        );

        for (let playerIndex = playerIds.length; playerIndex < MAX_PLAYERS; playerIndex++) {
            const emptySlot = document.createElement("li");
            emptySlot.className = "empty-player";
            emptySlot.textContent = "Ожидание игрока";

            this.playerListElement.appendChild(emptySlot);
        }
    }

    private CreatePlayerElement(playerId: string, playerIndex: number, localPlayerId: string | null): HTMLLIElement {
        const isLocalPlayer = playerId === localPlayerId;
        const listItem = document.createElement("li");
        listItem.className = "player";
        const avatarContainer = document.createElement("div");
        avatarContainer.className = "player-avatar";
        const avatarImage = document.createElement("img");
        avatarImage.src = "/images/avatar.png";
        avatarImage.alt = "Аватар игрока";

        avatarContainer.appendChild(avatarImage);

        const playerInformation = document.createElement("div");

        playerInformation.className = "player-info";
        const playerName = document.createElement("strong");

        playerName.textContent = isLocalPlayer ? "Вы" : `Игрок ${playerIndex + 1}`;
        const readyStatus = document.createElement("span");

        readyStatus.className = "ready";
        readyStatus.textContent = "Готов";

        playerInformation.append(playerName, readyStatus);
        const playerRole = document.createElement("span");
        playerRole.className = "player-role";

        playerRole.textContent = isLocalPlayer && this.hostMode ? "Капитан" : "Игрок";
        listItem.append(avatarContainer, playerInformation, playerRole);

        return listItem;
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
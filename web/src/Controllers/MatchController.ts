import { Game } from "../Core/Game";
import { GameServerGateway } from "../Network/GameServerGateway";
import type { GameOverPayload, GameStartedPayload}  from "../Network/ServerContracts";
import { MusicManager } from "../Services/MusicManager";
import { AppView } from "../Views/AppView";

export class MatchController {
    private application: Game;
    private gameServerGateway: GameServerGateway;
    private musicManager: MusicManager;
    private appView: AppView;
    private matchTimerElement: HTMLElement;

    public constructor(application: Game, gameServerGateway: GameServerGateway, musicManager: MusicManager, appView: AppView) {
        this.application = application;
        this.gameServerGateway = gameServerGateway;
        this.musicManager = musicManager;
        this.appView = appView;
        this.matchTimerElement = this.GetRequiredElement("match-timer");
    }

    public Initialize(): void {
        this.gameServerGateway.SetGameStartedHandler(this.HandleGameStarted);
        this.gameServerGateway.SetGameOverHandler(this.HandleGameOver);
        this.gameServerGateway.SetMatchTimerHandler(this.HandleMatchTimerChanged);
    }

    private HandleGameStarted = (payload: GameStartedPayload): void => {
        void this.StartMatch(payload);
    };

    private HandleGameOver = (payload: GameOverPayload): void => {
        this.application.StopMatch();
        this.appView.ShowScreen(payload.win ? "win" : "lose");
        void this.musicManager.PlayMusic("menu");
    };

    private HandleMatchTimerChanged = (remainingSeconds: number): void => {
        this.matchTimerElement.textContent = this.FormatMatchTime(remainingSeconds);
    };

    private async StartMatch(payload: GameStartedPayload): Promise<void> {
        this.matchTimerElement.textContent = "00:00";
        await this.musicManager.PlayMusic("game");
        this.application.StartMatch(payload);
        this.appView.ShowScreen("game");
    }

    private FormatMatchTime(totalSeconds: number): string {
        const safeSeconds = Math.max(0, Math.floor(totalSeconds));
        const minutes = Math.floor(safeSeconds / 60);
        const seconds = safeSeconds % 60;

        return (
            `${String(minutes).padStart(2, "0")}:` + String(seconds).padStart(2, "0")
        );
    }

    private GetRequiredElement(elementId: string): HTMLElement {
        const element = document.getElementById(elementId);

        if (!element) {
            throw new Error(`Не найден элемент #${elementId}`);
        }

        return element;
    }
}
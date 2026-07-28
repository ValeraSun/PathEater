import type { GameStartedPayload } from "../Network/ServerContracts";
import { GameServerGateway } from "../Network/GameServerGateway";
import { MusicManager } from "../Services/MusicManager";
import { GameView } from "../Views/GameView";
import { GameSession } from "./GameSession";

export class Game {
    private gameServerGateway: GameServerGateway;
    private gameView: GameView;
    private gameSession: GameSession;
    private musicManager: MusicManager;

    public constructor(gameServerGateway: GameServerGateway, gameView: GameView, gameSession: GameSession, musicManager: MusicManager) {
        this.gameServerGateway = gameServerGateway;
        this.gameView = gameView;
        this.gameSession = gameSession;
        this.musicManager = musicManager;
    }

    public InitializeView(gameContainer: HTMLElement, healthBarContainer: HTMLElement): void {
        this.gameView.Initialize(gameContainer, healthBarContainer);
    }

    public async ConnectToServer(): Promise<void> {
        await this.gameServerGateway.Connect();
    }

    public StartMatch(payload: GameStartedPayload): void {
        this.gameSession.Start(payload);
    }

    public StopMatch(): void {
        this.gameSession.Stop();
    }

    public GetGameServerGateway(): GameServerGateway {
        return this.gameServerGateway;
    }

    public GetMusicManager(): MusicManager {
        return this.musicManager;
    }
}
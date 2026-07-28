import { GameServerGateway } from "../Network/GameServerGateway";
import { MusicManager } from "../Services/MusicManager";
import { AppView } from "../Views/AppView";
import { LobbyView } from "../Views/LobbyView";
import { RoomView } from "../Views/RoomView";

export class LobbyController {
    private gameServerGateway: GameServerGateway;
    private musicManager: MusicManager;
    private appView: AppView;
    private roomView: RoomView;
    private lobbyView: LobbyView;

    public constructor(
        gameServerGateway: GameServerGateway,
        musicManager: MusicManager,
        appView: AppView,
        roomView: RoomView,
        lobbyView: LobbyView
    ) {
        this.gameServerGateway = gameServerGateway;
        this.musicManager = musicManager;
        this.appView = appView;
        this.roomView = roomView;
        this.lobbyView = lobbyView;
    }

    public Initialize(): void {
        this.lobbyView.SetLeaveHandler(this.HandleLeaveLobby);
        this.lobbyView.SetStartGameHandler(this.HandleStartGame);
        this.gameServerGateway.SetRoomPlayersHandler(this.HandleRoomPlayersChanged);
    }

    private HandleLeaveLobby = (): void => {
        this.gameServerGateway.ExitRoom();
        this.lobbyView.SetHostMode(false);
        this.lobbyView.SetStatus("");
        this.lobbyView.SetStartButtonEnabled(true);
        this.roomView.SetControlsEnabled(true);
        this.roomView.SetStatus("");
        void this.musicManager.PlayMusic("menu");
        this.appView.ShowScreen("room");
    };

    private HandleStartGame = (): void => {
        if (!this.lobbyView.IsHostMode()) {
            return;
        }

        this.lobbyView.SetStartButtonEnabled(false);
        this.lobbyView.SetStatus("Запуск игры");
        this.gameServerGateway.StartGame();
    };

    private HandleRoomPlayersChanged = (playerIds: string[]): void => {
        this.lobbyView.RenderPlayers(playerIds, this.gameServerGateway.localPlayerId);
    };
}
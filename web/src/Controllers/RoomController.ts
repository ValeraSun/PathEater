import { GameServerGateway } from "../Network/GameServerGateway";
import { MusicManager } from "../Services/MusicManager";
import { AppView } from "../Views/AppView";
import { LobbyView } from "../Views/LobbyView";
import { RoomView } from "../Views/RoomView";

export class RoomController {
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
        this.roomView.SetCreateRoomHandler(this.HandleCreateRoom);
        this.roomView.SetJoinRoomHandler(this.HandleJoinRoom);
        this.roomView.SetBackHandler(this.HandleBackToMenu);
    }

    private HandleCreateRoom =
        async (): Promise<void> => {
            this.roomView.SetControlsEnabled(false);
            this.roomView.SetStatus("Создаём комнату");

            try {
                const room = await this.gameServerGateway.CreateRoom();
                this.lobbyView.SetHostMode(true);
                this.lobbyView.SetStatus(`Код комнаты: ${room.roomId}`);
                this.appView.ShowScreen("lobby");
            } catch (error: unknown) {
                this.ShowRequestError(error, "Не удалось создать комнату");
            }
        };

    private HandleJoinRoom =
        async (): Promise<void> => {
            const roomCode = this.roomView.GetRoomCode();

            if (!roomCode) {
                this.roomView.SetStatus("Введите код комнаты");
                this.roomView.FocusRoomCodeInput();
                return;
            }

            this.roomView.SetControlsEnabled(false);
            this.roomView.SetStatus(`Подключение к комнате ${roomCode}...`);

            try {
                await this.gameServerGateway.JoinRoom(roomCode);

                this.lobbyView.SetHostMode(false);
                this.lobbyView.SetStatus("Ожидание запуска игры командиром...");
                this.appView.ShowScreen("lobby");
            } catch (error: unknown) {
                this.ShowRequestError(error, "Комната не найдена");
            }
        };

    private HandleBackToMenu = (): void => {
        void this.musicManager.PlayMusic("menu");
        this.appView.ShowScreen("menu");
    };

    private ShowRequestError(error: unknown, fallbackMessage: string): void {
        const errorMessage = error instanceof Error ? error.message : fallbackMessage;
        this.roomView.SetStatus(errorMessage);
        this.roomView.SetControlsEnabled(true);
    }

    private gameServerGateway: GameServerGateway;
    private musicManager: MusicManager;
    private appView: AppView;
    private roomView: RoomView;
    private lobbyView: LobbyView;
}
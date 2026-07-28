import { EntityMessageHandler } from "../Core/EntityMessageHandler";
import type { GameOverPayload, GameStartedPayload, PlayerStatePayload, RoomInfoPayload } from "./ServerContracts";
import { IsDeleteEntityPayload, IsEntityInfo, IsGameOverPayload, IsGameStartedPayload, IsMatchTimerPayload, IsRoomInfoPayload, IsRoomPlayersPayload } from "./ServerValidators";
import { WebSocketClient } from "../Services/WebSocketClient";

const REQUEST_TIMEOUT_MILLISECONDS = 5000;

export class GameServerGateway {
    public localPlayerId: string | null = null;
    private webSocketClient: WebSocketClient;
    private entityMessageHandler: EntityMessageHandler;
    private listenersInitialized = false;
    private roomPlayersHandler: ((playerIds: string[]) => void) | null = null;
    private gameStartedHandler: ((payload: GameStartedPayload) => void) | null = null;
    private gameOverHandler: ((payload: GameOverPayload) => void) | null = null;
    private matchTimerHandler: ((seconds: number) => void) | null = null;

    public constructor(entityMessageHandler: EntityMessageHandler) {
        this.entityMessageHandler = entityMessageHandler;
        this.webSocketClient = new WebSocketClient();
        this.InitializeListeners();
    }

    public async Connect(): Promise<void> {
        const webSocketProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const serverUrl = `${webSocketProtocol}//` + `${window.location.host}/ws`;

        await this.webSocketClient.connect(serverUrl);
    }

    public SetRoomPlayersHandler(handler: (playerIds: string[]) => void): void {
        this.roomPlayersHandler = handler;
    }

    public SetGameStartedHandler(handler: (payload: GameStartedPayload) => void): void {
        this.gameStartedHandler = handler;
    }

    public SetGameOverHandler(handler: (payload: GameOverPayload) => void): void {
        this.gameOverHandler = handler;
    }

    public SetMatchTimerHandler(handler: (seconds: number) => void): void {
        this.matchTimerHandler = handler;
    }

    public async CreateRoom():
        Promise<RoomInfoPayload> {
        this.webSocketClient.send("createRoom",{});
        const response = await this.WaitForResponse("SuccessCreateRoom");

        if (!IsRoomInfoPayload(response)) {
            throw new Error(
                "Сервер вернул некорректные данные"
            );
        }

        if (response.error) {
            throw new Error(response.error);
        }
        this.RememberLocalPlayer(response);

        return response;
    }

    public async JoinRoom(roomId: string): Promise<RoomInfoPayload> {
        this.webSocketClient.send("joinRoom", {roomID: roomId});

        const response = await this.WaitForResponse("SuccessJoinRoom");

        if (!IsRoomInfoPayload(response)) {
            throw new Error(
                "Сервер вернул некорректные данные"
            );
        }

        if (response.error) {
            throw new Error(response.error);
        }
        this.RememberLocalPlayer(response);

        return response;
    }

    public ExitRoom(): void {
        this.webSocketClient.send("exitRoom", {});
        this.localPlayerId = null;
        this.entityMessageHandler.Clear();
    }

    public StartGame(): void {
        this.webSocketClient.send("createGameSession", {});
    }

    public SendPlayerState(playerState: PlayerStatePayload): void {
        this.webSocketClient.send("playerState", playerState);
    }

    private InitializeListeners(): void {
        if (this.listenersInitialized) {
            return;
        }

        this.listenersInitialized = true;

        this.webSocketClient.on("CreateEntity",
            (payload: unknown): void => {
                if (!IsEntityInfo(payload)) {
                    console.log("Некорректное создание сущности:", payload);
                    return;
                }

                this.entityMessageHandler.CreateEntity(payload);
            }
        );

        this.webSocketClient.on("UpdateEntity",
            (payload: unknown): void => {
                if (!IsEntityInfo(payload)) {
                    console.log("Некорректное обновление сущности:", payload);
                    return;
                }

                this.entityMessageHandler.UpdateEntity(payload);
            }
        );

        this.webSocketClient.on("DeleteEntity",
            (payload: unknown): void => {
                if (!IsDeleteEntityPayload(payload)) {
                    return;
                }
                this.entityMessageHandler.RemoveEntity(payload.id);
            }
        );

        this.webSocketClient.on("RoomPlayers",
            (payload: unknown): void => {
                if (!IsRoomPlayersPayload(payload)) {
                    return;
                }
                this.roomPlayersHandler?.(payload.playerIds);
            }
        );

        this.webSocketClient.on("GameStarted",
            (payload: unknown): void => {
                if (!IsGameStartedPayload(payload)) {
                    console.log("Некорректные данные начала игры:", payload);
                    return;
                }

                this.localPlayerId = payload.playerId;
                this.entityMessageHandler.SetLocalPlayerId(payload.playerId);
                this.gameStartedHandler?.(payload);
            }
        );

        this.webSocketClient.on("GameOver",
            (payload: unknown): void => {
                if (!IsGameOverPayload(payload)) {
                    console.log("Некорректные данные окончания игры:", payload);
                    return;
                }
                this.gameOverHandler?.(payload);
            }
        );

        this.webSocketClient.on("Time",
            (payload: unknown): void => {
                if (!IsMatchTimerPayload(payload)) {
                    return;
                }
                this.matchTimerHandler?.(payload.time);
            }
        );
    }

    private RememberLocalPlayer(roomInformation: RoomInfoPayload): void {
        this.localPlayerId = roomInformation.playerId;
        this.entityMessageHandler.SetLocalPlayerId(roomInformation.playerId);
    }

    private WaitForResponse(eventName: string): Promise<unknown> {
        return new Promise(
            (
                resolve: (value: unknown) => void,
                reject: (reason?: unknown) => void
            ): void => {
                const timeoutId = window.setTimeout(
                    (): void => {
                        reject(
                            new Error(
                                "Сервер не отвечает"
                            )
                        );
                    },
                    REQUEST_TIMEOUT_MILLISECONDS
                );

                this.webSocketClient.once(
                    eventName,
                    (payload: unknown): void => {
                        window.clearTimeout(timeoutId);
                        resolve(payload);
                    }
                );
            }
        );
    }
}
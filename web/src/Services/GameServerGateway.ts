import { EntityManager } from "./EntityManager";
import { WebSocketClient } from "./WebSocketClient";
import type { PlayerStatePayload } from "../Controllers/PlayerController";

export interface EntityInfo {
    id: string;
    type: string;
    data: unknown;
}

export interface DeleteEntityPayload {
    id: string;
}

export interface SnapshotPayload {
    entities: EntityInfo[];
}

export interface RoomInfoPayload {
    roomId: string;
    playerId: string;
    error?: string;
}

export interface GameStartedPayload {
    playerId: string;
    spawn: { x: number; y: number; z: number };
}

const REQUEST_TIMEOUT_MS = 5000;

function timeout(ms: number): Promise<never> {
    return new Promise((_, reject) =>
        setTimeout(() => reject(new Error("Сервер не отвечает")), ms)
    );
}

export class GameServerGateway {
    private wsClient: WebSocketClient;
    private entityManager: EntityManager;

    public localPlayerId: string | null = null;
    public onGameStarted: ((payload: GameStartedPayload) => void) | null = null;

    constructor(entityManager: EntityManager) {
        this.wsClient = new WebSocketClient();
        this.entityManager = entityManager;
    }

    public async ConnectToServer(): Promise<void> {
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const url = `${protocol}//${window.location.host}/ws`;

        await this.wsClient.connect(url);
    }

    public InitListeners(): void {
        this.wsClient.on("CreateEntity", (payload: unknown) => {
            if (!this.isValidEntityInfo(payload)) return;
            this.entityManager.CreateEntity(payload.id, payload.type, payload.data as any);
        });

        this.wsClient.on("UpdateEntity", (payload: unknown) => {
            if (!this.isValidEntityInfo(payload)) return;
            this.entityManager.UpdateEntity(payload.id, payload.type, payload.data as any);
        });

        this.wsClient.on("DeleteEntity", (payload: unknown) => {
            if (!this.isValidDeletePayload(payload)) return;
            this.entityManager.DeleteEntity(payload.id);
        });

        this.wsClient.on("Snapshot", (payload: unknown) => {
            if (!this.isValidSnapshotPayload(payload)) return;

            this.entityManager.ApplySnapshot(
                payload.entities.map(entity => ({
                    id: entity.id,
                    type: entity.type,
                    data: entity.data as any
                }))
            );
        });

        this.wsClient.on("GameStarted", (payload: GameStartedPayload) => {
            this.localPlayerId = payload.playerId;
            this.entityManager.SetLocalPlayerId(payload.playerId);
            this.onGameStarted?.(payload);
        });

        this.wsClient.on("error", (payload: any) => {
            console.error("Ошибка от сервера:", payload?.message ?? payload);
        });
    }

    private waitFor(event: string): Promise<any> {
        const answer = new Promise(resolve => {
            this.wsClient.once(event, resolve);
        });

        return Promise.race([answer, timeout(REQUEST_TIMEOUT_MS)]);
    }

    public async CreateRoom(): Promise<RoomInfoPayload> {
        this.wsClient.send("createRoom", {});
        const info = await this.waitFor("SuccessCreateRoom");

        if (info.error) throw new Error(info.error);

        this.rememberPlayer(info);
        return info;
    }

    public async JoinRoom(roomId: string): Promise<RoomInfoPayload> {
        this.wsClient.send("joinRoom", { roomID: roomId });
        const info = await this.waitFor("SuccessJoinRoom");

        if (info.error) throw new Error(info.error);

        this.rememberPlayer(info);
        return info;
    }

    public ExitRoom(): void {
        this.wsClient.send("exitRoom", {});

        this.localPlayerId = null;
    }

    public ExitMenu(): void {
        this.wsClient.send("exitMenu", {});
    }

    public DeleteRoom(roomId: string): void {
        this.wsClient.send("deleteRoom", {
            roomID: roomId
        });
    }

    public StartGame(): void {
        this.wsClient.send("createGameSession", {});
    }

    public SendPlayerState(state: PlayerStatePayload): void {
        this.wsClient.send("playerState", state);
    }

    private rememberPlayer(info: RoomInfoPayload): void {
        this.localPlayerId = info.playerId;
        this.entityManager.SetLocalPlayerId(info.playerId);
    }

    private isValidEntityInfo(payload: any): payload is EntityInfo {
        return payload
            && typeof payload === "object"
            && typeof payload.id === "string"
            && typeof payload.type === "string"
            && "data" in payload;
    }

    private isValidDeletePayload(payload: any): payload is DeleteEntityPayload {
        return payload
            && typeof payload === "object"
            && typeof payload.id === "string";
    }

    private isValidSnapshotPayload(payload: any): payload is SnapshotPayload {
        return payload
            && typeof payload === "object"
            && Array.isArray(payload.entities);
    }
}
import { EntityManager } from "./EntityManager";
import { WebSocketClient } from "./WebSocketClient";
import { EntityParser } from "./EntityParser";

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

export class GameServerGateway {
    constructor(
        private wsClient: WebSocketClient,
        private entityManager: EntityManager
    ) {
        this.initListeners();
    }

    private initListeners() {
        this.wsClient.on("CreateEntity", (payload: unknown) => {
            if (!this.isValidEntityInfo(payload)) return; 

            try {
                const parsedData = EntityParser.Parse(payload.type, payload.data);
                this.entityManager.CreateEntity(payload.id, payload.type, parsedData);
            } catch (error) {
                console.error("Ошибка парсинга CreateEntity:", error);
            }
        });

        this.wsClient.on("UpdateEntity", (payload: unknown) => {
            if (!this.isValidEntityInfo(payload)) return;
            
            try {
                const parsedData = EntityParser.Parse(payload.type, payload.data);
                this.entityManager.UpdateEntity(payload.id, payload.type, parsedData);
            } catch (error) {
                console.error("Ошибка парсинга UpdateEntity:", error);
            }
        });

        this.wsClient.on("DeleteEntity", (payload: unknown) => {
            if (!this.isValidDeletePayload(payload)) return;
            
            this.entityManager.DeleteEntity(payload.id);
        });

        this.wsClient.on("Snapshot", (payload: unknown) => {
            if (!this.isValidSnapshotPayload(payload)) return;

            try {
                const parsedEntities = payload.entities.map(entity => ({
                    id: entity.id,
                    type: entity.type,
                    data: EntityParser.Parse(entity.type, entity.data)
                }));
                
                this.entityManager.ApplySnapshot(parsedEntities);
            } catch (error) {
                console.error("Ошибка применения Snapshot:", error);
            }
        });
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
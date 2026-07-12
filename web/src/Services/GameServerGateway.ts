export class GameServerGateway {
    private wsClient: WebSocketClient;
    private entityManager: EntityManager;

    constructor(wsClient: WebSocketClient, entityManager: EntityManager) {
        this.wsClient = wsClient;
        this.entityManager = entityManager;

        // Подписываемся на CreateEntity
        this.wsClient.on("CreateEntity", (payload: EntityInfo) => {
            this.entityManager.CreateEntity(payload.id, payload.type, payload.data);
        });

        // Подписываемся на UpdateEntity
        this.wsClient.on("UpdateEntity", (payload: EntityInfo) => {
            this.entityManager.UpdateEntity(payload.id, payload.type, payload.data);
        });

        // Подписываемся на DeleteEntity
        this.wsClient.on("DeleteEntity", (payload: DeleteEntityPayload) => {
            this.entityManager.DeleteEntity(payload.id);
        });

        // Подписываемся на snapshot
        this.wsClient.on("snapshot", (payload: SnapshotPayload) => {
            this.entityManager.ApplySnapshot(payload.entities);
        });
    }
}
import type { ClientMessage, ServerMessage, Vector3D} from "../Models/NetworkMessages";
import { EntityManager } from "./EntityManager";

export class NetworkManager 
{
    private socket: WebSocket | null = null;
    private entityManager: EntityManager;

    constructor(entityManager: EntityManager) 
    {
        this.entityManager = entityManager;
    }

    public Connect(): void 
    {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";

        const url = window.location.port === "5173"
                    ? "ws://localhost:8080/ws"
                    : `${protocol}://${window.location.host}/ws`;

        this.socket = new WebSocket(url);

        this.socket.onopen = () => {
            console.log("WebSocket connected:", url);
        };

        this.socket.onmessage = (event: MessageEvent) => {
            const message = JSON.parse(event.data) as ServerMessage;
            this.HandleMessage(message);
        };

        this.socket.onclose = () => {
            console.log("WebSocket closed");
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket error:", error);
        };
    }

    public CreateGameSession(): void 
    {
         this.Send(
            {action: "createGameSession"}
        )
    }
    public CreateRoom(): void 
    {
        this.Send(
            {action: "createRoom"}
        )
    }

    public SendMove(position: Vector3D, direction: Vector3D): void 
    {
        this.Send({
            action: "move",
            position,
            direction
        });
    }

    public SendInteract(targetId: string): void 
    {
        this.Send({
            type: "action",
            action: "interact",
            targetId
        });
    }

    private Send(message: ClientMessage): void 
    {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) 
        {
            return;
        }

        this.socket.send(JSON.stringify(message));
    }

    private HandleMessage(message: ServerMessage): void 
    {
        switch (message.type) 
        {
            case "entity_create":
                this.entityManager.CreateEntity(
                    message.id,
                    message.kind,
                    message.position,
                    message.rotationY ?? 0
                );
                break;

            case "entity_update":
                this.entityManager.UpdateEntity(
                    message.id,
                    message.position,
                    message.rotationY
                );
                break;

            case "entity_delete":
                this.entityManager.DeleteEntity(message.id);
                break;

            case "snapshot":
                this.entityManager.ApplySnapshot(message.entities);
                break;
        }
    }
}
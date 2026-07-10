import type { ClientMessage, ServerMessage, Vector3D } from "../Models/NetworkMessages";

import { EntityManager } from "./EntityManager";

type NavigationState = 
{
    ship: {
        x: number;
        z: number;
        rotationY: number;
        hp: number;
    };

    asteroids: Array<{
        id: string;
        x: number;
        z: number;
    }>;

    monsters: Array<{
        id: string;
        x: number;
        z: number;
    }>;
};

export class NetworkManager 
{
    private socket: WebSocket | null = null;
    private entityManager: EntityManager;

    private computerStateHandler:
        | ((
              computerId: string,
              lockedBy: string | null
          ) => void)
        | null = null;

    private navigationStateHandler:
        | ((state: NavigationState) => void)
        | null = null;

    public constructor(entityManager: EntityManager) 
    {
        this.entityManager = entityManager;
    }

    public Connect(): void 
    {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";

        const url = window.location.port === "5173" ? "ws://localhost:8080/ws" : `${protocol}://${window.location.host}/ws`;

        this.socket = new WebSocket(url);

        this.socket.onopen = (): void => { console.log("WebSocket connected:", url); };

        this.socket.onmessage = (
            event: MessageEvent
        ): void => {
            try {
                const message = JSON.parse(
                    String(event.data)
                ) as ServerMessage;

                this.HandleMessage(message);
            } catch (error) {
                console.error(
                    "Некорректное сообщение от сервера:",
                    event.data,
                    error
                );
            }
        };

        this.socket.onclose = (): void => {
            console.log("WebSocket closed");
        };

        this.socket.onerror = (
            error: Event
        ): void => {
            console.error("WebSocket error:", error);
        };
    }

    public SendMove(position: Vector3D, rotationY: number): void 
    {
        this.Send({
            type: "action",
            action: "move",
            position,
            rotationY
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

    public SendUseComputer(computerId: string): void 
    {
        this.Send({
            type: "use_computer",
            computerId
        });
    }

    public SendReleaseComputer(computerId: string): void 
    {
        this.Send({
            type: "release_computer",
            computerId
        });
    }

    public SendShipInput(computerId: string, inputX: number, inputZ: number): void 
    {
        this.Send({
            type: "ship_input",
            computerId,
            inputX,
            inputZ
        });
    }

    public OnComputerState(
        handler: (
            computerId: string,
            lockedBy: string | null
        ) => void
    ): void {
        this.computerStateHandler = handler;
    }

    public OnNavigationState(handler: (state: NavigationState) => void): void 
    {
        this.navigationStateHandler = handler;
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

            case "computer_state":
                this.computerStateHandler?.(
                    message.computerId,
                    message.lockedBy
                );
                break;

            case "navigation_state":
                this.navigationStateHandler?.({
                    ship: message.ship,
                    asteroids: message.asteroids,
                    monsters: message.monsters
                });
                break;
        }
    }
}
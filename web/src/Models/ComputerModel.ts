import * as THREE from "three";

export class ComputerModel {
    public id: string;
    public interactionPosition: THREE.Vector3;

    private controllingPlayerId: string | null = null;

    public constructor(id: string, interactionPosition: THREE.Vector3) {
        this.id = id;
        this.interactionPosition = interactionPosition.clone();
    }

    public SetControllingPlayer(playerId: string | null): void {
        this.controllingPlayerId = playerId;
    }

    public IsAvailable(): boolean {
        return this.controllingPlayerId === null;
    }

    public IsControlledBy(playerId: string): boolean {
        return this.controllingPlayerId === playerId;
    }

    public GetControllingPlayerId(): string | null {
        return this.controllingPlayerId;
    }
}
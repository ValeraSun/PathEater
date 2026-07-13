
import * as THREE from "three"
export interface PlayerData {
    position: THREE.Vector3;
}

export interface MonsterData {
    position: THREE.Vector3;
}

export class EntityParser {
    public static Parse(type: string, data: unknown): any {
        switch (type) {
            case "player":
                return data as PlayerData; 
            case "monster":
                return data as MonsterData;
            default:
                throw new Error(`Неизвестный тип сущности: ${type}`);
        }
    }
}
export interface Vector3D {
    x: number;
    y: number;
    z: number;
}

export interface PlayerStatePayload {
    move_front: boolean;
    move_left: boolean;
    move_right: boolean;
    move_back: boolean;
    interact: boolean;
    attack: boolean;
    direction: Vector3D;
}

export interface EntityTransformData {
    position?: Vector3D;
    rotation?: Vector3D;
    health?: number;
    dead?: boolean;
    attacking?: boolean;
}

export interface AsteroidStateData {
    position?: Vector3D;
    radius?: number;
    destroyed?: boolean;
}

export interface MonsterStateData {
    position?: Vector3D;
    rotation?: Vector3D;
    health?: number;
    died?: boolean;
}

export interface EntityCreateInfo {
    id: string;
    type: string;
    data: unknown;
}

export interface EntityUpdateInfo {
    id: string;
    type: string;
    data: unknown;
}

export interface DeleteEntityPayload {
    id: string;
}

export interface RoomInfoPayload {
    roomId: string;
    playerId: string;
    error?: string;
}

export interface GameStartedPayload {
    playerId: string;
    spawn: Vector3D;
}

export interface GameOverPayload {
    win: boolean;
    status: number;
}

export interface MatchTimerPayload {
    time: number;
}

export interface RoomPlayersPayload {
    playerIds: string[];
}

export interface ShipWireData {
    baggage_status?: number;
    weapon_direction?: Vector3D;
    health?: number;
}
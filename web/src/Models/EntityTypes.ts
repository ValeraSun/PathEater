export type EntityType = "player" | "wall" | "room";

export type Vector3Data = {
    x: number;
    y: number;
    z: number;
};

export type PlayerData = {
    position: Vector3Data;
    rotationY: number;
    hp: number;
    name: string;
};

export type WallData = {
    position: Vector3Data;
    size: Vector3Data;
    rotationY: number;
};

export type RoomData = {
    name: string;
    position: Vector3Data;
    size: Vector3Data;
};

export type EntityDataMap = {
    player: PlayerData;
    wall: WallData;
    room: RoomData;
};
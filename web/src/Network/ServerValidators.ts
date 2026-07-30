import type {
    AsteroidStateData,
    DeleteEntityPayload,
    EntityCreateInfo,
    EntityTransformData,
    BulletStateData,
    GameOverPayload,
    GameStartedPayload,
    MatchTimerPayload,
    MonsterStateData,
    RoomInfoPayload,
    RoomPlayersPayload,
    ShipWireData,
    Vector3D,
    HoleStateData,
    DoorStateData
} from "./ServerContracts";


export interface BreakdownStateData {
    position?: Vector3D;
    radius?: number;
    normal?: Vector3D;
}

function IsObject(value: unknown): value is Record<string, unknown> {
    return value !== null && typeof value === "object";
}

function IsFiniteNumber(value: unknown): value is number {
    return typeof value === "number" && Number.isFinite(value);
}

function IsBoolean(value: unknown): value is boolean {
    return typeof value === "boolean";
}

export function IsVector3D(value: unknown): value is Vector3D {
    if (!IsObject(value)) return false;
    return IsFiniteNumber(value.x) && IsFiniteNumber(value.y) && IsFiniteNumber(value.z);
}

type FieldCheck = (value: unknown) => boolean;

function HasValidOptionalFields(value: unknown, schema: Record<string, FieldCheck>): boolean {
    if (!IsObject(value)) return false;

    for (const field in schema) {
        const fieldValue = value[field];
        if (fieldValue !== undefined && !schema[field](fieldValue)) {
            return false;
        }
    }

    return true;
}

export function IsEntityTransformData(value: unknown): value is EntityTransformData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
        rotation: IsVector3D,
        health: IsFiniteNumber,
        dead: IsBoolean,
        attacking: IsBoolean
    });
}

export function IsAsteroidStateData(value: unknown): value is AsteroidStateData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
        radius: IsFiniteNumber,
        destroyed: IsBoolean
    });
}

export function IsBulletData(value: unknown): value is BulletStateData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
    });
}

export function IsMonsterStateData(value: unknown): value is MonsterStateData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
        rotation: IsVector3D,
        health: IsFiniteNumber,
        died: IsBoolean
    });
}

export function IsShipWireData(value: unknown): value is ShipWireData {
    return HasValidOptionalFields(value, {
        baggage_status: IsFiniteNumber,
        weapon_direction: IsVector3D,
        health: IsFiniteNumber
    });
}

export function IsDoorStateData(value: unknown): value is DoorStateData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
        rotation: IsVector3D,
        isOpen: IsBoolean
    });
}

export function IsBreakdownStateData(value: unknown): value is BreakdownStateData {
    return HasValidOptionalFields(value, {
        position: IsVector3D,
        radius: IsFiniteNumber,
        normal: IsVector3D
    });
}

export function IsEntityInfo(value: unknown): value is EntityCreateInfo {
    if (!IsObject(value)) return false;
    return typeof value.id === "string" && typeof value.type === "string" && "data" in value;
}

export function IsDeleteEntityPayload(value: unknown): value is DeleteEntityPayload {
    return IsObject(value) && typeof value.id === "string";
}

export function IsRoomInfoPayload(value: unknown): value is RoomInfoPayload {
    if (!IsObject(value)) return false;
    if (typeof value.roomId !== "string" || typeof value.playerId !== "string") return false;
    return value.error === undefined || typeof value.error === "string";
}

export function IsRoomPlayersPayload(value: unknown): value is RoomPlayersPayload {
    if (!IsObject(value) || !Array.isArray(value.playerIds)) return false;
    return value.playerIds.every((id: unknown) => typeof id === "string");
}

export function IsGameStartedPayload(value: unknown): value is GameStartedPayload {
    if (!IsObject(value)) return false;
    return typeof value.playerId === "string" && IsVector3D(value.spawn);
}

export function IsGameOverPayload(value: unknown): value is GameOverPayload {
    if (!IsObject(value)) return false;
    return typeof value.win === "boolean" && IsFiniteNumber(value.status);
}

export function IsMatchTimerPayload(value: unknown): value is MatchTimerPayload {
    return IsObject(value) && IsFiniteNumber(value.time);
}
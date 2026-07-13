import type { PlayerData, RoomData, Vector3Data, WallData } from "./EntityTypes";

function IsObject(value: unknown): value is Record<string, unknown> 
{
    return typeof value === "object" && value !== null;
}

function ParseVector3(value: unknown): Vector3Data | null 
{
    if (!IsObject(value)) 
    {
        return null;
    }

    const x = value["x"];
    const y = value["y"];
    const z = value["z"];

    if (typeof x !== "number" || typeof y !== "number" || typeof z !== "number") 
    {
        return null;
    }

    return { x, y, z };
}

export function ParsePlayerData(value: unknown): PlayerData | null 
{
    if (!IsObject(value)) 
    {
        return null;
    }

    const position = ParseVector3(value["position"]);

    const rotationY = value["rotationY"];
    const hp = value["hp"];
    const name = value["name"];

    if (!position || typeof rotationY !== "number" || typeof hp !== "number" || typeof name !== "string")
    {
        return null;
    }

    return {
        position,
        rotationY,
        hp,
        name
    };
}

export function ParseWallData(value: unknown): WallData | null 
{
    if (!IsObject(value)) 
    {
        return null;
    }

    const position = ParseVector3(value["position"]);
    const size = ParseVector3(value["size"]);
    const rotationY = value["rotationY"];

    if (!position || !size || typeof rotationY !== "number") 
    {
        return null;
    }

    return {
        position,
        size,
        rotationY
    };
}

export function ParseRoomData(value: unknown): RoomData | null 
{
    if (!IsObject(value)) 
    {
        return null;
    }

    const name = value["name"];
    const position = ParseVector3(value["position"]);
    const size = ParseVector3(value["size"]);

    if (typeof name !== "string" || !position || !size) 
    {
        return null;
    }

    return {
        name,
        position,
        size
    };
}
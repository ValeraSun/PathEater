export interface PlayerData
{
    position: {
        x: number;
        y: number;
        z: number;
    };

    rotation: {
        x: number;
        y: number;
        z: number;
    };

    health: number;
}

export interface AlienData {
    position: {
        x: number;
        y: number;
        z: number;
    };

    rotation: {
        x: number;
        y: number;
        z: number;
    };
    health: number;
    died: boolean;
    attacking: boolean;
}

export interface ShipData
{
    baggageStatus: number;
    health: number;
}

export class EntityParser
{
    public static Parse(type: string, data: unknown): unknown
    {
        switch (type)
        {
            case "player":
                return data as PlayerData;
            case "alien":
                return data as AlienData;
            case "ship":
                return data as ShipData;
            default:
                return data;
        }
    }
}
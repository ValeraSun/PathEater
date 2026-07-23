import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";
import { ShipView, type ShipStateData } from "../Views/ShipView";
import { type NavigationDisplayState } from "../Views/ComputerView";

export interface EntityTransformData {
    position?: {
        x: number;
        y: number;
        z: number;
    };

    rotation?: {
        x: number;
        y: number;
        z: number;
    };

    health: number;
}

interface Vec2 {
    x: number;
    z: number;
}

interface ShipRadarData {
    weapon_direction?: { x: number; y: number; z: number };
    health?: number;
}

interface AsteroidData {
    position?: { x: number; y: number; z: number };
    radius?: number;
    destroyed?: boolean;
}

interface CosmoAlienData {
    position?: { x: number; y: number; z: number };
    rotation?: { x: number; y: number; z: number };
    health?: number;
    died?: boolean;
}

interface EntityRecord {
    id: string;
    type: string;
    object: THREE.Object3D;
    targetPosition: THREE.Vector3;
    targetRotationY: number;
}

const NETWORK_LERP_SPEED = 12;

export class EntityManager
{
    private readonly scene: THREE.Scene;
    private readonly entities = new Map<string, EntityRecord>();

    private localPlayerId: string | null = null;
    private shipView: ShipView | null = null;
    private onLocalPlayerUpdate: ((data: EntityTransformData) => void) | null = null;
    private shipHp = 100;
    private shipWeaponDirection: Vec2 = { x: 0, z: -1 };
    private asteroids = new Map<string, Vec2>();
    private monsters = new Map<string, Vec2>();
    private onRadarChanged: ((state: NavigationDisplayState) => void) | null = null;

    public SetRadarChangedHandler(handler: (state: NavigationDisplayState) => void): void
    {
        this.onRadarChanged = handler;
        this.EmitRadarChanged();
    }

    private EmitRadarChanged(): void
    {
        if (!this.onRadarChanged) return;

        const rotationY = Math.atan2(
            this.shipWeaponDirection.x,
            this.shipWeaponDirection.z
        );

        this.onRadarChanged({
            ship: { x: 0, z: 0, rotationY, hp: this.shipHp },
            asteroids: Array.from(this.asteroids, ([id, pos]) => ({ id, x: pos.x, z: pos.z })),
            monsters: Array.from(this.monsters, ([id, pos]) => ({ id, x: pos.x, z: pos.z }))
        });
    }

    public constructor(scene: THREE.Scene)
    {
        this.scene = scene;
    }

    public SetLocalPlayerUpdateHandler(handler: (data: EntityTransformData) => void): void
    {
        this.onLocalPlayerUpdate = handler;
    }

    public SetShipView(shipView: ShipView): void
    {
        this.shipView = shipView;
    }

    public SetLocalPlayerId(id: string): void
    {
        this.localPlayerId = id;
        this.DeleteEntity(id);
    }

    public CreateEntity(id: string, type: string, data: unknown): void
    {
        if (type === "ship")
        {
            this.UpdateShip(data);
            return;
        }

        if (type === "ship")
        {
            this.UpdateShip(data);
            return;
        }

        if (type === "asteroid")
        {
            this.UpdateAsteroid(id, data);
            return;
        }

        if (type === "alien")
        {
            this.UpdateMonster(id, data);
            return;
        }

        if (id === this.localPlayerId)
        {
            this.onLocalPlayerUpdate?.(data as EntityTransformData);

            return;
        }

        const existing = this.entities.get(id);

        if (existing)
        {
            this.UpdateEntity(id, type, data);

            return;
        }

        const object = this.CreateObject(type);

        if (!object)
        {
            console.warn(`Невозможно создать сущность неизвестного типа: ${type}`);

            return;
        }

        const transformData = data as EntityTransformData;

        this.ApplyTransform(object, transformData);

        this.scene.add(object);
        this.entities.set(id, {
            id,
            type,
            object,
            targetPosition: object.position.clone(),
            targetRotationY: object.rotation.y
        });
    }

    public UpdateEntity(id: string, type: string, data: unknown): void
    {
         if (type === "ship")
        {
            this.UpdateShip(data);
            return;
        }

        if (type === "asteroid")
        {
            this.UpdateAsteroid(id, data);
            return;
        }

        if (type === "alien")
        {
            this.UpdateMonster(id, data);
            return;
        }

        if (id === this.localPlayerId)
        {
            this.onLocalPlayerUpdate?.(data as EntityTransformData);

            return;
        }

        const entity = this.entities.get(id);

        if (!entity)
        {
            this.CreateEntity(id, type, data);

            return;
        }

        this.SetTarget(entity, data as EntityTransformData);
    }

    public DeleteEntity(id: string): void
    {
        const entity = this.entities.get(id);

        if (!entity)
        {
            return;
        }

        if (this.asteroids.delete(id))
        {
            this.EmitRadarChanged();
            return;
        }

        if (this.monsters.delete(id))
        {
            this.EmitRadarChanged();
            return;
        }

        this.scene.remove(entity.object);
        this.DisposeObject(entity.object);
        this.entities.delete(id);
    }

    public Clear(): void
    {
        for (const id of this.entities.keys())
        {
            this.DeleteEntity(id);
        }
    }

    public Update(dt: number): void
    {
        const interpolation =
            1 -
            Math.exp(
                -NETWORK_LERP_SPEED * dt
            );

        for (
            const entity of
            this.entities.values()
        )
        {
            entity.object.position.lerp(
                entity.targetPosition,
                interpolation
            );

            entity.object.rotation.y = this.LerpAngle(
                entity.object.rotation.y,
                entity.targetRotationY,
                interpolation
            );
        }
    }

    public GetEntity(
        id: string
    ): THREE.Object3D | null
    {
        return (
            this.entities.get(id)?.object ??
            null
        );
    }

    private UpdateShip(data: unknown): void
    {
        if (!this.shipView)
        {
            console.warn("ShipView не установлен в EntityManager");
            return;
        }

        if (!this.IsShipStateData(data))
        {
            console.warn("Получено некорректное состояние корабля", data);
            return;
        }

        this.shipView.UpdateState(data);

        const shipData = data as ShipRadarData;

        if (typeof shipData.health === "number")
        {
            this.shipHp = shipData.health;
        }

        if (shipData.weapon_direction)
        {
            this.shipWeaponDirection = {
                x: shipData.weapon_direction.x,
                z: shipData.weapon_direction.z
            };
        }

        this.EmitRadarChanged();
    }

    private IsShipStateData(data: unknown): data is ShipStateData
    {
        if (!data || typeof data !== "object")
        {
            return false;
        }

        const state = data as Record<string, unknown>;

        const baggageIsValid =
            state.baggage_status === undefined ||
            typeof state.baggage_status === "number";

        const healthIsValid =
            state.health === undefined ||
            typeof state.health === "number";

        return baggageIsValid && healthIsValid;
    }

    private UpdateAsteroid(id: string, data: unknown): void
    {
        if (!this.IsAsteroidData(data))
        {
            console.warn("Получены некорректные данные астероида", data);
            return;
        }

        if (data.destroyed)
        {
            this.asteroids.delete(id);
            this.EmitRadarChanged();
            return;
        }

        if (data.position)
        {
            this.asteroids.set(id, { x: data.position.x, z: data.position.z });
            this.EmitRadarChanged();
        }
    }

    private UpdateMonster(id: string, data: unknown): void
    {
        if (!this.IsCosmoAlienData(data))
        {
            console.warn("Получены некорректные данные пришельца", data);
            return;
        }

        if (data.died)
        {
            this.monsters.delete(id);
            this.EmitRadarChanged();
            return;
        }

        if (data.position)
        {
            this.monsters.set(id, { x: data.position.x, z: data.position.z });
            this.EmitRadarChanged();
        }
    }

    private IsAsteroidData(data: unknown): data is AsteroidData
    {
        if (!data || typeof data !== "object") return false;
        const d = data as Record<string, unknown>;
        return (
            (d.position === undefined || typeof d.position === "object") &&
            (d.destroyed === undefined || typeof d.destroyed === "boolean")
        );
    }

    private IsCosmoAlienData(data: unknown): data is CosmoAlienData
    {
        if (!data || typeof data !== "object") return false;
        const d = data as Record<string, unknown>;
        return (
            (d.position === undefined || typeof d.position === "object") &&
            (d.died === undefined || typeof d.died === "boolean")
        );
    }

    private SetTarget(entity: EntityRecord, data: EntityTransformData): void
    {
        if (data.position)
        {
            entity.targetPosition.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }

        if (data.rotation)
        {
            entity.targetRotationY = this.DirectionToYaw(data.rotation);
        }
    }

    private CreateObject(
        type: string
    ): THREE.Object3D | null
    {
        switch (type)
        {
            case "player":
            {
                const playerView =
                    new PlayerView();

                return playerView.mesh;
            }

            case "monster":
                return this.CreateBox(
                    0xff0000
                );

            case "door":
                return this.CreateBox(
                    0x4444ff
                );

            case "cargo":
                return this.CreateBox(
                    0xffaa00
                );

            default:
                return null;
        }
    }

    private ApplyTransform(object: THREE.Object3D, data: EntityTransformData): void
    {
        if (data.position)
        {
            object.position.set(data.position.x, data.position.y, data.position.z);
        }

        if (data.rotation)
        {
            object.rotation.y = this.DirectionToYaw(data.rotation);
        }
    }

    private DirectionToYaw(dir: { x: number; y: number; z: number }): number
    {
        return Math.atan2(dir.x, dir.z);
    }

    private LerpAngle(current: number, target: number, t: number): number {
        const delta = Math.atan2(Math.sin(target - current), Math.cos(target - current));
        return current + delta * t;
    }

    private DisposeObject(
        object: THREE.Object3D
    ): void
    {
        object.traverse(child =>
        {
            if (!(child instanceof THREE.Mesh))
            {
                return;
            }

            child.geometry.dispose();

            if (Array.isArray(child.material))
            {
                for (
                    const material of
                    child.material
                )
                {
                    material.dispose();
                }
            }
            else
            {
                child.material.dispose();
            }
        });
    }

    private CreateBox(
        color: number
    ): THREE.Mesh
    {
        return new THREE.Mesh(
            new THREE.BoxGeometry(
                1,
                1,
                1
            ),
            new THREE.MeshStandardMaterial({
                color
            })
        );
    }
}
import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";
import { type NavigationDisplayState } from "../Views/ComputerView";
import { ShipView, type ShipStateData } from "../Views/ShipView";
import { AlienView } from "../Views/AlienView";
import { CreateHole } from "../Views/BreakdownView";

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

export interface AlienStateData extends EntityTransformData {
    dead?: boolean;
    attacking?: boolean;
}

interface Vec2 {
    x: number;
    y: number;
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

interface AnimatedEntityView {
    mesh: THREE.Object3D;
    AdvanceAnimation(dt: number): void;
    SetMoving?(isMoving: boolean): void;
}

interface EntityRecord {
    id: string;
    type: string;
    object: THREE.Object3D;
    animatedView: AnimatedEntityView | null;
    targetPosition: THREE.Vector3;
    targetRotationY: number;
}

interface CreatedEntityObject {
    object: THREE.Object3D;
    animatedView: AnimatedEntityView | null;
}

const NETWORK_LERP_SPEED = 12;

export class EntityManager
{
    private readonly scene: THREE.Scene;
    private readonly entities = new Map<string, EntityRecord>();

    private localPlayerId: string | null = null;
    private shipView: ShipView | null = null;
    private onLocalPlayerUpdate: ((data: EntityTransformData) => void) | null = null;
    private onPlayersChanged: ((playerIds: string[]) => void) | null = null;
    private onRadarChanged: ((state: NavigationDisplayState) => void) | null = null;

    private shipHp = 100;
    private shipWeaponDirection: Vec2 = { x: 0, y: -1 };
    private asteroids = new Map<string, Vec2>();
    private monsters = new Map<string, Vec2>();

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
        this.EmitPlayersChanged();
    }

    public SetPlayersChangedHandler(handler: (playerIds: string[]) => void): void
    {
        this.onPlayersChanged = handler;
        this.EmitPlayersChanged();
    }

    public SetRadarChangedHandler(handler: (state: NavigationDisplayState) => void): void
    {
        this.onRadarChanged = handler;
        this.EmitRadarChanged();
    }

    public CreateEntity(id: string, type: string, data: unknown): void
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

        if (type === "cosmoAlien")
        {
            this.UpdateMonster(id, data);
            return;
        }

        if (type === "breakdown")
        {
            console.log("CreateEntity breakdown:", id, data);
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

        const created = this.CreateObject(type);

        if (!created)
        {
            console.warn(`Невозможно создать сущность неизвестного типа: ${type}`);
            return;
        }

        const transformData = data as EntityTransformData;

        this.ApplyTransform(created.object, transformData);
        this.ApplyEntityState(type, created.animatedView, data);

        this.scene.add(created.object);

        this.entities.set(id, {
            id,
            type,
            object: created.object,
            animatedView: created.animatedView,
            targetPosition: created.object.position.clone(),
            targetRotationY: created.object.rotation.y
        });

        if (type === "player")
        {
            this.EmitPlayersChanged();
        }
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

        if (type === "cosmoAlien")
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
        this.ApplyEntityState(entity.type, entity.animatedView, data);
    }

    public DeleteEntity(id: string): void
    {
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

        const entity = this.entities.get(id);

        if (!entity)
        {
            return;
        }

        const wasPlayer = entity.type === "player";

        this.scene.remove(entity.object);
        this.DisposeObject(entity.object);
        this.entities.delete(id);

        if (wasPlayer)
        {
            this.EmitPlayersChanged();
        }
    }

    public Clear(): void
    {
        for (const id of this.entities.keys())
        {
            this.DeleteEntity(id);
        }

        this.EmitPlayersChanged();
    }

    public Update(dt: number): void
{
    const interpolation = 1 - Math.exp(-NETWORK_LERP_SPEED * dt);

    for (const entity of this.entities.values())
    {
        if (entity.type === "player")
        {
            const distanceToTarget = entity.object.position.distanceTo(entity.targetPosition);
            const isMoving = distanceToTarget > 0.01;
            entity.animatedView?.SetMoving?.(isMoving);
        }

        entity.animatedView?.AdvanceAnimation(dt);

        entity.object.position.lerp(entity.targetPosition, interpolation);

        entity.object.rotation.y = this.LerpAngle(
            entity.object.rotation.y,
            entity.targetRotationY,
            interpolation
        );
    }
}

    public GetEntity(id: string): THREE.Object3D | null
    {
        return this.entities.get(id)?.object ?? null;
    }

    private EmitPlayersChanged(): void
    {
        const playerIds: string[] = [];

        if (this.localPlayerId)
        {
            playerIds.push(this.localPlayerId);
        }

        for (const entity of this.entities.values())
        {
            if (entity.type === "player")
            {
                playerIds.push(entity.id);
            }
        }

        this.onPlayersChanged?.(playerIds);
    }

    private EmitRadarChanged(): void
    {
        if (!this.onRadarChanged)
        {
            return;
        }

        const rotationY = Math.atan2(
            this.shipWeaponDirection.x,
            this.shipWeaponDirection.y
        );

        this.onRadarChanged({
            ship: { x: 0, y: 0, rotationY, hp: this.shipHp },
            asteroids: Array.from(this.asteroids, ([id, pos]) => ({ id, x: pos.x, y: pos.y })),
            monsters: Array.from(this.monsters, ([id, pos]) => ({ id, x: pos.x, y: pos.y }))
        });
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
                y: shipData.weapon_direction.y
            };
        }

        this.EmitRadarChanged();
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
           // console.dir(data)
            this.asteroids.set(id, { x: data.position.x, y: data.position.y });
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
            this.monsters.set(id, { x: data.position.x, y: data.position.z });
            this.EmitRadarChanged();
        }
    }

    private ApplyEntityState(type: string, animatedView: AnimatedEntityView | null, data: unknown): void
    {
        if (type !== "alien" && type !== "monster")
        {
            return;
        }

        if (!(animatedView instanceof AlienView))
        {
            return;
        }

        if (!this.IsAlienStateData(data))
        {
            console.warn("Получено некорректное состояние пришельца:", data);
            return;
        }

        const died = data.dead ?? false;
        const attacking = data.attacking ?? false;

        animatedView.mesh.visible = !died;

        if (died)
        {
            return;
        }

        animatedView.SetAttacking(attacking);
    }

    private IsShipStateData(data: unknown): data is ShipStateData
    {
        if (!data || typeof data !== "object")
        {
            return false;
        }

        const state = data as Record<string, unknown>;

        const baggageIsValid = state.baggageStatus === undefined || typeof state.baggageStatus === "number";
        const healthIsValid = state.health === undefined || typeof state.health === "number";

        return baggageIsValid && healthIsValid;
    }

    private IsAsteroidData(data: unknown): data is AsteroidData
    {
        if (!data || typeof data !== "object")
        {
            return false;
        }

        const d = data as Record<string, unknown>;

        return (
            (d.position === undefined || typeof d.position === "object") &&
            (d.destroyed === undefined || typeof d.destroyed === "boolean")
        );
    }

    private IsCosmoAlienData(data: unknown): data is CosmoAlienData
    {
        if (!data || typeof data !== "object")
        {
            return false;
        }

        const d = data as Record<string, unknown>;

        return (
            (d.position === undefined || typeof d.position === "object") &&
            (d.died === undefined || typeof d.died === "boolean")
        );
    }

    private IsAlienStateData(data: unknown): data is AlienStateData
    {
        if (!data || typeof data !== "object")
        {
            return false;
        }

        const state = data as Record<string, unknown>;

        const attackingIsValid = state.attacking === undefined || typeof state.attacking === "boolean";
        const diedIsValid = state.dead === undefined || typeof state.dead === "boolean";
        const healthIsValid = state.health === undefined || typeof state.health === "number";

        return attackingIsValid && diedIsValid && healthIsValid;
    }

    private SetTarget(entity: EntityRecord, data: EntityTransformData): void
    {
        if (data.position)
        {
            entity.targetPosition.set(data.position.x, data.position.y, data.position.z);
        }

        if (data.rotation)
        {
            entity.targetRotationY = this.DirectionToYaw(data.rotation);
        }
    }

    private CreateObject(type: string): CreatedEntityObject | null
    {
        switch (type)
        {
            case "player":
            {
                const playerView = new PlayerView();
                return {
                    object: playerView.mesh,
                    animatedView: playerView
                };
            }

            case "alien":
            case "monster":
            {
                const alienView = new AlienView();
                return {
                    object: alienView.mesh,
                    animatedView: alienView
                };
            }

            case "door":
            {
                return {
                    object: this.CreateBox(0x4444ff),
                    animatedView: null
                };
            }

            case "cargo":
            {
                return {
                    object: this.CreateBox(0xffaa00),
                    animatedView: null
                };
            }
            case "breakdown":
            {
                return {
                    object: CreateHole(),
                    animatedView: null
                };
            }

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

    private LerpAngle(current: number, target: number, t: number): number
    {
        const delta = Math.atan2(Math.sin(target - current), Math.cos(target - current));
        return current + delta * t;
    }

    private DisposeObject(object: THREE.Object3D): void
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
                for (const material of child.material)
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

    private CreateBox(color: number): THREE.Mesh
    {
        return new THREE.Mesh(
            new THREE.BoxGeometry(1, 1, 1),
            new THREE.MeshStandardMaterial({ color })
        );
    }
}
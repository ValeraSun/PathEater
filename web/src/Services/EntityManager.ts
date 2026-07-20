import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";
import {
    ShipView,
    type ShipStateData
} from "../Views/ShipView";

export interface EntityTransformData {
    position?: {
        x: number;
        y: number;
        z: number;
    };

    rotationY?: number;
    health: number;
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

            entity.object.rotation.y =
                this.LerpAngle(
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
            console.warn(
                "ShipView не установлен в EntityManager"
            );

            return;
        }

        if (!this.IsShipStateData(data))
        {
            console.warn(
                "Получено некорректное состояние корабля",
                data
            );

            return;
        }

        this.shipView.UpdateState(data);
    }

    private IsShipStateData(
        data: unknown
    ): data is ShipStateData
    {
        if (
            !data ||
            typeof data !== "object"
        )
        {
            return false;
        }

        const state =
            data as Record<string, unknown>;

        const baggageIsValid =
            state.baggageStatus === undefined ||
            typeof state.baggageStatus ===
                "number";

        const healthIsValid =
            state.health === undefined ||
            typeof state.health === "number";

        return baggageIsValid && healthIsValid;
    }

    private SetTarget(
        entity: EntityRecord,
        data: EntityTransformData
    ): void
    {
        if (data.position)
        {
            entity.targetPosition.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }

        if (
            typeof data.rotationY ===
            "number"
        )
        {
            entity.targetRotationY =
                data.rotationY;
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

    private ApplyTransform(
        object: THREE.Object3D,
        data: EntityTransformData
    ): void
    {
        if (data.position)
        {
            object.position.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }

        if (
            typeof data.rotationY ===
            "number"
        )
        {
            object.rotation.y =
                data.rotationY;
        }
    }

    private LerpAngle(
        current: number,
        target: number,
        interpolation: number
    ): number
    {
        const delta = Math.atan2(
            Math.sin(target - current),
            Math.cos(target - current)
        );

        return (
            current +
            delta * interpolation
        );
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
import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";
import {ShipView, type ShipStateData } from "../Views/ShipView";
import { AlienView } from "../Views/AlienView";

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
    died?: boolean;
    attacking?: boolean;
}

interface AnimatedEntityView {
    mesh: THREE.Object3D;
    AdvanceAnimation(dt: number): void;
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

    public SetPlayersChangedHandler(handler: (playerIds: string[]) => void): void {
        this.onPlayersChanged = handler;
        this.EmitPlayersChanged();
    }

    private EmitPlayersChanged(): void {
        const playerIds: string[] = [];

        if (this.localPlayerId) {
            playerIds.push(this.localPlayerId);
        }

        for (const entity of this.entities.values()) {
            if (entity.type === "player") {
                playerIds.push(entity.id);
            }
        }

        this.onPlayersChanged?.(playerIds);
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
        this.EmitPlayersChanged();
    }

    public CreateEntity(
        id: string,
        type: string,
        data: unknown
    ): void {
        if (type === "ship") {
            this.UpdateShip(data);
            return;
        }

        if (id === this.localPlayerId) {
            this.onLocalPlayerUpdate?.(
                data as EntityTransformData
            );
            return;
        }

        const existing = this.entities.get(id);

        if (existing) {
            this.UpdateEntity(id, type, data);
            return;
        }

        const created = this.CreateObject(type);

        if (!created) {
            console.warn(
                `Невозможно создать сущность неизвестного типа: ${type}`
            );
            return;
        }

        const transformData =
            data as EntityTransformData;

        this.ApplyTransform(
            created.object,
            transformData
        );

        this.ApplyEntityState(
            type,
            created.animatedView,
            data
        );

        this.scene.add(created.object);

        this.entities.set(id, {
            id,
            type,
            object: created.object,
            animatedView: created.animatedView,
            targetPosition:
                created.object.position.clone(),
            targetRotationY:
                created.object.rotation.y
        });

        if (type === "player") {
            this.EmitPlayersChanged();
        }
    }

    public UpdateEntity(id: string, type: string, data: unknown): void {
        if (type === "ship") {
            this.UpdateShip(data);
            return;
        }

        if (id === this.localPlayerId) {
            this.onLocalPlayerUpdate?.(data as EntityTransformData);
            return;
        }

        const entity = this.entities.get(id);

        if (!entity) {
            this.CreateEntity(id, type, data);
            return;
        }

        this.SetTarget(entity, data as EntityTransformData);
        this.ApplyEntityState(entity.type, entity.animatedView, data);
    }

    public DeleteEntity(id: string): void
    {
        const entity = this.entities.get(id);

        if (!entity)
        {
            return;
        }

        const wasPlayer = entity.type === "player";

        this.scene.remove(entity.object);
        this.DisposeObject(entity.object);
        this.entities.delete(id);

        if (wasPlayer) {
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
        const interpolation = 1 - Math.exp( -NETWORK_LERP_SPEED * dt);

        for (const entity of this.entities.values())
        {
            entity.animatedView?.AdvanceAnimation(dt);

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

    public GetEntity( id: string): THREE.Object3D | null
    {
        return (
            this.entities.get(id)?.object ?? null
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
            console.warn("Получено некорректное состояние корабля",data);
            return;
        }

        this.shipView.UpdateState(data);
    }

    private ApplyEntityState(type: string, animatedView: AnimatedEntityView | null, data: unknown): void {
        if (type !== "alien" && type !== "monster") {
            return;
        }

        if (!(animatedView instanceof AlienView)) {
            return;
        }

        if (!this.IsAlienStateData(data)) {
            console.warn("Получено некорректное состояние пришельца:", data);
            return;
        }

        const died = data.died ?? false;
        const attacking = data.attacking ?? false;
        animatedView.mesh.visible = !died;

        if (died) {
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

    private IsAlienStateData(data: unknown): data is AlienStateData {
        if (!data || typeof data !== "object") {
            return false;
        }

        const state = data as Record<string, unknown>;

        const attackingIsValid = state.attacking === undefined || typeof state.attacking === "boolean";
        const diedIsValid = state.died === undefined || typeof state.died === "boolean";
        const healthIsValid = state.health === undefined || typeof state.health === "number";

        return (
            attackingIsValid && diedIsValid && healthIsValid
        );
    }

    private SetTarget(entity: EntityRecord, data: EntityTransformData): void {
        if (data.position) {
            entity.targetPosition.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }

        if (data.rotation) {
            entity.targetRotationY =this.DirectionToYaw(data.rotation);
        }
    }

    private CreateObject(type: string): CreatedEntityObject | null {
        switch (type) {
            case "player": {
                const playerView = new PlayerView();
                return {
                    object: playerView.mesh,
                    animatedView: playerView
                };
            }
            case "alien":
            case "monster": {
                const alienView = new AlienView();
                return {
                    object: alienView.mesh,
                    animatedView: alienView
                };
            }
            case "door": {
                const object = this.CreateBox(0x4444ff);
                return {
                    object,
                    animatedView: null
                };
            }
            case "cargo": {
                const object = this.CreateBox(0xffaa00);
                return {
                    object,
                    animatedView: null
                };
            }
            default:
                return null;
        }
    }

    private ApplyTransform(object: THREE.Object3D, data: EntityTransformData): void {
        if (data.position) {
            object.position.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }

        if (data.rotation) {
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

    private CreateBox( color: number): THREE.Mesh
    {
        return new THREE.Mesh(
            new THREE.BoxGeometry(1, 1, 1),
            new THREE.MeshStandardMaterial({ color})
        );
    }
}
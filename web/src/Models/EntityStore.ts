import * as THREE from "three";
import type { EntityTransformData, Vector3D } from "../Network/ServerContracts";
import { type EntityModel } from "./EntityModel";

export class EntityStore {
    private entitiesById = new Map<string, EntityModel>();

    public CreateEntity(entityId: string, entityType: string, transformData: EntityTransformData): EntityModel {
        const initialPosition = new THREE.Vector3();

        if (transformData.position) {
            initialPosition.set(transformData.position.x, transformData.position.y, transformData.position.z);
        }

        const initialRotationY = transformData.rotation ? this.ConvertDirectionToYaw(transformData.rotation): 0;

        const entity: EntityModel = {
            id: entityId,
            type: entityType,

            position: initialPosition,
            targetPosition: initialPosition.clone(),

            rotationY: initialRotationY,
            targetRotationY: initialRotationY,

            health: transformData.health,
            dead: transformData.dead,
            attacking: transformData.attacking
        };

        this.entitiesById.set(entityId, entity);
        return entity;
    }

    public UpdateEntity(entityId: string, transformData: EntityTransformData): EntityModel | null {
        const entity = this.entitiesById.get(entityId);

        if (!entity) {
            return null;
        }

        if (transformData.position) {
            entity.targetPosition.set(transformData.position.x, transformData.position.y, transformData.position.z);
        }

        if (transformData.rotation) {
            entity.targetRotationY = this.ConvertDirectionToYaw(transformData.rotation);
        }

        if (typeof transformData.health === "number") {
            entity.health = transformData.health;
        }

        if (typeof transformData.dead === "boolean") {
            entity.dead = transformData.dead;
        }

        if (typeof transformData.attacking === "boolean") {
            entity.attacking = transformData.attacking;
        }

        if (typeof transformData.isOpen === "boolean") {
            console.log("передано состояние")
            entity.isOpen = transformData.isOpen;
            console.log(transformData.isOpen);
        }
        return entity;
    }

    public GetEntity(entityId: string): EntityModel | null {
        return this.entitiesById.get(entityId) ?? null;
    }

    public GetAllEntities(): Iterable<EntityModel> {
        return this.entitiesById.values();
    }

    public RemoveEntity(entityId: string): boolean {
        return this.entitiesById.delete(entityId);
    }

    public Clear(): void {
        this.entitiesById.clear();
    }

    private ConvertDirectionToYaw(direction: Vector3D): number {
        return Math.atan2(direction.x, direction.z);
    }
}
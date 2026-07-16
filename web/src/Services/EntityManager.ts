import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";

export interface EntityTransformData {
    position?: {
        x: number;
        y: number;
        z: number;
    };

    rotationY?: number;
}

interface EntityRecord {
    id: string;
    type: string;
    object: THREE.Object3D;
    targetPosition: THREE.Vector3;
    targetRotationY: number;
}

interface SnapshotEntity {
    id: string;
    type: string;
    data: EntityTransformData;
}

const NETWORK_LERP_SPEED = 12;

export class EntityManager {
    private scene: THREE.Scene;
    private entities = new Map<string, EntityRecord>();
    private localPlayerId: string | null = null;
    private onLocalPlayerUpdate: ((data: EntityTransformData) => void) | null = null;
    public SetLocalPlayerUpdateHandler(handler: (data: EntityTransformData) => void): void {
        this.onLocalPlayerUpdate = handler;
    }

    public constructor(scene: THREE.Scene) {
        this.scene = scene;
    }

    public SetLocalPlayerId(id: string): void {
        this.localPlayerId = id;
        this.DeleteEntity(id);
    }

    public CreateEntity(id: string, type: string, data: EntityTransformData): void {
        if (id === this.localPlayerId) {
            this.onLocalPlayerUpdate?.(data);
            return;
        }
 
        const existing = this.entities.get(id);
 
        if (existing) {
            this.UpdateEntity(id, type, data);
            return;
        }
 
        const object = this.createObject(type);
 
        if (!object) {
            console.warn(`Невозможно создать сущность неизвестного типа: ${type}`);
            return;
        }
 
        this.applyTransform(object, data);
        this.scene.add(object);
 
        this.entities.set(id, {
            id,
            type,
            object,
            targetPosition: object.position.clone(),
            targetRotationY: object.rotation.y
        });
    }

    public UpdateEntity(id: string, type: string, data: EntityTransformData): void {
        if (id === this.localPlayerId) {
            this.onLocalPlayerUpdate?.(data);
            return;
        }
 
        const entity = this.entities.get(id);
 
        if (!entity) {
            this.CreateEntity(id, type, data);
            return;
        }
 
        if (entity.type !== type) {
            this.DeleteEntity(id);
            this.CreateEntity(id, type, data);
            return;
        }
 
        this.setTarget(entity, data);
    }

    public DeleteEntity(id: string): void {
        const entity = this.entities.get(id);
 
        if (!entity) {
            return;
        }
 
        this.scene.remove(entity.object);
        this.disposeObject(entity.object);
        this.entities.delete(id);
    }

    public ApplySnapshot(entities: SnapshotEntity[]): void {
        const receivedIds = new Set(
            entities.map(entity => entity.id)
        );

        for (const id of [...this.entities.keys()]) {
            if (!receivedIds.has(id)) {
                this.DeleteEntity(id);
            }
        }

        for (const entity of entities) {
            this.CreateEntity(
                entity.id,
                entity.type,
                entity.data
            );
        }
    }

    public Clear(): void 
    {
        for (const id of [...this.entities.keys()]) {
            this.DeleteEntity(id);
        }
    }

    public Update(dt: number): void {
        const t = 1 - Math.exp(-NETWORK_LERP_SPEED * dt);
 
        for (const entity of this.entities.values()) {
            entity.object.position.lerp(entity.targetPosition, t);
 
            entity.object.rotation.y = this.lerpAngle(
                entity.object.rotation.y,
                entity.targetRotationY,
                t
            );
        }
    }

    public GetEntity(id: string): THREE.Object3D | null {
        return this.entities.get(id)?.object ?? null;
    }

    private setTarget(entity: EntityRecord, data: EntityTransformData): void {
        if (data.position) {
            entity.targetPosition.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }
 
        if (typeof data.rotationY === "number") {
            entity.targetRotationY = data.rotationY;
        }
    }
 
    private createObject(type: string): THREE.Object3D | null {
        switch (type) {
            case "player": {
                const playerView = new PlayerView();
                return playerView.mesh;
            }
            case "monster":
                return this.CreateBox(0xff0000);
            case "door":
                return this.CreateBox(0x4444ff);
            case "cargo":
                return this.CreateBox(0xffaa00);
            default:
                return null;
        }
    }

    private applyTransform(
        object: THREE.Object3D,
        data: EntityTransformData
    ): void {
        if (data.position) {
            object.position.set(
                data.position.x,
                data.position.y,
                data.position.z
            );
        }
        if (typeof data.rotationY === "number") {
            object.rotation.y = data.rotationY;
        }
    }

    private lerpAngle(a: number, b: number, t: number): number {
        const delta = Math.atan2(Math.sin(b - a), Math.cos(b - a));
        return a + delta * t;
    }
 
    private disposeObject(object: THREE.Object3D): void {
        object.traverse(child => {
            if (!(child instanceof THREE.Mesh)) {
                return;
            }
            child.geometry.dispose();
            if (Array.isArray(child.material)) {
                for (const material of child.material) {
                    material.dispose();
                }
            } else {
                child.material.dispose();
            }
        });
    }

    private CreateBox(color: number): THREE.Mesh {
        return new THREE.Mesh(
            new THREE.BoxGeometry(1, 1, 1),
            new THREE.MeshStandardMaterial({ color })
        );
    }
}
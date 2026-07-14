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
}

interface SnapshotEntity {
    id: string;
    type: string;
    data: EntityTransformData;
}

export class EntityManager {
    private scene: THREE.Scene;
    private entities = new Map<string, EntityRecord>();

    public constructor(scene: THREE.Scene) {
        this.scene = scene;
    }

    public CreateEntity(id: string, type: string, data: EntityTransformData): void 
    {
        if (this.entities.has(id)) 
        {
            this.UpdateEntity(id, type, data);
            return;
        }

        const object = this.createObject(type);

        if (!object) 
        {
            console.warn(`Невозможно создать сущность неизвестного типа: ${type}`);
            return;
        }

        this.applyTransform(object, data);

        this.scene.add(object);

        this.entities.set(id, {
            id,
            type,
            object
        });
    }

    public UpdateEntity(id: string, type: string, data: EntityTransformData): void 
    {
        const entity = this.entities.get(id);

        if (!entity) 
        {
            this.CreateEntity(id, type, data);
            return;
        }

        this.applyTransform(entity.object, data);
    }

    public DeleteEntity(id: string): void 
    {
        const entity = this.entities.get(id);

        if (!entity) 
        {
            return;
        }

        this.scene.remove(entity.object);
        this.disposeObject(entity.object);
        this.entities.delete(id);
    }

    public ApplySnapshot(entities: SnapshotEntity[]): void 
    {
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

    private createObject(type: string): THREE.Object3D | null 
    {
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

    private disposeObject(object: THREE.Object3D): void 
    {
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
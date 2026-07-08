import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";
import type { Vector3D } from "../Models/NetworkMessages";

type EntityKind = "player" | "monster" | "door" | "cargo";

type EntityRecord = {
    id: string;
    kind: EntityKind;
    object: THREE.Object3D;
};

export class EntityManager 
{
    private scene: THREE.Scene;
    private entities = new Map<string, EntityRecord>();

    constructor(scene: THREE.Scene) 
    {
        this.scene = scene;
    }

    public CreateEntity(id: string, kind: EntityKind, position: Vector3D, rotationY = 0): void 
    {
        if (this.entities.has(id)) 
        {
            this.UpdateEntity(id, position, rotationY);
            return;
        }

        let object: THREE.Object3D;

        switch (kind) 
        {
            case "player": 
            {
                const playerView = new PlayerView();
                object = playerView.mesh;
                break;
            }

            case "monster": 
            {
                object = this.CreateBox(0xff0000);
                break;
            }

            case "door": 
            {
                object = this.CreateBox(0x4444ff);
                break;
            }

            case "cargo": 
            {
                object = this.CreateBox(0xffaa00);
                break;
            }

            default: return;
        }

        object.position.set(position.x, position.y, position.z);
        object.rotation.y = rotationY;

        this.scene.add(object);

        this.entities.set(id, {
            id,
            kind,
            object
        });
    }

    public UpdateEntity(id: string, position?: Vector3D, rotationY?: number): void 
    {
        const entity = this.entities.get(id);

        if (!entity) 
        {
            return;
        }

        if (position) 
        {
            entity.object.position.set(position.x, position.y, position.z);
        }

        if (rotationY !== undefined) 
        {
            entity.object.rotation.y = rotationY;
        }
    }

    public DeleteEntity(id: string): void 
    {
        const entity = this.entities.get(id);

        if (!entity) 
        {
            return;
        }

        this.scene.remove(entity.object);
        this.entities.delete(id);
    }

    public ApplySnapshot(
        entities: Array<{
            id: string;
            kind: EntityKind;
            position: Vector3D;
            rotationY?: number;
        }>
    ): void 
    {
        for (const entity of entities) 
        {
            this.CreateEntity(
                entity.id,
                entity.kind,
                entity.position,
                entity.rotationY ?? 0
            );
        }
    }

    private CreateBox(color: number): THREE.Mesh 
    {
        return new THREE.Mesh(
            new THREE.BoxGeometry(1, 1, 1),
            new THREE.MeshStandardMaterial({ color })
        );
    }
}
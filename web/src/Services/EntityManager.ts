import * as THREE from "three";
import { PlayerView } from "../Views/PlayerView";

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

    public CreateEntity(type:string, id: string, data: unknown): void 
    {
        if (this.entities.has(id)) 
        {
            this.UpdateEntity(type, id, data);
            return;
        }

        let object: THREE.Object3D;

        switch (type) 
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

        object.position.set(data.position.x, data.position.y, data.position.z);
        object.rotation.y = data.rotationY;

        this.scene.add(object);

        this.entities.set(id, {
            id,
            type,
            object
        });
    }

    public UpdateEntity(type: string, id: string, data: unknown): void 
    {
        const entity = this.entities.get(id);

        if (!entity) 
        {
            return;
        }

        if (data.position) 
        {
            entity.object.position.set(data.position.x, data.position.y, data.position.z);
        }

        if (data.rotationY !== undefined) 
        {
            entity.object.rotation.y = data.rotationY;
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
            type: string;
            data: unknown;
        }>
    ): void 
    {
        for (const entity of entities) 
        {
            this.CreateEntity(
                entity.id,
                entity.type,
                entity.data.position,
                entity.data.rotationY ?? 0
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
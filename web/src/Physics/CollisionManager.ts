import * as THREE from "three";
import { OBB, CollisionLayer, type CollisionLayerValue } from "./OBB";
import { UniformGrid } from "./UniformGrid";
import { RESOLVE_ITERATIONS, SKIN, MAX_STEP, GRID_CELL_SIZE, EPSILON, EPSILON_SQ } from "./PhysicsConst";

type JsonCollider = 
{
    id: string;
    type: string;
    center: { x: number; y: number; z: number };
    size?: { x: number; y: number; z: number };
    halfExtents?: { x: number; y: number; z: number };
    rotation?: { x: number; y: number; z: number; w: number };
    layer?: string;
};

export class CollisionManager 
{
    private staticColliders: OBB[] = [];
    private dynamicColliders: OBB[] = [];

    private grid = new UniformGrid(GRID_CELL_SIZE);

    private stepVec = new THREE.Vector3();
    private pushOut = new THREE.Vector3();
    private slideVec = new THREE.Vector3();
    private startPos = new THREE.Vector3();

    public async LoadShipColliders(path: string): Promise<void> 
    {
        const response = await fetch(path);
        const data = await response.json();

        for (const item of data.colliders as JsonCollider[]) 
        {
            const center = new THREE.Vector3(item.center.x, item.center.y, item.center.z);

            let halfExtents: THREE.Vector3;

            if (item.halfExtents) 
            {
                halfExtents = new THREE.Vector3(
                    item.halfExtents.x,
                    item.halfExtents.y,
                    item.halfExtents.z,
                );
            } 
            else if (item.size) 
            {
                halfExtents = new THREE.Vector3(
                    item.size.x / 2,
                    item.size.y / 2,
                    item.size.z / 2,
                );
            } 
            else 
            {
                continue;
            }

            const quaternion = item.rotation
                ? new THREE.Quaternion(item.rotation.x, item.rotation.y, item.rotation.z, item.rotation.w,)
                : new THREE.Quaternion();

            const obb = new OBB(item.id, center, halfExtents, quaternion, CollisionLayer.Static);

            this.staticColliders.push(obb);
            this.grid.Insert(obb);
        }
    }

    public AddDynamic(obb: OBB): void 
    {
        this.dynamicColliders.push(obb);
    }

    public RemoveDynamic(obb: OBB): void 
    {
        const i = this.dynamicColliders.indexOf(obb);
        if (i !== -1) this.dynamicColliders.splice(i, 1);
    }

    private Resolve(body: OBB, mask: CollisionLayerValue): THREE.Vector3 | null 
    {
        let lastNormal: THREE.Vector3 | null = null;

        for (let iteration = 0; iteration < RESOLVE_ITERATIONS; iteration++) 
        {
            let deepestDepth = 0;
            let deepestNormal: THREE.Vector3 | null = null;

            const candidates = this.grid.QueryBody(body);

            for (const collider of candidates) 
            {
                if ((collider.layer & mask) === 0) continue;

                const result = body.Intersects(collider);

                if (!result.hit) continue;

                if (result.depth > deepestDepth) 
                {
                    deepestDepth = result.depth;
                    deepestNormal = result.normal.clone();
                }
            }

            for (const collider of this.dynamicColliders) 
            {
                if (collider === body) continue;
                if ((collider.layer & mask) === 0) continue;

                const result = body.Intersects(collider);

                if (!result.hit) continue;

                if (result.depth > deepestDepth) 
                {
                    deepestDepth = result.depth;
                    deepestNormal = result.normal.clone();
                }
            }

            if (deepestNormal === null) break;

            this.pushOut.copy(deepestNormal).multiplyScalar(deepestDepth + SKIN);
            body.center.add(this.pushOut);
            body.UpdateAxes();

            lastNormal = deepestNormal;
        }

        return lastNormal;
    }

    public MoveAndSlide(body: OBB, displacement: THREE.Vector3, mask: CollisionLayerValue): THREE.Vector3 
    {
        this.startPos.copy(body.center);

        const remaining = displacement.clone();
        const totalLength = remaining.length();

        if (totalLength < EPSILON) 
        {
            this.Resolve(body, mask);
            return new THREE.Vector3();
        }

        const steps = Math.max(1, Math.ceil(totalLength / MAX_STEP));

        for (let i = 0; i < steps; i++) 
        {
            if (remaining.lengthSq() < EPSILON_SQ) break;

            this.stepVec.copy(remaining).multiplyScalar(1 / (steps - i));

            body.center.add(this.stepVec);
            body.UpdateAxes();

            remaining.sub(this.stepVec);

            const normal = this.Resolve(body, mask);

            if (normal !== null) 
            {
                const into = remaining.dot(normal);

                if (into < 0) 
                {
                    this.slideVec.copy(normal).multiplyScalar(into);
                    remaining.sub(this.slideVec);
                }
            }
        }

        return new THREE.Vector3().subVectors(body.center, this.startPos);
    }

    public IsFree(body: OBB, mask: CollisionLayerValue): boolean 
    {
        const candidates = this.grid.QueryBody(body);

        for (const collider of candidates) 
        {
            if ((collider.layer & mask) === 0) continue;
            if (body.Intersects(collider).hit) return false;
        }

        return true;
    }
    // УДАЛИТЬ ПЕРЕД ЗАЩИТОЙ
    public AddDebugHelpers(scene: THREE.Scene): void 
    {
        for (const collider of this.staticColliders) 
        {
            const geometry = new THREE.BoxGeometry(
                collider.halfExtents.x * 2,
                collider.halfExtents.y * 2,
                collider.halfExtents.z * 2,
            );

            const mesh = new THREE.Mesh(
                geometry,
                new THREE.MeshBasicMaterial({ color: 0xff0000, wireframe: true }),
            );

            mesh.position.copy(collider.center);
            mesh.quaternion.copy(collider.quaternion);

            scene.add(mesh);
        }
    }
}
import * as THREE from "three";
import { EPSILON_SQ } from "./PhysicsConst";

export const CollisionLayer = {
    None: 0,
    Static: 1 << 0,
    Player: 1 << 1,
    Item: 1 << 2,
    Monster: 1 << 3,
} as const;

export type CollisionLayerValue = typeof CollisionLayer[keyof typeof CollisionLayer];

export type CollisionResult = {
    hit: boolean;
    normal: THREE.Vector3;
    depth: number;
};

export class OBB 
{
    public id: string;
    public center: THREE.Vector3;
    public halfExtents: THREE.Vector3;
    public quaternion: THREE.Quaternion;
    public layer: CollisionLayerValue;
    public axes: THREE.Vector3[] = [];

    constructor(
        id: string,
        center: THREE.Vector3,
        halfExtents: THREE.Vector3,
        quaternion: THREE.Quaternion,
        layer: CollisionLayerValue
    ) 
    {
        this.id = id;
        this.center = center;
        this.halfExtents = halfExtents;
        this.quaternion = quaternion;
        this.layer = layer;

        this.UpdateAxes();
    }

    public static FromAABB(
        id: string,
        center: THREE.Vector3,
        size: THREE.Vector3,
        layer: CollisionLayerValue
    ): OBB 
    {
        return new OBB(
            id,
            center,
            size.clone().multiplyScalar(0.5),
            new THREE.Quaternion(),
            layer
        );
    }

    public UpdateAxes(): void 
    {
        this.axes = [
            new THREE.Vector3(1, 0, 0).applyQuaternion(this.quaternion).normalize(),
            new THREE.Vector3(0, 1, 0).applyQuaternion(this.quaternion).normalize(),
            new THREE.Vector3(0, 0, 1).applyQuaternion(this.quaternion).normalize(),
        ];
    }

    private Project(axis: THREE.Vector3): { min: number; max: number } 
    {
        const radius =
            Math.abs(this.halfExtents.x * axis.dot(this.axes[0])) +
            Math.abs(this.halfExtents.y * axis.dot(this.axes[1])) +
            Math.abs(this.halfExtents.z * axis.dot(this.axes[2]));

        const centerProjection = this.center.dot(axis);

        return {
            min: centerProjection - radius,
            max: centerProjection + radius,
        };
    }

    public Intersects(other: OBB): CollisionResult 
    {
        const axesToTest: THREE.Vector3[] = [];

        axesToTest.push(...this.axes);
        axesToTest.push(...other.axes);

        for (const a of this.axes) 
        {
            for (const b of other.axes) 
            {
                const cross = new THREE.Vector3().crossVectors(a, b);

                if (cross.lengthSq() > EPSILON_SQ) 
                {
                    axesToTest.push(cross.normalize());
                }
            }
        }

        let minOverlap = Infinity;
        const smallestAxis = new THREE.Vector3();

        for (const axis of axesToTest) 
        {
            const p1 = this.Project(axis);
            const p2 = other.Project(axis);

            const overlap = Math.min(p1.max, p2.max) - Math.max(p1.min, p2.min);

            if (overlap <= 0) 
            {
                return {
                    hit: false,
                    normal: new THREE.Vector3(),
                    depth: 0,
                };
            }

            if (overlap < minOverlap) 
            {
                minOverlap = overlap;
                smallestAxis.copy(axis);
            }
        }

        const direction = new THREE.Vector3().subVectors(this.center, other.center);

        if (direction.dot(smallestAxis) < 0) {
            smallestAxis.negate();
        }

        return {
            hit: true,
            normal: smallestAxis,
            depth: minOverlap,
        };
    }
}
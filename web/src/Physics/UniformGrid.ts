import * as THREE from "three";
import { OBB } from "./OBB";

export class UniformGrid 
{
    private cellSize: number;
    private cells = new Map<string, OBB[]>();

    private tmpAABB = new THREE.Box3();
    private result: OBB[] = [];
    private seen = new Set<OBB>();

    constructor(cellSize: number) 
    {
        this.cellSize = cellSize;
    }

    public Clear(): void 
    {
        this.cells.clear();
    }

    private CellIndex(value: number): number 
    {
        return Math.floor(value / this.cellSize);
    }

    private Key(x: number, y: number, z: number): string 
    {
        return `${x}_${y}_${z}`;
    }

    private ComputeAABB(obb: OBB, target: THREE.Box3): void 
    {
        const h = obb.halfExtents;
        const a = obb.axes;

        const ex =
            Math.abs(a[0].x) * h.x + Math.abs(a[1].x) * h.y + Math.abs(a[2].x) * h.z;

        const ey =
            Math.abs(a[0].y) * h.x + Math.abs(a[1].y) * h.y + Math.abs(a[2].y) * h.z;

        const ez =
            Math.abs(a[0].z) * h.x + Math.abs(a[1].z) * h.y + Math.abs(a[2].z) * h.z;

        target.min.set(obb.center.x - ex, obb.center.y - ey, obb.center.z - ez);
        target.max.set(obb.center.x + ex, obb.center.y + ey, obb.center.z + ez);
    }

    public Insert(obb: OBB): void 
    {
        this.ComputeAABB(obb, this.tmpAABB);

        const minX = this.CellIndex(this.tmpAABB.min.x);
        const minY = this.CellIndex(this.tmpAABB.min.y);
        const minZ = this.CellIndex(this.tmpAABB.min.z);

        const maxX = this.CellIndex(this.tmpAABB.max.x);
        const maxY = this.CellIndex(this.tmpAABB.max.y);
        const maxZ = this.CellIndex(this.tmpAABB.max.z);

        for (let x = minX; x <= maxX; x++) 
        {
            for (let y = minY; y <= maxY; y++) 
            {
                for (let z = minZ; z <= maxZ; z++) 
                {
                    const key = this.Key(x, y, z);

                    let list = this.cells.get(key);

                    if (list === undefined) 
                    {
                        list = [];
                        this.cells.set(key, list);
                    }

                    list.push(obb);
                }
            }
        }
    }

    public QueryBody(body: OBB): OBB[] 
    {
        this.ComputeAABB(body, this.tmpAABB);
        return this.QueryAABB(this.tmpAABB);
    }

    public QueryAABB(aabb: THREE.Box3): OBB[] 
    {
        this.result.length = 0;
        this.seen.clear();

        const minX = this.CellIndex(aabb.min.x);
        const minY = this.CellIndex(aabb.min.y);
        const minZ = this.CellIndex(aabb.min.z);

        const maxX = this.CellIndex(aabb.max.x);
        const maxY = this.CellIndex(aabb.max.y);
        const maxZ = this.CellIndex(aabb.max.z);

        for (let x = minX; x <= maxX; x++) 
        {
            for (let y = minY; y <= maxY; y++) 
            {
                for (let z = minZ; z <= maxZ; z++) 
                {
                    const list = this.cells.get(this.Key(x, y, z));

                    if (list === undefined) continue;

                    for (const obb of list) 
                    {
                        if (this.seen.has(obb)) continue;

                        this.seen.add(obb);
                        this.result.push(obb);
                    }
                }
            }
        }

        return this.result;
    }

    public get CellCount(): number 
    {
        return this.cells.size;
    }
}
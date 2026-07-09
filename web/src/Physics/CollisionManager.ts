import * as THREE from "three";

export class CollisionManager 
{
    private colliders: THREE.Box3[] = [];

    public async LoadShipColliders(): Promise<void> 
    {
        const response = await fetch("/data/ship_wall_colliders.json");
        const data = await response.json();

        for (const collider of data.colliders) 
        {
            const center = new THREE.Vector3(
                collider.center.x,
                collider.center.y,
                collider.center.z
            );

            const size = new THREE.Vector3(
                collider.size.x,
                collider.size.y,
                collider.size.z
            );

            const box = new THREE.Box3().setFromCenterAndSize(center, size);

            this.colliders.push(box);
        }

        console.log("Ship colliders loaded:", this.colliders.length);
    }

    public CanMove(playerBox: THREE.Box3): boolean 
    {
        for (const collider of this.colliders) 
        {
            if (playerBox.intersectsBox(collider)) 
            {
                return false;
            }
        }

        return true;
    }
    
    //убрать перед защитой
    public AddDebugHelpers(scene: THREE.Scene): void 
    {
        for (const collider of this.colliders) {
            const helper = new THREE.Box3Helper(collider, 0xff0000);
            scene.add(helper);
        }
    }
}
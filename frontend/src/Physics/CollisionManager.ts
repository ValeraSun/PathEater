import * as THREE from "three";

export class CollisionManager 
{
    private colliders: THREE.Box3[] = [];

    public AddCollider(box: THREE.Box3): void 
    {
        this.colliders.push(box);
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
}
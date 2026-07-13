import * as THREE from "three";
import { OBB, CollisionLayer } from "../Physics/OBB";
import { LOCAL_PLAYER_ID, PLAYER_COLLIDER_SIZE, PLAYER_SPAWN_POSITION } from "../Config/PlayerConfig";

export class PlayerModel 
{
    public body = OBB.FromAABB(
        LOCAL_PLAYER_ID,
        PLAYER_SPAWN_POSITION,
        PLAYER_COLLIDER_SIZE,
        CollisionLayer.Player
    );

    public get position(): THREE.Vector3 
    {
        return this.body.center;
    }

    speed = 3;
}
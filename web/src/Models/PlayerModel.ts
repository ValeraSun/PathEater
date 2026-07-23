import * as THREE from "three";
import { PLAYER_SPAWN_POSITION } from "../Config/PlayerConfig";

export class PlayerModel {
    public readonly position: THREE.Vector3;
    public speed = 3;

    public constructor(spawnPosition: THREE.Vector3 = PLAYER_SPAWN_POSITION) {
        this.position = spawnPosition.clone();
    }
}
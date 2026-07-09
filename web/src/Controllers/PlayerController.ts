import * as THREE from "three";
import { PlayerModel } from "../Models/PlayerModel";
import { PlayerView } from "../Views/PlayerView";
import { InputController } from "./InputController";
import { CollisionManager } from "../Physics/CollisionManager";
import { CollisionLayer } from "../Physics/OBB";
import { NetworkManager } from "../Services/NetworkManager";
import { DEFAULT_SPEED_Y, START_SEND_TIME_VALUE } from "../Config/GameConfig";
import { NETWORK_SEND_INTERVAL } from "../Config/NetworkConfig";

export class PlayerController 
{
    private model: PlayerModel;
    private view: PlayerView;
    private input: InputController;
    private camera: THREE.PerspectiveCamera;
    private networkManager: NetworkManager;
    private collisionManager: CollisionManager;
    private lastSendTime = START_SEND_TIME_VALUE;

    constructor(model: PlayerModel, view: PlayerView,input: InputController, camera: THREE.PerspectiveCamera, networkManager: NetworkManager, collisionManager: CollisionManager) 
    {
        this.model = model;
        this.view = view;
        this.input = input;
        this.camera = camera;
        this.networkManager = networkManager;
        this.collisionManager = collisionManager;
    }

    public Update(dt: number): void 
    {
        const forward = new THREE.Vector3();

        this.camera.getWorldDirection(forward);
        forward.y = DEFAULT_SPEED_Y;
        forward.normalize();

        const right = new THREE.Vector3().crossVectors(forward, this.camera.up).normalize();

        const wish = new THREE.Vector3();

        if (this.input.IsKeyDown("KeyW")) 
        {
            wish.add(forward);
        }

        if (this.input.IsKeyDown("KeyS")) 
        {
            wish.sub(forward);
        }

        if (this.input.IsKeyDown("KeyD")) 
        {
            wish.add(right);
        }

        if (this.input.IsKeyDown("KeyA")) 
        {
            wish.sub(right);
        }

        const movedInput = wish.lengthSq() > 0;

        if (movedInput) 
        {
            wish.normalize().multiplyScalar(this.model.speed);
        }

        const displacement = wish.multiplyScalar(dt);

        const PLAYER_MASK = CollisionLayer.Static | CollisionLayer.Item;

        this.collisionManager.MoveAndSlide(this.model.body, displacement, PLAYER_MASK);

        this.view.Update(this.model.position);

        const now = performance.now();

        if (movedInput && now - this.lastSendTime > NETWORK_SEND_INTERVAL) 
        {
            this.networkManager.SendMove(this.camera.position, this.camera.rotation);

            this.lastSendTime = now;
        }
    }
}
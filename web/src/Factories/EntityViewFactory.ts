import * as THREE from "three";
import { MeshFactory } from "./MeshFactory";
import { AlienView } from "../Views/AlienView";
import { PlayerView } from "../Views/PlayerView";

export interface AnimatedEntityView {
    AdvanceAnimation(deltaTime: number): void;
    SetMoving?(moving: boolean): void;
    SetAttacking?(attacking: boolean): void;
}

export interface CreatedEntityView {
    object: THREE.Object3D;
    animatedView: AnimatedEntityView | null;
    ownsResources: boolean;
}

export class EntityViewFactory {
    public CreateView(entityType: string): CreatedEntityView | null {
        switch (entityType) {
            case "player": {
                const playerView = new PlayerView();

                return {
                    object: playerView.mesh,
                    animatedView: playerView,
                    ownsResources: false
                };
            }

            case "alien":
            case "monster": {
                const alienView = new AlienView();

                return {
                    object: alienView.mesh,
                    animatedView: alienView,
                    ownsResources: false
                };
            }

            case "door":
            case "cargo":
                return {
                    object: MeshFactory.CreateBox(1, 1, 1),
                    animatedView: null,
                    ownsResources: true
                };

            default:
                return null;
        }
    }
}
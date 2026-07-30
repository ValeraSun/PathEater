import * as THREE from "three";
import { MeshFactory } from "./MeshFactory";
import { AlienView } from "../Views/AlienView";
import { PlayerView } from "../Views/PlayerView";
import { DoorView } from "../Views/DoorView";
import { CreateHole } from "../Views/BreakdownView";

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

            case "door": {
    const doorView = new DoorView();
    return {
        object: doorView.mesh,        // группа, содержащая модель
        animatedView: doorView,       // сам DoorView (реализует AnimatedEntityView)
        ownsResources: false
    };
}
            case "cargo":
                return {
                    object: MeshFactory.CreateBox(1, 1, 1),
                    animatedView: null,
                    ownsResources: true
                };
case "hole": {
    const mesh = CreateHole();
     mesh.scale.set(1, 1, 1);
    return {
        object: mesh,
        animatedView: null, // анимация не нужна
        ownsResources: true // мы создали новый меш
    };
}
            default:
                return null;
        }
    }
}
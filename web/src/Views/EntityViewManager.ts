import * as THREE from "three";
import type { EntityModel } from "../Models/EntityModel";
import { EntityViewFactory, type AnimatedEntityView } from "../Factories/EntityViewFactory";

interface EntityViewRecord {
    object: THREE.Object3D;
    animatedView: AnimatedEntityView | null;
    ownsResources: boolean;
}

const POSITION_INTERPOLATION_SPEED = 12;
const MOVEMENT_DISTANCE_THRESHOLD = 0.01;

export class EntityViewManager {
    private scene: THREE.Scene;
    private entityViewFactory: EntityViewFactory;
    private viewsByEntityId = new Map<string, EntityViewRecord>();

    public constructor(scene: THREE.Scene, entityViewFactory: EntityViewFactory) {
        this.scene = scene;
        this.entityViewFactory = entityViewFactory;
    }

    public CreateEntity(entityModel: EntityModel): void {
        if (this.viewsByEntityId.has(entityModel.id)) {
            return;
        }

        const createdView = this.entityViewFactory.CreateView(entityModel.type);

        if (!createdView) {
            console.warn(`Неизвестный тип сущности: ${entityModel.type}`);
            return;
        }

        createdView.object.position.copy(entityModel.position);
        createdView.object.rotation.y = entityModel.rotationY;

        this.scene.add(createdView.object);

        this.viewsByEntityId.set(entityModel.id, {
            object: createdView.object,
            animatedView: createdView.animatedView,
            ownsResources: createdView.ownsResources
        });
    }

    public UpdateEntity(entityModel: EntityModel, deltaTime: number): void {
        const entityView = this.viewsByEntityId.get(entityModel.id);

        if (!entityView) {
            return;
        }

        const interpolation = 1 - Math.exp(-POSITION_INTERPOLATION_SPEED * deltaTime);
        const distanceToTarget = entityView.object.position.distanceTo(entityModel.targetPosition);

        entityView.animatedView?.SetMoving?.(distanceToTarget > MOVEMENT_DISTANCE_THRESHOLD);
        entityView.animatedView?.SetAttacking?.(entityModel.attacking ?? false);
        entityView.animatedView?.AdvanceAnimation(deltaTime);

        entityView.object.position.lerp(entityModel.targetPosition, interpolation);
        entityView.object.rotation.y = this.InterpolateAngle(
            entityView.object.rotation.y,
            entityModel.targetRotationY,
            interpolation
        );

        entityView.object.visible = !(entityModel.dead ?? false);
    }

    public RemoveEntity(entityId: string): void {
        const entityView = this.viewsByEntityId.get(entityId);

        if (!entityView) {
            return;
        }

        this.scene.remove(entityView.object);

        if (entityView.ownsResources) {
            this.DisposeOwnedObject(entityView.object);
        }

        this.viewsByEntityId.delete(entityId);
    }

    public Clear(): void {
        for (const entityId of [...this.viewsByEntityId.keys()]) {
            this.RemoveEntity(entityId);
        }
    }

    private InterpolateAngle(currentAngle: number, targetAngle: number, interpolation: number): number {
        const angleDifference = Math.atan2(
            Math.sin(targetAngle - currentAngle),
            Math.cos(targetAngle - currentAngle)
        );

        return currentAngle + angleDifference * interpolation;
    }

    private DisposeOwnedObject(object: THREE.Object3D): void {
        object.traverse((child: THREE.Object3D): void => {
            if (!(child instanceof THREE.Mesh)) {
                return;
            }

            child.geometry.dispose();

            if (Array.isArray(child.material)) {
                for (const material of child.material) {
                    material.dispose();
                }
            } else {
                child.material.dispose();
            }
        });
    }
}
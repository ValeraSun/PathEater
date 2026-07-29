import { EntityStore } from "../Models/EntityStore";
import { RadarModel } from "../Models/RadarModel";
import type { EntityCreateInfo, EntityTransformData, EntityUpdateInfo } from "../Network/ServerContracts";
import { IsDoorStateData, IsEntityTransformData, IsAsteroidStateData, IsMonsterStateData, IsShipWireData, IsBulletData } from "../Network/ServerValidators";
import type { DoorStateData } from "../Network/ServerContracts";
import { EntityViewManager } from "../Views/EntityViewManager";
import { ShipView } from "../Views/ShipView";

export class EntityMessageHandler {
    private entityStore: EntityStore;
    private entityViewManager: EntityViewManager;
    private radarModel: RadarModel;
    private shipView: ShipView;
    private localPlayerId: string | null = null;
    private localPlayerStateHandler: ((playerState: EntityTransformData) => void) | null = null;

    public constructor(entityStore: EntityStore, entityViewManager: EntityViewManager, radarModel: RadarModel, shipView: ShipView) {
        this.entityStore = entityStore;
        this.entityViewManager = entityViewManager;
        this.radarModel = radarModel;
        this.shipView = shipView;
    }

    public SetLocalPlayerId(playerId: string): void {
        this.localPlayerId = playerId;
        this.entityStore.RemoveEntity(playerId);
        this.entityViewManager.RemoveEntity(playerId);
    }

    public SetLocalPlayerStateHandler(handler: (playerState: EntityTransformData) => void): void {
        this.localPlayerStateHandler = handler;
    }

    public CreateEntity(entityInformation: EntityCreateInfo): void {
        if (entityInformation.id === this.localPlayerId) {
            this.ApplyLocalPlayerState(entityInformation.data);
            return;
        }

        if (entityInformation.type === "asteroid") {
            this.UpdateAsteroid(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "bullet") {
            this.UpdateBullet(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "cosmoAlien") {
            this.UpdateRadarMonster(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "ship") {
            this.UpdateShip(entityInformation.data);
            return;
        }

        if (entityInformation.type === "door") {
        this.UpdateDoor(entityInformation.id, entityInformation.data);
        return;
        }

        if (!IsEntityTransformData(entityInformation.data)) {
            console.warn("Получены некорректные данные сущности:", entityInformation);
            return;
        }


        const entityModel = this.entityStore.CreateEntity(
            entityInformation.id,
            entityInformation.type,
            entityInformation.data
        );

        this.entityViewManager.CreateEntity(entityModel);
    }

    public UpdateEntity(entityInformation: EntityUpdateInfo): void {
        if (entityInformation.id === this.localPlayerId) {
            this.ApplyLocalPlayerState(entityInformation.data);
            return;
        }

        if (entityInformation.type === "asteroid") {
            this.UpdateAsteroid(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "bullet") {
            this.UpdateBullet(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "cosmoAlien") {
            this.UpdateRadarMonster(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "ship") {
            this.UpdateShip(entityInformation.data);
            return;
        } 
        
        if (entityInformation.type === "door" || entityInformation.type === "spacedoor") {
            this.UpdateDoor(entityInformation.id, entityInformation.data);
            return;
        }

        if (!IsEntityTransformData(entityInformation.data)) {
            console.warn("Получены некорректные данные обновления сущности:", entityInformation);
            return;
        }

        const entityModel = this.entityStore.UpdateEntity(entityInformation.id, entityInformation.data);

        if (!entityModel) {
            this.CreateEntity(entityInformation);
        }
    }

    public RemoveEntity(entityId: string): void {
        this.radarModel.RemoveAsteroid(entityId);
        this.radarModel.RemoveBullet(entityId)
        this.radarModel.RemoveMonster(entityId);
        this.entityViewManager.RemoveEntity(entityId);
        this.entityStore.RemoveEntity(entityId);
    }

    public UpdateEntityViews(deltaTime: number): void {
        for (const entityModel of this.entityStore.GetAllEntities()) {
            this.entityViewManager.UpdateEntity(entityModel, deltaTime);
        }
    }

    public Clear(): void {
        this.entityViewManager.Clear();
        this.entityStore.Clear();
        this.radarModel.Clear();
    }

    private ApplyLocalPlayerState(data: unknown): void {
        if (!IsEntityTransformData(data)) {
            console.warn("Получено некорректное состояние локального игрока:", data);
            return;
        }

        this.localPlayerStateHandler?.(data);
    }

    private UpdateAsteroid(entityId: string, data: unknown): void {
        if (!IsAsteroidStateData(data)) {
            console.warn("Получено некорректное состояние астероида:", data);
            return;
        }

        if (data.destroyed) {
            this.radarModel.RemoveAsteroid(entityId);
            return;
        }

        if (!data.position) {
            return;
        }

        this.radarModel.UpdateAsteroid(
            entityId,
            {
                x: data.position.x,
                y: data.position.y
            },
            data.radius
        );
    }

    private UpdateBullet(entityId: string, data: unknown): void {
        if (!IsBulletData(data)) {
            console.warn("Получено некорректное состояние астероида:", data);
            return;
        }

        if (!data.position) {
            return;
        }

        this.radarModel.UpdateBullet(
            entityId,
            {
                x: data.position.x,
                y: data.position.y
            },
        );
    }

    private UpdateShip(data: unknown): void {
        if (!IsShipWireData(data)) {
            console.log("Получено некорректное состояние корабля:", data);
            return;
        }

        this.shipView.UpdateState(data);

        this.radarModel.UpdateShip(
            data.health,
            data.weapon_direction
                ? { x: data.weapon_direction.x, y: data.weapon_direction.y }
                : undefined
        );
    }

    private UpdateRadarMonster(entityId: string, data: unknown): void {
        if (!IsMonsterStateData(data)) {
            console.warn("Получено некорректное состояние монстра:", data);
            return;
        }

        if (data.died) {
            this.radarModel.RemoveMonster(entityId);
            return;
        }

        if (!data.position) {
            return;
        }

        this.radarModel.UpdateMonster(entityId, {
            x: data.position.x,
            y: data.position.y
        });
    }

private HandleDoor(entityId: string, data: unknown): void {
    if (!IsDoorStateData(data)) {
        console.warn("Получено некорректное состояние двери:", data);
        return;
    }

    // Получаем или создаем модель
    let entityModel = this.entityStore.GetEntity(entityId);
    if (!entityModel) {
        entityModel = this.entityStore.CreateEntity(entityId, "door", data);
        this.entityViewManager.CreateEntity(entityModel);
    } else {
        entityModel = this.entityStore.UpdateEntity(entityId, data);
    }

    // Обновляем вьюху
    const view = this.entityViewManager.GetEntityView(entityId);
    if (view?.object) {
        // Ищем SpaceDoorView
        let doorView: SpaceDoorView | null = null;
        if (view.object instanceof SpaceDoorView) {
            doorView = view.object;
        } else {
            // Ищем в детях
            view.object.children.forEach(child => {
                if (child instanceof SpaceDoorView) {
                    doorView = child;
                }
            });
        }

        if (doorView) {
            // Обновляем позицию
            if (data.position) {
                doorView.mesh.position.set(data.position.x, data.position.y, data.position.z);
            }
            if (data.isOpen !== undefined) {
                doorView.SetState(data.isOpen, data.openProgress);
            }
        }
    }
}

// Добавляем метод Update для обновления анимаций дверей
// Вызывайте его из UpdateEntityViews или отдельно
public UpdateDoor(deltaTime: number): void {
    for (const entity of this.entityStore.GetAllEntities()) {
        if (entity.type === "door") {
            const view = this.entityViewManager.GetEntityView(entity.id);
            if (view?.object) {
                let doorView: DoorView | null = null;
                if (view.object instanceof DoorView) {
                    doorView = view.object;
                } else {
                    view.object.children.forEach(child => {
                        if (child instanceof DoorView) {
                            doorView = child;
                        }
                    });
                }
                doorView?.Update(deltaTime);
            }
        }
    }
}
}
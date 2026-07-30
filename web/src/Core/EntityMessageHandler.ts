import { EntityStore } from "../Models/EntityStore";
import { RadarModel } from "../Models/RadarModel";
import type { EntityCreateInfo, EntityTransformData, EntityUpdateInfo } from "../Network/ServerContracts";
import { IsDoorStateData, IsEntityTransformData, IsAsteroidStateData, IsMonsterStateData, IsShipWireData, IsBulletData } from "../Network/ServerValidators";
import type { DoorStateData } from "../Network/ServerContracts";
import { EntityViewManager } from "../Views/EntityViewManager";
import { ShipView } from "../Views/ShipView";
import { IsHoleStateData } from "../Network/ServerValidators";
import type { HoleStateData } from "../Network/ServerContracts";
import { DoorView } from "../Views/DoorView";

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

if (entityInformation.type === "hole") {
    this.UpdateHole(entityInformation.id, entityInformation.data);
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
        
        if (entityInformation.type === "hole") {
    this.UpdateHole(entityInformation.id, entityInformation.data);
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



private UpdateDoor(entityId: string, data: unknown): void {
    console.log("[Door] UpdateDoor called", entityId, data);
    if (!IsDoorStateData(data)) {
        console.warn("[Door] Некорректные данные двери:", data);
        return;
    }

    let entityModel = this.entityStore.GetEntity(entityId);
    if (!entityModel) {
        console.log("[Door] Creating new door entity", entityId);
        entityModel = this.entityStore.CreateEntity(entityId, "door", data);
        this.entityViewManager.CreateEntity(entityModel);
    } else {
        entityModel = this.entityStore.UpdateEntity(entityId, data);
    }

    const view = this.entityViewManager.GetEntityView(entityId);
    console.log("[Door] View for door:", view);
    if (view?.animatedView && view.animatedView instanceof DoorView) {
        const doorView = view.animatedView;
        console.log("[Door] DoorView found, applying data");
        if (data.position) {
            doorView.mesh.position.set(data.position.x, data.position.y, data.position.z);
            console.log("[Door] Position set to", data.position);
        }
        if (data.rotation) {
            doorView.mesh.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
            console.log("[Door] Rotation set to", data.rotation);
        }
        if (data.isOpen !== undefined) {
            doorView.SetState(data.isOpen, data.openProgress);
            console.log("[Door] State set to open=", data.isOpen, "progress=", data.openProgress);
        }
    } else {
        console.warn("[Door] DoorView not found for entity", entityId);
    }
}

private UpdateHole(entityId: string, data: unknown): void {
    if (!IsHoleStateData(data)) {
        console.warn("Некорректные данные поломки:", data);
        return;
    }

    // Получаем или создаём модель
    let entityModel = this.entityStore.GetEntity(entityId);
    if (!entityModel) {
        entityModel = this.entityStore.CreateEntity(entityId, "hole", data);
        this.entityViewManager.CreateEntity(entityModel);
    } else {
        entityModel = this.entityStore.UpdateEntity(entityId, data);
    }

    // Обновляем вьюху
    const view = this.entityViewManager.GetEntityView(entityId);
    if (view?.object) {
        if (data.position) {
            view.object.position.set(data.position.x, data.position.y, data.position.z);
        }
        if (data.rotation) {
            view.object.rotation.set(data.rotation.x, data.rotation.y, data.rotation.z);
        }
        if (data.radius) {
            // Масштабируем меш (круг) – радиус влияет на размер
            view.object.scale.set(data.radius, data.radius, data.radius);
        }
    }
}

}


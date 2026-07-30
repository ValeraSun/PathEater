import { EntityStore } from "../Models/EntityStore";
import { RadarModel } from "../Models/RadarModel";
import type { EntityCreateInfo, EntityTransformData, EntityUpdateInfo } from "../Network/ServerContracts";
import {  IsEntityTransformData, IsAsteroidStateData, IsMonsterStateData, IsShipWireData, IsBulletData, IsDoorStateData, IsBreakdownStateData } from "../Network/ServerValidators";
import { EntityViewManager } from "../Views/EntityViewManager";
import { ShipView } from "../Views/ShipView";

export class EntityMessageHandler {
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

    public SetMonsterDamagedHandler(handler: () => void): void {
        this.monsterDamagedHandler = handler;
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

         if (entityInformation.type === "hole" || entityInformation.type === "breakdown") {
            this.UpdateBreakdown(entityInformation.id, entityInformation.data);
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

        if (entityInformation.type === "door") {
            this.UpdateDoor(entityInformation.id, entityInformation.data);
            return;
        }

        if (entityInformation.type === "hole" || entityInformation.type === "breakdown") {
            this.UpdateBreakdown(entityInformation.id, entityInformation.data);
            return;
        }

        if (!IsEntityTransformData(entityInformation.data)) {
            console.warn("Получены некорректные данные обновления сущности:", entityInformation);
            return;
        }

        const previousHealth = this.entityStore.GetEntity(entityInformation.id)?.health;
        const entityModel = this.entityStore.UpdateEntity(entityInformation.id, entityInformation.data);

        if (!entityModel) {
            this.CreateEntity(entityInformation);
            return;
        }

        this.CheckMonsterDamage(entityInformation.type, previousHealth, entityModel.health);
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

    private CheckMonsterDamage(entityType: string, previousHealth: number | undefined, currentHealth: number | undefined): void {
        const isMonster = entityType === "alien" || entityType === "monster";

        if (!isMonster) {
            return;
        }

        if (typeof previousHealth !== "number" || typeof currentHealth !== "number") {
            return;
        }

        if (currentHealth < previousHealth) {
            this.monsterDamagedHandler?.();
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

    private entityStore: EntityStore;
    private entityViewManager: EntityViewManager;
    private radarModel: RadarModel;
    private shipView: ShipView;
    private localPlayerId: string | null = null;
    private localPlayerStateHandler: ((playerState: EntityTransformData) => void) | null = null;
    private monsterDamagedHandler: (() => void) | null = null;
    private UpdateDoor(entityId: string, data: unknown): void {
    if (!IsDoorStateData(data)) {
        console.warn("UpdateDoor: некорректные данные", data);
        return;
    }

    let entityModel = this.entityStore.GetEntity(entityId);
    if (!entityModel) {
        // Создаём модель с временной позицией, если position нет
        const tempPos = data.position || { x: 0, y: 0, z: 0 };
        entityModel = this.entityStore.CreateEntity(entityId, "door", data);
        entityModel.position.set(tempPos.x, tempPos.y, tempPos.z);
        entityModel.targetPosition.copy(entityModel.position);
        // Всегда создаём вьюху
        this.entityViewManager.CreateEntity(entityModel);
        console.log("UpdateDoor: создана новая дверь", entityId);
    } else {
        if (data.position) {
            entityModel.position.set(data.position.x, data.position.y, data.position.z);
            entityModel.targetPosition.copy(entityModel.position);
        }
        this.entityStore.UpdateEntity(entityId, data);
    }
}

  private UpdateBreakdown(entityId: string, data: unknown): void {
        // Создаем валидатор для данных поломки
        if (!IsBreakdownStateData(data)) {
            console.warn("UpdateBreakdown: некорректные данные", data);
            return;
        }

        // Получаем или создаем модель
        let entityModel = this.entityStore.GetEntity(entityId);
        const anyData = data as any;
        
        if (!entityModel) {
            // Создаем модель с позицией
            const tempPos = anyData.position || { x: 0, y: 0, z: 0 };
            entityModel = this.entityStore.CreateEntity(
                entityId, 
                "breakdown", 
                { position: tempPos }
            );
            entityModel.position.set(tempPos.x, tempPos.y, tempPos.z);
            entityModel.targetPosition.copy(entityModel.position);
            
            // Создаем вьюху с данными
            this.entityViewManager.CreateEntity(entityModel, anyData);
        } else {
            // Обновляем позицию
            if (anyData.position) {
                entityModel.position.set(anyData.position.x, anyData.position.y, anyData.position.z);
                entityModel.targetPosition.copy(entityModel.position);
            }
        }

        // Обновляем вьюху
        const view = this.entityViewManager.GetEntityView(entityId);
        if (view?.object && view.object instanceof THREE.Group) {
            // Если есть позиция
            if (anyData.position) {
                view.object.position.set(anyData.position.x, anyData.position.y, anyData.position.z);
            }
            
            // Если есть радиус
            if (anyData.radius) {
                // Ищем Mesh в группе и обновляем его
                view.object.children.forEach(child => {
                    if (child instanceof THREE.Mesh) {
                        const scale = anyData.radius / 1.2; // 1.2 - базовый радиус
                        child.scale.set(scale, scale, scale);
                    }
                });
            }
            
            // Ориентация на стене - если передана нормаль
            if (anyData.normal) {
                const normal = new THREE.Vector3(
                    anyData.normal.x,
                    anyData.normal.y,
                    anyData.normal.z
                );
                const up = new THREE.Vector3(0, 0, 1);
                const quaternion = new THREE.Quaternion().setFromUnitVectors(up, normal.clone().normalize());
                view.object.quaternion.copy(quaternion);
            }
        }
    }
}




interface Vector2D {
    x: number;
    y: number;
}

interface AsteroidRadarEntry {
    position: Vector2D;
    radius: number;
}

export interface RadarDisplayState {
    ship: {
        x: number;
        y: number;
        rotationY: number;
        hp: number;
    };
    asteroids: Array<{
        id: string;
        x: number;
        y: number;
        radius: number;
    }>;
    monsters: Array<{
        id: string;
        x: number;
        y: number;
    }>;
    bullets: Array<{
        id: string;
        x: number;
        y: number;
    }>;
}

const DEFAULT_ASTEROID_RADIUS = 90;

export class RadarModel {
    private shipHealth = 100;
    private shipWeaponDirection: Vector2D = { x: 0, y: -1 };
    private asteroidsById = new Map<string, AsteroidRadarEntry>();
    private monstersById = new Map<string, Vector2D>();
    private bulletsById = new Map<string, Vector2D>();
    

    private changedHandler: ((radarState: RadarDisplayState) => void) | null = null;

    public SetChangedHandler(handler: (radarState: RadarDisplayState) => void): void {
        this.changedHandler = handler;
        this.EmitChanged();
    }

    public UpdateShip(health?: number, weaponDirection?: Vector2D): void {
        if (typeof health === "number") {
            this.shipHealth = health;
        }

        if (weaponDirection) {
            this.shipWeaponDirection = { ...weaponDirection };
        }

        this.EmitChanged();
    }

    public UpdateAsteroid(entityId: string, position: Vector2D, radius = DEFAULT_ASTEROID_RADIUS): void {
        this.asteroidsById.set(entityId, {
            position: { ...position },
            radius
        });

        this.EmitChanged();
    }

    public UpdateBullet(entityId: string, position: Vector2D): void {
        this.bulletsById.set(entityId, { ...position });
        this.EmitChanged();
    }

    public RemoveAsteroid(entityId: string): boolean {
        const removed = this.asteroidsById.delete(entityId);

        if (removed) {
            this.EmitChanged();
        }

        return removed;
    }

    public UpdateMonster(entityId: string, position: Vector2D): void {
        this.monstersById.set(entityId, { ...position });
        this.EmitChanged();
    }

    public RemoveMonster(entityId: string): boolean {
        const removed = this.monstersById.delete(entityId);

        if (removed) {
            this.EmitChanged();
        }

        return removed;
    }

    public RemoveBullet(entityId: string): boolean {
        const removed = this.bulletsById.delete(entityId);

        if (removed) {
            this.EmitChanged();
        }

        return removed;
    }

    public Clear(): void {
        this.shipHealth = 100;
        this.shipWeaponDirection = { x: 0, y: -1 };
        this.asteroidsById.clear();
        this.monstersById.clear();
        this.EmitChanged();
    }

    private EmitChanged(): void {
        if (!this.changedHandler) {
            return;
        }

        const rotationY = Math.atan2(this.shipWeaponDirection.x, this.shipWeaponDirection.y);

        this.changedHandler({
            ship: {
                x: 0,
                y: 0,
                rotationY,
                hp: this.shipHealth
            },
            asteroids: Array.from(this.asteroidsById, ([id, asteroid]) => ({
                id,
                x: asteroid.position.x,
                y: asteroid.position.y,
                radius: asteroid.radius
            })),
            monsters: Array.from(this.monstersById, ([id, position]) => ({
                id,
                x: position.x,
                y: position.y
            })),
            bullets: Array.from(this.bulletsById, ([id, position]) => ({
                id,
                x: position.x,
                y: position.y
            }))

        });
    }
}
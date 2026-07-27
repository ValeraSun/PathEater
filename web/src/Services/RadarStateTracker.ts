import type { NavigationDisplayState } from "../Views/ComputerView";

interface Vec2 {
    x: number;
    y: number;
}

interface AsteroidRadarEntry {
    position: Vec2;
    radius: number;
}

const DEFAULT_ASTEROID_RADIUS = 90.0;

export class RadarStateTracker {
    private shipHp = 100;
    private shipWeaponDirection: Vec2 = { x: 0, y: -1 };
    private readonly asteroids = new Map<string, AsteroidRadarEntry>();
    private readonly monsters = new Map<string, Vec2>();

    private onChanged: ((state: NavigationDisplayState) => void) | null = null;

    public SetChangedHandler(handler: (state: NavigationDisplayState) => void): void {
        this.onChanged = handler;
        this.Emit();
    }

    public UpdateShip(hp: number | undefined, weaponDirection: Vec2 | undefined): void {
        if (typeof hp === "number") {
            this.shipHp = hp;
        }

        if (weaponDirection) {
            this.shipWeaponDirection = weaponDirection;
        }

        this.Emit();
    }

    public UpdateAsteroid(id: string, position: Vec2 | null, radius?: number): void {
        if (!position) {
            this.asteroids.delete(id);
        } else {
            this.asteroids.set(id, { position, radius: radius ?? DEFAULT_ASTEROID_RADIUS });
        }

        this.Emit();
    }

    public RemoveAsteroid(id: string): boolean {
        const removed = this.asteroids.delete(id);
        if (removed) this.Emit();
        return removed;
    }

    public UpdateMonster(id: string, position: Vec2 | null): void {
        if (!position) {
            this.monsters.delete(id);
        } else {
            this.monsters.set(id, position);
        }

        this.Emit();
    }

    public RemoveMonster(id: string): boolean {
        const removed = this.monsters.delete(id);
        if (removed) this.Emit();
        return removed;
    }

    private Emit(): void {
        if (!this.onChanged) return;

        const rotationY = Math.atan2(this.shipWeaponDirection.x, this.shipWeaponDirection.y);

        this.onChanged({
            ship: { x: 0, y: 0, rotationY, hp: this.shipHp },
            asteroids: Array.from(this.asteroids, ([id, a]) => ({
                id,
                x: a.position.x,
                y: a.position.y,
                radius: a.radius
            })),
            monsters: Array.from(this.monsters, ([id, pos]) => ({ id, x: pos.x, y: pos.y }))
        });
    }
}
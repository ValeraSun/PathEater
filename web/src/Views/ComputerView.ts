import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

export type NavigationDisplayState = {
    ship: {
        x: number;
        z: number;
        rotationY: number;
        hp: number;
    };

    asteroids: Array<{
        id: string;
        x: number;
        z: number;
    }>;

    monsters: Array<{
        id: string;
        x: number;
        z: number;
    }>;
};

export class ComputerView 
{
    private id = "bridge_computer_1";
    private group = new THREE.Group();
    private canvas: HTMLCanvasElement;
    private context: CanvasRenderingContext2D;
    private texture: THREE.CanvasTexture;

    private screenMesh: THREE.Mesh<THREE.PlaneGeometry, THREE.MeshBasicMaterial>;

    private lockedBy: string | null = null;

    public constructor() 
    {
        this.canvas = document.createElement("canvas");
        this.canvas.width = 1024;
        this.canvas.height = 512;

        const context = this.canvas.getContext("2d");

        if (!context) 
        {
            throw new Error("Не удалось создать контекст экрана компьютера");
        }

        this.context = context;
        this.texture = new THREE.CanvasTexture(this.canvas);
        this.texture.colorSpace = THREE.SRGBColorSpace;
        this.texture.minFilter = THREE.LinearFilter;
        this.texture.magFilter = THREE.LinearFilter;
        this.screenMesh = this.CreateScreenMesh();
        this.DrawInitialScreen();
        this.LoadComputerModel();
        this.group.add(this.screenMesh);
    }

    private LoadComputerModel(): void 
    {
        const loader = new GLTFLoader();

        loader.load("/models/bridge_console.glb",
            (gltf) => {
                const model = gltf.scene;
                model.position.set(0, 0, 0);
                model.scale.set(1, 1, 1);
                this.group.add(model);
                console.log("Bridge console loaded");
            },
            undefined,
            (error) => {
                console.error(
                    "Ошибка загрузки bridge_console.glb:",
                    error
                );
            }
        );
    }

    private CreateScreenMesh(): THREE.Mesh<THREE.PlaneGeometry, THREE.MeshBasicMaterial> 
    {
        const geometry = new THREE.PlaneGeometry(6, 2.5);

        const material = new THREE.MeshBasicMaterial({
            map: this.texture,
            side: THREE.DoubleSide,
            toneMapped: false,

            polygonOffset: true,
            polygonOffsetFactor: -2,
            polygonOffsetUnits: -2
        });

        const screen = new THREE.Mesh(geometry, material);
        screen.position.set(41.5, 3, -9);
        screen.rotation.y = -Math.PI / 2;
        screen.renderOrder = 10;

        return screen;
    }

    private DrawInitialScreen(): void 
    {
        const width = this.canvas.width;
        const height = this.canvas.height;

        this.context.fillStyle = "#020b12";
        this.context.fillRect(0, 0, width, height);

        this.context.strokeStyle = "#00d9ff";
        this.context.lineWidth = 6;
        this.context.strokeRect(12, 12, width - 24, height - 24);

        this.context.fillStyle = "#00e5ff";
        this.context.font = "bold 40px monospace";
        this.context.fillText("NAVIGATION SYSTEM", 55, 75);
        this.context.font = "28px monospace";
        this.context.fillText("WAITING FOR SERVER DATA...", 55, 140);
        this.texture.needsUpdate = true;
    }

    public UpdateDisplay(state: NavigationDisplayState): void 
    {
        const width = this.canvas.width;
        const height = this.canvas.height;

        this.context.fillStyle = "#020b12";
        this.context.fillRect(0, 0, width, height);

        this.context.strokeStyle = "#00d9ff";
        this.context.lineWidth = 6;
        this.context.strokeRect(12, 12, width - 24, height - 24);
        this.context.fillStyle = "#00e5ff";
        this.context.font = "bold 36px monospace";
        this.context.fillText(
            "SHIP RADAR",
            45,
            60
        );

        this.context.fillStyle = "#ffffff";
        this.context.font = "26px monospace";
        this.context.fillText(
            `SHIP HP: ${Math.round(state.ship.hp)}`,
            45,
            105
        );

        const centerX = width / 2;
        const centerY = height / 2 + 25;
        const mapScale = 8;

        this.DrawShip(centerX, centerY, state.ship.rotationY);

        for (const asteroid of state.asteroids) 
        {
            const screenX = centerX + (asteroid.x - state.ship.x) * mapScale;
            const screenY = centerY + (asteroid.z - state.ship.z) * mapScale;

            this.DrawAsteroid(screenX, screenY);
        }

        for (const monster of state.monsters) 
        {
            const screenX = centerX + (monster.x - state.ship.x) * mapScale;
            const screenY = centerY + (monster.z - state.ship.z) * mapScale;

            this.DrawMonster(screenX, screenY);
        }

        this.context.font = "22px monospace";

        this.context.fillStyle = "#ffad33";
        this.context.fillText(
            `ASTEROIDS: ${state.asteroids.length}`,
            45,
            height - 45
        );

        this.context.fillStyle = "#ff4055";
        this.context.fillText(
            `MONSTERS: ${state.monsters.length}`,
            330,
            height - 45
        );

        this.context.fillStyle =
            this.lockedBy === null
                ? "#45ff88"
                : "#ff4055";

        this.context.fillText(
            this.lockedBy === null
                ? "CONTROL: AVAILABLE"
                : "CONTROL: OCCUPIED",
            675,
            height - 45
        );

        this.texture.needsUpdate = true;
    }

    private DrawShip(x: number, y: number, rotationY: number): void 
    {
        this.context.save();

        this.context.translate(x, y);
        this.context.rotate(-rotationY);

        this.context.fillStyle = "#45ff88";
        this.context.strokeStyle = "#ffffff";
        this.context.lineWidth = 3;

        this.context.beginPath();
        this.context.moveTo(28, 0);
        this.context.lineTo(-18, -18);
        this.context.lineTo(-10, 0);
        this.context.lineTo(-18, 18);
        this.context.closePath();

        this.context.fill();
        this.context.stroke();

        this.context.restore();
    }

    private DrawAsteroid(x: number, y: number): void 
    {
        if (!this.IsInsideScreen(x, y)) 
        {
            return;
        }

        this.context.fillStyle = "#ffad33";
        this.context.strokeStyle = "#fff2d5";
        this.context.lineWidth = 2;

        this.context.beginPath();
        this.context.arc(x, y, 10, 0, Math.PI * 2);
        this.context.fill();
        this.context.stroke();
    }

    private DrawMonster(x: number, y: number): void 
    {
        if (!this.IsInsideScreen(x, y)) 
            {
            return;
        }

        this.context.fillStyle = "#ff4055";

        this.context.fillRect(x - 9, y - 9, 18, 18);
    }

    private IsInsideScreen(x: number, y: number): boolean 
    {
        return (
            x >= 25 &&
            x <= this.canvas.width - 25 &&
            y >= 120 &&
            y <= this.canvas.height - 75
        );
    }

    public SetLockedBy(playerId: string | null): void 
    {
        this.lockedBy = playerId;
    }

    public IsFree(): boolean 
    {
        return this.lockedBy === null;
    }

    public IsOwnedBy(playerId: string): boolean 
    {
        return this.lockedBy === playerId;
    }

    public GetLockedBy(): string | null 
    {
        return this.lockedBy;
    }

    public GetId(): string 
    {
        return this.id;
    }

    public GetInteractionPosition(): THREE.Vector3 
    {
        return new THREE.Vector3(39, 1, -9);
    }

    public GetObject(): THREE.Group 
    {
        return this.group;
    }
}
import * as THREE from "three";
import { WindowResize } from "../Utils/WindowResize";
import { ComputerView } from "./ComputerView";
import { HealthbarView } from "./HealthbarView";
import { PlayerView } from "./PlayerView";
import { ShipView } from "./ShipView";
import { SpaceView } from "./SpaceView";
import {ButtonView} from "./ButtonView"

export class GameView {
    private scene: THREE.Scene;
    private camera: THREE.PerspectiveCamera;
    private renderer: THREE.WebGLRenderer;
    private shipView: ShipView;
    private spaceView: SpaceView;
    private computerView: ComputerView;
    

    private healthBarView: HealthbarView | null = null;
    private initialized = false;

    public constructor() {
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        this.renderer = new THREE.WebGLRenderer({antialias: true});
        this.shipView = new ShipView();
        this.spaceView = new SpaceView();
        this.computerView = new ComputerView();
    }

    public Initialize(gameContainer: HTMLElement, healthBarContainer: HTMLElement): void {
        if (this.initialized) {
            return;
        }

        this.initialized = true;

        this.renderer.setSize(window.innerWidth, window.innerHeight);
        gameContainer.appendChild(this.renderer.domElement);

        this.healthBarView = new HealthbarView(healthBarContainer);

        this.camera.position.set(0, 5, 12);
        this.camera.lookAt(0, 0, 1);

        this.AddLights();
        const buttons = this.InitButtons()
        
        buttons.forEach(element => {
             this.scene.add(element.GetObject());
        });

        this.scene.add(this.spaceView.GetObject());
        this.scene.add(this.shipView.GetObject());
        this.scene.add(this.computerView.GetObject());

        WindowResize.Handle(this.camera, this.renderer);
    }

    public Render(): void {
        this.renderer.render(this.scene, this.camera);
    }

    public AttachPlayerView(playerView: PlayerView): void {
        this.scene.add(playerView.mesh);
    }

    public RemovePlayerView(playerView: PlayerView): void {
        this.scene.remove(playerView.mesh);
    }

    public SetPlayerHealth(health: number): void {
        this.healthBarView?.SetHealth(health);
    }

    public ShowPlayerHealth(): void {
        this.healthBarView?.Show();
    }

    public HidePlayerHealth(): void {
        this.healthBarView?.Hide();
    }

    public GetScene(): THREE.Scene {
        return this.scene;
    }

    public GetCamera(): THREE.PerspectiveCamera {
        return this.camera;
    }

    public GetRendererDomElement(): HTMLCanvasElement {
        return this.renderer.domElement;
    }

    public GetShipView(): ShipView {
        return this.shipView;
    }

    public GetComputerView(): ComputerView {
        return this.computerView;
    }

    private AddLights(): void {
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
        const hemisphereLight = new THREE.HemisphereLight(0xffffff, 0x444466, 1);
        const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
        directionalLight.position.set(10, 10, 10);
        this.scene.add(ambientLight);
        this.scene.add(hemisphereLight);
        this.scene.add(directionalLight);
    }

    private InitButtons(): ButtonView[] {
        const buttons: ButtonView[] = [
            new ButtonView(23.4, 2.7, -17.5),
            new ButtonView(26.59, 2.7, -17.5),
            new ButtonView(29.8, 2.7, -17.5),
        ];

        return buttons;
    }
}


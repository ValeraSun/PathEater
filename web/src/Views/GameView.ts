import * as THREE from "three";
import { ShipView } from "./ShipView";
import { SpaceView } from "./SpaceView";
import { WindowResize } from "../Utils/WindowResize";
import { PlayerView } from "./PlayerView";
import { ComputerView } from "./ComputerView";

export class GameView 
{
    public Render(): void 
    {
        this.renderer.render(this.scene, this.camera);
    }

    public Init(): void 
    {
        this.renderer.setSize(window.innerWidth, window.innerHeight);
        const container = document.getElementById("app");
        container?.appendChild(this.renderer.domElement);

        this.camera.position.set(0, 5, 12);
        this.camera.lookAt(0, 0, 1);

        this.addLights();

        this.scene.add(this.spaceView.GetObject());
        this.scene.add(this.shipView.GetObject());
        this.scene.add(this.computerView.GetObject());

        WindowResize.Handle(this.camera, this.renderer);
    }

    public AttachPlayerView(playerView: PlayerView): void {
        this.scene.add(playerView.mesh);
    }

    public GetCamera(): THREE.PerspectiveCamera 
    {
        return this.camera;
    }

    public GetScene(): THREE.Scene 
    {
        return this.scene;
    }

    public GetRendererDomElement(): HTMLCanvasElement 
    {
        return this.renderer.domElement;
    }

    public GetComputerView(): ComputerView 
    {
        return this.computerView;
    }

    private scene: THREE.Scene;
    private camera: THREE.PerspectiveCamera;
    private renderer: THREE.WebGLRenderer;

    private shipView: ShipView;
    private spaceView: SpaceView;
    private computerView: ComputerView;

    constructor() 
    {
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
        this.renderer = new THREE.WebGLRenderer({antialias: true,});
        this.shipView = new ShipView();
        this.spaceView = new SpaceView();
        this.computerView = new ComputerView();
    }

    private addLights(): void 
    {
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
        this.scene.add(ambientLight);

        const hemisphereLight = new THREE.HemisphereLight(0xffffff, 0x444466, 1.0);

        this.scene.add(hemisphereLight);

        const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
        directionalLight.position.set(10, 10, 10);
        this.scene.add(directionalLight);
    }  
}
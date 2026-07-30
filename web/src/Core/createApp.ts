import * as THREE from "three";
import { EntityMessageHandler } from "./EntityMessageHandler";
import { Game } from "./Game";
import { GameSession } from "./GameSession";
import { InputController } from "../Controllers/InputController";
import { EntityViewFactory } from "../Factories/EntityViewFactory";
import { ComputerModel } from "../Models/ComputerModel";
import { EntityStore } from "../Models/EntityStore";
import { RadarModel } from "../Models/RadarModel";
import { GameServerGateway } from "../Network/GameServerGateway";
import { MusicManager } from "../Services/MusicManager";
import { EntityViewManager } from "../Views/EntityViewManager";
import { GameView } from "../Views/GameView";
import { InteractionView } from "../Views/InteractionView";

export function CreateApplication(): Game {
    const gameView = new GameView();
    const inputController = new InputController();
    const entityStore = new EntityStore();
    const radarModel = new RadarModel();
    const entityViewFactory = new EntityViewFactory();

    const entityViewManager = new EntityViewManager(
        gameView.GetScene(),
        entityViewFactory
    );

    const entityMessageHandler = new EntityMessageHandler(
        entityStore,
        entityViewManager,
        radarModel, 
        gameView.GetShipView()
    );

    const gameServerGateway = new GameServerGateway(entityMessageHandler);
    const computerModel = new ComputerModel("bridge_computer_1",new THREE.Vector3(39, 1, -9));
    const interactionView = new InteractionView();
    const musicManager = new MusicManager();

    radarModel.SetChangedHandler((radarState): void => {
        gameView.GetComputerView().UpdateDisplay(radarState);
    });

    const gameSession = new GameSession(
        gameView,
        inputController,
        gameServerGateway,
        entityMessageHandler,
        computerModel,
        interactionView
    );

    entityMessageHandler.SetLocalPlayerStateHandler((playerState): void => {
        gameSession.ApplyLocalPlayerState(playerState);
    });

    return new Game(gameServerGateway, gameView, gameSession, musicManager);
}
import { LobbyController } from "./Controllers/LobbyController";
import { MatchController } from "./Controllers/MatchController";
import { MenuController } from "./Controllers/MenuController";
import { RoomController } from "./Controllers/RoomController";
import { CreateApplication } from "./Core/createApp";
import { AppView } from "./Views/AppView";
import { LobbyView } from "./Views/LobbyView";
import { MenuView } from "./Views/MenuView";
import { RoomView } from "./Views/RoomView";

const application = CreateApplication();

const appView = new AppView();
const menuView = new MenuView();
const roomView = new RoomView();
const lobbyView = new LobbyView();

const gameServerGateway = application.GetGameServerGateway();
const musicManager = application.GetMusicManager();

application.InitializeView(appView.GetGameContainer(), appView.GetHealthBarContainer());

const menuController = new MenuController(application, appView, menuView);
const roomController = new RoomController(gameServerGateway, musicManager, appView, roomView, lobbyView);
const lobbyController = new LobbyController(gameServerGateway, musicManager, appView, roomView, lobbyView);
const matchController = new MatchController(application, gameServerGateway,musicManager, appView);

menuController.Initialize();
roomController.Initialize();
lobbyController.Initialize();
matchController.Initialize();

document.addEventListener(
    "click",
    (): void => {
        void musicManager.PlayMusic("menu");
    },
    {
        once: true
    }
);

appView.ShowScreen("menu");
import { Game } from "../Core/Game";
import { AppView } from "../Views/AppView";
import { MenuView } from "../Views/MenuView";

export class MenuController {
    private application: Game;
    private appView: AppView;
    private menuView: MenuView;
    private connectedToServer = false;

    public constructor(application: Game, appView: AppView, menuView: MenuView) {
        this.application = application;
        this.appView = appView;
        this.menuView = menuView;
    }

    public Initialize(): void {
        this.menuView.SetPlayHandler(this.HandlePlayButtonClick);
    }

    private HandlePlayButtonClick =
        async (): Promise<void> => {
            if (this.connectedToServer) {
                this.appView.ShowScreen("room");
                return;
            }

            this.menuView.SetConnectionPending(true);

            try {
                await this.application.ConnectToServer();
                this.connectedToServer = true;
                this.menuView.SetConnectionPending(false);
                this.appView.ShowScreen("room");
            } catch (error: unknown) {
                //this.menuView.ShowConnectionError(error);
            }
        };
}
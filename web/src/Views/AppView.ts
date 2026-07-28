export type ScreenName = "menu" | "room" | "lobby" | "game" | "lose" | "win";

export class AppView {
    private screens: Record<ScreenName, HTMLElement>;
    private gameContainer: HTMLElement;
    private healthBarContainer: HTMLElement;

    public constructor() {
        this.screens = {
            menu: this.GetRequiredElement("menu-screen"),
            room: this.GetRequiredElement("room-screen"),
            lobby: this.GetRequiredElement("lobby-screen"),
            game: this.GetRequiredElement("game-screen"),
            lose: this.GetRequiredElement("lose-screen"),
            win: this.GetRequiredElement("win-screen")
        };

        this.gameContainer = this.GetRequiredElement("game-container");
        this.healthBarContainer = this.GetRequiredElement("healthbar-container");
    }

    public ShowScreen(screenName: ScreenName): void {
        for (const screen of Object.values(this.screens)) {
            screen.classList.remove("screen--active");
        }

        this.screens[screenName].classList.add("screen--active");

        document.body.classList.remove(
            "background--menu",
            "background--room",
            "background--lobby",
            "background--game",
            "background--lose",
            "background--win",
            "background--blurred"
        );

        document.body.classList.add(`background--${screenName}`);
        const shouldBlurBackground = screenName === "room" || screenName === "lobby";
        document.body.classList.toggle("background--blurred", shouldBlurBackground);
    }

    public GetGameContainer(): HTMLElement {
        return this.gameContainer;
    }

    public GetHealthBarContainer(): HTMLElement {
        return this.healthBarContainer;
    }

    private GetRequiredElement(elementId: string): HTMLElement {
        const element = document.getElementById(elementId);

        if (!element) {
            throw new Error(`Не найден элемент #${elementId}`);
        }

        return element;
    }
}
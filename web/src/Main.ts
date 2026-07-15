import { Game } from "./Core/Game";

type ScreenName = "menu" | "room" | "lobby" | "game";

const screens: Record<ScreenName, HTMLElement> = {
    menu: getElement("menu-screen"),
    room: getElement("room-screen"),
    lobby: getElement("lobby-screen"),
    game: getElement("game-screen")
};

const playButton = getButton("play-button");
const createRoomButton = getButton("create-room-button");
const joinRoomButton = getButton("join-room-button");
const roomBackButton = getButton("room-back-button");
const leaveLobbyButton = getButton("leave-lobby-button");
const startGameButton = getButton("start-game-button");

const roomCodeInput = getInput("room-code-input");

const menuStatus = getElement("menu-status");
const roomStatus = getElement("room-status");
const lobbyStatus = getElement("lobby-status");
const gameContainer = getElement("game-container");

const game = Game.GetInstance();
const gateway = game.GetGateway();

let connected = false;
let isHost = false;
let roomId: string | null = null;

function showScreen(name: ScreenName): void {
    for (const screen of Object.values(screens)) {
        screen.classList.remove("screen--active");
    }

    screens[name].classList.add("screen--active");
}

playButton.addEventListener("click", async () => {
    if (connected) {
        showScreen("room");
        return;
    }

    playButton.disabled = true;
    menuStatus.textContent = "Идет подключение к серверу";

    try {
        await game.Connect();

        connected = true;
        showScreen("room");
    } catch (error) {
        menuStatus.textContent =
            error instanceof Error
                ? error.message
                : "Сервер недоступен";

        playButton.disabled = false;
    }
});

createRoomButton.addEventListener("click", async () => {
    setRoomButtonsDisabled(true);
    roomStatus.textContent = "Создаем комнату";

    try {
        const room = await gateway.CreateRoom();

        roomId = room.roomId;
        isHost = true;

        startGameButton.hidden = false;
        lobbyStatus.textContent = `Код комнаты: ${roomId}`;

        showScreen("lobby");
    } catch (error) {
        roomStatus.textContent =
            error instanceof Error
                ? error.message
                : "Не удалось создать комнату";

        setRoomButtonsDisabled(false);
    }
});

joinRoomButton.addEventListener("click", async () => {
    const code = roomCodeInput.value.trim().toUpperCase();

    if (!code) {
        roomStatus.textContent = "Введите код комнаты";
        roomCodeInput.focus();
        return;
    }

    setRoomButtonsDisabled(true);
    roomStatus.textContent = `Подключение к комнате ${code}…`;

    try {
        const room = await gateway.JoinRoom(code);

        roomId = room.roomId;
        isHost = false;

        startGameButton.hidden = true;
        lobbyStatus.textContent =
            "Ожидание запуска игры командиром…";

        showScreen("lobby");
    } catch (error) {
        roomStatus.textContent =
            error instanceof Error
                ? error.message
                : "Комната не найдена";

        setRoomButtonsDisabled(false);
    }
});

roomBackButton.addEventListener("click", () => {
    showScreen("menu");
});

leaveLobbyButton.addEventListener("click", () => {
    roomId = null;
    isHost = false;

    roomStatus.textContent = "";
    setRoomButtonsDisabled(false);

    showScreen("room");
});

startGameButton.addEventListener("click", () => {
    if (!isHost) {
        return;
    }

    startGameButton.disabled = true;
    lobbyStatus.textContent = "Запуск игры…";

    gateway.StartGame();
});

gateway.onGameStarted = async payload => {
    await game.StartMatch(payload);

    showScreen("game");
    game.MountTo(gameContainer);
};

roomCodeInput.addEventListener("input", () => {
    roomCodeInput.value = roomCodeInput.value
        .toUpperCase()
        .replace(/[^A-Z0-9]/g, "");
});

function setRoomButtonsDisabled(disabled: boolean): void {
    createRoomButton.disabled = disabled;
    joinRoomButton.disabled = disabled;
    roomCodeInput.disabled = disabled;
}

function getElement(id: string): HTMLElement {
    const element = document.getElementById(id);

    if (!element) {
        throw new Error(`Не найден элемент #${id}`);
    }

    return element;
}

function getButton(id: string): HTMLButtonElement {
    const element = document.getElementById(id);

    if (!(element instanceof HTMLButtonElement)) {
        throw new Error(`Не найдена кнопка #${id}`);
    }

    return element;
}

function getInput(id: string): HTMLInputElement {
    const element = document.getElementById(id);

    if (!(element instanceof HTMLInputElement)) {
        throw new Error(`Не найден input #${id}`);
    }

    return element;
}

showScreen("menu");
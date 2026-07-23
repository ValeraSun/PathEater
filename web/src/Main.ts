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
const musicManager = game.GetMusicManager();

document.addEventListener(
    "click",
    () => {
        void musicManager.PlayMusic("menu");
    },
    { once: true }
);

let connected = false;
let isHost = false;
let roomId: string | null = null;
const playersCount = getElement("players-count");
const playerList = getElement("player-list");

const MAX_PLAYERS = 4;

gateway.SetRoomPlayersHandler(playerIds => {
    playersCount.textContent = `${playerIds.length}/${MAX_PLAYERS}`;
    renderPlayers(playerIds);
});

function renderPlayers(playerIds: string[]): void {
    playerList.replaceChildren();

    playerIds.forEach((playerId, index) => {
        const isLocal = playerId === gateway.localPlayerId;

        const item = document.createElement("li");
        item.className = "player";

        const avatar = document.createElement("div");
        avatar.className = "player-avatar";

        const image = document.createElement("img");
        image.src = "/images/avatar.png";
        image.alt = "Аватар игрока";
        avatar.appendChild(image);

        const info = document.createElement("div");
        info.className = "player-info";

        const name = document.createElement("strong");
        name.textContent = isLocal
            ? "Вы"
            : `Игрок ${index + 1}`;

        const ready = document.createElement("span");
        ready.className = "ready";
        ready.textContent = "Готов";

        info.append(name, ready);

        const role = document.createElement("span");
        role.className = "player-role";
        role.textContent =
            isLocal && isHost
                ? "Капитан"
                : "Игрок";

        item.append(avatar, info, role);
        playerList.appendChild(item);
    });

    for (
        let index = playerIds.length;
        index < MAX_PLAYERS;
        index++
    ) {
        const empty = document.createElement("li");
        empty.className = "empty-player";
        empty.textContent = "Ожидание игрока";

        playerList.appendChild(empty);
    }
}

function showScreen(name: ScreenName): void {
    for (const screen of Object.values(screens)) {
        screen.classList.remove("screen--active");
    }

    screens[name].classList.add("screen--active");

    document.body.classList.toggle(
        "background--blurred",
        name !== "menu"
    );
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
        menuStatus.textContent = "";
        playButton.disabled = false;
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
    const code = roomCodeInput.value.trim();

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
    void musicManager.PlayMusic("menu");
    showScreen("menu");
});

leaveLobbyButton.addEventListener("click", () => {
    gateway.ExitRoom();

    roomId = null;
    isHost = false;

    roomStatus.textContent = "";
    lobbyStatus.textContent = "";

    startGameButton.hidden = true;
    startGameButton.disabled = false;

    setRoomButtonsDisabled(false);
    void musicManager.PlayMusic("menu");

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
    await musicManager.PlayMusic("game");
    await game.StartMatch(payload);

    showScreen("game");
    game.MountTo(gameContainer);
};

roomCodeInput.addEventListener("input", () => {
    roomCodeInput.value = roomCodeInput.value;
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
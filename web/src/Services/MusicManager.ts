export type MusicTrack = "menu" | "game";

interface MusicConfig {
    src: string;
    volume: number;
}

const MUSIC_TRACKS: Record<MusicTrack, MusicConfig> = {
    menu: {
        src: "/music/menu.mp3",
        volume: 0.4
    },

    game: {
        src: "/music/game.mp3",
        volume: 0.3
    }
};

const FADE_STEP_MS = 40;
const FADE_STEP_VOLUME = 0.04;

export class MusicManager {
    private readonly music: HTMLAudioElement;

    private currentTrack: MusicTrack | null = null;
    private fadeTimer: number | null = null;

    public constructor() {
        this.music = new Audio();

        this.music.loop = true;
        this.music.preload = "auto";
    }

    public async PlayMusic(track: MusicTrack): Promise<void> {
        if (this.currentTrack === track) {
            if (this.music.paused) {
                await this.TryPlay();
            }

            return;
        }

        const config = MUSIC_TRACKS[track];

        this.StopFade();

        this.music.pause();
        this.music.currentTime = 0;
        this.music.src = config.src;
        this.music.volume = config.volume;
        this.currentTrack = track;

        await this.TryPlay();
    }

    public StopMusic(): void {
        this.StopFade();
        this.music.pause();
        this.music.currentTime = 0;
        this.currentTrack = null;
    }

    public PauseMusic(): void {
        this.music.pause();
    }

    public async ResumeMusic(): Promise<void> {
        if (!this.currentTrack) return;
        await this.TryPlay();
    }

    public FadeOut(onFinished?: () => void): void {
        this.StopFade();

        this.fadeTimer = window.setInterval(() => {
            const nextVolume = this.music.volume - FADE_STEP_VOLUME;

            if (nextVolume <= 0) {
                this.music.volume = 0;
                this.music.pause();
                this.StopFade();
                onFinished?.();
                return;
            }

            this.music.volume = nextVolume;
        }, FADE_STEP_MS);
    }

    public SetVolume(volume: number): void {
        this.music.volume = Math.min(1, Math.max(0, volume));
    }

    public SetMuted(muted: boolean): void {
        this.music.muted = muted;
    }

    private async TryPlay(): Promise<void> {
        try {
            await this.music.play();
        } catch (error) {
            console.warn("Браузер заблокировал автоматическое воспроизведение музыки:",error);
        }
    }

    private StopFade(): void {
        if (this.fadeTimer === null) {
            return;
        }
        window.clearInterval(this.fadeTimer);
        this.fadeTimer = null;
    }
}
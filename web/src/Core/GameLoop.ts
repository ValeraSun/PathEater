import { MAX_DELTA_TIME, MILLISECONDS_IN_SECOND } from "../Config/GameConfig";

export interface Updatable {
    Update(deltaTime: number): void;
}

export class GameLoop {
    private systems: Updatable[];
    private renderFrame: (deltaTime: number) => void;
    private finishFrame: (() => void) | null;

    private animationFrameId: number | null = null;
    private previousFrameTime = 0;
    private running = false;

    public constructor(systems: Updatable[], renderFrame: (deltaTime: number) => void, finishFrame?: () => void) {
        this.systems = systems;
        this.renderFrame = renderFrame;
        this.finishFrame = finishFrame ?? null;
    }

    public Start(): void {
        if (this.running) {
            return;
        }

        this.running = true;
        this.previousFrameTime = performance.now();
        this.animationFrameId = requestAnimationFrame(this.HandleAnimationFrame);
    }

    public Stop(): void {
        this.running = false;

        if (this.animationFrameId !== null) {
            cancelAnimationFrame(this.animationFrameId);
            this.animationFrameId = null;
        }
    }

    public IsRunning(): boolean {
        return this.running;
    }

    private readonly HandleAnimationFrame = (currentTime: number): void => {
        if (!this.running) {
            return;
        }

        const deltaTime = Math.min(
            (currentTime - this.previousFrameTime) / MILLISECONDS_IN_SECOND,
            MAX_DELTA_TIME
        );

        this.previousFrameTime = currentTime;

        for (const system of this.systems) {
            system.Update(deltaTime);
        }

        this.renderFrame(deltaTime);
        this.finishFrame?.();
        this.animationFrameId = requestAnimationFrame(this.HandleAnimationFrame);
    };
}
export class HealthbarView
{
    private readonly image: HTMLImageElement;

    public constructor()
    {
        this.image = document.createElement("img");
        this.image.style.position = "fixed";
        this.image.style.left = "10px";
        this.image.style.top = "40px";
        this.image.style.width = "550px";
        this.image.style.zIndex = "1000";
        this.image.style.display = "none";
        this.image.style.opacity = "0.7";

        document.body.appendChild(this.image);

        this.SetHealth(100);
    }

    public SetHealth(health: number): void
    {
        let state = 0;

        if (health >= 100)
            state = 100;
        else if (health >= 80)
            state = 80;
        else if (health >= 60)
            state = 60;
        else if (health >= 40)
            state = 40;
        else if (health >= 20)
            state = 20;

        this.image.src =
            `/images/health_${state}.png`;
    }

    public Show(): void
    {
        this.image.style.display = "block";
    }

    public Hide(): void
    {
        this.image.style.display = "none";
    }
}
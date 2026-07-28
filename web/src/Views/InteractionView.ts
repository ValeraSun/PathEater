export class InteractionView 
{
    private element: HTMLDivElement;

    public constructor() 
    {
        this.element = document.createElement("div");
        this.element.textContent = "Нажмите E, чтобы управлять кораблём";
        this.element.style.position = "fixed";
        this.element.style.left = "50%";
        this.element.style.bottom = "80px";
        this.element.style.transform = "translateX(-50%)";
        this.element.style.padding = "12px 20px";
        this.element.style.background = "rgba(0, 15, 25, 0.85)";
        this.element.style.border = "1px solid #00e5ff";
        this.element.style.color = "#00e5ff";
        this.element.style.fontFamily = "monospace";
        this.element.style.fontSize = "18px";
        this.element.style.zIndex = "1000";
        this.element.style.display = "none";

        document.body.appendChild(this.element);
    }

    public Show(text = "Нажмите E, чтобы управлять кораблём"): void 
    {
        this.element.textContent = text;
        this.element.style.display = "block";
    }

    public Hide(): void 
    {
        this.element.style.display = "none";
    }
}
export class InputController
{
    keys = new Set<string>();

    constructor()
    { 
        window.addEventListener("keydown", (event) => {
            if (!this.keys.has(event.code)) 
            {
                this.pressedOnce.add(event.code);
            }

            this.keys.add(event.code);
        });
        window.addEventListener("keyup",(e)=>{ this.keys.delete(e.code); }); 
    }

    public IsKeyDown(code: string): boolean 
    {
        return this.keys.has(code);
    }

    public WasPressedOnce(code: string): boolean 
    {
        if (!this.pressedOnce.has(code)) 
        {
            return false;
        }

        this.pressedOnce.delete(code);
        return true;
    }

    private pressedOnce = new Set<string>();
}
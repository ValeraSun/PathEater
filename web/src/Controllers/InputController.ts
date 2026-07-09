export class InputController
{
    keys = new Set<string>();

    constructor()
    { 
        window.addEventListener("keydown",(e)=>{ this.keys.add(e.code); }); 
        window.addEventListener("keyup",(e)=>{ this.keys.delete(e.code); }); 
    }

    public IsKeyDown(code: string): boolean 
    {
        return this.keys.has(code);
    }
}
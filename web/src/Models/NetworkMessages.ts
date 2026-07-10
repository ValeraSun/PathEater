export type Vector3D = { x: number; y: number; z: number; };

export type ClientMessage =
   {
      cmd: "createRoom";
      payload: {}
   }
   | 
   {
      cmd: "createGameSession";
      payload: {}
   }
   |   
   {
      cmd: "move";
      payload: {
         position: Vector3D;
         direction: Vector3D;
      }
   }
   |  
   {
      cmd: "interact";
      payload: {
          targetId: string;
       };
   }

export type ServerMessage =
    | {
          type: "entity_create";
          id: string;
          kind: "player" | "monster" | "door" | "cargo";
          position: Vector3D;
          rotationY?: number;
      }
    | {
          type: "entity_update";
          id: string;
          position?: Vector3D;
          rotationY?: number;
      }
    | {
          type: "entity_delete";
          id: string;
      }
    | {
          type: "snapshot";
          entities: Array<{
              id: string;
              kind: "player" | "monster" | "door" | "cargo";
              position: Vector3D;
              rotationY?: number;
          }>;
      };
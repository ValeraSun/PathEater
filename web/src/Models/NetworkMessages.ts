export type Vector3D = {
    x: number;
    y: number;
    z: number;
};

export type EntityKind =
    | "player"
    | "monster"
    | "door"
    | "cargo";

export type NetworkEntity = {
    id: string;
    kind: EntityKind;
    position: Vector3D;
    rotationY?: number;
};

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
      cmd: "playerState";
      payload: {
          move_front: boolean;
          move_left: boolean;
          move_right: boolean;
          move_back: boolean; 
          interact: boolean;
          attack: boolean;  
          direction: Vector3D;
       };
   }

export type ServerMessage =
    | {
          type: "entity_create";
          id: string;
          kind: EntityKind;
          position: Vector3D;
          rotation: Vector3D;
      }
    | {
          type: "entity_update";
          id: string;
          position?: Vector3D;
          rotation?: number;
      }
    | {
          type: "entity_delete";
          id: string;
      }
    | {
          type: "snapshot";
          entities: NetworkEntity[];
      }
    | {
          type: "computer_state";
          computerId: string;
          lockedBy: string | null;
      }
    | {
          type: "navigation_state";
          ship: {
              x: number;
              z: number;
              rotationY: number;
              hp: number;
          };
          asteroids: Array<{
              id: string;
              x: number;
              z: number;
          }>;
          monsters: Array<{
              id: string;
              x: number;
              z: number;
          }>;
      };
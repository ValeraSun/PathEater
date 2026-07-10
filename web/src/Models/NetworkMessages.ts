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
    | {
          type: "action";
          action: "move";
          position: Vector3D;
          rotationY: number;
      }
    | {
          type: "action";
          action: "interact";
          targetId: string;
      }
    | {
          type: "use_computer";
          computerId: string;
      }
    | {
          type: "release_computer";
          computerId: string;
      }
    | {
          type: "ship_input";
          computerId: string;

          inputX: number;
          inputZ: number;
      };

export type ServerMessage =
    | {
          type: "entity_create";
          id: string;
          kind: EntityKind;
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
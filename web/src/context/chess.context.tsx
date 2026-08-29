import {
  createContext,
  useContext,
  useEffect,
  useRef,
} from "react"


export interface ChallengeRequestPayload {
  to_user_id: string;
  from_user_id?: string;
}

export interface ChallengeAcceptPayload {
  from_user_id: string;
  to_user_id?: string;
}

export interface GameStartPayload {
  game_id: string;
  white_player: string;
  black_player: string;
  your_color: "White" | "Black";
}

export interface MakeMovePayload {
  game_id: string;
  move: string;
}

export interface MoveMadePayload {
  game_id: string;
  move: string;
  fen: string;
}

export interface GameOverPayload {
  game_id: string;
  result: "WhiteWins" | "BlackWins" | "NoOutCome" | string;
  reason: string;
}

export interface ErrorPayload {
  message: string;
}

// Map mapping Event Names -> Payload Types
export interface SocketEventMap {
  presence: string[];
  challenge_request: ChallengeRequestPayload;
  challenge_accept: ChallengeAcceptPayload;
  game_start: GameStartPayload;
  make_move: MakeMovePayload;
  move_made: MoveMadePayload;
  game_over: GameOverPayload;
  error: ErrorPayload;
}

export type EventType = keyof SocketEventMap;

export interface SocketEnvelope<K extends EventType = EventType> {
  type: K;
  payload: SocketEventMap[K];
}

export type InternalHandler = (payload: unknown) => void;


export type EventHandler<T extends EventType> = (payload: SocketEventMap[T]) => void

// --- Context Type Definition ---
export interface ChessSocketContextType {
  isConnected: boolean;
  onlineUsers: string[];
  send: <K extends EventType>(type: K, payload: SocketEventMap[K]) => void;
  sendChallenge: (toUserId: string) => void;
  acceptChallenge: (fromUserId: string) => void;
  sendMove: (gameId: string, move: string) => void;
  subscribe: <K extends EventType>(
    event: K,
    handler: EventHandler<K>
  ) => () => void;
}

export const ChessSocketContext = createContext<ChessSocketContextType | null>(
  null
)

// --- Custom Hooks ---

export const useChessSocket = (): ChessSocketContextType => {
  const context = useContext(ChessSocketContext);
  if (!context) {
    throw new Error("useChessSocket must be used within a ChessSocketProvider");
  }
  return context;
};

// 2. Specialized hook to subscribe to any socket event with automatic cleanup
// Automatic payload inference based on the passed EventType
export const useSocketEvent = <K extends EventType>(
  event: K,
  handler: EventHandler<K>
): void => {
  const { subscribe } = useChessSocket();
  const handlerRef = useRef<EventHandler<K>>(handler);

  useEffect(() => {
    handlerRef.current = handler;
  }, [handler]);

  useEffect(() => {
    const callback: EventHandler<K> = (data) => {
      if (handlerRef.current) {
        handlerRef.current(data);
      }
    };

    const unsubscribe = subscribe(event, callback);
    return () => {
      unsubscribe();
    };
  }, [event, subscribe]);
};

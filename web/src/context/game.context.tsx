import { createContext, useContext, useEffect, useRef } from "react"

export type PlayerStatus = "online" | "idle"

export interface UserPayload {
  id: string
  username?: string
  firstName?: string
  lastName?: string
}
export interface PresencePayload {
  presence_type: "add_user" | "remove_user"
  user_data: UserPayload
}

export interface PlayerStatusPayload {
  userId: string
  status: PlayerStatus
}

export interface ChallengePayload {
  from_user_data: UserPayload
  to_user_id?: string
}

export interface ChallengeAcceptPayload {
  from_user_id: string
  to_user_id?: string
}

export interface GameStartPayload {
  game_id: string
}


export interface MakeMovePayload {
  game_id: string
  move: string
}

export interface MoveMadePayload {
  id: string
  turn: "white" | "black"
  move: string
}

export interface GameOverPayload {
  game_id: string
  result: "WhiteWins" | "BlackWins" | "NoOutCome" | string
  reason: string
}

export interface ErrorPayload {
  message: string
}


// Map mapping Event Names -> Payload Types
export interface SocketEventMap {
  // Client events
  make_move: MakeMovePayload
  challenge_request: ChallengePayload
  challenge_accept: ChallengePayload

  // Server event
  move_made: MoveMadePayload
  challenge: ChallengePayload
  game_start: GameStartPayload
  presence: PresencePayload
  game_over: GameOverPayload
  error: ErrorPayload
}

export type EventType = keyof SocketEventMap

export interface SocketEnvelope<K extends EventType = EventType> {
  type: K
  payload: SocketEventMap[K]
}

export type InternalHandler = (payload: unknown) => void

export type EventHandler<T extends EventType> = (
  payload: SocketEventMap[T]
) => void

// --- Context Type Definition ---
export interface ChessSocketContextType {
  isConnected: boolean
  send: <K extends EventType>(type: K, payload: SocketEventMap[K]) => void
  sendChallenge: (toUserId: string) => void
  acceptChallenge: (fromUserId: string) => void
  sendMove: (gameId: string, move: string) => void
  subscribe: <K extends EventType>(
    event: K,
    handler: EventHandler<K>
  ) => () => void
}

export const GameSocketContext = createContext<ChessSocketContextType | null>(
  null
)

// --- Custom Hooks ---

export const useChessContext = (): ChessSocketContextType => {
  const context = useContext(GameSocketContext)
  if (!context) {
    throw new Error("useChessSocket must be used within a ChessSocketProvider")
  }
  return context
}

// 2. Specialized hook to subscribe to any socket event with automatic cleanup
// Automatic payload inference based on the passed EventType
export const useSocketEvent = <K extends EventType>(
  event: K,
  handler: EventHandler<K>
): void => {
  const { subscribe } = useChessContext()
  const handlerRef = useRef<EventHandler<K>>(handler)

  useEffect(() => {
    handlerRef.current = handler
  }, [handler])

  useEffect(() => {
    const callback: EventHandler<K> = (data) => {
      if (handlerRef.current) {
        handlerRef.current(data)
      }
    }

    const unsubscribe = subscribe(event, callback)
    return () => {
      unsubscribe()
    }
  }, [event, subscribe])
}

import useAuthContext from "@/context/auth.context"
import {
  GameSocketContext,
  type EventHandler,
  type EventType,
  type InternalHandler,
  type SocketEnvelope,
  type SocketEventMap,
} from "@/context/game.context"
import useUserContext from "@/context/user.context"
import { useCallback, useEffect, useRef, useState, type ReactNode } from "react"


const WS_ADDR = import.meta.env.VITE_SERVER_URL || ""

export interface ChessSocketProviderProps {
  children: ReactNode
}

export const GameSocketProvider = ({
  children,
}: ChessSocketProviderProps) => {


  const { client } = useAuthContext()
  const { user } = useUserContext()
  const userId = user.id
  const [isConnected, setIsConnected] = useState<boolean>(false)

  const wsRef = useRef<WebSocket | null>(null)

  // Registry to store event listeners: Map<EventType, Set<Handler>>
  const listenersRef = useRef<Map<EventType, Set<InternalHandler>>>(new Map())

  // Subscribe to an event, returns an unsubscribe cleanup function
  const subscribe = useCallback(
    <K extends EventType>(event: K, handler: EventHandler<K>) => {
      if (!listenersRef.current.has(event)) {
        listenersRef.current.set(event, new Set())
      }

      const handlerSet = listenersRef.current.get(event)
      const internalFn: InternalHandler = (payload: unknown) => {
        handler(payload as SocketEventMap[K])
      }

      handlerSet?.add(internalFn)

      return () => {
        handlerSet?.delete(internalFn)
      }
    },
    []
  )

  // Generic message sender
  const send = useCallback(
    <K extends EventType>(type: K, payload: SocketEventMap[K]) => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        const message: SocketEnvelope<K> = { type, payload }
        wsRef.current.send(JSON.stringify(message))
      } else {
        console.warn(`[WebSocket] Cannot send "${type}". Socket not connected.`)
      }
    },
    []
  )
  const sendChallenge = useCallback(
    (toUserId: string) => {
      send("challenge_request", { to_user_id: toUserId })
    },
    [send]
  )

  const acceptChallenge = useCallback(
    (fromUserId: string) => {
      send("challenge_accept", { from_user_id: fromUserId })
    },
    [send]
  )

  const sendMove = useCallback(
    (gameId: string, move: string) => {
      send("make_move", { game_id: gameId, move })
    },
    [send]
  )
  // WebSocket lifecycle management
  useEffect(() => {
    if (!userId) return

    let isMounted = true
    let ws: WebSocket | null = null

    const initConn = async () => {
      try {
        // 1. Fetch short-lived single-use ticket
        const res = await client.post<{ ticket: string }>("/api/v1/tickets")
        const ticket = res.data.ticket

        // If the component unmounted while awaiting the HTTP request, abort
        if (!isMounted) return

        // 2. Connect using the one-time ticket
        const wsUrl = `ws://${WS_ADDR}/ws?ticket=${encodeURIComponent(ticket)}`
        ws = new WebSocket(wsUrl)
        wsRef.current = ws

        ws.onopen = () => {
          if (!isMounted) return
          setIsConnected(true)
        }

        ws.onclose = () => {
          if (!isMounted) return
          setIsConnected(false)
        }

        ws.onerror = (err) => {
          console.error("[WebSocket Error]", err)
        }

        ws.onmessage = (messageEvent: MessageEvent) => {
          try {
            const envelope: SocketEnvelope = JSON.parse(messageEvent.data)
            const { type, payload } = envelope


            const handlers = listenersRef.current.get(type)
            if (handlers) {
              handlers.forEach((handler) => handler(payload))
            }
          } catch (err) {
            console.error(
              "[WebSocket] Failed to parse message:",
              messageEvent.data,
              err
            )
          }
        }
      } catch (err) {
        console.error("[WebSocket] Failed to acquire ticket or connect:", err)
      }
    }

    // 3. Trigger async connection
    initConn()

    // 4. Safe cleanup
    return () => {
      isMounted = false
      if (ws) {
        ws.close()
      } else if (wsRef.current) {
        wsRef.current.close()
      }
      wsRef.current = null
    }
  }, [userId, client])

  return (
    <GameSocketContext.Provider
      value={{
        isConnected,
        send,
        sendChallenge,
        acceptChallenge,
        sendMove,
        subscribe,
      }}
    >
      {children}
    </GameSocketContext.Provider>
  )
}

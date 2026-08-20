import { THEMES, type BoardTheme } from "@/types/board-theme.types"
import type { BoardResponse, MoveMessage, Position } from "@/types/game.types"
import { getPossibleMoves } from "@/utils/game/rules"
import { getPieceImageUrl } from "@/utils/game/pieces"
import { useEffect, useRef, useState } from "react"
import { useParams } from "react-router"

export default function GamePage() {
  const { gameId } = useParams()
  console.log(gameId)
  const [board, setBoard] = useState<string[][]>(
    Array(8).fill(Array(8).fill("."))
  )
  const [turn, setTurn] = useState<"WHITE" | "BLACK">("WHITE")
  const [selected, setSelected] = useState<Position | null>(null)
  const [theme, setTheme] = useState<BoardTheme>(THEMES.greenLichess)
  const [possibleMoves, setPossibleMoves] = useState<Position[]>([])

  // Explicitly type the WebSocket ref
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    // Initialize WebSocket connection
    const socket = new WebSocket("ws://localhost:8080/ws")
    ws.current = socket

    socket.onmessage = (event: MessageEvent) => {
      try {
        const data: BoardResponse = JSON.parse(event.data)
        setBoard(data.board)
        setTurn(data.turn)
        // Reset the selection
        setSelected(null)
        setPossibleMoves([])
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err)
      }
    }

    socket.onerror = (error) => {
      console.error("WebSocket Error:", error)
    }

    // Cleanup on unmount
    return () => {
      socket.close()
    }
  }, [])

  const handleSquareClick = (row: number, col: number): void => {
    // If clicking a square that is in possibleMoves, execute the move
    const isValidTarget = possibleMoves.some(
      (m) => m.row === row && m.col === col
    )

    if (selected && isValidTarget) {
      const movePayload: MoveMessage = {
        from: selected,
        to: { row, col },
      }

      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify(movePayload))
      }

      setSelected(null)
      setPossibleMoves([])
      return
    }

    // Otherwise, check if clicking a valid owned piece to select it
    const clickedPiece = board[row][col]
    if (clickedPiece !== ".") {
      const isWhitePiece = clickedPiece === clickedPiece.toUpperCase()
      if (
        (turn === "WHITE" && isWhitePiece) ||
        (turn === "BLACK" && !isWhitePiece)
      ) {
        const pos = { row, col }
        setSelected(pos)
        // Calculate possible moves locally
        const moves = getPossibleMoves(board, pos, turn)
        setPossibleMoves(moves)
        return
      }
    }

    // Deselect if clicking an empty/invalid square
    setSelected(null)
    setPossibleMoves([])
  }

  const resetBoard = () => {}

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        fontFamily: "sans-serif",
      }}
    >
      <h2>Golang Real-Time Chess</h2>

      {/* Theme Picker */}
      <div style={{ display: "flex", gap: "10px", alignItems: "center" }}>
        <label>Theme:</label>
        <select
          onChange={(e) => setTheme(THEMES[e.target.value])}
          style={{ padding: "6px 12px", borderRadius: "4px", fontSize: "14px" }}
        >
          {Object.keys(THEMES).map((key) => (
            <option key={key} value={key}>
              {THEMES[key].name}
            </option>
          ))}
        </select>
      </div>

      <h3>
        Current Turn:{" "}
        <span style={{ color: turn === "WHITE" ? "#4CAF50" : "#E91E63" }}>
          {turn}
        </span>
      </h3>

      {/* Grid Board */}

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(8, 64px)",
          gridTemplateRows: "repeat(8, 64px)",
          border: "4px solid #222",
          borderRadius: "4px",
          overflow: "hidden",
          boxShadow: "0 8px 24px rgba(0,0,0,0.2)",
        }}
      >
        {board.map((row, rIdx) =>
          row.map((cell, cIdx) => {
            const isDark = (rIdx + cIdx) % 2 === 1
            const isSelected = selected?.row === rIdx && selected?.col === cIdx
            const isPossibleMove = possibleMoves.some(
              (m) => m.row === rIdx && m.col === cIdx
            )
            const isCapture = isPossibleMove && cell !== "."
            const pieceUrl = getPieceImageUrl(cell)

            const cursorType =
              cell === "."
                ? "default"
                : (turn === "BLACK" && cell !== cell.toLowerCase()) ||
                    (turn === "WHITE" &&
                      cell !== "." &&
                      cell !== cell.toUpperCase())
                  ? "not-allowed"
                  : "pointer"

            return (
              <div
                key={`${rIdx}-${cIdx}`}
                onClick={() => handleSquareClick(rIdx, cIdx)}
                style={{
                  width: "64px",
                  height: "64px",
                  backgroundColor: isSelected
                    ? theme.selectedSquare
                    : isDark
                      ? theme.darkSquare
                      : theme.lightSquare,
                  display: "flex",
                  alignItems: "center",
                  cursor: cursorType,
                  justifyContent: "center",
                  userSelect: "none",
                  position: "relative",
                }}
              >
                {/* Piece Image */}
                {pieceUrl && (
                  <img
                    src={pieceUrl}
                    alt={cell}
                    style={{
                      width: "85%",
                      height: "85%",
                      pointerEvents: "none",
                      zIndex: 1,
                    }}
                  />
                )}

                {/* Move Indicator: Small Dot for empty square */}
                {isPossibleMove && !isCapture && (
                  <div
                    style={{
                      position: "absolute",
                      width: "20px",
                      height: "20px",
                      borderRadius: "50%",
                      backgroundColor: "rgba(0, 0, 0, 0.25)",
                      pointerEvents: "none",
                      zIndex: 2,
                    }}
                  />
                )}

                {/* Capture Indicator: Hollow Circle around target piece */}
                {isCapture && (
                  <div
                    style={{
                      position: "absolute",
                      width: "54px",
                      height: "54px",
                      borderRadius: "50%",
                      border: "5px solid rgba(0, 0, 0, 0.25)",
                      boxSizing: "border-box",
                      pointerEvents: "none",
                      zIndex: 2,
                    }}
                  />
                )}
              </div>
            )
          })
        )}
      </div>
      <div
        style={{
          padding: "8px 16px",
        }}
      >
        <button
          onClick={resetBoard}
          style={{
            backgroundColor: "red",
            outline: "none",
            color: "white",
            padding: "4px 8px",
            borderRadius: "4px",
            fontSize: "24px",
          }}
        >
          Reset
        </button>
      </div>
      <p>Click a piece, then click a target square to execute a move!</p>
    </div>
  )
}

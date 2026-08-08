import { useEffect, useState, useRef } from "react";
import {
  type Position,
  type MoveMessage,
  type BoardResponse,
} from "./types/chess";
import { THEMES, type BoardTheme } from "./types/themes";
import { getPieceImageUrl } from "./utils/pieces";

// // Mapping string pieces to unicode characters
// const PIECE_SYMBOLS: Record<string, string> = {
//   r: "♜",
//   n: "♞",
//   b: "♝",
//   q: "♛",
//   k: "♚",
//   p: "♟",
//   R: "♖",
//   N: "♘",
//   B: "♗",
//   Q: "♕",
//   K: "♔",
//   P: "♙",
//   ".": "",
// };

export default function App() {
  const [board, setBoard] = useState<string[][]>(
    Array(8).fill(Array(8).fill(".")),
  );
  const [turn, setTurn] = useState<"WHITE" | "BLACK">("WHITE");
  const [selected, setSelected] = useState<Position | null>(null);

  const [theme, setTheme] = useState<BoardTheme>(THEMES.greenLichess);

  // Explicitly type the WebSocket ref
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    // Initialize WebSocket connection
    const socket = new WebSocket("ws://localhost:8080/ws");
    ws.current = socket;

    socket.onmessage = (event: MessageEvent) => {
      try {
        const data: BoardResponse = JSON.parse(event.data);
        setBoard(data.board);
        setTurn(data.turn);
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
      }
    };

    socket.onerror = (error) => {
      console.error("WebSocket Error:", error);
    };

    // Cleanup on unmount
    return () => {
      socket.close();
    };
  }, []);

  const handleSquareClick = (row: number, col: number): void => {
    if (!selected) {
      // First click: Select piece if the square is not empty
      if (board[row][col] !== ".") {
        setSelected({ row, col });
      }
    } else {
      // Second click: Send move payload to Go backend
      const movePayload: MoveMessage = {
        from: selected,
        to: { row, col },
      };

      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify(movePayload));
      }

      setSelected(null); // Reset selection
    }
  };

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

      <div
        style={{
          display: "grid",
          gridTemplateColumns: "repeat(8, 60px)",
          gridTemplateRows: "repeat(8, 60px)",
          border: "3px solid #333",
        }}
      >
        {board.map((row, rIdx) =>
          row.map((cell, cIdx) => {
            const isDark = (rIdx + cIdx) % 2 === 1;
            const isSelected = selected?.row === rIdx && selected?.col === cIdx;

            const pieceUrl = getPieceImageUrl(cell);

            return (
              <div
                key={`${rIdx}-${cIdx}`}
                onClick={() => handleSquareClick(rIdx, cIdx)}
                style={{
                  width: "60px",
                  height: "60px",
                  backgroundColor: isSelected
                    ? theme.selectedSquare
                    : isDark
                      ? theme.darkSquare
                      : theme.lightSquare,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  fontSize: "36px",
                  cursor: "pointer",
                  userSelect: "none",
                }}
              >
                {pieceUrl && (
                  <img
                    src={pieceUrl}
                    alt={cell}
                    style={{
                      width: "85%",
                      height: "85%",
                      pointerEvents: "none",
                    }}
                  />
                )}
              </div>
            );
          }),
        )}
      </div>
      <p>Click a piece, then click a target square to execute a move!</p>
    </div>
  );
}

import { ChessBoard } from "@/components/game/chess-board"
import { GameOverDialog } from "@/components/game/game-dialogue"
import { type MoveRecord } from "@/components/game/sidebar"
import useAuthContext from "@/context/auth.context"
import {
    useChessContext,
    useSocketEvent,
    type GameOverPayload,
    type MoveMadePayload
} from "@/context/game.context"
import useUserContext from "@/context/user.context"
import type { User } from "@/types/user.types"
import React, { useCallback, useEffect, useState } from "react"
import { useParams } from "react-router"

const INITIAL_BOARD = [
  ["r", "n", "b", "q", "k", "b", "n", "r"],
  ["p", "p", "p", "p", "p", "p", "p", "p"],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  ["P", "P", "P", "P", "P", "P", "P", "P"],
  ["R", "N", "B", "Q", "K", "B", "N", "R"],
]

export const GamePage: React.FC = () => {
  const { gameId } = useParams<{ gameId: string }>()
  const { sendMove } = useChessContext()
  const {user}= useUserContext()

  const { client } = useAuthContext()

  const [board, setBoard] = useState<string[][]>(INITIAL_BOARD)
  const [currentTurn, setCurrentTurn] = useState<"white" | "black">(
    "white"
  )
  const [selectedSquare, setSelectedSquare] = useState<string | null>(null)
  const [availableMoves, setAvailableMoves] = useState<string[]>([])
  const [lastMove, setLastMove] = useState<{ from: string; to: string } | null>(
    null
  )
  const [moveHistory, setMoveHistory] = useState<MoveRecord[]>([])

  const [gameOverData, setGameOverData] = useState<GameOverPayload | null>(null)

  const [playerColor, setPlayerColor] = useState<"white"|"black">("white")

  // Parse the state of the board
  const parseBoardState = useCallback((rawBoard: string): string[][] => {
    const emptyBoard = () => Array.from({ length: 8 }, () => Array(8).fill("."))

    if (!rawBoard || typeof rawBoard !== "string") {
      return emptyBoard()
    }

    // 1. Split into row strings, trimming trailing whitespace or carriage returns
    const rawRows = rawBoard.trim().split(/\r?\n/)

    if (rawRows.length !== 8) {
      console.error(
        `Invalid board row count: expected 8, received ${rawRows.length}`
      )
      return emptyBoard()
    }

    const parsedBoard: string[][] = []

    for (let r = 0; r < 8; r++) {
      const rowStr = rawRows[r].trim()

      // Check if the row is space-separated ("r n b q ...") or continuous ("rnbq...")
      let cells: string[]
      if (rowStr.includes(" ")) {
        cells = rowStr.split(/\s+/)
      } else {
        cells = rowStr.split("")
      }

      if (cells.length !== 8) {
        console.error(
          `Invalid column count at row ${r}: expected 8, received ${cells.length}`
        )
        return emptyBoard()
      }

      parsedBoard.push(cells)
    }

    return parsedBoard
  }, [])

  useEffect(() => {
    const fetchGameData = async () => {
      const { data, status } = await client.get<{
        id: string
        white_player: User
        black_player: User
        turn: "white" | "black"
        board: string
      }>(`/api/v1/games/${gameId}`)

      if (status === 200) {
        const { board, turn, black_player, white_player } = data
        setCurrentTurn(turn)
        setBoard(parseBoardState(board))
        const currPlayerCol = black_player.id === user.id? "black":"white"
        setPlayerColor(currPlayerCol)
      }
    }
    fetchGameData()
  },[client,gameId,parseBoardState, user])

  // 1. Listen for move confirmations from server
  useSocketEvent("move_made", (data: MoveMadePayload) => {

    // const from = data.move.substring(0, 2)
    // const to = data.move.substring(2, 4)
    // setLastMove({ from, to })

    // Update move records
    // setMoveHistory((prev) => {
    //   const isWhiteMove = currentTurn === "White"
    //   if (isWhiteMove) {
    //     return [...prev, { turnNumber: prev.length + 1, white: data.move }]
    //   }
    //   const updated = [...prev]
    //   if (updated.length > 0) {
    //     updated[updated.length - 1].black = data.move
    //   }
    //   return updated
    // })

    // Toggle turn
    setCurrentTurn(data.turn)
    setSelectedSquare(null)
    setAvailableMoves([])
    const newBoard = parseBoardState(data.fen)
    setBoard(newBoard)
  })

  // 2. Listen for game over
  useSocketEvent("game_over", (data: GameOverPayload) => {
    if (data.game_id === gameId) {
      setGameOverData(data)
    }
  })

  // const applyMoveToLocalBoard = (from: string, to: string) => {
  //   const fromR = 8 - parseInt(from[1], 10);
  //   const fromC = from.charCodeAt(0) - 97;
  //   const toR = 8 - parseInt(to[1], 10);
  //   const toC = to.charCodeAt(0) - 97;

  //   setBoard((prev) => {
  //     const next = prev.map((row) => [...row]);
  //     const piece = next[fromR][fromC];
  //     next[toR][toC] = piece;
  //     next[fromR][fromC] = ".";
  //     return next;
  //   });
  // };

  // Parse FEN or execute local board update

  const handleSquareClick = (square: string) => {
    if (currentTurn !== playerColor || !gameId) return

    // Move to clicked destination
    if (selectedSquare && selectedSquare !== square) {
      const moveUci = `${selectedSquare}${square}`
      sendMove(gameId, moveUci)
      setSelectedSquare(null)
      setAvailableMoves([])
      return
    }

    // Select piece
    const r = 8 - parseInt(square[1], 10)
    const c = square.charCodeAt(0) - 97
    const piece = board[r][c]

    if (!piece || piece === ".") {
      setSelectedSquare(null)
      return
    }

    const isWhitePiece = piece === piece.toUpperCase()
    if (
      (playerColor === "white" && isWhitePiece) ||
      (playerColor === "black" && !isWhitePiece)
    ) {
      setSelectedSquare(square)
    }
  }

  const handleResign = () => {
    if (window.confirm("Are you sure you want to resign?")) {
      // Dispatch resign event
      //
    }
  }

  return (
    <div className="flex min-h-[calc(100vh-80px)] flex-col items-center justify-center gap-8 bg-background p-6 lg:flex-row">
      {/* 8x8 Chessboard */}
      <div className="flex w-full max-w-[620px] items-center justify-center">
        <ChessBoard
          board={board}
          playerColor={playerColor}
          selectedSquare={selectedSquare}
          availableMoves={availableMoves}
          lastMove={lastMove}
          isMyTurn={currentTurn === playerColor}
          onSquareClick={handleSquareClick}
        />
      </div>

      {/* Move History & Match Details */}
      <div className="h-[620px] w-full max-w-sm">
        {/*<GameSidebar
          gameId={gameId || ""}
          whitePlayer={whitePlayer}
          blackPlayer={blackPlayer}
          playerColor={playerColor}
          currentTurn={currentTurn}
          moves={moveHistory}
          onResign={handleResign}
        />*/}
      </div>

      {/* Game Over Popup */}
      {gameOverData && (
        <GameOverDialog
          open={!!gameOverData}
          result={gameOverData.result}
          reason={gameOverData.reason}
          playerColor={playerColor}
        />
      )}
    </div>
  )
}

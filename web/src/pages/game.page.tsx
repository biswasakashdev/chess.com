import React, { useState, useEffect } from "react";
import { useParams, useLocation } from "react-router";
import {
  useChessSocket,
  useSocketEvent,
  type MoveMadePayload,
  type GameOverPayload,
} from "@/context/game.context";
import { ChessBoard } from "@/components/game/chess-board";
import { GameSidebar, type MoveRecord } from "@/components/game/sidebar";
import { GameOverDialog } from "@/components/game/game-dialogue";

const INITIAL_BOARD = [
  ["r", "n", "b", "q", "k", "b", "n", "r"],
  ["p", "p", "p", "p", "p", "p", "p", "p"],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  [".", ".", ".", ".", ".", ".", ".", "."],
  ["P", "P", "P", "P", "P", "P", "P", "P"],
  ["R", "N", "B", "Q", "K", "B", "N", "R"],
];

export const GamePage: React.FC = () => {
  const { gameId } = useParams<{ gameId: string }>();
  const location = useLocation();
  const { sendMove } = useChessSocket();

  // Color passed via router state or default
  const playerColor: "White" | "Black" = location.state?.color || "White";
  const whitePlayer: string = location.state?.whitePlayer || "White Player";
  const blackPlayer: string = location.state?.blackPlayer || "Black Player";

  const [board, setBoard] = useState<string[][]>(INITIAL_BOARD);
  const [currentTurn, setCurrentTurn] = useState<"White" | "Black">("White");
  const [selectedSquare, setSelectedSquare] = useState<string | null>(null);
  const [availableMoves, setAvailableMoves] = useState<string[]>([]);
  const [lastMove, setLastMove] = useState<{ from: string; to: string } | null>(null);
  const [moveHistory, setMoveHistory] = useState<MoveRecord[]>([]);

  const [gameOverData, setGameOverData] = useState<GameOverPayload | null>(null);

  // 1. Listen for move confirmations from server
  useSocketEvent("move_made", (data: MoveMadePayload) => {
    if (data.game_id !== gameId) return;

    const from = data.move.substring(0, 2);
    const to = data.move.substring(2, 4);
    setLastMove({ from, to });

    // Parse FEN or execute local board update
    applyMoveToLocalBoard(from, to);

    // Update move records
    setMoveHistory((prev) => {
      const isWhiteMove = currentTurn === "White";
      if (isWhiteMove) {
        return [...prev, { turnNumber: prev.length + 1, white: data.move }];
      }
      const updated = [...prev];
      if (updated.length > 0) {
        updated[updated.length - 1].black = data.move;
      }
      return updated;
    });

    // Toggle turn
    setCurrentTurn((t) => (t === "White" ? "Black" : "White"));
    setSelectedSquare(null);
    setAvailableMoves([]);
  });

  // 2. Listen for game over
  useSocketEvent("game_over", (data: GameOverPayload) => {
    if (data.game_id === gameId) {
      setGameOverData(data);
    }
  });

  const applyMoveToLocalBoard = (from: string, to: string) => {
    const fromR = 8 - parseInt(from[1], 10);
    const fromC = from.charCodeAt(0) - 97;
    const toR = 8 - parseInt(to[1], 10);
    const toC = to.charCodeAt(0) - 97;

    setBoard((prev) => {
      const next = prev.map((row) => [...row]);
      const piece = next[fromR][fromC];
      next[toR][toC] = piece;
      next[fromR][fromC] = ".";
      return next;
    });
  };

  const handleSquareClick = (square: string) => {
    if (currentTurn !== playerColor || !gameId) return;

    // Move to clicked destination
    if (selectedSquare && selectedSquare !== square) {
      const moveUci = `${selectedSquare}${square}`;
      sendMove(gameId, moveUci);
      setSelectedSquare(null);
      setAvailableMoves([]);
      return;
    }

    // Select piece
    const r = 8 - parseInt(square[1], 10);
    const c = square.charCodeAt(0) - 97;
    const piece = board[r]?.[c];

    if (!piece || piece === ".") {
      setSelectedSquare(null);
      return;
    }

    const isWhitePiece = piece === piece.toUpperCase();
    if (
      (playerColor === "White" && isWhitePiece) ||
      (playerColor === "Black" && !isWhitePiece)
    ) {
      setSelectedSquare(square);
    }
  };

  const handleResign = () => {
    if (window.confirm("Are you sure you want to resign?")) {
      // Dispatch resign event
    }
  };

  return (
    <div className="flex flex-col lg:flex-row items-center justify-center gap-8 min-h-[calc(100vh-80px)] p-6 bg-background">
      {/* 8x8 Chessboard */}
      <div className="flex justify-center items-center w-full max-w-[620px]">
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
      <div className="w-full max-w-sm h-[620px]">
        <GameSidebar
          gameId={gameId || ""}
          whitePlayer={whitePlayer}
          blackPlayer={blackPlayer}
          playerColor={playerColor}
          currentTurn={currentTurn}
          moves={moveHistory}
          onResign={handleResign}
        />
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
  );
};

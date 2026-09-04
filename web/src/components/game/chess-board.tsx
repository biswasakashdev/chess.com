import React from "react";
import { getPieceImageUrl } from "@/lib/game/pieces";

interface ChessBoardProps {
  board: string[][]; // 8x8 array of characters
  playerColor: "white" | "black";
  selectedSquare: string | null;
  availableMoves: string[]; // e.g. ["e3", "e4"]
  lastMove: { from: string; to: string } | null;
  isMyTurn: boolean;
  onSquareClick: (square: string) => void;
}

export const ChessBoard: React.FC<ChessBoardProps> = ({
  board,
  playerColor,
  selectedSquare,
  availableMoves,
  lastMove,
  isMyTurn,
  onSquareClick,
}) => {
  const files = ["a", "b", "c", "d", "e", "f", "g", "h"];
  const ranks = ["8", "7", "6", "5", "4", "3", "2", "1"];

  // Orient files/ranks based on perspective
  const displayRanks = playerColor === "black" ? [...ranks].reverse() : ranks;
  const displayFiles = playerColor === "black" ? [...files].reverse() : files;

  return (
    <div className="relative aspect-square w-full max-w-[620px] select-none rounded-lg border-4 border-muted/80 shadow-2xl bg-[#769656] overflow-hidden">
      <div className="grid grid-cols-8 grid-rows-8 h-full w-full">
        {displayRanks.map((rank) =>
          displayFiles.map((file) => {
            const square = `${file}${rank}`;

            // Map rank & file back to original 0-indexed [r][c] board matrix
            const r = 8 - parseInt(rank, 10);
            const c = file.charCodeAt(0) - 97;
            const pieceCode = board[r][c] || ".";

            const isLightSquare = (r + c) % 2 === 0;
            const isSelected = selectedSquare === square;
            const isTarget = availableMoves.includes(square);
            const isLastMoveSquare =
              lastMove?.from === square || lastMove?.to === square;

            return (
              <button
                key={square}
                type="button"
                onClick={() => onSquareClick(square)}
                className={`relative flex items-center justify-center transition-colors focus:outline-none ${
                  isLightSquare ? "bg-[#eeeed2]" : "bg-[#769656]"
                } ${isLastMoveSquare ? "bg-amber-200/70 dark:bg-amber-500/40" : ""}`}
              >
                {/* Selected square highlight */}
                {isSelected && (
                  <div className="absolute inset-0 bg-yellow-400/50 ring-2 ring-yellow-400" />
                )}

                {/* Move Hint Indicator */}
                {isTarget && (
                  <div
                    className={`absolute z-20 rounded-full ${
                      pieceCode
                        ? "h-full w-full border-4 border-black/20"
                        : "h-4 w-4 bg-black/25"
                    }`}
                  />
                )}

                {/* Piece Image */}
                {pieceCode !== "." && (
                  <img
                    src={getPieceImageUrl(pieceCode)}
                    alt={pieceCode}
                    draggable={false}
                    className={`z-10 h-4/5 w-4/5 object-contain transition-transform duration-150 ${
                      isSelected ? "scale-110" : "hover:scale-105"
                    } ${!isMyTurn ? "cursor-default" : "cursor-pointer"}`}
                  />
                )}

                {/* Board Edge Coordinates */}
                {file === (playerColor === "black" ? "h" : "a") && (
                  <span
                    className={`absolute top-0.5 left-1 text-[10px] font-bold ${
                      isLightSquare ? "text-[#769656]" : "text-[#eeeed2]"
                    }`}
                  >
                    {rank}
                  </span>
                )}
                {rank === (playerColor === "black" ? "8" : "1") && (
                  <span
                    className={`absolute bottom-0.5 right-1 text-[10px] font-bold ${
                      isLightSquare ? "text-[#769656]" : "text-[#eeeed2]"
                    }`}
                  >
                    {file}
                  </span>
                )}
              </button>
            );
          })
        )}
      </div>
    </div>
  );
};

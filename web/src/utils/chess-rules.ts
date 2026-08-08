import { type Position } from "../types/chess";

// Helper to check if coordinates are inside the 8x8 grid
function inBounds(row: number, col: number): boolean {
  return row >= 0 && row < 8 && col >= 0 && col < 8;
}

// Helper to determine piece ownership
function isWhite(piece: string): boolean {
  return piece !== "." && piece === piece.toUpperCase();
}

function isBlack(piece: string): boolean {
  return piece !== "." && piece === piece.toLowerCase();
}

function isSameTeam(piece1: string, piece2: string): boolean {
  if (piece1 === "." || piece2 === ".") return false;
  return (
    (isWhite(piece1) && isWhite(piece2)) || (isBlack(piece1) && isBlack(piece2))
  );
}

export function getPossibleMoves(
  board: string[][],
  selected: Position,
  currentTurn: "WHITE" | "BLACK",
): Position[] {
  const { row, col } = selected;
  const piece = board[row]?.[col];

  if (!piece || piece === ".") return [];

  // Enforce turn check
  const pieceIsWhite = isWhite(piece);
  if (
    (currentTurn === "WHITE" && !pieceIsWhite) ||
    (currentTurn === "BLACK" && pieceIsWhite)
  ) {
    return [];
  }

  const moves: Position[] = [];
  const pLower = piece.toLowerCase();

  // Helper for sliding pieces (Rook, Bishop, Queen)
  const addRayMoves = (directions: [number, number][]) => {
    for (const [dr, dc] of directions) {
      let r = row + dr;
      let c = col + dc;
      while (inBounds(r, c)) {
        const target = board[r][c];
        if (target === ".") {
          moves.push({ row: r, col: c });
        } else {
          if (!isSameTeam(piece, target)) {
            moves.push({ row: r, col: c }); // Capture enemy piece
          }
          break; // Blocked by piece
        }
        r += dr;
        c += dc;
      }
    }
  };

  switch (pLower) {
    case "p": {
      // Pawns
      const direction = pieceIsWhite ? -1 : 1; // White moves up (-1), Black moves down (+1)
      const startRow = pieceIsWhite ? 6 : 1;

      // 1-step forward
      const fRow = row + direction;
      if (inBounds(fRow, col) && board[fRow][col] === ".") {
        moves.push({ row: fRow, col });

        // 2-steps forward from starting position
        const f2Row = row + 2 * direction;
        if (row === startRow && board[f2Row][col] === ".") {
          moves.push({ row: f2Row, col });
        }
      }

      // Diagonal Captures
      for (const dc of [-1, 1]) {
        const cCol = col + dc;
        if (inBounds(fRow, cCol)) {
          const target = board[fRow][cCol];
          if (target !== "." && !isSameTeam(piece, target)) {
            moves.push({ row: fRow, col: cCol });
          }
        }
      }
      break;
    }

    case "n": {
      // Knight (L-shapes)
      const knightOffsets = [
        [-2, -1],
        [-2, 1],
        [-1, -2],
        [-1, 2],
        [1, -2],
        [1, 2],
        [2, -1],
        [2, 1],
      ];
      for (const [dr, dc] of knightOffsets) {
        const r = row + dr;
        const c = col + dc;
        if (inBounds(r, c) && !isSameTeam(piece, board[r][c])) {
          moves.push({ row: r, col: c });
        }
      }
      break;
    }

    case "b":
      // Bishop (Diagonals)
      addRayMoves([
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
      ]);
      break;

    case "r":
      // Rook (Orthogonals)
      addRayMoves([
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
      ]);
      break;

    case "q":
      // Queen (Diagonals + Orthogonals)
      addRayMoves([
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
      ]);
      break;

    case "k": {
      // King (1-step in any direction)
      const kingOffsets = [
        [-1, -1],
        [-1, 0],
        [-1, 1],
        [0, -1],
        [0, 1],
        [1, -1],
        [1, 0],
        [1, 1],
      ];
      for (const [dr, dc] of kingOffsets) {
        const r = row + dr;
        const c = col + dc;
        if (inBounds(r, c) && !isSameTeam(piece, board[r][c])) {
          moves.push({ row: r, col: c });
        }
      }
      break;
    }
  }

  return moves;
}

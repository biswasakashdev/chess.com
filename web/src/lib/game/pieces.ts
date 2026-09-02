// Mapping piece notation to standard SVG filenames
import pawnLight from "@/assets/pieces/Chess_plt45.svg"
import pawnDark from "@/assets/pieces/Chess_pdt45.svg"
import rookDark from "@/assets/pieces/Chess_rdt45.svg"
import rookLight from "@/assets/pieces/Chess_rlt45.svg"
import kingDark from "@/assets/pieces/Chess_kdt45.svg"
import kingLight from "@/assets/pieces/Chess_klt45.svg"
import nightDark from "@/assets/pieces/Chess_ndt45.svg"
import nightLight from "@/assets/pieces/Chess_nlt45.svg"
import bishopLight from "@/assets/pieces/Chess_blt45.svg"
import bishopDark from "@/assets/pieces/Chess_bdt45.svg"
import queenDark from "@/assets/pieces/Chess_qdt45.svg"
import queenLight from "@/assets/pieces/Chess_qlt45.svg"

export function getPieceImageUrl(piece: string): string | undefined {
  if (piece === ".") return undefined

  const isWhite = piece === piece.toUpperCase()
  const colorPrefix = isWhite ? "w" : "b"

  // Wikipedia standard piece set hosted via Wikimedia
  const pieceMap: Record<string, string> = {
    pw: pawnLight,
    rw: rookLight,
    nw: nightLight,
    bw: bishopLight,
    qw: queenLight,
    kw: kingLight,
    pb: pawnDark,
    rb: rookDark,
    nb: nightDark,
    bb: bishopDark,
    qb: queenDark,
    kb: kingDark,
  }

  const pieceName = pieceMap[`${piece.toLowerCase()}${colorPrefix}`]

  return pieceName
}

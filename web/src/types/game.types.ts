export interface Position {
  row: number
  col: number
}

export interface MoveMessage {
  from: Position
  to: Position
}

export interface BoardResponse {
  board: string[][]
  turn: "WHITE" | "BLACK"
}

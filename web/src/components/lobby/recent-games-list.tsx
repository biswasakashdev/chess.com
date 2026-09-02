import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"
import React from "react"

interface RecentGame {
  id: string
  opponent: string
  result: "win" | "loss" | "draw"
  ratingChange: string
  mode: string
  moves: number
  date: string
}

const mockGames: RecentGame[] = [
  {
    id: "1",
    opponent: "LevyRoz",
    result: "win",
    ratingChange: "+8",
    mode: "10 min",
    moves: 34,
    date: "2 hrs ago",
  },
  {
    id: "2",
    opponent: "AnnaCramling",
    result: "loss",
    ratingChange: "-9",
    mode: "3 min",
    moves: 42,
    date: "Yesterday",
  },
  {
    id: "3",
    opponent: "DanyaN",
    result: "draw",
    ratingChange: "+1",
    mode: "5 min",
    moves: 68,
    date: "2 days ago",
  },
]

export const RecentGamesTable: React.FC = () => {


  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg font-semibold">Match History</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Opponent</TableHead>
              <TableHead>Outcome</TableHead>
              <TableHead>Mode</TableHead>
              <TableHead>Moves</TableHead>
              <TableHead>Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {mockGames.map((game) => (
              <TableRow key={game.id}>
                <TableCell className="font-medium">{game.opponent}</TableCell>
                <TableCell>
                  <Badge
                    variant={
                      game.result === "win"
                        ? "default"
                        : game.result === "loss"
                          ? "destructive"
                          : "secondary"
                    }
                  >
                    {game.result.toUpperCase()} ({game.ratingChange})
                  </Badge>
                </TableCell>
                <TableCell>{game.mode}</TableCell>
                <TableCell>{game.moves}</TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {game.date}
                </TableCell>

              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  )
}

import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { PlayCircle, Clock } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

interface ActiveGame {
  id: string
  opponent: string
  opponentRating: number
  timeControl: string
  isMyTurn: boolean
  timeLeft: string
}

const activeMock: ActiveGame[] = [
  {
    id: "g1",
    opponent: "MagnusC",
    opponentRating: 2850,
    timeControl: "10+0 Rapid",
    isMyTurn: true,
    timeLeft: "06:42",
  },
  {
    id: "g2",
    opponent: "HikaruN",
    opponentRating: 2820,
    timeControl: "3+2 Blitz",
    isMyTurn: false,
    timeLeft: "01:18",
  },
]

export const ActiveGamesList: React.FC = () => {
  return (
    <Card className="h-full">
      <CardHeader className="flex flex-row items-center justify-between pb-3">
        <CardTitle className="flex items-center gap-2 text-lg font-semibold">
          Active Games
          <Badge variant="secondary">{activeMock.length}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {activeMock.length === 0 ? (
          <p className="py-4 text-center text-sm text-muted-foreground">
            No games in progress.
          </p>
        ) : (
          activeMock.map((game) => (
            <div
              key={game.id}
              className="flex items-center justify-between rounded-lg border bg-card p-3 transition-colors hover:bg-accent/50"
            >
              <div className="flex items-center space-x-3">
                <Avatar className="h-9 w-9">
                  <AvatarImage
                    src={`https://avatar.vercel.sh/${game.opponent}`}
                  />
                  <AvatarFallback>{game.opponent.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="text-sm font-semibold">
                      {game.opponent}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      ({game.opponentRating})
                    </span>
                  </div>
                  <div className="flex items-center space-x-2 text-xs text-muted-foreground">
                    <span>{game.timeControl}</span>
                    <span>•</span>
                    <span className="flex items-center">
                      <Clock className="mr-1 h-3 w-3" />
                      {game.timeLeft}
                    </span>
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-2">
                {game.isMyTurn ? (
                  <Badge className="animate-pulse border-emerald-500/20 bg-emerald-500/15 text-emerald-600 dark:text-emerald-400">
                    Your Turn
                  </Badge>
                ) : (
                  <Badge variant="outline" className="text-muted-foreground">
                    Waiting
                  </Badge>
                )}
                <Button size="sm" variant="default">
                  <PlayCircle className="mr-1 h-4 w-4" /> Resume
                </Button>
              </div>
            </div>
          ))
        )}
      </CardContent>
    </Card>
  )
}

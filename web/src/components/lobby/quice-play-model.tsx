import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Flame, Plus, Timer, Zap, type LucideProps } from "lucide-react"
import React, { useState } from "react"

type GameType = "blt" | "bltz" | "rpd" | "cls"

const timeControls: {
  name: string
  type: GameType
  time: string
  icon: React.ForwardRefExoticComponent<
    Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>
  >
  desc: string
}[] = [
  {
    name: "Classical",
    type: "cls",
    time: "30 min",
    icon: Timer,
    desc: "Tournament prep",
  },

  {
    name: "Rapid",
    type: "rpd",
    time: "10 min",
    icon: Timer,
    desc: "Thoughtful play",
  },

  {
    name: "Bullet",
    type: "blt",
    time: "1 min",
    icon: Flame,
    desc: "Fast & chaotic",
  },
  {
    name: "Blitz",
    type: "bltz",
    time: "3 min",
    icon: Zap,
    desc: "Standard blitz",
  },
]

export const QuickPlayModal: React.FC = () => {
  const [gameType, setGameType] = useState<GameType>("cls")
  const createGame = () => {}

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="lg" className="w-full shadow-md sm:w-auto">
          <Plus className="mr-2 h-5 w-5" /> Create New Game
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-120">
        <DialogHeader>
          <DialogTitle className="text-xl">Choose Time Control</DialogTitle>
        </DialogHeader>

        <div className="grid grid-cols-2 gap-3 py-4">
          {timeControls.map((tc) => (
            <Button
              key={tc.name}
              variant={gameType === tc.type ? "outline" : "ghost"}

              onClick={() => setGameType(tc.type)}
              className="flex h-20 cursor-pointer flex-col items-center justify-center space-y-1"
            >
              <div className="flex items-center space-x-1 font-semibold">
                <tc.icon className="h-4 w-4 text-primary" />
                <span>{tc.name}</span>
              </div>
              <span className="text-xs text-muted-foreground">{tc.time}</span>
            </Button>
          ))}
        </div>

        <Button className="mt-4 w-full" size="lg" onClick={createGame}>
          Start Matchmaking
        </Button>
      </DialogContent>
    </Dialog>
  )
}

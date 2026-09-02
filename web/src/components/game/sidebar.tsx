import React from "react";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { Flag, Swords } from "lucide-react";

export interface MoveRecord {
  turnNumber: number;
  white?: string;
  black?: string;
}

interface GameSidebarProps {
  gameId: string;
  whitePlayer: string;
  blackPlayer: string;
  playerColor: "White" | "Black";
  currentTurn: "White" | "Black";
  moves: MoveRecord[];
  onResign: () => void;
}

export const GameSidebar: React.FC<GameSidebarProps> = ({
  whitePlayer,
  blackPlayer,
  playerColor,
  currentTurn,
  moves,
  onResign,
}) => {
  return (
    <Card className="flex flex-col h-full w-full max-w-sm border-border shadow-md">
      {/* Player Header */}
      <CardHeader className="pb-3 space-y-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-semibold flex items-center gap-2">
            <Swords className="h-4 w-4 text-primary" /> Live Match
          </CardTitle>
          <Badge
            variant={currentTurn === playerColor ? "default" : "outline"}
            className={currentTurn === playerColor ? "bg-primary text-primary-foreground" : ""}
          >
            {currentTurn === playerColor ? "Your Turn" : "Opponent's Turn"}
          </Badge>
        </div>

        {/* Players Card List */}
        <div className="space-y-2 pt-1">
          <div className="flex items-center justify-between rounded-md p-2 bg-muted/40">
            <div className="flex items-center gap-2">
              <Avatar className="h-7 w-7">
                <AvatarFallback className="text-xs bg-black text-white font-bold">B</AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">{blackPlayer}</span>
            </div>
            {currentTurn === "Black" && <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />}
          </div>

          <div className="flex items-center justify-between rounded-md p-2 bg-muted/40">
            <div className="flex items-center gap-2">
              <Avatar className="h-7 w-7">
                <AvatarFallback className="text-xs bg-white text-black border border-border font-bold">W</AvatarFallback>
              </Avatar>
              <span className="text-sm font-medium">{whitePlayer}</span>
            </div>
            {currentTurn === "White" && <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />}
          </div>
        </div>
      </CardHeader>

      <Separator />

      {/* Move History Table */}
      <CardContent className="flex-1 p-0">
        <div className="grid grid-cols-3 px-4 py-2 text-xs font-semibold text-muted-foreground border-b border-border bg-muted/20">
          <span>#</span>
          <span>White</span>
          <span>Black</span>
        </div>
        <ScrollArea className="h-80 px-4">
          {moves.length === 0 ? (
            <div className="flex items-center justify-center h-48 text-xs text-muted-foreground">
              No moves played yet
            </div>
          ) : (
            <div className="divide-y divide-border/40 py-1">
              {moves.map((m) => (
                <div key={m.turnNumber} className="grid grid-cols-3 py-1.5 text-xs font-mono">
                  <span className="text-muted-foreground">{m.turnNumber}.</span>
                  <span className="font-semibold text-foreground">{m.white || ""}</span>
                  <span className="font-semibold text-foreground">{m.black || ""}</span>
                </div>
              ))}
            </div>
          )}
        </ScrollArea>
      </CardContent>

      <Separator />

      {/* Footer Controls */}
      <CardFooter className="p-3">
        <Button variant="outline" size="sm" onClick={onResign} className="w-full text-destructive hover:bg-destructive/10">
          <Flag className="h-4 w-4 mr-2" /> Resign Game
        </Button>
      </CardFooter>
    </Card>
  );
};

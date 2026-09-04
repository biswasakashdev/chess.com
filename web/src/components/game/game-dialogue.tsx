import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Trophy, Home } from "lucide-react";
import { useNavigate } from "react-router";

interface GameOverDialogProps {
  open: boolean;
  result: string;
  reason: string;
  playerColor: "white" | "black";
}

export const GameOverDialog: React.FC<GameOverDialogProps> = ({
  open,
  result,
  reason,
  playerColor,
}) => {
  const navigate = useNavigate();

  const isWin =
    (result === "WhiteWins" && playerColor === "white") ||
    (result === "BlackWins" && playerColor === "black");
  const isDraw = result === "Draw" || result === "NoOutCome";

  return (
    <Dialog open={open}>
      <DialogContent className="sm:max-w-md text-center">
        <DialogHeader className="items-center space-y-3">
          <div className={`p-3 rounded-full ${isWin ? "bg-amber-100 text-amber-600" : "bg-muted text-muted-foreground"}`}>
            <Trophy className="h-8 w-8" />
          </div>
          <DialogTitle className="text-xl font-bold">
            {isDraw ? "Game Drawn" : isWin ? "Victory!" : "Defeat"}
          </DialogTitle>
          <DialogDescription className="text-sm">
            {result} • {reason}
          </DialogDescription>
        </DialogHeader>

        <DialogFooter className="sm:justify-center mt-4">
          <Button onClick={() => navigate("/home")} className="gap-2">
            <Home className="h-4 w-4" /> Back to Lobby
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

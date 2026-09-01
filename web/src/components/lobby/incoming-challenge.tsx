import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  type ChallengeRequestPayload,
  useChessSocket,
  useSocketEvent,
} from "@/context/game.context";
import useUserContext from "@/context/user.context";
import { Check, Clock, Swords, X } from "lucide-react";
import { useEffect, useState } from "react";

const CHALLENGE_TIMEOUT_SEC = 10;

// --- 1. Bottom-Right Floating Incoming Challenge Toast ---

export const IncomingChallengeToast = () => {

  const { user } = useUserContext()
  const currentUserId = user.id
  const { acceptChallenge } = useChessSocket();
  const [activeChallenge, setActiveChallenge] = useState<{
    fromUserId: string;
    expiresAt: number;
  } | null>(null);
  const [timeLeft, setTimeLeft] = useState<number>(CHALLENGE_TIMEOUT_SEC);

  // Listen for incoming challenges
  useSocketEvent("challenge_request", (payload: ChallengeRequestPayload) => {
    if (!payload.from_user_id || payload.from_user_id === currentUserId) return;

    // Trigger toast for 30 seconds
    setActiveChallenge({
      fromUserId: payload.from_user_id,
      expiresAt: Date.now() + CHALLENGE_TIMEOUT_SEC * 1000,
    });
    setTimeLeft(CHALLENGE_TIMEOUT_SEC);
  });

  // Auto-dismiss & Countdown Timer Loop
  useEffect(() => {
    if (!activeChallenge) return;

    const timer = setInterval(() => {
      const remainingMs = activeChallenge.expiresAt - Date.now();
      const remainingSec = Math.max(0, Math.ceil(remainingMs / 1000));

      setTimeLeft(remainingSec);

      if (remainingSec <= 0) {
        setActiveChallenge(null);
      }
    }, 200);

    return () => clearInterval(timer);
  }, [activeChallenge]);

  if (!activeChallenge) return null;

  const progressPercent = (timeLeft / CHALLENGE_TIMEOUT_SEC) * 100;

  const handleAccept = () => {
    acceptChallenge(activeChallenge.fromUserId);
    setActiveChallenge(null);
  };

  const handleDecline = () => {
    setActiveChallenge(null);
  };

  return (
    <div className="fixed bottom-5 right-5 z-50 w-80 animate-in fade-in slide-in-from-bottom-5 duration-300">
      <Card className="shadow-2xl border-primary/20 bg-background/95 backdrop-blur-md overflow-hidden">
        <Progress value={progressPercent} className="h-1 rounded-none bg-muted" />

        <CardContent className="p-4 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Avatar className="h-10 w-10 border border-border">
                <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">
                  {activeChallenge.fromUserId.substring(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div>
                <p className="text-xs text-muted-foreground font-medium flex items-center gap-1">
                  <Swords className="h-3 w-3 text-primary" /> Game Challenge
                </p>
                <p className="text-sm font-semibold leading-tight">
                  {activeChallenge.fromUserId}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-1 text-xs font-mono text-muted-foreground bg-muted px-2 py-1 rounded-md">
              <Clock className="h-3 w-3" />
              <span>{timeLeft}s</span>
            </div>
          </div>

          <div className="flex items-center gap-2 pt-1">
            <Button
              size="sm"
              variant="outline"
              onClick={handleDecline}
              className="w-1/2 h-8 text-xs"
            >
              <X className="h-3.5 w-3.5 mr-1 text-muted-foreground" />
              Decline
            </Button>
            <Button
              size="sm"
              onClick={handleAccept}
              className="w-1/2 h-8 text-xs bg-primary hover:bg-primary/90"
            >
              <Check className="h-3.5 w-3.5 mr-1" />
              Accept
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

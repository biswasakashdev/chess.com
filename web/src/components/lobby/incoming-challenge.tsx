import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useChessContext, useSocketEvent, type ChallengePayload } from "@/context/game.context";
import useUserContext from "@/context/user.context";
import { getCapitalise, getInitialsFullName } from "@/utils/user-utils";
import { Check, Clock, Swords, X } from "lucide-react";
import React, { useCallback, useEffect, useState } from "react";

interface ChallengeItem {
  fromUserId: string;
  fromUsername: string;
  fromUserFullname: string;
  expiresAt: number;
}

const CHALLENGE_TIMEOUT_SEC = 15;

// Individual toast item with its own countdown loop
const ChallengeToastItem = React.memo(
  ({
    challenge,
    onAccept,
    onDecline,
    onExpire,
  }: {
    challenge: ChallengeItem;
    onAccept: (fromUserId: string) => void;
    onDecline: (challengeId: string) => void;
    onExpire: (challengeId: string) => void;
  }) => {
    const [timeLeft, setTimeLeft] = useState<number>(() =>
      Math.max(0, Math.ceil((challenge.expiresAt - Date.now()) / 1000))
    );

    useEffect(() => {
      const interval = setInterval(() => {
        const remainingMs = challenge.expiresAt - Date.now();
        const remainingSec = Math.max(0, Math.ceil(remainingMs / 1000));

        setTimeLeft(remainingSec);

        if (remainingSec <= 0) {
          clearInterval(interval);
          onExpire(challenge.fromUserId);
        }
      }, 200);

      return () => clearInterval(interval);
    }, [challenge.expiresAt, challenge.fromUserId, onExpire]);

    const progressPercent = (timeLeft / CHALLENGE_TIMEOUT_SEC) * 100;

    return (
      <Card className="shadow-lg border-primary/20 bg-background/95 backdrop-blur-md overflow-hidden animate-in fade-in slide-in-from-bottom-3 duration-200">
        <Progress
          value={progressPercent}
          className="h-1 rounded-none bg-muted transition-all duration-200"
        />

        <CardContent className="p-3 space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <Avatar className="h-9 w-9 border border-border">
                <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">
                  {getInitialsFullName(challenge.fromUserFullname)}
                </AvatarFallback>
              </Avatar>
              <div className="max-w-[140px] truncate">
                <p className="text-[11px] text-muted-foreground font-medium flex items-center gap-1">
                  <Swords className="h-3 w-3 text-primary" /> Accept
                </p>
                <p className="text-sm font-semibold truncate leading-tight">
                  {challenge.fromUserId}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-1 text-[11px] font-mono text-muted-foreground bg-muted px-2 py-0.5 rounded-md">
              <Clock className="h-3 w-3" />
              <span>{timeLeft}s</span>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => onDecline(challenge.fromUserId)}
              className="w-1/2 h-7 text-xs"
            >
              <X className="h-3.5 w-3.5 mr-1 text-muted-foreground" />

            </Button>
            <Button
              size="sm"
              onClick={() => onAccept(challenge.fromUserId)}
              className="w-1/2 h-7 text-xs bg-primary hover:bg-primary/90"
            >
              <Check className="h-3.5 w-3.5 mr-1" />
              Accept
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }
);

ChallengeToastItem.displayName = "ChallengeToastItem";

export const IncomingChallengeToast = () => {
  const { user } = useUserContext();
  const { acceptChallenge } = useChessContext();
  const [challenges, setChallenges] = useState<ChallengeItem[]>([]);

  // Socket listener for new incoming challenges
  useSocketEvent("challenge", (payload: ChallengePayload) => {

    const {  username, id, firstName,lastName} = payload.from_user_data
    if (!id || id === user.id) return;


    if (username && id && firstName && lastName) {
      setChallenges((prev) => {
        // Prevent duplicates from the same sender if one is already active
        const exists = prev.some(
          (c) => c.fromUserId === id
        );
        if (exists) return prev;

        return [
          ...prev,
          {
            fromUserId: id,
            fromUserFullname: getCapitalise(firstName, lastName),
            fromUsername: username,
            expiresAt: Date.now() + CHALLENGE_TIMEOUT_SEC * 1000,
          },
        ];
      });
    }
  });

  const removeChallenge = useCallback((userId: string) => {
    setChallenges((prev) => prev.filter((c) => c.fromUserId !== userId));
  }, []);

  const handleAccept = useCallback(
    (fromUserId: string) => {
      acceptChallenge(fromUserId);
      removeChallenge(fromUserId);
    },
    [acceptChallenge, removeChallenge]
  );

  const handleDecline = useCallback(
    (challengeId: string) => {
      removeChallenge(challengeId);
    },
    [removeChallenge]
  );

  if (challenges.length === 0) return null;

  return (
    <aside
      aria-label="Incoming Challenges"
      className="fixed bottom-5 right-5 z-50 w-84 pointer-events-none"
    >
      <ScrollArea className="max-h-[calc(100vh-6rem)] pr-2 pointer-events-auto">
        <div className="flex flex-col gap-2.5 pb-1">
          {challenges.map((challenge) => (
            <ChallengeToastItem
              challenge={challenge}
              onAccept={handleAccept}
              onDecline={handleDecline}
              onExpire={removeChallenge}
            />
          ))}
        </div>
      </ScrollArea>
    </aside>
  );
};

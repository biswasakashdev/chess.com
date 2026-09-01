import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Circle, Swords, UserPlus } from "lucide-react"

export interface UserRowProps {
  firstName: string
  lastName: string
  username: string
  id: string
  status?: "online" | "in-game" | "idle"
  rating?: number
  challenge?: boolean
}

export const UserRow = ({
  firstName,
  lastName,
  username,
  rating,
  status,
  challenge = false,
}: UserRowProps) => {
  const fullName = firstName+" "+ lastName
  const sendChallengeHandler = () => {
    // Send challenge
  }

  const sendRequestHandler = () => {
    // Send request
  }
  return (
    <div className="flex items-center justify-between rounded-md p-2 transition-colors hover:bg-accent">
      <div className="flex items-center space-x-3">
        <div className="relative">
          <Avatar className="h-8 w-8">
            <AvatarFallback>
              <span>{firstName[0] + lastName[0]}</span>
            </AvatarFallback>
          </Avatar>
          {status && (
            <Circle
              className={`absolute right-0 bottom-0 h-2.5 w-2.5 fill-current ${
                status === "online"
                  ? "text-emerald-500"
                  : status === "in-game"
                    ? "text-amber-500"
                    : "text-muted-foreground"
              }`}
            />
          )}
        </div>

        <div>
          <div className="text-xs font-semibold">{fullName}</div>
          <div className="text-[11px] text-muted-foreground">
            {username} • {rating} • {status && <span> {status}</span>}
          </div>
        </div>
      </div>
      <div className="flex space-x-2">
        {challenge ? (
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7"
            title="Challenge"
            onClick={sendChallengeHandler}
          >
            <Swords className="h-3.5 w-3.5" />
          </Button>
        ) : (
          <Button size="sm" variant="outline" onClick={sendRequestHandler}>
            <UserPlus className="mr-1 h-3.5 w-3.5" /> Add
          </Button>
        )}
      </div>
    </div>
  )
}

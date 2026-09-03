import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { useChessContext } from "@/context/game.context"
import { getCapitalise, getInitials } from "@/utils/user-utils"
import { Circle, Swords, UserPlus } from "lucide-react"

export interface UserRowProps {
  firstName: string
  lastName: string
  username: string
  id: string
  rating?: number
  challenge?: boolean
  isActive?: boolean
  sendRequestHandler?: () => void
}

export const UserRow = ({
  id,
  firstName,
  lastName,
  username,
  rating,
  isActive,
  challenge = false,
  sendRequestHandler,
}: UserRowProps) => {
  const fullName = getCapitalise(firstName, lastName)
  const {sendChallenge}= useChessContext()
  return (
    <div className="flex items-center justify-between rounded-md p-2 transition-colors hover:bg-accent">
      <div className="flex items-center space-x-3">
        <div className="relative">
          <Avatar className="h-8 w-8">
            <AvatarFallback>
              <span>{getInitials(firstName, lastName)}</span>
            </AvatarFallback>
          </Avatar>

          <Circle
            className={`absolute right-0 bottom-0 h-2.5 w-2.5 fill-current ${
              isActive ? "text-emerald-500" : "text-muted-foreground"
            }`}
          />
        </div>

        <div>
          <div className="text-xs font-semibold">{fullName}</div>
          <div className="text-[11px] text-muted-foreground">
            @{username} {rating && <span>• {rating}</span>}{" "}
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
            onClick={() => sendChallenge(id)}
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

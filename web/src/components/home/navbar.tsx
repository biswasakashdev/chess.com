import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { User as UserType } from "@/types/user.types"
import { LogOutIcon, User } from "lucide-react"
import { Avatar, AvatarFallback } from "../ui/avatar"
import { Button } from "../ui/button"
import { getCapitalise, getInitials } from "@/utils/user-utils"

export const Navbar = ({
  user,
  logoutHandler,
}: {
  user?: UserType
  logoutHandler?: () => void
}) => {
  return (
    <header className="sticky top-0 z-20 flex h-20 items-center justify-between border-b bg-background/95 px-6 py-4 backdrop-blur">
      <img src="/logo.png" alt="go-chess" className="h-full object-fill" />
      <div className="flex w-fit items-center space-x-3">
        {user && (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <div className="w-fit space-x-3">

              <Button variant="ghost" size="icon" className="rounded-full">
                <Avatar>
                  <AvatarFallback>{getInitials(user.firstName,user.lastName)}</AvatarFallback>
                </Avatar>
              </Button>
              <span className="cursor-pointer">@{user.username}</span>
              </div>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-fit">
              <DropdownMenuGroup>
                <DropdownMenuItem className="px-2">
                  <User />
                  {getCapitalise(user.firstName, user.lastName)}
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={logoutHandler}>
                <LogOutIcon />
                Sign Out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )}
      </div>
    </header>
  )
}

import { QuickPlayModal } from "@/components/lobby/quice-play-model"
import useUserContext from "@/context/user.context"
import { Avatar, AvatarFallback } from "../ui/avatar"
import { Item, ItemContent, ItemMedia, ItemTitle } from "../ui/item"

export const Navbar = () => {
  const { user } = useUserContext()
  const initials = user
    ? `${user.firstName[0]}${user.lastName[0]}`.toUpperCase()
    : undefined

  return (
    <header className="sticky top-0 z-20 flex h-20 items-center justify-between border-b bg-background/95 px-6 py-4 backdrop-blur">
      <img src="/logo.png" alt="go-chess" className="h-full object-fill" />
      <div className="flex w-fit items-center space-x-3">
        <Item className="flex flex-row" size={"sm"}>
          <ItemMedia variant={"icon"}>
            <Avatar>
              <AvatarFallback>{initials}</AvatarFallback>
            </Avatar>
          </ItemMedia>
          <ItemContent>
            <ItemTitle>{user.username}</ItemTitle>
          </ItemContent>
        </Item>
        <QuickPlayModal />
      </div>
    </header>
  )
}

import useUserContext from "@/context/user.context"
import { AvatarFallback, Avatar } from "../ui/avatar"
import useAuthContext from "@/context/auth.context"

export const Navbar = () => {
  const { user } = useUserContext()
  const { authorization } = useAuthContext()
  const initials = authorization
    ? `${user.firstName[0] + user.lastName[0]}`.toUpperCase()
    : undefined
  return (
    <header className="sticky top-0 z-50 flex w-full justify-center border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container flex h-16 items-center justify-between px-4 sm:px-8">
        <img src="/logo.png" alt="logo" className="h-full" />
        <div>
          {authorization && (
            <Avatar>
              <AvatarFallback>{initials}</AvatarFallback>
            </Avatar>
          )}
        </div>
      </div>
    </header>
  )
}

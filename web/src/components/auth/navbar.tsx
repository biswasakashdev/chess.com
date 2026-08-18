import useAuthContext from "@/context/auth.context"
import { Avatar, AvatarFallback } from "../ui/avatar"
import { useEffect, useState } from "react"
import type { User } from "@/types/user.types"
import { useNavigate } from "react-router"

export const Navbar = () => {
  const { authorization, client } = useAuthContext()
  const navigate = useNavigate()

  const [user, setUser] = useState<User | undefined>(undefined)

  useEffect(() => {
    const fetchUser = async () => {
      if (!authorization) return

      const res = await client.get("/api/v1/users", {
        headers: {
          Authorization: `Bearer ${authorization}`,
        },
      })

      if (res.status !== 200) {
        console.error("Failed to fetch the user.")
      }
      setUser(res.data)
    }
    fetchUser()
  }, [authorization, navigate])

  const initials = user
    ? `${user.firstName[0]}${user.lastName[0]}`.toUpperCase()
    : undefined

  return (
    <header className="sticky top-0 z-50 flex w-full justify-center border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container flex h-16 items-center justify-between px-4 sm:px-8">
        <img src="/logo.png" alt="logo" className="h-full" />
        <div>
          {user && (
            <Avatar>
              <AvatarFallback>{initials}</AvatarFallback>
            </Avatar>
          )}
        </div>
      </div>
    </header>
  )
}

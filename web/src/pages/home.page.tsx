import { Navbar } from "@/components/home/navbar"
import { Toaster } from "@/components/ui/sonner"
import { GameSocketProvider } from "@/context-provider/game.provider"
import useAuthContext from "@/context/auth.context"
import { UserContext } from "@/context/user.context"
import type { User } from "@/types/user.types"
import { useEffect, useState } from "react"
import { Outlet, useNavigate } from "react-router"



export default function HomePage() {
  const { authorization, client, updateAuthorization } = useAuthContext()
  const navigate = useNavigate()

  const [user, setUser] = useState<User | undefined>(undefined)

  useEffect(() => {
    const fetchUser = async () => {
      if (!authorization) navigate("/auth")

      const res = await client.get("/api/v1/users")

      if (res.status !== 200) {
        updateAuthorization(undefined)
        navigate("/auth")
        return
      }

      setUser(res.data)
    }
    fetchUser()
  }, [authorization, navigate, client, updateAuthorization])

  if (!user) {
    return
  }

  return (
    <UserContext.Provider
      value={{
        user: user,
      }}
    >
      <GameSocketProvider>

      <div className="flex min-h-screen flex-col bg-background text-foreground">
        <Navbar />
        <Outlet />
      </div>
      <Toaster/>
      </GameSocketProvider>
    </UserContext.Provider>
  )
}

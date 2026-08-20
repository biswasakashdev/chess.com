import { Navbar } from "@/components/home/navbar"
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

      if (res.status === 401) {
        updateAuthorization(undefined)
        navigate("/auth")
        return
      }

      if (res.status !== 200) {
        console.error("Failed to fetch the user.")
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
      <div className="flex min-h-screen flex-col bg-background text-foreground">
        <Navbar />
        <Outlet />
      </div>
    </UserContext.Provider>
  )
}

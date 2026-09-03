
import useAuthContext from "@/context/auth.context"
import { UserContext } from "@/context/user.context"
import type { User } from "@/types/user.types"
import { useEffect, useState } from "react"
import { useNavigate } from "react-router"



export default function UserProvider({children}:{children:React.ReactNode}) {
  const { authorization, client, updateAuthorization } = useAuthContext()
  const navigate = useNavigate()

  const [user, setUser] = useState<User | undefined>(undefined)

  useEffect(() => {
    const fetchUser = async () => {
      if (!authorization) navigate("/auth")

      const res = await client.get("/api/v1/users/profile")

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
      {children}
    </UserContext.Provider>
  )
}

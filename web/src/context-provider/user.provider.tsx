import useAuthContext from "@/context/auth.context"
import { UserContext } from "@/context/user.context"
import type { User } from "@/types/user.types"
import { axiosInstance } from "@/utils/axios.config"
import { useEffect, useState } from "react"
import { redirect } from "react-router"

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { authorization } = useAuthContext()

  const [user, setUser] = useState<User | undefined>(undefined)

  if (!authorization) {
    throw redirect("/auth")
  }

  const instance = axiosInstance.create({
    headers: {
      Authorization: `Bearer ${authorization}`,
    },
  })

  useEffect(() => {
    const fetchUser = async () => {
      const res = await instance.get("/api/v1/users")

      if (res.status !== 200) {
        console.error("Failed to fetch the user.")
      }

      setUser(res.data)
    }
    fetchUser()
  }, [instance])

  if (!user) {
    return
  }

  return (
    <UserContext.Provider value={{ client: instance, user: user }}>
      {children}
    </UserContext.Provider>
  )
}

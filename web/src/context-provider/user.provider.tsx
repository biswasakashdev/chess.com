import useAuthContext from "@/context/auth.context"
import { UserContext } from "@/context/user.context"
import axios from "axios"

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { authorization } = useAuthContext()

  const instance = axios.create({
    headers: {
      Authorization: `Bearer ${authorization}`,
    },
  })

  return (
    <UserContext.Provider value={{ client: instance }}>
      {children}
    </UserContext.Provider>
  )
}

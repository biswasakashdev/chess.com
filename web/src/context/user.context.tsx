import type { User } from "@/types/user.types"
import { createContext, useContext } from "react"

export const UserContext = createContext<{
  user: User
}>({
  user: {
    firstName: "",
    lastName: "",
    username: "",
    id: "",
  },
})

export default function useUserContext() {
  return useContext(UserContext)
}

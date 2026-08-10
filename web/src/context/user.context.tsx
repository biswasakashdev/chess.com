import type { User } from "@/types/user.types"
import type { AxiosInstance } from "axios"
import axios from "axios"
import { createContext, useContext } from "react"

export const UserContext = createContext<{
  client: AxiosInstance
  user: User
}>({
  client: axios,
  user: {
    firstName: "",
    lastName: "",
    email: "",
    avatar: "",
  },
})

export default function useUserContext() {
  return useContext(UserContext)
}

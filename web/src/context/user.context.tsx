import type { AxiosInstance } from "axios"
import axios from "axios"
import { createContext, useContext } from "react"

export const UserContext = createContext<{
  client: AxiosInstance
}>({
  client: axios,
})

export default function useUserContext() {
  return useContext(UserContext)
}

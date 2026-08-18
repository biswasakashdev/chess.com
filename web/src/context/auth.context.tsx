import type { AxiosInstance } from "axios"
import axios from "axios"
import { createContext, useContext } from "react"

export const AuthContext = createContext<{
  authorization: string | undefined
  updateAuthorization: (auth: string | undefined) => void
  clearAuthorization: () => void
  client: AxiosInstance
}>({
  authorization: undefined,
  updateAuthorization: () => {},
  clearAuthorization: () => {},
  client: axios,
})
const useAuthContext = () => {
  return useContext(AuthContext)
}

export default useAuthContext

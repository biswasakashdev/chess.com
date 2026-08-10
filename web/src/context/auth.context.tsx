import type { Authorization } from "@/types/user.types"
import { createContext, useContext } from "react"

export const AuthContext = createContext<{
  authorization: Authorization | undefined
  updateAuthorization: (auth: Authorization | undefined) => void
}>({
  authorization: undefined,
  updateAuthorization: () => {},
})
const useAuthContext = () => {
  return useContext(AuthContext)
}

export default useAuthContext

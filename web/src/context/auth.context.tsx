import { createContext, useContext } from "react"

export const AuthContext = createContext<{
  authorization: string | undefined
  updateAuthorization: (auth: string | undefined) => void
  clearAuthorization: () => void
}>({
  authorization: undefined,
  updateAuthorization: () => {},
  clearAuthorization: () => {},
})
const useAuthContext = () => {
  return useContext(AuthContext)
}

export default useAuthContext

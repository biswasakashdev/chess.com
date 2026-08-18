import { AuthContext } from "@/context/auth.context"
import { useState } from "react"

const AUTH_TOKEN = "authtoken"

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const localToken = localStorage.getItem(AUTH_TOKEN) || undefined
  const [authorization, setAuthorization] = useState<string | undefined>(
    localToken
  )

  const updateAuthorization = (auth: string | undefined) => {
    setAuthorization(auth)
  }

  const clearAuthorization = () => {
    setAuthorization(undefined)
    localStorage.removeItem(AUTH_TOKEN)
  }

  return (
    <AuthContext.Provider
      value={{
        authorization: authorization,
        updateAuthorization: updateAuthorization,
        clearAuthorization: clearAuthorization,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

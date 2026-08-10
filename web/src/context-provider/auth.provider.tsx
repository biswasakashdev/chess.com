import { fetchAuthorization } from "@/api/auth.api"
import { AuthContext } from "@/context/auth.context"
import type { Authorization } from "@/types/user.types"
import { useEffect, useState } from "react"

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [authorization, setAuthorization] = useState<Authorization | undefined>(
    undefined
  )

  const updateAuthorization = (auth: Authorization | undefined) => {
    setAuthorization(auth)
  }

  useEffect(() => {
    const getAuthorization = async () => {
      const auth = await fetchAuthorization()
      if (auth) {
        setAuthorization({ token: auth })
      }
    }
    getAuthorization()
  }, [])

  return (
    <AuthContext.Provider
      value={{
        authorization: authorization,
        updateAuthorization: updateAuthorization,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

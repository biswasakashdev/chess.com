import { AuthContext } from "@/context/auth.context"
import axios from "axios"
import { useState } from "react"

const AUTH_TOKEN = "authtoken"

const url = import.meta.env.DEV ? "/backend" : ""
const axiosInstance = axios.create({
  baseURL: url,
})

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const localToken = localStorage.getItem(AUTH_TOKEN) || undefined
  const [authorization, setAuthorization] = useState<string | undefined>(
    localToken
  )

  const updateAuthorization = (auth: string | undefined) => {
    setAuthorization(auth)
    if (auth) localStorage.setItem(AUTH_TOKEN, auth)
  }

  const clearAuthorization = () => {
    setAuthorization(undefined)
    localStorage.removeItem(AUTH_TOKEN)
  }

  const instance = axiosInstance.create({
    baseURL: url,
    headers: {
      Authorization: authorization ? `Bearer ${authorization}` : undefined,
    },
  })

  return (
    <AuthContext.Provider
      value={{
        authorization: authorization,
        updateAuthorization: updateAuthorization,
        clearAuthorization: clearAuthorization,
        client: instance,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

import { fetchAuthorization } from "@/api/auth.api"
import useAuthContext from "@/context/auth.context"
import { UserContext } from "@/context/user.context"
import type { User } from "@/types/user.types"
import { axiosInstance } from "@/utils/axios.config"
import type { AxiosError } from "axios"
import { useEffect, useState } from "react"
import { redirect } from "react-router"

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { authorization, updateAuthorization } = useAuthContext()

  const [user, setUser] = useState<User | undefined>(undefined)

  if (!authorization) {
    throw redirect("/auth")
  }

  const instance = axiosInstance.create({
    headers: {
      Authorization: `Bearer ${authorization.token}`,
    },
  })

  // Retry if token expired.
  instance.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const config = error.config

      // Only retry on network errors or 5xx server errors
      const shouldRetry =
        !error.response ||
        (error.response.status >= 500 && error.response.status < 600)

      if (!shouldRetry || !config) {
        return Promise.reject(error)
      }

      const auth = await fetchAuthorization()
      if (!auth) {
        return Promise.reject(error)
      }

      config.headers.Authorization = `Bearer ${auth}`
      updateAuthorization({ token: auth })
      return instance(config)
    }
  )

  useEffect(() => {
    const fetchUser = async () => {
      const res = await instance.get("/api/v1/users")

      if (res.status !== 200) {
        console.error("Failed to fetch the user.")
      }

      setUser(res.data)
    }
    fetchUser()
  }, [instance])

  if (!user) {
    return
  }

  return (
    <UserContext.Provider value={{ client: instance, user: user }}>
      {children}
    </UserContext.Provider>
  )
}

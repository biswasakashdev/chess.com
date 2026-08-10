import { axiosInstance } from "@/utils/axios.config"

export async function fetchAuthorization(): Promise<string | undefined> {
  const res = await axiosInstance.get("/api/v1/auth")

  if (res.status !== 200) {
    return undefined
  }
  return res.data.token
}

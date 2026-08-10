import axios from "axios"
console.log(import.meta.env)
export const axiosInstance = axios.create({
  baseURL: "/backend",
})

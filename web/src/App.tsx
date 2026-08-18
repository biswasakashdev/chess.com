import { Outlet } from "react-router"
import { AuthProvider } from "./context-provider/auth.provider"
import { Navbar } from "./components/auth/navbar"

export function App() {
  return (
    <AuthProvider>
      {/* 1. Header */}
      <Navbar />
      <Outlet />
    </AuthProvider>
  )
}

export default App

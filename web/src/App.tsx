import { Outlet } from "react-router"
import { Toaster } from "./components/ui/sonner"
import { AuthProvider } from "./context-provider/auth.provider"

export function App() {
  return (
    <AuthProvider>
      <Outlet />
      <Toaster />
    </AuthProvider>
  )
}

export default App

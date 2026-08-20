import { Outlet } from "react-router"
import { AuthProvider } from "./context-provider/auth.provider"

export function App() {
  return (
    <AuthProvider>
      <Outlet />
    </AuthProvider>
  )
}

export default App

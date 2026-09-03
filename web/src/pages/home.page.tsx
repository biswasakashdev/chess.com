import { GameSocketProvider } from "@/context-provider/game.provider"
import UserProvider from "@/context-provider/user.provider"
import { Outlet } from "react-router"

export default function HomePage() {
  return (
    <UserProvider>
      <GameSocketProvider>
        <Outlet />
      </GameSocketProvider>
    </UserProvider>
  )
}

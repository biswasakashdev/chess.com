import AuthPage from "@/pages/auth.page"
import HomePage from "@/pages/home.page"
import { createBrowserRouter, Navigate } from "react-router"
import App from "./App"
import LobbyPage from "./pages/lobby.page"
import { GamePage } from "./pages/game.page"

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      {
        path: "home",
        Component: HomePage,
        children: [
          {
            index: true,
            Component: LobbyPage,
          },
          {
            path: ":gameId",
            Component: GamePage,
          },
        ],
      },
      {
        path: "/auth",
        Component: AuthPage,
      },
      {
        path: "/",
        element: <Navigate replace={true} to="/home" />,
      },
    ],
  },
])

export default router

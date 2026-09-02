import App from "@/App"
import AuthPage from "@/pages/auth.page"
import {GamePage} from "@/pages/game.page"
import { createBrowserRouter, Navigate } from "react-router"
import LobbyPage from "./pages/lobby.page"
import HomePage from "@/pages/home.page"

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      {
        // The routes under home are secured and not accessed without authorization.
        path: "home",
        Component: HomePage,
        children: [
          {
            index: true,
            Component: LobbyPage,
          },
          {
            path: ":gameId/game",
            Component: GamePage,
          },
        ],
      },
      {
        path: "auth",
        Component: AuthPage,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate replace={true} to="/home" />,
  },
])

export default router

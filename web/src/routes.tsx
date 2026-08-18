import { createBrowserRouter, Navigate } from "react-router"
import App from "./App"
import AuthPage from "./pages/auth.page"
import HomePage from "./pages/home.page"

export const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      {
        index: true,
        Component: HomePage,
      },
      {
        path: "/auth",
        Component: AuthPage,
      },
    ],
  },
  {
    path: "*",
    element: <Navigate replace={true} to="/" />,
  },
])

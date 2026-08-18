import { Outlet } from "react-router"

export default function HomePage() {
  return (
    <div>
      <div>Welcome to go chess.com</div>
      <Outlet />
    </div>
  )
}

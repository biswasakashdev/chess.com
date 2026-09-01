import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import useAuthContext from "@/context/auth.context"
import type { PlayerStatus } from "@/context/game.context"
import { ArrowUpRight } from "lucide-react"
import { useEffect, useState } from "react"
import { UserRow } from "../user-row"

interface ActiveUser {
  id: string
  username: string
  firstName: string
  lastName: string
  rating: number
  status: PlayerStatus
}

export function ActiveUsersSection() {
  const { client } = useAuthContext()
  const [activeUsers, setActiveUsers] = useState<ActiveUser[]>([])

  useEffect(() => {
    const fetchOnlineUsers = async () => {
      const res = await client.get("/api/v1/users", {
        params: {
          status: "online",
        },
      })
      if (res.status === 200) {
        setActiveUsers(res.data)
      }
    }
    fetchOnlineUsers()
  }, [client])


  return (
    <Card className="w-full max-w-2xl border-border/60 shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <CardTitle className="text-xl font-semibold tracking-tight">
              Active Friends
            </CardTitle>
            <Badge
              variant="secondary"
              className="px-2 py-0.5 text-xs font-medium"
            >
              {activeUsers.length} Online
            </Badge>
          </div>
        </div>
        <Button variant="outline" size="sm" className="hidden gap-1 sm:flex">
          View All
          <ArrowUpRight className="h-4 w-4" />
        </Button>
      </CardHeader>

      <CardContent className="grid gap-3">
        {activeUsers.map((row) => (
          <UserRow
            id={row.id}
            firstName={row.firstName}
            lastName={row.lastName}
            username={row.username}
            rating={row.rating}
            key={row.id}
          />
        ))}
      </CardContent>
    </Card>
  )
}

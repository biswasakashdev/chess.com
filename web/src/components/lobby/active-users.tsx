import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import useAuthContext from "@/context/auth.context"
import {
  useSocketEvent,
  type PresencePayload
} from "@/context/game.context"
import { ArrowUpRight } from "lucide-react"
import { useEffect, useState } from "react"
import { UserRow } from "../user-row"

interface ActiveUser {
  id: string
  username: string
  firstName: string
  lastName: string
  rating: number
}

export function ActiveUsersSection() {
  const { client } = useAuthContext()
  const [activeUsers, setActiveUsers] = useState<ActiveUser[]>([])

  useEffect(() => {
    const fetchOnlineUsers = async () => {
      const res = await client.get("/api/v1/friends", {
        params: {
          type: "online",
        },
      })
      if (res.status === 200) {
        setActiveUsers(res.data)
      }
    }
    fetchOnlineUsers()
  }, [client])

  useSocketEvent("presence", (payLoad: PresencePayload) => {
    console.log("Event revieved...")
    if (payLoad.presence_type === "add_user" && payLoad.user_data) {
      const userData = payLoad.user_data
      setActiveUsers((pre) => {
        return [
          ...pre,
          {
            username: userData.username,
            firstName: userData.username,
            lastName: userData.username,
            id: userData.id,
            rating: userData.rating,
          },
        ]
      })
    } else if (
      payLoad.presence_type === "remove_user" &&
      payLoad.remove_user_id
    ) {
      setActiveUsers((pre) => {
        return [...pre.filter((us) => us.id !== payLoad.remove_user_id)]
      })
    }
  })

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
            isActive
            challenge
            key={row.id}
          />
        ))}
      </CardContent>
    </Card>
  )
}

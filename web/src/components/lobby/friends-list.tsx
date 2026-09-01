import { UserRow } from "@/components/user-row"
import useAuthContext from "@/context/auth.context"
import {
  useSocketEvent,
  type PlayerStatus,
  type PlayerStatusPayload
} from "@/context/game.context"
import { useEffect, useState } from "react"

interface Friend {
  id: string
  username: string
  rating: number
  firstName: string
  lastName: string
  status: PlayerStatus
}

export default function FriendsList() {
  const [activeFriends, setActiveFriends] = useState<Friend[]>([])
  const { client } = useAuthContext()

  useEffect(() => {
    const fetchData = async () => {
      const res = await client.get("/api/v1/users/friends")
      if (res.status === 200) {
        setActiveFriends(res.data)
      }
      setActiveFriends([])
    }

    fetchData()
  }, [client])


  useSocketEvent("player_status", (payload:PlayerStatusPayload) => {
    setActiveFriends((pre) => {
      return pre.map((frnd) => {
        if (frnd.id === payload.userId) {
          return { ...frnd, status: payload.status }
        }
        return { ...frnd }
      })
    })
  })

  return (
    <>
      {activeFriends.map((row) => (
        <UserRow
          id={row.id}
          firstName={row.firstName}
          lastName={row.lastName}
          username={row.username}
          rating={row.rating}
          challenge
          status={row.status}
          key={row.id}
        />
      ))}
    </>
  )
}

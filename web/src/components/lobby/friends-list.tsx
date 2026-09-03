import { UserRow } from "@/components/user-row"
import useAuthContext from "@/context/auth.context"
import { useSocketEvent, type PresencePayload } from "@/context/game.context"
import { useEffect, useState } from "react"

interface Friend {
  id: string
  username: string
  rating: number
  firstName: string
  lastName: string
  isActive: boolean
}

export default function FriendsList() {
  const [friendsList, setFriendsList] = useState<Friend[]>([])
  const { client } = useAuthContext()

  useEffect(() => {
    const fetchData = async () => {
      const res = await client.get("/api/v1/friends")
      if (res.status === 200) {
        setFriendsList(res.data)
      }
    }

    fetchData()
  }, [client])

  useSocketEvent("presence", (payLoad: PresencePayload) => {
    if (payLoad.presence_type === "add_user") {
      setFriendsList((pre) => {
        return [
          ...pre.map((itm) => {
            if (itm.id === payLoad.user_data.id) {
              return { ...itm, isActive: true }
            }
            return { ...itm }
          }),
        ]
      })
    } else if (
      payLoad.presence_type === "remove_user"
    ) {
      setFriendsList((pre) => {
        return [
          ...pre.map((us) => {
            if (us.id === payLoad.user_data.id) {
              return { ...us, isActive: false }
            }
            return { ...us }
          }),
        ]
      })
    }
  })

  return (
    <>
      {friendsList.map((row) => (
        <UserRow
          id={row.id}
          firstName={row.firstName}
          lastName={row.lastName}
          username={row.username}
          rating={row.rating}
          isActive={row.isActive}
          key={row.id}
          challenge
        />
      ))}
    </>
  )
}

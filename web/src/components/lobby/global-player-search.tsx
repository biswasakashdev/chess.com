import { Input } from "@/components/ui/input"
import useAuthContext from "@/context/auth.context"
import { Search } from "lucide-react"
import { useEffect, useState } from "react"
import { UserRow } from "../user-row"

interface PlayerList {
  id: string
  username: string
  rating: number
  firstName: string
  lastName: string
}

export const SearchPlayers = () => {
  const { client } = useAuthContext()
  const [query, setQuery] = useState("")

  const [users, setUsers] = useState<PlayerList[]>([])

  useEffect(() => {
    if (query.length < 3) {
      return
    }
    const fetechData = async () => {
      const res = await client.get("/api/v1/users", {
        params: {
          search: query,
        },
      })
      if (res.status !== 200) {
        return
      }
      setUsers(res.data)
    }
    const timeOut = setTimeout(() => {
      fetechData()
    }, 500)

    return () => {
      clearTimeout(timeOut)
    }
  }, [query, client])

  return (
    <div className="space-y-4">
      <div className="relative">
        <Search className="absolute top-3 left-3 h-4 w-4 text-muted-foreground" />
        <Input
          placeholder="Search global players by username..."
          className="pl-9"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>

      {query.length > 1 && (
        <div className="space-y-2 rounded-lg border bg-card p-2 shadow-sm">
          {users.map((row) => {
            return (
              <UserRow
                id={row.id}
                firstName={row.firstName}
                lastName={row.lastName}
                username={row.username}
                rating={row.rating}
                key={row.id}
              />
            )
          })}
        </div>
      )}
    </div>
  )
}

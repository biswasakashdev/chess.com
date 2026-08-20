import React, { useState } from "react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Search, Swords, UserPlus } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export const GlobalPlayerSearch: React.FC = () => {
  const [query, setQuery] = useState("")

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
          <div className="flex items-center justify-between rounded-md p-2 transition-colors hover:bg-accent">
            <div className="flex items-center space-x-3">
              <Avatar className="h-8 w-8">
                <AvatarImage src={`https://avatar.vercel.sh/${query}`} />
                <AvatarFallback>
                  {query.slice(0, 2).toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <div>
                <p className="text-sm font-semibold">{query}</p>
                <p className="text-xs text-muted-foreground">Rating: 1840</p>
              </div>
            </div>
            <div className="flex space-x-2">
              <Button size="sm" variant="outline">
                <UserPlus className="mr-1 h-3.5 w-3.5" /> Add
              </Button>
              <Button size="sm">
                <Swords className="mr-1 h-3.5 w-3.5" /> Challenge
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

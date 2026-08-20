import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Swords, Circle } from "lucide-react"

interface Friend {
  id: string
  name: string
  rating: number
  status: "in-game" | "online" | "idle"
}

const friends: Friend[] = [
  { id: "1", name: "GarryK", rating: 2812, status: "in-game" },
  { id: "2", name: "AnishGiri", rating: 2760, status: "online" },
  { id: "3", name: "FabiC", rating: 2805, status: "idle" },
]

export const SocialSidebar: React.FC = () => {
  return (
    <Card className="h-full rounded-none border-l lg:rounded-lg">
      <CardHeader className="pb-3">
        <CardTitle className="text-base font-semibold">Community</CardTitle>
      </CardHeader>
      <CardContent className="px-3">
        <Tabs defaultValue="friends" className="w-full">
          <TabsList className="mb-4 grid w-full grid-cols-2">
            <TabsTrigger value="friends">
              Friends ({friends.length})
            </TabsTrigger>
            <TabsTrigger value="online">Global (1,248)</TabsTrigger>
          </TabsList>

          <TabsContent value="friends" className="space-y-3">
            {friends.map((friend) => (
              <div
                key={friend.id}
                className="flex items-center justify-between rounded-md p-2 transition-colors hover:bg-accent/60"
              >
                <div className="flex items-center space-x-2.5">
                  <div className="relative">
                    <Avatar className="h-8 w-8">
                      <AvatarImage
                        src={`https://avatar.vercel.sh/${friend.name}`}
                      />
                      <AvatarFallback>{friend.name.slice(0, 2)}</AvatarFallback>
                    </Avatar>
                    <Circle
                      className={`absolute right-0 bottom-0 h-2.5 w-2.5 fill-current ${
                        friend.status === "online"
                          ? "text-emerald-500"
                          : friend.status === "in-game"
                            ? "text-amber-500"
                            : "text-muted-foreground"
                      }`}
                    />
                  </div>
                  <div>
                    <div className="text-xs font-semibold">{friend.name}</div>
                    <div className="text-[11px] text-muted-foreground">
                      {friend.rating} • {friend.status}
                    </div>
                  </div>
                </div>

                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7"
                  title="Challenge"
                >
                  <Swords className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))}
          </TabsContent>

          <TabsContent value="online">
            <p className="py-4 text-center text-xs text-muted-foreground">
              Matchmaking queue is active. Search players to invite directly.
            </p>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import React from "react"
import FriendsList from "./friends-list"
import { SearchPlayers } from "@/components/lobby/global-player-search"

export const SocialSidebar: React.FC = () => {
  return (
    <Card className="h-full rounded-none border-l lg:rounded-lg">
      <CardHeader className="pb-3">
        <CardTitle className="text-base font-semibold">Community</CardTitle>
      </CardHeader>
      <CardContent className="px-3">
        <Tabs defaultValue="friends" className="w-full">
          <TabsList className="mb-4 grid w-full grid-cols-2">
            <TabsTrigger value="friends">Friends</TabsTrigger>
            <TabsTrigger value="online">Global</TabsTrigger>
          </TabsList>

          <TabsContent value="friends" className="space-y-3">
            <FriendsList />
          </TabsContent>

          <TabsContent value="online">
            <SearchPlayers />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}

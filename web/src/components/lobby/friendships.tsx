import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import useAuthContext from "@/context/auth.context"
import { getCapitalise, getInitials } from "@/utils/user-utils"
import type { AxiosInstance } from "axios"
import axios from "axios"
import {
  Check,
  Send,
  UserCheck,
  X
} from "lucide-react"
import { useEffect, useState } from "react"

// --- Types ---

export interface Requests {
  id: string
  username: string
  firstName: string
  lastName: string
  rating: number
}

interface UpdateFriendshipRequest {
  target_user_id: string
  action: "accept" | "block" | "unblock" | "cancel"
}

export const FriendshipTabs = () => {
  return (
    <Card className="mx-auto w-full max-w-xl border-border shadow-md">
      <CardHeader className="pb-4">
        <CardTitle className="flex items-center gap-2 text-xl font-bold">
          <UserCheck className="h-5 w-5 text-primary" /> Friend Requests &
          Privacy
        </CardTitle>
        <CardDescription>
          Manage incoming challenges, sent invitations
        </CardDescription>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Tabs System */}
        <Tabs defaultValue="requests" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger
              value="requests"
              className="flex items-center gap-1.5 text-xs sm:text-sm"
            >
              <UserCheck className="h-4 w-4" />
              <span>Requests</span>
            </TabsTrigger>

            <TabsTrigger
              value="sent"
              className="flex items-center gap-1.5 text-xs sm:text-sm"
            >
              <Send className="h-3.5 w-3.5" />
              <span>Sent</span>
            </TabsTrigger>


          </TabsList>

          {/* --- Tab 1: Incoming Requests --- */}
          <TabsContent value="requests" className="mt-4">
            <RequestsSection />
          </TabsContent>

          {/* --- Tab 2: Sent Requests --- */}
          <TabsContent value="sent" className="mt-4">
            <SentRequestSection />
          </TabsContent>

        </Tabs>
      </CardContent>
    </Card>
  )
}

export const RequestsSection = () => {
  const [incoming, setIncoming] = useState<Requests[]>([])
  const { client } = useAuthContext()

  // Trigger rendering
  const [updateState, setUpdateState] = useState({})

  const handleAccept = async (userId: string): Promise<void> => {
    try {
      const payload: UpdateFriendshipRequest = {
        target_user_id: userId,
        action: "accept",
      }

      const response = await client.patch("/api/v1/friends", payload)

      // Optimistically update your local UI state here if needed
      console.log("Friend request accepted:", response.data)
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error(
          "Failed to accept friend request:",
          error.response?.data || error.message
        )
      } else {
        console.error("Unexpected error:", error)
      }
    }

    setUpdateState({})
  }

  const handleDecline = async (userId: string) => {
    try {
      // Declining an incoming request matches the "cancel" action on the friendship state
      const payload: UpdateFriendshipRequest = {
        target_user_id: userId,
        action: "cancel",
      }

      const response = await client.post("/api/v1/friends/requests", payload)

      // Optimistically update your local UI state here if needed
      console.log("Friend request declined:", response.data)
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error(
          "Failed to decline friend request:",
          error.response?.data || error.message
        )
      } else {
        console.error("Unexpected error:", error)
      }
    }
    setUpdateState({})
  }

  useEffect(() => {
    ;(async function () {
      const data = await fetchRequests(client, "pending")
      setIncoming(data)
    })()
  }, [client,updateState])

  return (
    <ScrollArea className="h-[300px] pr-3">
      {incoming.length === 0 ? (
        <div className="flex h-[200px] flex-col items-center justify-center space-y-1 text-sm text-muted-foreground">
          <UserCheck className="h-8 w-8 text-muted-foreground/50" />
          <p>No incoming friend requests</p>
        </div>
      ) : (
        <div className="space-y-2.5">
          {incoming.map((req) => (
            <div
              key={req.id}
              className="flex items-center justify-between rounded-lg border border-border bg-card p-3 transition-colors hover:bg-muted/40"
            >
              <div className="flex items-center gap-3">
                <Avatar className="h-9 w-9 border border-border">
                  <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                    {/* Derives 2-letter initials from name, falls back to username */}
                    {getInitials(req.firstName, req.lastName)}
                  </AvatarFallback>
                </Avatar>

                <div>
                  <div className="flex flex-col gap-2">
                    {/* Full Name */}
                    <span className="text-sm leading-none font-semibold">
                      {getCapitalise(req.firstName, req.lastName)}
                    </span>

                    <span className="text-xs leading-none text-muted-foreground">
                      @{req.username}
                    </span>

                    {/* Rating */}
                    {req.rating !== undefined && (
                      <Badge
                        variant="outline"
                        className="h-4 px-1.5 py-0 font-mono text-[10px] leading-none"
                      >
                        {req.rating}
                      </Badge>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-1.5">
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => handleDecline(req.id)}
                  className="h-8 px-2.5 text-xs text-muted-foreground hover:text-destructive"
                >
                  <X className="mr-1 h-3.5 w-3.5" />
                </Button>
                <Button
                  size="sm"
                  onClick={() => handleAccept(req.id)}
                  className="h-8 bg-primary px-3 text-xs hover:bg-primary/90"
                >
                  <Check className="mr-1 h-3.5 w-3.5" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </ScrollArea>
  )
}


export const SentRequestSection = () => {
  const [sent, setSent] = useState<Requests[]>([])
  const { client } = useAuthContext()

  // Trigger rendering
  const [ updateState, setUpdateState] = useState({})

  const handleCancelSent = async (userId: string) => {
    try {
      // Declining an incoming request matches the "cancel" action on the friendship state
      const payload: UpdateFriendshipRequest = {
        target_user_id: userId,
        action: "cancel",
      }

      const response = await client.patch("/api/v1/friends", payload)

      // Optimistically update your local UI state here if needed
      console.log("Friend request declined:", response.data)
    } catch (error) {
      if (axios.isAxiosError(error)) {
        console.error(
          "Failed to decline friend request:",
          error.response?.data || error.message
        )
      } else {
        console.error("Unexpected error:", error)
      }
    }
    setUpdateState({})
  }

  useEffect(() => {
    ;(async function () {
      const data = await fetchRequests(client, "sent")
      setSent(data)
    })()
  }, [client,updateState])

  return (
    <ScrollArea className="h-[300px] pr-3">
      {sent.length === 0 ? (
        <div className="flex h-[200px] flex-col items-center justify-center space-y-1 text-sm text-muted-foreground">
          <Send className="h-8 w-8 text-muted-foreground/50" />
          <p>No outgoing requests pending</p>
        </div>
      ) : (
        <div className="space-y-2.5">
          {sent.map((item) => (
            <div
              key={item.id}
              className="flex items-center justify-between rounded-lg border border-border bg-card p-3 transition-colors hover:bg-muted/40"
            >
              <div className="flex items-center gap-3">
                <Avatar className="h-9 w-9 border border-border">
                  <AvatarFallback className="bg-primary/10 text-xs font-semibold text-primary">
                    {/* Derives 2-letter initials from name, falls back to username */}
                    {getInitials(item.firstName, item.lastName)}
                  </AvatarFallback>
                </Avatar>

                <div>
                  <div className="flex items-center gap-2">
                    {/* Full Name */}
                    <span className="text-sm leading-none font-semibold">
                      {getCapitalise(item.firstName, item.lastName)}
                    </span>

                    <span className="text-xs leading-none text-muted-foreground">
                      @{item.username}
                    </span>

                    {/* Rating */}
                    {item.rating !== undefined && (
                      <Badge
                        variant="outline"
                        className="h-4 px-1.5 py-0 font-mono text-[10px] leading-none"
                      >
                        {item.rating}
                      </Badge>
                    )}
                  </div>
                </div>
              </div>

              <Button
                size="sm"
                variant="ghost"
                onClick={() => handleCancelSent(item.id)}
                className="h-8 px-3 text-xs text-destructive hover:bg-destructive/10"
              >
                Cancel
              </Button>
            </div>
          ))}
        </div>
      )}
    </ScrollArea>
  )
}
const fetchRequests = async (
  client: AxiosInstance,
  type: "sent" | "blocked" | "pending"
) => {
  const { status, data } = await client.get("/api/v1/friends/requests", {
    params: {
      type,
    },
  })

  if (status === 200) {
    return data
  }
}

interface UpdateFriendshipRequest {
  target_user_id: string
  action: "accept" | "block" | "unblock" | "cancel"
}

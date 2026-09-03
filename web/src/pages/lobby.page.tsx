import { Navbar } from "@/components/home/navbar"
import { ActiveGamesList } from "@/components/lobby/active-game-list"
import { ActiveUsersSection } from "@/components/lobby/active-users"
import { FriendshipTabs } from "@/components/lobby/friendships"
import { RecentGamesTable } from "@/components/lobby/recent-games-list"
import { SocialSidebar } from "@/components/lobby/social-sidebar"
import { Button } from "@/components/ui/button"
import { Toaster } from "@/components/ui/sonner"
import useAuthContext from "@/context/auth.context"
import {
  useChessContext,
  useSocketEvent,
  type ChallengePayload,
  type GameStartPayload,
} from "@/context/game.context"
import useUserContext from "@/context/user.context"
import { getCapitalise } from "@/utils/user-utils"
import { Swords } from "lucide-react"
import { toast } from "sonner"

export const LobbyPage = () => {
  const { acceptChallenge } = useChessContext()
  const { user } = useUserContext()
  const{ updateAuthorization }= useAuthContext()

  useSocketEvent("game_start", (payload: GameStartPayload) => {
    const gameId = payload.game_id
    console.log("The created game id is: ",gameId)
    // navigate(`/game/${gameId}`)
  })

  useSocketEvent("challenge_request", (payload: ChallengePayload) => {
    const { username, firstName,lastName,id } = payload.from_user_data

    console.log("Some message came...")
    if (username && id && firstName && lastName) {
      toast(
        <div>
          <div>
            <div className="text-xs font-semibold">{getCapitalise(firstName, lastName)}</div>
            <div className="text-[11px] text-muted-foreground">
              @{username}
            </div>
          </div>
        </div>,
        {
          action: (
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7"
              title="Challenge"
              onClick={() => acceptChallenge(id)}
            >
              <Swords className="h-3.5 w-3.5" />
            </Button>
          ),
        }
      )
    }
  })

  return (
    <>
      <Navbar user={user} logoutHandler={()=> updateAuthorization(undefined)}/>
      {/* Main Grid Layout */}
      <main className="mx-auto grid w-full max-w-[1600px] flex-1 grid-cols-1 gap-6 p-6 lg:grid-cols-4">
        {/* Left / Center 3-column Primary Area */}
        <div className="space-y-2 lg:col-span-3">
          {/*<UserStatsCard
            username="GrandmasterDev"
            rating={1942}
            winRate={58}
            streak={4}
            totalGames={482}
          />*/}

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <ActiveGamesList />
            <RecentGamesTable />
            <ActiveUsersSection />
            <FriendshipTabs />
          </div>
        </div>

        {/* Right 1-column Social / Online Panel */}
        <div className="lg:col-span-1">
          <SocialSidebar />
        </div>

        <Toaster />
      </main>
    </>
  )
}

export default LobbyPage

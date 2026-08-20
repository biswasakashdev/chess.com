import { ActiveGamesList } from "@/components/lobby/active-game-list"
import { GlobalPlayerSearch } from "@/components/lobby/global-playler-search"
import { RecentGamesTable } from "@/components/lobby/recent-games-list"
import { SocialSidebar } from "@/components/lobby/social-sidebar"
import { UserStatsCard } from "@/components/lobby/user-status-card"

export const LobbyPage = () => {
  return (
    <>
      {/* Main Grid Layout */}
      <main className="mx-auto grid w-full max-w-[1600px] flex-1 grid-cols-1 gap-6 p-6 lg:grid-cols-4">
        {/* Left / Center 3-column Primary Area */}
        <div className="space-y-6 lg:col-span-3">
          <UserStatsCard
            username="GrandmasterDev"
            rating={1942}
            winRate={58}
            streak={4}
            totalGames={482}
          />

          <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div className="space-y-4">
              <GlobalPlayerSearch />
              <ActiveGamesList />
            </div>
            <div className="h-full">
              <RecentGamesTable />
            </div>
          </div>
        </div>

        {/* Right 1-column Social / Online Panel */}
        <div className="lg:col-span-1">
          <SocialSidebar />
        </div>
      </main>
    </>
  )
}

export default LobbyPage

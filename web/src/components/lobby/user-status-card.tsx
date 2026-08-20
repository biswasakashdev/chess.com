import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Trophy, Flame, Swords, TrendingUp } from "lucide-react"
import { motion } from "framer-motion"

interface UserStatsProps {
  username: string
  rating: number
  winRate: number
  streak: number
  totalGames: number
}

export const UserStatsCard: React.FC<UserStatsProps> = ({
  rating,
  winRate,
  streak,
  totalGames,
}) => {
  const stats = [
    {
      title: "Rapid Rating",
      value: rating,
      icon: Trophy,
      trend: "+24 this week",
    },
    {
      title: "Win Rate",
      value: `${winRate}%`,
      icon: TrendingUp,
      trend: "54W / 36L / 10D",
    },
    {
      title: "Current Streak",
      value: `${streak} Wins`,
      icon: Flame,
      trend: "Best: 8",
    },
    {
      title: "Total Games",
      value: totalGames,
      icon: Swords,
      trend: "Active since 2025",
    },
  ]

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, idx) => (
        <motion.div
          key={stat.title}
          initial={{ opacity: 0, y: 15 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: idx * 0.05 }}
        >
          <Card className="transition-colors hover:border-primary/40">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <stat.icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="mt-1 text-xs text-muted-foreground">{stat.trend}</p>
            </CardContent>
          </Card>
        </motion.div>
      ))}
    </div>
  )
}

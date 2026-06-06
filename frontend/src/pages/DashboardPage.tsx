import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useCurrentUser } from '@/features/auth/hooks/useCurrentUser'

export function DashboardPage() {
  const { data: user } = useCurrentUser()

  return (
    <div className="space-y-12">
      <section className="space-y-3">
        <p className="text-sm font-medium text-muted-foreground">Dashboard</p>
        <h2 className="text-3xl font-bold tracking-normal">
          Welcome, {user?.name ?? 'learner'}
        </h2>
        <p className="max-w-2xl text-muted-foreground">
          Your protected learning workspace is ready for the next vocabulary
          tasks.
        </p>
      </section>

      <Card>
        <CardHeader>
          <CardTitle>Auth is connected</CardTitle>
          <CardDescription>
            This protected page is available only with a valid access token.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-3 text-sm text-muted-foreground sm:grid-cols-3">
          <div className="rounded-xl border border-border bg-muted p-4">
            <p className="font-semibold text-foreground">Target band</p>
            <p>{user?.targetBand.toFixed(1)}</p>
          </div>
          <div className="rounded-xl border border-border bg-muted p-4">
            <p className="font-semibold text-foreground">Level</p>
            <p>
              {user?.level} - {user?.levelTitle}
            </p>
          </div>
          <div className="rounded-xl border border-border bg-muted p-4">
            <p className="font-semibold text-foreground">Current streak</p>
            <p>{user?.currentStreak ?? 0} days</p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

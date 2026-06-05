import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'

type DashboardPlaceholderProps = {
  eyebrow: string
  title: string
  description: string
}

export function DashboardPlaceholder({
  eyebrow,
  title,
  description,
}: DashboardPlaceholderProps) {
  return (
    <div className="space-y-12">
      <section className="space-y-3">
        <p className="text-sm font-medium text-muted-foreground">{eyebrow}</p>
        <h2 className="text-3xl font-bold tracking-normal">{title}</h2>
        <p className="max-w-2xl text-muted-foreground">{description}</p>
      </section>

      <Card>
        <CardHeader>
          <CardTitle>Layout placeholder</CardTitle>
          <CardDescription>
            This screen is connected to the dashboard navigation without adding
            feature data yet.
          </CardDescription>
        </CardHeader>
        <CardContent className="text-sm text-muted-foreground">
          Future task content can be added here while keeping the same page
          chrome.
        </CardContent>
      </Card>
    </div>
  )
}

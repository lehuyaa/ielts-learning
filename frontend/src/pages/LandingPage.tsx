import { zodResolver } from '@hookform/resolvers/zod'
import { BookOpen, CheckCircle2, Sparkles } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { z } from 'zod'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

const roadmapItems = [
  'Band-based roadmap',
  'Smart flashcards',
  'Quiz practice',
  'Progress tracking',
]

const demoFormSchema = z.object({
  email: z.string().email('Enter a valid email address.'),
})

type DemoFormValues = z.infer<typeof demoFormSchema>

export function LandingPage() {
  const form = useForm<DemoFormValues>({
    resolver: zodResolver(demoFormSchema),
    defaultValues: {
      email: '',
    },
  })

  function handleDemoSubmit() {
    form.reset()
  }

  return (
    <main className="min-h-screen bg-background text-foreground">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-6 py-8">
        <nav className="flex items-center justify-between">
          <a className="flex items-center gap-3 font-semibold" href="/">
            <span className="grid size-10 place-items-center rounded-2xl bg-primary text-lg font-bold text-primary-foreground">
              <BookOpen className="size-5" aria-hidden="true" />
            </span>
            <span>LexPath</span>
          </a>

          <Button asChild variant="outline">
            <a href="/">Start Learning</a>
          </Button>
        </nav>

        <div className="grid flex-1 items-center gap-12 py-16 lg:grid-cols-[1.05fr_0.95fr]">
          <div className="max-w-2xl">
            <Badge className="mb-5" variant="outline">
              <Sparkles className="size-3.5" aria-hidden="true" />
              IELTS vocabulary learning platform
            </Badge>

            <h1 className="text-5xl font-bold leading-tight tracking-normal text-slate-950 md:text-6xl">
              Build vocabulary for your target IELTS band.
            </h1>

            <p className="mt-6 max-w-xl text-lg leading-8 text-slate-600">
              Follow structured lessons, review words with spaced repetition,
              and check progress with quizzes designed for IELTS topics.
            </p>

            <div className="mt-8 flex flex-col gap-3 sm:flex-row">
              <Button asChild size="lg">
                <a href="/">Start Learning</a>
              </Button>
              <Button asChild size="lg" variant="outline">
                <a href="/">View Roadmap</a>
              </Button>
            </div>
          </div>

          <Card className="shadow-xl shadow-slate-200/70">
            <CardHeader>
              <CardDescription>Today&apos;s word</CardDescription>
              <CardTitle className="text-3xl">Sustainable</CardTitle>
              <div className="flex gap-2 pt-2">
                <Badge>Band 7</Badge>
                <Badge variant="success">adjective</Badge>
              </div>
            </CardHeader>

            <CardContent className="grid gap-6">
              <div className="grid gap-3">
                {roadmapItems.map((item) => (
                  <div
                    className="flex items-center justify-between rounded-xl border border-border bg-muted px-4 py-3"
                    key={item}
                  >
                    <span className="font-medium text-slate-700">{item}</span>
                    <CheckCircle2
                      className="size-5 text-success"
                      aria-hidden="true"
                    />
                  </div>
                ))}
              </div>

              <Form {...form}>
                <form
                  className="grid gap-4 rounded-xl border border-border bg-background p-4"
                  onSubmit={form.handleSubmit(handleDemoSubmit)}
                >
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Email</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="alex@example.com"
                            type="email"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          A small form check for the shared UI setup.
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <Button type="submit" variant="secondary">
                    Check Components
                  </Button>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>
      </section>
    </main>
  )
}

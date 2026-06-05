const roadmapItems = [
  'Band-based roadmap',
  'Smart flashcards',
  'Quiz practice',
  'Progress tracking',
]

export function LandingPage() {
  return (
    <main className="min-h-screen bg-slate-50 text-slate-900">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-6 py-8">
        <nav className="flex items-center justify-between">
          <a className="flex items-center gap-3 font-semibold" href="/">
            <span className="grid size-10 place-items-center rounded-2xl bg-indigo-600 text-lg font-bold text-white">
              L
            </span>
            <span>LexPath</span>
          </a>

          <a
            className="rounded-xl border border-slate-200 bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm transition hover:border-indigo-200 hover:text-indigo-600"
            href="/"
          >
            Start Learning
          </a>
        </nav>

        <div className="grid flex-1 items-center gap-12 py-16 lg:grid-cols-[1.05fr_0.95fr]">
          <div className="max-w-2xl">
            <p className="mb-5 inline-flex rounded-full border border-indigo-100 bg-white px-4 py-2 text-sm font-medium text-indigo-600 shadow-sm">
              IELTS vocabulary learning platform
            </p>

            <h1 className="text-5xl font-bold leading-tight tracking-normal text-slate-950 md:text-6xl">
              Build vocabulary for your target IELTS band.
            </h1>

            <p className="mt-6 max-w-xl text-lg leading-8 text-slate-600">
              Follow structured lessons, review words with spaced repetition,
              and check progress with quizzes designed for IELTS topics.
            </p>

            <div className="mt-8 flex flex-col gap-3 sm:flex-row">
              <a
                className="rounded-xl bg-indigo-600 px-5 py-3 text-center font-semibold text-white shadow-lg shadow-indigo-200 transition hover:bg-indigo-700"
                href="/"
              >
                Start Learning
              </a>
              <a
                className="rounded-xl border border-slate-200 bg-white px-5 py-3 text-center font-semibold text-slate-700 shadow-sm transition hover:border-indigo-200 hover:text-indigo-600"
                href="/"
              >
                View Roadmap
              </a>
            </div>
          </div>

          <div className="rounded-2xl border border-slate-200 bg-white p-6 shadow-xl shadow-slate-200/70">
            <div className="rounded-xl bg-indigo-600 p-6 text-white">
              <p className="text-sm font-medium text-indigo-100">Today&apos;s word</p>
              <h2 className="mt-8 text-4xl font-bold tracking-normal">Sustainable</h2>
              <p className="mt-2 text-indigo-100">adjective · Band 7</p>
            </div>

            <div className="mt-6 grid gap-3">
              {roadmapItems.map((item) => (
                <div
                  className="flex items-center justify-between rounded-xl border border-slate-100 bg-slate-50 px-4 py-3"
                  key={item}
                >
                  <span className="font-medium text-slate-700">{item}</span>
                  <span className="size-2 rounded-full bg-emerald-500" />
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>
    </main>
  )
}

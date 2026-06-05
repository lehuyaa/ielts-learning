type SectionHeaderProps = {
  eyebrow: string
  title: string
  description: string
}

export function SectionHeader({ eyebrow, title, description }: SectionHeaderProps) {
  return (
    <div className="mx-auto max-w-3xl space-y-3 text-center">
      <p className="text-sm font-semibold uppercase tracking-normal text-primary">
        {eyebrow}
      </p>
      <h2 className="text-3xl font-bold leading-tight tracking-normal text-slate-950 md:text-4xl">
        {title}
      </h2>
      <p className="text-base leading-6 text-muted-foreground md:text-lg md:leading-8">
        {description}
      </p>
    </div>
  )
}

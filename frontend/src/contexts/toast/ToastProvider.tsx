import { CheckCircle2, Info, X, XCircle } from 'lucide-react'
import type React from 'react'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'

import { cn } from '@/lib/utils'

import {
  ToastContext,
  type ToastContextValue,
  type ToastInput,
  type ToastTone,
} from './toastContext'

type ToastItem = ToastInput & {
  id: number
}

const toastStyles: Record<
  ToastTone,
  {
    icon: typeof CheckCircle2
    iconClassName: string
  }
> = {
  success: {
    icon: CheckCircle2,
    iconClassName: 'text-emerald-600',
  },
  error: {
    icon: XCircle,
    iconClassName: 'text-destructive',
  },
  info: {
    icon: Info,
    iconClassName: 'text-primary',
  },
}

type ToastProviderProps = {
  children: React.ReactNode
}

export function ToastProvider({ children }: ToastProviderProps) {
  const [toasts, setToasts] = useState<ToastItem[]>([])

  const dismissToast = useCallback((id: number) => {
    setToasts((current) => current.filter((toast) => toast.id !== id))
  }, [])

  const showToast = useCallback((toast: ToastInput) => {
    const id = Date.now() + Math.floor(Math.random() * 1000)
    setToasts((current) => [...current, { ...toast, id }])
  }, [])

  useEffect(() => {
    if (toasts.length === 0) {
      return
    }

    const timers = toasts.map((toast) =>
      window.setTimeout(() => {
        dismissToast(toast.id)
      }, 4000),
    )

    return () => {
      timers.forEach((timer) => window.clearTimeout(timer))
    }
  }, [dismissToast, toasts])

  const value = useMemo<ToastContextValue>(
    () => ({
      showToast,
    }),
    [showToast],
  )

  return (
    <ToastContext.Provider value={value}>
      {children}
      {typeof document !== 'undefined'
        ? createPortal(
            <div className="pointer-events-none fixed right-4 top-4 z-[100] flex w-[min(360px,calc(100vw-2rem))] flex-col gap-3">
              {toasts.map((toast) => {
                const tone = toast.tone ?? 'info'
                const style = toastStyles[tone]
                const Icon = style.icon

                return (
                  <section
                    className="pointer-events-auto rounded-2xl border border-border bg-white p-4 shadow-lg"
                    key={toast.id}
                    role="status"
                  >
                    <div className="flex items-start gap-3">
                      <Icon
                        className={cn('mt-0.5 size-5 shrink-0', style.iconClassName)}
                        aria-hidden="true"
                      />
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-semibold text-foreground">
                          {toast.title}
                        </p>
                        {toast.description ? (
                          <p className="mt-1 text-sm text-muted-foreground">
                            {toast.description}
                          </p>
                        ) : null}
                      </div>
                      <button
                        aria-label="Dismiss notification"
                        className="cursor-pointer rounded-full p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                        onClick={() => dismissToast(toast.id)}
                        type="button"
                      >
                        <X className="size-4" aria-hidden="true" />
                      </button>
                    </div>
                  </section>
                )
              })}
            </div>,
            document.body,
          )
        : null}
    </ToastContext.Provider>
  )
}

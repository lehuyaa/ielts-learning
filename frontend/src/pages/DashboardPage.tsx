import { APIError } from '@/api/api'
import { DashboardOverview } from '@/features/dashboard/DashboardOverview'
import {
  loadingDashboardViewModel,
  mapDashboardSummaryToViewModel,
} from '@/features/dashboard/mapDashboardSummary'
import { useDashboardSummary } from '@/features/dashboard/hooks/useDashboardSummary'

export function DashboardPage() {
  const dashboardQuery = useDashboardSummary()
  const dashboard = dashboardQuery.data
    ? mapDashboardSummaryToViewModel(dashboardQuery.data)
    : loadingDashboardViewModel

  return (
    <DashboardOverview
      dashboard={dashboard}
      errorMessage={getDashboardErrorMessage(dashboardQuery.error)}
      isEmpty={
        !dashboardQuery.isLoading &&
        !dashboardQuery.error &&
        !dashboardQuery.data
      }
      isLoading={dashboardQuery.isLoading}
      onRetry={() => {
        void dashboardQuery.refetch()
      }}
    />
  )
}

function getDashboardErrorMessage(error: Error | null) {
  if (!error) {
    return null
  }

  return error instanceof APIError
    ? error.message
    : 'Unable to load the dashboard right now.'
}

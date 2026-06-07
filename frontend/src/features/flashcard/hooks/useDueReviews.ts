import { useQuery } from "@tanstack/react-query";

import { getDueReviews } from "@/api/flashcard";
import type { DueReviewsQueryParams } from "@/types/flashcard";

export function dueReviewsQueryKey(params?: DueReviewsQueryParams) {
  return ["dueReviews", params ?? {}] as const;
}

export function useDueReviews(params?: DueReviewsQueryParams) {
  return useQuery({
    queryKey: dueReviewsQueryKey(params),
    queryFn: () => getDueReviews(params),
    staleTime: 60 * 1000,
  });
}

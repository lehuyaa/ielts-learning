import { api, unwrapData } from "@/api/api";
import type { APIResponse } from "@/types/api";
import type {
  DueReviewsQueryParams,
  DueReviewsResponse,
  LessonFlashcardsResponse,
  ReviewFlashcardRequest,
  ReviewFlashcardResponse,
} from "@/types/flashcard";

export async function getLessonFlashcards(lessonId: string | number) {
  const response = await api.get<APIResponse<LessonFlashcardsResponse>>(
    `/lessons/${lessonId}/flashcards`,
  );

  return unwrapData<LessonFlashcardsResponse>(response);
}

export async function getDueReviews(params?: DueReviewsQueryParams) {
  const response = await api.get<APIResponse<DueReviewsResponse>>(
    "/reviews/due",
    { params },
  );

  return unwrapData<DueReviewsResponse>(response);
}

export async function reviewFlashcard(payload: ReviewFlashcardRequest) {
  const response = await api.post<APIResponse<ReviewFlashcardResponse>>(
    "/flashcards/review",
    payload,
  );

  return unwrapData<ReviewFlashcardResponse>(response);
}

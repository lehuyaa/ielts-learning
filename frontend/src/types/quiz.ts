export type QuizQuestionType =
  | "MEANING_CHOICE"
  | "WORD_CHOICE"
  | "USAGE_CHOICE"
  | "SENTENCE_CHOICE"
  | "MULTIPLE_CHOICE";

export type QuizOption = {
  id: number;
  label: string;
  text: string;
};

export type QuizQuestion = {
  id: number;
  type: QuizQuestionType;
  prompt: string;
  points: number;
  options: QuizOption[];
};

export type LessonQuiz = {
  id: number;
  title: string;
  requiredScore: number;
  timeLimitSeconds: number | null;
};

export type LessonQuizResponse = {
  lesson: LessonQuiz;
  questions: QuizQuestion[];
};

export type QuizAnswerInput = {
  questionId: number;
  optionId: number;
};

export type SubmitQuizRequest = {
  answers: QuizAnswerInput[];
};

export type CheckQuizAnswerRequest = {
  questionId: number;
  optionId: number;
};

export type CheckQuizAnswerResponse = {
  questionId: number;
  selectedOptionId: number;
  correctOptionId: number;
  isCorrect: boolean;
  explanation: string;
  earnedPoints: number;
};

export type QuizResultItem = {
  questionId: number;
  selectedOptionId: number | null;
  correctOptionId: number;
  isCorrect: boolean;
  explanation: string;
  pointsAwarded: number;
};

export type SubmitQuizResponse = {
  attemptId: number;
  lessonId: number;
  score: number;
  requiredScore: number;
  passed: boolean;
  correctCount: number;
  totalQuestions: number;
  earnedXp: number;
  completedAt: string | null;
  results: QuizResultItem[];
};

export type QuizSessionQuestion = QuizQuestion;

export type QuizSessionViewModel = {
  lessonId: number;
  lessonTitle: string;
  topicTitle: string;
  band: string;
  requiredScore: number;
  timeLimitSeconds: number;
  questions: QuizSessionQuestion[];
};

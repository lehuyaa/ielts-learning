export type MockQuizOption = {
  id: string
  label: 'A' | 'B' | 'C' | 'D'
  text: string
}

export type MockQuizQuestion = {
  id: string
  vocabularyId: number
  type: 'MULTIPLE_CHOICE'
  prompt: string
  options: MockQuizOption[]
  correctOptionId: string
  explanation: string
  points: number
}

export type MockQuiz = {
  lessonId: number
  lessonTitle: string
  topicTitle: string
  band: string
  requiredScore: number
  timeLimitSeconds: number
  questions: MockQuizQuestion[]
}

export const mockQuiz: MockQuiz = {
  lessonId: 1,
  lessonTitle: 'Environment Vocabulary',
  topicTitle: 'Environment',
  band: 'Band 7',
  requiredScore: 80,
  timeLimitSeconds: 30,
  questions: [
    {
      id: 'q1',
      vocabularyId: 1,
      type: 'MULTIPLE_CHOICE',
      prompt: 'What does "sustainable" mean?',
      correctOptionId: 'q1-a',
      explanation:
        'The correct answer is: "Able to be maintained without depleting natural resources".',
      points: 20,
      options: [
        {
          id: 'q1-a',
          label: 'A',
          text: 'Able to be maintained without depleting natural resources',
        },
        {
          id: 'q1-b',
          label: 'B',
          text: 'Relating to the study of living organisms',
        },
        {
          id: 'q1-c',
          label: 'C',
          text: 'Causing harm to the environment',
        },
        {
          id: 'q1-d',
          label: 'D',
          text: 'A method of rapid economic growth',
        },
      ],
    },
    {
      id: 'q2',
      vocabularyId: 2,
      type: 'MULTIPLE_CHOICE',
      prompt: 'Which word means "to make a harmful effect less severe"?',
      correctOptionId: 'q2-c',
      explanation:
        'Mitigate means to reduce the severity or impact of a problem.',
      points: 20,
      options: [
        { id: 'q2-a', label: 'A', text: 'Proliferate' },
        { id: 'q2-b', label: 'B', text: 'Deteriorate' },
        { id: 'q2-c', label: 'C', text: 'Mitigate' },
        { id: 'q2-d', label: 'D', text: 'Consume' },
      ],
    },
    {
      id: 'q3',
      vocabularyId: 3,
      type: 'MULTIPLE_CHOICE',
      prompt: 'Which sentence uses "proliferate" correctly?',
      correctOptionId: 'q3-c',
      explanation:
        'The correct answer is: "Online misinformation has proliferated rapidly since 2016."',
      points: 20,
      options: [
        {
          id: 'q3-a',
          label: 'A',
          text: 'Scientists proliferate the new vaccine to ensure safety.',
        },
        {
          id: 'q3-b',
          label: 'B',
          text: 'The government plans to proliferate education spending.',
        },
        {
          id: 'q3-c',
          label: 'C',
          text: 'Online misinformation has proliferated rapidly since 2016.',
        },
        {
          id: 'q3-d',
          label: 'D',
          text: 'Doctors proliferate with their medical opinions on the matter.',
        },
      ],
    },
    {
      id: 'q4',
      vocabularyId: 4,
      type: 'MULTIPLE_CHOICE',
      prompt: 'What is the best synonym for "resilience"?',
      correctOptionId: 'q4-b',
      explanation:
        'Resilience is the ability to recover after difficulty, so "endurance" is the closest option.',
      points: 20,
      options: [
        { id: 'q4-a', label: 'A', text: 'Convenience' },
        { id: 'q4-b', label: 'B', text: 'Endurance' },
        { id: 'q4-c', label: 'C', text: 'Accuracy' },
        { id: 'q4-d', label: 'D', text: 'Isolation' },
      ],
    },
    {
      id: 'q5',
      vocabularyId: 5,
      type: 'MULTIPLE_CHOICE',
      prompt: 'Which phrase is a natural collocation with "infrastructure"?',
      correctOptionId: 'q5-a',
      explanation:
        '"Transport infrastructure" is a common IELTS collocation for cities and development.',
      points: 20,
      options: [
        { id: 'q5-a', label: 'A', text: 'Transport infrastructure' },
        { id: 'q5-b', label: 'B', text: 'Infrastructure emotion' },
        { id: 'q5-c', label: 'C', text: 'Infrastructure opinion' },
        { id: 'q5-d', label: 'D', text: 'Silent infrastructure' },
      ],
    },
    {
      id: 'q6',
      vocabularyId: 6,
      type: 'MULTIPLE_CHOICE',
      prompt: 'In IELTS writing, "accountability" is most useful for topics about...',
      correctOptionId: 'q6-d',
      explanation:
        'Accountability is especially useful when discussing government, leadership, and public responsibility.',
      points: 20,
      options: [
        { id: 'q6-a', label: 'A', text: 'Weather patterns' },
        { id: 'q6-b', label: 'B', text: 'Personal hobbies' },
        { id: 'q6-c', label: 'C', text: 'Food preparation' },
        { id: 'q6-d', label: 'D', text: 'Government responsibility' },
      ],
    },
  ],
}

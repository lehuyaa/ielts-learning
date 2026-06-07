export type LessonVocabularyItem = {
  id: string
  word: string
  ipa: string
  partOfSpeech: string
  shortDefinition: string
  definition: string
  example: string
  difficulty: 'Beginner' | 'Intermediate' | 'Advanced'
  band: string
  learned: boolean
}

export type MockLesson = {
  id: string
  title: string
  topic: string
  bandRange: string
  description: string
  estimatedMinutes: number
  requiredScore: number
  xpReward: number
  score: number | null
  bestScore: number | null
  vocabulary: LessonVocabularyItem[]
}

export const mockLesson: MockLesson = {
  id: '1',
  title: 'Education Vocabulary',
  topic: 'Education',
  bandRange: '5.5 - 7.0',
  description:
    'Build precise IELTS vocabulary for education, learning environments, academic pressure, and study habits.',
  estimatedMinutes: 25,
  requiredScore: 80,
  xpReward: 350,
  score: null,
  bestScore: null,
  vocabulary: [
    {
      id: 'accommodate',
      word: 'Accommodate',
      ipa: "/əˈkɒm.ə.deɪt/",
      partOfSpeech: 'verb',
      shortDefinition: 'To provide space or housing for; to adapt or adjust to...',
      definition: 'To provide space or housing for; to adapt or adjust to meet needs.',
      example:
        'Universities must accommodate students with diverse learning needs.',
      difficulty: 'Intermediate',
      band: 'Band 6',
      learned: true,
    },
    {
      id: 'adversely',
      word: 'Adversely',
      ipa: "/ˈaed.vɜːs.li/",
      partOfSpeech: 'adverb',
      shortDefinition: 'In a way that is harmful or unfavorable.',
      definition: 'In a way that has a negative or harmful effect.',
      example: 'Budget cuts may adversely affect classroom resources.',
      difficulty: 'Intermediate',
      band: 'Band 6.5',
      learned: true,
    },
    {
      id: 'analytical',
      word: 'Analytical',
      ipa: "/ˌæn.əˈlɪt.ɪ.kəl/",
      partOfSpeech: 'adjective',
      shortDefinition: 'Relating to analysis; using systematic examination of...',
      definition: 'Using careful reasoning and systematic examination.',
      example: 'Academic writing often requires analytical thinking.',
      difficulty: 'Advanced',
      band: 'Band 7',
      learned: true,
    },
    {
      id: 'curriculum',
      word: 'Curriculum',
      ipa: "/kəˈrɪk.jə.ləm/",
      partOfSpeech: 'noun',
      shortDefinition: 'The subjects comprising a course of study in a school o...',
      definition: 'The subjects and learning activities included in a course.',
      example: 'The curriculum includes academic vocabulary development.',
      difficulty: 'Beginner',
      band: 'Band 5.5',
      learned: true,
    },
    {
      id: 'deteriorate',
      word: 'Deteriorate',
      ipa: "/dɪˈtɪə.ri.ə.reɪt/",
      partOfSpeech: 'verb',
      shortDefinition: 'To become progressively worse in quality or...',
      definition: 'To become worse over time.',
      example: 'Study habits can deteriorate without regular practice.',
      difficulty: 'Intermediate',
      band: 'Band 6.5',
      learned: true,
    },
    {
      id: 'endeavour',
      word: 'Endeavour',
      ipa: "/ɪnˈdev.ər/",
      partOfSpeech: 'noun/verb',
      shortDefinition: 'A serious attempt or effort to do something.',
      definition: 'A serious attempt or sustained effort to achieve a goal.',
      example: 'Preparing for IELTS is a demanding academic endeavour.',
      difficulty: 'Advanced',
      band: 'Band 7',
      learned: true,
    },
    {
      id: 'facilitate',
      word: 'Facilitate',
      ipa: "/fəˈsɪl.ɪ.teɪt/",
      partOfSpeech: 'verb',
      shortDefinition: 'To make an action or process easier or more...',
      definition: 'To make a process easier or more effective.',
      example: 'Clear feedback can facilitate faster vocabulary growth.',
      difficulty: 'Intermediate',
      band: 'Band 6.5',
      learned: true,
    },
  ],
}

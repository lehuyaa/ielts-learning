export type RoadmapLessonStatus =
  | 'locked'
  | 'unlocked'
  | 'in-progress'
  | 'completed'

export type RoadmapLesson = {
  id: string
  title: string
  duration: string
  wordCount: number
  status: RoadmapLessonStatus
}

export type RoadmapTopic = {
  id: string
  title: string
  description: string
  progress: number
  icon: string
  completedLessons: number
  totalLessons: number
  lessons: RoadmapLesson[]
}

export type RoadmapBand = {
  id: string
  band: string
  title: string
  description: string
  progress: number
  lessonCount: number
  topicCount: number
  topics: RoadmapTopic[]
}

export const roadmapBands: RoadmapBand[] = [
  {
    id: 'band-5',
    band: '5.0',
    title: 'Foundation Vocabulary',
    description:
      'Build reliable words for everyday IELTS topics and short academic responses.',
    progress: 72,
    lessonCount: 12,
    topicCount: 3,
    topics: [
      {
        id: 'education-basics',
        title: 'Education',
        description: 'School, study habits, exams, and learning environments.',
        icon: '🎓',
        progress: 100,
        completedLessons: 8,
        totalLessons: 8,
        lessons: [
          {
            id: 'lesson-5-education-1',
            title: 'School Systems',
            duration: '12 min',
            wordCount: 18,
            status: 'completed',
          },
          {
            id: 'lesson-5-education-2',
            title: 'Study Habits',
            duration: '10 min',
            wordCount: 16,
            status: 'completed',
          },
          {
            id: 'lesson-5-education-3',
            title: 'Exam Pressure',
            duration: '14 min',
            wordCount: 20,
            status: 'completed',
          },
        ],
      },
      {
        id: 'daily-life',
        title: 'Health',
        description: 'Common situations for speaking and general writing tasks.',
        icon: '🏥',
        progress: 80,
        completedLessons: 6,
        totalLessons: 8,
        lessons: [
          {
            id: 'lesson-5-life-1',
            title: 'Routines',
            duration: '11 min',
            wordCount: 15,
            status: 'completed',
          },
          {
            id: 'lesson-5-life-2',
            title: 'Home and Family',
            duration: '13 min',
            wordCount: 18,
            status: 'in-progress',
          },
          {
            id: 'lesson-5-life-3',
            title: 'Free Time',
            duration: '12 min',
            wordCount: 17,
            status: 'unlocked',
          },
        ],
      },
      {
        id: 'society-basics',
        title: 'Society',
        description: 'People, communities, culture, and public life.',
        icon: '👥',
        progress: 60,
        completedLessons: 5,
        totalLessons: 8,
        lessons: [
          {
            id: 'lesson-5-society-1',
            title: 'Community Life',
            duration: '12 min',
            wordCount: 18,
            status: 'completed',
          },
          {
            id: 'lesson-5-society-2',
            title: 'Social Habits',
            duration: '13 min',
            wordCount: 19,
            status: 'unlocked',
          },
        ],
      },
    ],
  },
  {
    id: 'band-6',
    band: '6.0',
    title: 'Topic Fluency',
    description:
      'Strengthen vocabulary for common IELTS themes with clearer examples and collocations.',
    progress: 38,
    lessonCount: 40,
    topicCount: 4,
    topics: [
      {
        id: 'environment',
        title: 'Environment',
        description: 'Climate, conservation, pollution, and sustainable choices.',
        icon: '🌿',
        progress: 40,
        completedLessons: 4,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-6-environment-1',
            title: 'Climate Change',
            duration: '15 min',
            wordCount: 22,
            status: 'unlocked',
          },
          {
            id: 'lesson-6-environment-2',
            title: 'Sustainability',
            duration: '16 min',
            wordCount: 24,
            status: 'locked',
          },
          {
            id: 'lesson-6-environment-3',
            title: 'Urban Pollution',
            duration: '14 min',
            wordCount: 20,
            status: 'locked',
          },
        ],
      },
      {
        id: 'technology',
        title: 'Technology',
        description: 'Innovation, digital habits, automation, and online life.',
        icon: '💻',
        progress: 65,
        completedLessons: 7,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-6-technology-1',
            title: 'Digital Communication',
            duration: '13 min',
            wordCount: 19,
            status: 'unlocked',
          },
          {
            id: 'lesson-6-technology-2',
            title: 'Artificial Intelligence',
            duration: '17 min',
            wordCount: 25,
            status: 'locked',
          },
          {
            id: 'lesson-6-technology-3',
            title: 'Online Privacy',
            duration: '14 min',
            wordCount: 21,
            status: 'locked',
          },
        ],
      },
      {
        id: 'government',
        title: 'Government',
        description: 'Policy, law, leadership, and public institutions.',
        icon: '🏛️',
        progress: 20,
        completedLessons: 2,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-6-government-1',
            title: 'Public Policy',
            duration: '14 min',
            wordCount: 20,
            status: 'unlocked',
          },
          {
            id: 'lesson-6-government-2',
            title: 'Local Services',
            duration: '15 min',
            wordCount: 21,
            status: 'locked',
          },
        ],
      },
      {
        id: 'economy',
        title: 'Economy',
        description: 'Markets, employment, growth, and financial change.',
        icon: '📈',
        progress: 0,
        completedLessons: 0,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-6-economy-1',
            title: 'Economic Growth',
            duration: '16 min',
            wordCount: 22,
            status: 'locked',
          },
          {
            id: 'lesson-6-economy-2',
            title: 'Consumer Trends',
            duration: '14 min',
            wordCount: 20,
            status: 'locked',
          },
        ],
      },
    ],
  },
  {
    id: 'band-7',
    band: '7.0',
    title: 'Academic Precision',
    description:
      'Learn more precise vocabulary for argument, comparison, cause, and evaluation.',
    progress: 8,
    lessonCount: 16,
    topicCount: 4,
    topics: [
      {
        id: 'society',
        title: 'Society',
        description: 'Social change, communities, culture, and public policy.',
        icon: '👥',
        progress: 8,
        completedLessons: 0,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-7-society-1',
            title: 'Social Trends',
            duration: '16 min',
            wordCount: 24,
            status: 'locked',
          },
          {
            id: 'lesson-7-society-2',
            title: 'Public Services',
            duration: '18 min',
            wordCount: 26,
            status: 'locked',
          },
        ],
      },
      {
        id: 'work',
        title: 'Work',
        description: 'Employment, productivity, leadership, and career change.',
        icon: '💼',
        progress: 0,
        completedLessons: 0,
        totalLessons: 10,
        lessons: [
          {
            id: 'lesson-7-work-1',
            title: 'Workplace Change',
            duration: '15 min',
            wordCount: 23,
            status: 'locked',
          },
          {
            id: 'lesson-7-work-2',
            title: 'Leadership',
            duration: '17 min',
            wordCount: 24,
            status: 'locked',
          },
        ],
      },
    ],
  },
  {
    id: 'band-8',
    band: '8.0',
    title: 'Advanced Expression',
    description:
      'Practice flexible academic language for nuanced opinions and complex explanations.',
    progress: 0,
    lessonCount: 8,
    topicCount: 2,
    topics: [
      {
        id: 'advanced-arguments',
        title: 'Advanced Arguments',
        description: 'Concession, evaluation, implication, and synthesis.',
        icon: '🏆',
        progress: 0,
        completedLessons: 0,
        totalLessons: 8,
        lessons: [
          {
            id: 'lesson-8-arguments-1',
            title: 'Balanced Opinions',
            duration: '18 min',
            wordCount: 28,
            status: 'locked',
          },
          {
            id: 'lesson-8-arguments-2',
            title: 'Nuanced Conclusions',
            duration: '19 min',
            wordCount: 30,
            status: 'locked',
          },
        ],
      },
    ],
  },
]

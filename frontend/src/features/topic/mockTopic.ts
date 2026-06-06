export type TopicLessonStatus =
  | 'completed'
  | 'in-progress'
  | 'unlocked'
  | 'locked'

export type TopicLesson = {
  id: string
  title: string
  description: string
  wordCount: number
  estimatedMinutes: number
  xpReward: number
  status: TopicLessonStatus
  progressPercentage: number
  lockedReason: string | null
}

export type MockTopic = {
  id: string
  title: string
  icon: string
  band: string
  description: string
  progressPercentage: number
  completedLessons: number
  totalLessons: number
  totalXP: number
  lessons: TopicLesson[]
}

export const mockTopic: MockTopic = {
  id: '1',
  title: 'Technology',
  icon: '💻',
  band: 'Band 6.0',
  description:
    'Master technology-related vocabulary for IELTS Writing and Speaking tasks.',
  progressPercentage: 40,
  completedLessons: 4,
  totalLessons: 10,
  totalXP: 3400,
  lessons: [
    {
      id: '1',
      title: 'Digital Communication',
      description:
        'Essential terms for discussing online communication, social media, and digital platforms.',
      wordCount: 12,
      estimatedMinutes: 25,
      xpReward: 300,
      status: 'completed',
      progressPercentage: 100,
      lockedReason: null,
    },
    {
      id: '2',
      title: 'Internet & Connectivity',
      description:
        'Vocabulary for describing internet access, networking, and global connectivity issues.',
      wordCount: 10,
      estimatedMinutes: 20,
      xpReward: 250,
      status: 'completed',
      progressPercentage: 100,
      lockedReason: null,
    },
    {
      id: '3',
      title: 'Artificial Intelligence',
      description:
        'Key terms for automation, machine learning, algorithms, and intelligent systems.',
      wordCount: 14,
      estimatedMinutes: 28,
      xpReward: 350,
      status: 'completed',
      progressPercentage: 100,
      lockedReason: null,
    },
    {
      id: '4',
      title: 'Data Privacy & Security',
      description:
        'Key vocabulary for cyber security, data protection, and online privacy concerns.',
      wordCount: 14,
      estimatedMinutes: 28,
      xpReward: 350,
      status: 'completed',
      progressPercentage: 100,
      lockedReason: null,
    },
    {
      id: '5',
      title: 'E-commerce & Digital Economy',
      description:
        'Terms related to online shopping, digital payments, and the modern economy.',
      wordCount: 12,
      estimatedMinutes: 25,
      xpReward: 300,
      status: 'in-progress',
      progressPercentage: 65,
      lockedReason: null,
    },
    {
      id: '6',
      title: 'Mobile Technology',
      description:
        'Vocabulary for smartphones, mobile apps, and portable device innovations.',
      wordCount: 11,
      estimatedMinutes: 22,
      xpReward: 275,
      status: 'unlocked',
      progressPercentage: 0,
      lockedReason: null,
    },
    {
      id: '7',
      title: 'Software & Applications',
      description:
        'Essential terms for discussing software development, apps, and digital tools.',
      wordCount: 13,
      estimatedMinutes: 26,
      xpReward: 325,
      status: 'unlocked',
      progressPercentage: 0,
      lockedReason: null,
    },
    {
      id: '8',
      title: 'Cloud Computing',
      description:
        'Advanced vocabulary for cloud services, remote storage, and distributed systems.',
      wordCount: 10,
      estimatedMinutes: 20,
      xpReward: 250,
      status: 'locked',
      progressPercentage: 0,
      lockedReason: 'Score 850+ to unlock',
    },
    {
      id: '9',
      title: 'Innovation & Research',
      description:
        'Terms for technological advancement, R&D, and breakthrough discoveries.',
      wordCount: 16,
      estimatedMinutes: 32,
      xpReward: 450,
      status: 'locked',
      progressPercentage: 0,
      lockedReason: 'Score 900+ to unlock',
    },
    {
      id: '10',
      title: 'Future Technologies',
      description:
        'Speculative vocabulary for emerging tech: quantum computing, biotech, nanotech.',
      wordCount: 18,
      estimatedMinutes: 35,
      xpReward: 500,
      status: 'locked',
      progressPercentage: 0,
      lockedReason: 'Score 950+ to unlock',
    },
  ],
}

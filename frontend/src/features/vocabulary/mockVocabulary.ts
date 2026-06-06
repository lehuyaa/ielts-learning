export type VocabularyDifficulty = "Beginner" | "Intermediate" | "Advanced";

export type VocabularyStatus = "New" | "Learning" | "Review" | "Mastered";

export type RelatedForm = {
  word: string;
  partOfSpeech: string;
};

export type VocabularyExample = {
  sentence: string;
  note: string;
};

export type MockVocabulary = {
  id: string;
  word: string;
  slug: string;
  ipa: string;
  partOfSpeech: string;
  topic: string;
  band: string;
  bandScore: number;
  difficulty: VocabularyDifficulty;
  status: VocabularyStatus;
  frequency: "Medium" | "High" | "Very High";
  frequencyScore: number;
  masteryScore: number;
  reviewCount: number;
  lastReviewedAt: string;
  primaryMeaning: string;
  secondaryMeaning: string;
  meaningVi: string;
  meaningEn: string;
  shortDefinition: string;
  examples: VocabularyExample[];
  synonyms: string[];
  antonyms: string[];
  collocations: string[];
  ieltsUsage: string;
  relatedForms: RelatedForm[];
};

export const mockVocabulary: MockVocabulary[] = [
  {
    id: "sustainable",
    word: "Sustainable",
    slug: "sustainable",
    ipa: "/səˈsteɪ.nə.bəl/",
    partOfSpeech: "adjective",
    topic: "Environment",
    band: "Band 7+",
    bandScore: 7,
    difficulty: "Advanced",
    status: "Mastered",
    frequency: "Very High",
    frequencyScore: 92,
    masteryScore: 92,
    reviewCount: 8,
    lastReviewedAt: "2 days ago",
    primaryMeaning:
      "Able to be maintained at a certain rate or level; conserving an ecological balance by avoiding depletion of natural resources.",
    secondaryMeaning:
      "Capable of being upheld or defended over a long period of time.",
    meaningVi: "bền vững; có thể duy trì lâu dài",
    meaningEn:
      "Able to continue over time without damaging resources or creating long-term harm.",
    shortDefinition:
      "Able to continue without damaging resources or the environment.",
    examples: [
      {
        sentence:
          "Governments should invest in sustainable transport to reduce urban pollution.",
        note: "Useful for environment and city-planning essays.",
      },
      {
        sentence:
          "A sustainable lifestyle often requires reducing waste and using renewable energy.",
        note: "Natural collocation for IELTS Speaking Part 3.",
      },
    ],
    synonyms: ["eco-friendly", "renewable", "viable", "long-term"],
    antonyms: ["unsustainable", "wasteful", "short-lived"],
    collocations: [
      "sustainable development",
      "sustainable energy",
      "sustainable transport",
      "sustainable lifestyle",
    ],
    ieltsUsage:
      "High frequency in IELTS Writing Task 2 for environment, urban development, energy, and consumer behaviour questions.",
    relatedForms: [
      { word: "Sustainability", partOfSpeech: "noun" },
      { word: "Sustainably", partOfSpeech: "adverb" },
      { word: "Unsustainable", partOfSpeech: "adjective" },
    ],
  },
  {
    id: "innovation",
    word: "Innovation",
    slug: "innovation",
    ipa: "/ˌɪn.əˈveɪ.ʃən/",
    partOfSpeech: "noun",
    topic: "Technology",
    band: "Band 6.5",
    bandScore: 6.5,
    difficulty: "Intermediate",
    status: "Review",
    frequency: "High",
    frequencyScore: 84,
    masteryScore: 76,
    reviewCount: 5,
    lastReviewedAt: "today",
    primaryMeaning: "A new idea, method, product, or process.",
    secondaryMeaning:
      "The act of introducing something original or more effective.",
    meaningVi: "sự đổi mới; sáng kiến",
    meaningEn: "A new idea or method that improves how something is done.",
    shortDefinition: "A new idea, method, or product.",
    examples: [
      {
        sentence:
          "Technological innovation can improve access to education in remote areas.",
        note: "Strong for technology and education topics.",
      },
    ],
    synonyms: ["invention", "breakthrough", "improvement"],
    antonyms: ["stagnation", "tradition"],
    collocations: [
      "technological innovation",
      "drive innovation",
      "innovation hub",
    ],
    ieltsUsage:
      "Use it when discussing technology, business, education reform, and economic growth.",
    relatedForms: [
      { word: "Innovate", partOfSpeech: "verb" },
      { word: "Innovative", partOfSpeech: "adjective" },
      { word: "Innovator", partOfSpeech: "noun" },
    ],
  },
  {
    id: "curriculum",
    word: "Curriculum",
    slug: "curriculum",
    ipa: "/kəˈrɪk.jə.ləm/",
    partOfSpeech: "noun",
    topic: "Education",
    band: "Band 5.5",
    bandScore: 5.5,
    difficulty: "Beginner",
    status: "Learning",
    frequency: "High",
    frequencyScore: 80,
    masteryScore: 48,
    reviewCount: 2,
    lastReviewedAt: "yesterday",
    primaryMeaning:
      "The subjects and learning content included in a course of study.",
    secondaryMeaning:
      "A planned set of lessons or academic experiences in a school.",
    meaningVi: "chương trình học",
    meaningEn: "The subjects and content taught in a school or course.",
    shortDefinition: "The subjects included in a course of study.",
    examples: [
      {
        sentence:
          "A modern curriculum should include digital literacy and critical thinking.",
        note: "Good for education reform essays.",
      },
    ],
    synonyms: ["syllabus", "course content", "programme"],
    antonyms: ["extracurricular activity"],
    collocations: [
      "school curriculum",
      "national curriculum",
      "modern curriculum",
    ],
    ieltsUsage:
      "Useful for IELTS topics about education systems, schools, exams, and skills.",
    relatedForms: [
      { word: "Curricular", partOfSpeech: "adjective" },
      { word: "Extracurricular", partOfSpeech: "adjective" },
    ],
  },
  {
    id: "resilience",
    word: "Resilience",
    slug: "resilience",
    ipa: "/rɪˈzɪl.i.əns/",
    partOfSpeech: "noun",
    topic: "Psychology",
    band: "Band 7",
    bandScore: 7,
    difficulty: "Advanced",
    status: "New",
    frequency: "Medium",
    frequencyScore: 66,
    masteryScore: 12,
    reviewCount: 0,
    lastReviewedAt: "not reviewed",
    primaryMeaning:
      "The ability to recover quickly after difficulties or stress.",
    secondaryMeaning:
      "The capacity of a system or person to adapt to pressure.",
    meaningVi: "khả năng phục hồi; sự kiên cường",
    meaningEn: "The ability to recover from problems or adapt to challenges.",
    shortDefinition: "The ability to recover from difficulty.",
    examples: [
      {
        sentence:
          "Students need resilience to cope with academic pressure and failure.",
        note: "Useful for education and psychology questions.",
      },
    ],
    synonyms: ["adaptability", "toughness", "endurance"],
    antonyms: ["fragility", "weakness"],
    collocations: [
      "build resilience",
      "emotional resilience",
      "community resilience",
    ],
    ieltsUsage:
      "Strong word for mental health, education, society, and disaster recovery topics.",
    relatedForms: [
      { word: "Resilient", partOfSpeech: "adjective" },
      { word: "Resiliently", partOfSpeech: "adverb" },
    ],
  },
  {
    id: "urbanisation",
    word: "Urbanisation",
    slug: "urbanisation",
    ipa: "/ˌɜː.bən.aɪˈzeɪ.ʃən/",
    partOfSpeech: "noun",
    topic: "Society",
    band: "Band 7",
    bandScore: 7,
    difficulty: "Advanced",
    status: "Review",
    frequency: "High",
    frequencyScore: 82,
    masteryScore: 70,
    reviewCount: 4,
    lastReviewedAt: "3 days ago",
    primaryMeaning:
      "The process by which more people live in cities and towns.",
    secondaryMeaning: "The social and physical growth of urban areas.",
    meaningVi: "đô thị hóa",
    meaningEn: "The growth of cities as more people move to urban areas.",
    shortDefinition: "The growth of cities and urban populations.",
    examples: [
      {
        sentence:
          "Rapid urbanisation can place pressure on housing and public transport.",
        note: "Common in city and society essays.",
      },
    ],
    synonyms: ["city growth", "urban growth"],
    antonyms: ["ruralisation", "depopulation"],
    collocations: [
      "rapid urbanisation",
      "urbanisation rate",
      "urbanisation trend",
    ],
    ieltsUsage:
      "Useful for housing, transport, public services, environment, and society topics.",
    relatedForms: [
      { word: "Urbanise", partOfSpeech: "verb" },
      { word: "Urban", partOfSpeech: "adjective" },
    ],
  },
  {
    id: "diagnosis",
    word: "Diagnosis",
    slug: "diagnosis",
    ipa: "/ˌdaɪ.əɡˈnəʊ.sɪs/",
    partOfSpeech: "noun",
    topic: "Health",
    band: "Band 6.5",
    bandScore: 6.5,
    difficulty: "Intermediate",
    status: "Mastered",
    frequency: "Medium",
    frequencyScore: 64,
    masteryScore: 88,
    reviewCount: 7,
    lastReviewedAt: "1 day ago",
    primaryMeaning:
      "The identification of an illness or problem by examining symptoms.",
    secondaryMeaning:
      "A professional judgement about the cause of a condition.",
    meaningVi: "sự chẩn đoán",
    meaningEn: "The process of identifying an illness or problem.",
    shortDefinition: "Identification of an illness or problem.",
    examples: [
      {
        sentence:
          "Early diagnosis can improve treatment outcomes and reduce healthcare costs.",
        note: "Good for healthcare policy topics.",
      },
    ],
    synonyms: ["identification", "assessment"],
    antonyms: ["misdiagnosis"],
    collocations: [
      "early diagnosis",
      "accurate diagnosis",
      "medical diagnosis",
    ],
    ieltsUsage:
      "Useful for health, public services, technology in medicine, and prevention topics.",
    relatedForms: [
      { word: "Diagnose", partOfSpeech: "verb" },
      { word: "Diagnostic", partOfSpeech: "adjective" },
    ],
  },
];

export function findVocabularyById(vocabularyId: string | undefined) {
  if (!vocabularyId) {
    return undefined;
  }

  const vocabulary = mockVocabulary.find(
    (item) => item.id === vocabularyId || item.slug === vocabularyId,
  );

  if (vocabulary) {
    return vocabulary;
  }

  return mockVocabulary?.[0];
}

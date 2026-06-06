export type FlashcardRating = "again" | "hard" | "good" | "easy";

export type MockFlashcard = {
  id: string;
  vocabularyId: string;
  word: string;
  ipa: string;
  partOfSpeech: string;
  meaningEn: string;
  meaningVi: string;
  exampleSentence: string;
  synonyms: string[];
  collocations: string[];
  ieltsUsage: string;
  difficulty: "Beginner" | "Intermediate" | "Advanced";
  band: string;
  topicTitle: string;
  status: "New" | "Learning" | "Review" | "Mastered";
};

export const lessonFlashcards: MockFlashcard[] = [
  {
    id: "sustainable",
    vocabularyId: "1",
    word: "Sustainable",
    ipa: "/səˈsteɪ.nə.bəl/",
    partOfSpeech: "adjective",
    meaningEn:
      "Able to continue over time without damaging resources or creating long-term harm.",
    meaningVi: "bền vững; có thể duy trì lâu dài",
    exampleSentence:
      "Governments should invest in sustainable transport to reduce urban pollution.",
    synonyms: ["renewable", "viable", "eco-friendly", "long-term"],
    collocations: ["sustainable development", "sustainable energy"],
    ieltsUsage: "Useful for IELTS Writing Task 2 environment essays.",
    difficulty: "Advanced",
    band: "Band 7",
    topicTitle: "Environment",
    status: "Review",
  },
  {
    id: "innovation",
    vocabularyId: "2",
    word: "Innovation",
    ipa: "/ˌɪn.əˈveɪ.ʃən/",
    partOfSpeech: "noun",
    meaningEn: "A new idea, method, product, or process.",
    meaningVi: "sự đổi mới; sáng kiến",
    exampleSentence:
      "Technological innovation can improve access to education in remote areas.",
    synonyms: ["breakthrough", "invention", "improvement"],
    collocations: ["technological innovation", "drive innovation"],
    ieltsUsage: "Useful for technology, business, and education reform topics.",
    difficulty: "Intermediate",
    band: "Band 6.5",
    topicTitle: "Technology",
    status: "Learning",
  },
  {
    id: "resilience",
    vocabularyId: "3",
    word: "Resilience",
    ipa: "/rɪˈzɪl.i.əns/",
    partOfSpeech: "noun",
    meaningEn: "The ability to recover quickly after difficulties or stress.",
    meaningVi: "khả năng phục hồi; sự kiên cường",
    exampleSentence:
      "Students need resilience to cope with academic pressure and failure.",
    synonyms: ["adaptability", "endurance", "toughness"],
    collocations: ["build resilience", "emotional resilience"],
    ieltsUsage: "Strong for psychology, education, and social pressure topics.",
    difficulty: "Advanced",
    band: "Band 7",
    topicTitle: "Psychology",
    status: "New",
  },
  {
    id: "curriculum",
    vocabularyId: "4",
    word: "Curriculum",
    ipa: "/kəˈrɪk.jə.ləm/",
    partOfSpeech: "noun",
    meaningEn: "The subjects and content taught in a school or course.",
    meaningVi: "chương trình học",
    exampleSentence:
      "A modern curriculum should include digital literacy and critical thinking.",
    synonyms: ["syllabus", "course content", "programme"],
    collocations: ["school curriculum", "national curriculum"],
    ieltsUsage: "Useful for education systems and school reform questions.",
    difficulty: "Beginner",
    band: "Band 5.5",
    topicTitle: "Education",
    status: "Learning",
  },
  {
    id: "urbanisation",
    vocabularyId: "5",
    word: "Urbanisation",
    ipa: "/ˌɜː.bən.aɪˈzeɪ.ʃən/",
    partOfSpeech: "noun",
    meaningEn: "The growth of cities as more people move to urban areas.",
    meaningVi: "đô thị hóa",
    exampleSentence:
      "Rapid urbanisation can place pressure on housing and public transport.",
    synonyms: ["city growth", "urban growth"],
    collocations: ["rapid urbanisation", "urbanisation trend"],
    ieltsUsage: "Common in city, housing, and society essays.",
    difficulty: "Advanced",
    band: "Band 7",
    topicTitle: "Society",
    status: "Review",
  },
  {
    id: "diagnosis",
    vocabularyId: "6",
    word: "Diagnosis",
    ipa: "/ˌdaɪ.əɡˈnəʊ.sɪs/",
    partOfSpeech: "noun",
    meaningEn: "The identification of an illness or problem.",
    meaningVi: "sự chẩn đoán",
    exampleSentence:
      "Early diagnosis can improve treatment outcomes for serious diseases.",
    synonyms: ["identification", "assessment", "evaluation"],
    collocations: ["early diagnosis", "medical diagnosis"],
    ieltsUsage: "Useful for health and public healthcare topics.",
    difficulty: "Intermediate",
    band: "Band 6",
    topicTitle: "Health",
    status: "New",
  },
  {
    id: "accountability",
    vocabularyId: "7",
    word: "Accountability",
    ipa: "/əˌkaʊn.təˈbɪl.ə.ti/",
    partOfSpeech: "noun",
    meaningEn: "Responsibility for actions and their results.",
    meaningVi: "trách nhiệm giải trình",
    exampleSentence:
      "Public officials should have greater accountability for spending decisions.",
    synonyms: ["responsibility", "liability", "answerability"],
    collocations: ["public accountability", "ensure accountability"],
    ieltsUsage: "Strong for government, leadership, and policy discussions.",
    difficulty: "Advanced",
    band: "Band 7.5",
    topicTitle: "Government",
    status: "Review",
  },
  {
    id: "proliferate",
    vocabularyId: "8",
    word: "Proliferate",
    ipa: "/prəˈlɪf.ər.eɪt/",
    partOfSpeech: "verb",
    meaningEn: "To increase rapidly in number or amount.",
    meaningVi: "tăng nhanh; sinh sôi nảy nở",
    exampleSentence:
      "Online learning platforms have proliferated since digital access improved.",
    synonyms: ["multiply", "spread", "increase rapidly"],
    collocations: ["proliferate rapidly", "proliferate online"],
    ieltsUsage: "Useful for technology, media, and social change topics.",
    difficulty: "Advanced",
    band: "Band 7.5",
    topicTitle: "Technology",
    status: "Learning",
  },
  {
    id: "mitigate",
    vocabularyId: "9",
    word: "Mitigate",
    ipa: "/ˈmɪt.ɪ.ɡeɪt/",
    partOfSpeech: "verb",
    meaningEn: "To make a problem or harmful effect less severe.",
    meaningVi: "giảm nhẹ; làm dịu bớt",
    exampleSentence:
      "Governments can mitigate traffic congestion by improving public transport.",
    synonyms: ["reduce", "alleviate", "lessen"],
    collocations: ["mitigate risk", "mitigate the impact"],
    ieltsUsage: "High-value verb for problem-solution essays.",
    difficulty: "Advanced",
    band: "Band 7",
    topicTitle: "Global Issues",
    status: "Review",
  },
  {
    id: "infrastructure",
    vocabularyId: "10",
    word: "Infrastructure",
    ipa: "/ˈɪn.frəˌstrʌk.tʃər/",
    partOfSpeech: "noun",
    meaningEn: "The basic systems and services a society or organization needs.",
    meaningVi: "cơ sở hạ tầng",
    exampleSentence:
      "Investment in infrastructure is essential for sustainable economic growth.",
    synonyms: ["facilities", "systems", "public works"],
    collocations: ["transport infrastructure", "digital infrastructure"],
    ieltsUsage: "Useful for cities, transport, economy, and development topics.",
    difficulty: "Intermediate",
    band: "Band 6.5",
    topicTitle: "Economics",
    status: "Learning",
  },
];

export const reviewFlashcards = lessonFlashcards.filter((card) =>
  ["sustainable", "urbanisation", "accountability", "mitigate", "innovation"].includes(
    card.id,
  ),
);

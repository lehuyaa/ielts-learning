# Seed Data - IELTS Vocabulary Platform

## 1. Seed Data Goals

Seed data should make the app usable immediately after setup.

Minimum V1 seed:

- 1 course
- 4 band levels
- 8 topics
- 16 lessons
- 100 vocabulary items for first MVP
- Quiz questions for each lesson
- Achievements
- Demo user

---

## 2. Course

```json
{
  "title": "IELTS Vocabulary Roadmap",
  "slug": "ielts-vocabulary-roadmap",
  "description": "A structured IELTS vocabulary path from Band 5.0 to Band 8.0.",
  "isPublished": true,
  "orderIndex": 1
}
```

---

## 3. Band Levels

```json
[
  {
    "bandScore": 5.0,
    "title": "Band 5.0 Foundation Vocabulary",
    "description": "Essential words for everyday IELTS topics.",
    "orderIndex": 1
  },
  {
    "bandScore": 6.0,
    "title": "Band 6.0 Core IELTS Vocabulary",
    "description": "Common IELTS words for clearer speaking and writing.",
    "orderIndex": 2
  },
  {
    "bandScore": 7.0,
    "title": "Band 7.0 Academic Vocabulary",
    "description": "Stronger academic vocabulary for higher band answers.",
    "orderIndex": 3
  },
  {
    "bandScore": 8.0,
    "title": "Band 8.0 Advanced Vocabulary",
    "description": "Advanced collocations and precise academic language.",
    "orderIndex": 4
  }
]
```

---

## 4. Topics

Recommended topics:

```json
[
  "Education",
  "Technology",
  "Environment",
  "Health",
  "Crime",
  "Government",
  "Work",
  "Culture",
  "Travel",
  "Society"
]
```

For MVP, use 5 topics first:

```json
[
  "Education",
  "Technology",
  "Environment",
  "Health",
  "Work"
]
```

---

## 5. Lesson Structure

Each topic should have 2-5 lessons.

Example:

```json
{
  "topic": "Education",
  "lessons": [
    {
      "title": "Education Basics",
      "slug": "education-basics",
      "estimatedMinutes": 10,
      "requiredScore": 80
    },
    {
      "title": "Academic Education Vocabulary",
      "slug": "academic-education-vocabulary",
      "estimatedMinutes": 12,
      "requiredScore": 80
    }
  ]
}
```

---

## 6. Sample Vocabulary

### Education

```json
[
  {
    "word": "curriculum",
    "slug": "curriculum",
    "ipa": "/kəˈrɪkjələm/",
    "partOfSpeech": "noun",
    "meaningVi": "chương trình học",
    "meaningEn": "the subjects included in a course of study",
    "exampleSentence": "Schools should update their curriculum to include digital skills.",
    "exampleMeaningVi": "Các trường nên cập nhật chương trình học để bao gồm kỹ năng số.",
    "synonyms": ["syllabus", "course content"],
    "antonyms": [],
    "collocations": ["school curriculum", "national curriculum", "curriculum reform"],
    "ieltsUsage": "Useful for Education essays about school reform.",
    "difficulty": "MEDIUM"
  },
  {
    "word": "literacy",
    "slug": "literacy",
    "ipa": "/ˈlɪtərəsi/",
    "partOfSpeech": "noun",
    "meaningVi": "khả năng đọc viết",
    "meaningEn": "the ability to read and write",
    "exampleSentence": "Improving literacy rates should be a priority for developing countries.",
    "exampleMeaningVi": "Cải thiện tỷ lệ biết chữ nên là ưu tiên ở các nước đang phát triển.",
    "synonyms": ["reading ability", "basic education"],
    "antonyms": ["illiteracy"],
    "collocations": ["literacy rate", "digital literacy", "financial literacy"],
    "ieltsUsage": "Useful for Education, Society, and Technology topics.",
    "difficulty": "MEDIUM"
  },
  {
    "word": "tuition",
    "slug": "tuition",
    "ipa": "/tjuˈɪʃən/",
    "partOfSpeech": "noun",
    "meaningVi": "học phí / sự dạy kèm",
    "meaningEn": "money paid for education or instruction",
    "exampleSentence": "High tuition fees can prevent talented students from attending university.",
    "exampleMeaningVi": "Học phí cao có thể khiến sinh viên tài năng không thể học đại học.",
    "synonyms": ["school fee", "instruction"],
    "antonyms": [],
    "collocations": ["tuition fees", "private tuition"],
    "ieltsUsage": "Useful for essays about university education and inequality.",
    "difficulty": "EASY"
  }
]
```

### Environment

```json
[
  {
    "word": "sustainable",
    "slug": "sustainable",
    "ipa": "/səˈsteɪnəbl/",
    "partOfSpeech": "adjective",
    "meaningVi": "bền vững",
    "meaningEn": "able to continue over a period of time without damaging the environment",
    "exampleSentence": "Sustainable development is essential for future generations.",
    "exampleMeaningVi": "Phát triển bền vững là điều thiết yếu cho các thế hệ tương lai.",
    "synonyms": ["eco-friendly", "long-term", "renewable"],
    "antonyms": ["unsustainable"],
    "collocations": ["sustainable development", "sustainable growth", "sustainable lifestyle"],
    "ieltsUsage": "Very useful for Environment, Economy, and Government topics.",
    "difficulty": "MEDIUM"
  },
  {
    "word": "biodiversity",
    "slug": "biodiversity",
    "ipa": "/ˌbaɪəʊdaɪˈvɜːsəti/",
    "partOfSpeech": "noun",
    "meaningVi": "đa dạng sinh học",
    "meaningEn": "the variety of living things in an area",
    "exampleSentence": "Deforestation poses a serious threat to biodiversity.",
    "exampleMeaningVi": "Phá rừng gây ra mối đe dọa nghiêm trọng đối với đa dạng sinh học.",
    "synonyms": ["biological diversity"],
    "antonyms": [],
    "collocations": ["protect biodiversity", "loss of biodiversity"],
    "ieltsUsage": "Useful for essays about wildlife, forests, and environmental protection.",
    "difficulty": "HARD"
  },
  {
    "word": "emissions",
    "slug": "emissions",
    "ipa": "/ɪˈmɪʃənz/",
    "partOfSpeech": "noun",
    "meaningVi": "khí thải",
    "meaningEn": "gases released into the air",
    "exampleSentence": "Governments should introduce stricter policies to reduce carbon emissions.",
    "exampleMeaningVi": "Chính phủ nên đưa ra chính sách nghiêm ngặt hơn để giảm khí thải carbon.",
    "synonyms": ["pollutants", "exhaust gases"],
    "antonyms": [],
    "collocations": ["carbon emissions", "greenhouse gas emissions", "reduce emissions"],
    "ieltsUsage": "Useful for climate change and transport topics.",
    "difficulty": "MEDIUM"
  }
]
```

### Technology

```json
[
  {
    "word": "innovation",
    "slug": "innovation",
    "ipa": "/ˌɪnəˈveɪʃən/",
    "partOfSpeech": "noun",
    "meaningVi": "sự đổi mới",
    "meaningEn": "a new idea, method, or invention",
    "exampleSentence": "Technological innovation has transformed the way people communicate.",
    "exampleMeaningVi": "Đổi mới công nghệ đã thay đổi cách con người giao tiếp.",
    "synonyms": ["invention", "advancement", "breakthrough"],
    "antonyms": ["tradition"],
    "collocations": ["technological innovation", "drive innovation", "innovation in education"],
    "ieltsUsage": "Useful for Technology, Work, and Education topics.",
    "difficulty": "MEDIUM"
  },
  {
    "word": "automation",
    "slug": "automation",
    "ipa": "/ˌɔːtəˈmeɪʃən/",
    "partOfSpeech": "noun",
    "meaningVi": "tự động hóa",
    "meaningEn": "the use of machines or technology to do work without humans",
    "exampleSentence": "Automation may replace some repetitive jobs in the future.",
    "exampleMeaningVi": "Tự động hóa có thể thay thế một số công việc lặp đi lặp lại trong tương lai.",
    "synonyms": ["mechanization"],
    "antonyms": ["manual labor"],
    "collocations": ["workplace automation", "automation technology"],
    "ieltsUsage": "Useful for essays about jobs, AI, and economic change.",
    "difficulty": "MEDIUM"
  }
]
```

### Health

```json
[
  {
    "word": "sedentary",
    "slug": "sedentary",
    "ipa": "/ˈsedntri/",
    "partOfSpeech": "adjective",
    "meaningVi": "ít vận động",
    "meaningEn": "involving a lot of sitting and little physical activity",
    "exampleSentence": "A sedentary lifestyle can increase the risk of obesity and heart disease.",
    "exampleMeaningVi": "Lối sống ít vận động có thể làm tăng nguy cơ béo phì và bệnh tim.",
    "synonyms": ["inactive"],
    "antonyms": ["active"],
    "collocations": ["sedentary lifestyle", "sedentary habits"],
    "ieltsUsage": "Useful for Health and Work topics.",
    "difficulty": "HARD"
  },
  {
    "word": "well-being",
    "slug": "well-being",
    "ipa": "/ˌwel ˈbiːɪŋ/",
    "partOfSpeech": "noun",
    "meaningVi": "sức khỏe thể chất và tinh thần",
    "meaningEn": "the state of being healthy and happy",
    "exampleSentence": "Regular exercise has a positive impact on mental well-being.",
    "exampleMeaningVi": "Tập thể dục thường xuyên có tác động tích cực đến sức khỏe tinh thần.",
    "synonyms": ["health", "welfare"],
    "antonyms": [],
    "collocations": ["mental well-being", "emotional well-being", "improve well-being"],
    "ieltsUsage": "Useful for Health, Work-life balance, and Society topics.",
    "difficulty": "EASY"
  }
]
```

### Work

```json
[
  {
    "word": "productivity",
    "slug": "productivity",
    "ipa": "/ˌprɒdʌkˈtɪvəti/",
    "partOfSpeech": "noun",
    "meaningVi": "năng suất",
    "meaningEn": "the rate at which work is completed or goods are produced",
    "exampleSentence": "Flexible working hours can improve employee productivity.",
    "exampleMeaningVi": "Giờ làm linh hoạt có thể cải thiện năng suất của nhân viên.",
    "synonyms": ["efficiency", "output"],
    "antonyms": ["inefficiency"],
    "collocations": ["increase productivity", "employee productivity", "workplace productivity"],
    "ieltsUsage": "Useful for Work and Technology topics.",
    "difficulty": "MEDIUM"
  },
  {
    "word": "workforce",
    "slug": "workforce",
    "ipa": "/ˈwɜːkfɔːs/",
    "partOfSpeech": "noun",
    "meaningVi": "lực lượng lao động",
    "meaningEn": "all the people who work in a company, industry, or country",
    "exampleSentence": "A skilled workforce is essential for economic growth.",
    "exampleMeaningVi": "Một lực lượng lao động có kỹ năng là điều cần thiết cho tăng trưởng kinh tế.",
    "synonyms": ["employees", "labor force"],
    "antonyms": [],
    "collocations": ["skilled workforce", "global workforce", "workforce development"],
    "ieltsUsage": "Useful for Work, Economy, and Education topics.",
    "difficulty": "MEDIUM"
  }
]
```

---

## 7. Quiz Seed Example

For word:

```txt
sustainable
```

Question:

```json
{
  "type": "MEANING_CHOICE",
  "question": "What does 'sustainable' mean?",
  "options": [
    {
      "content": "Able to continue long-term without causing serious damage",
      "isCorrect": true
    },
    {
      "content": "Very expensive and difficult to buy",
      "isCorrect": false
    },
    {
      "content": "Related to ancient history",
      "isCorrect": false
    },
    {
      "content": "Only used by governments",
      "isCorrect": false
    }
  ],
  "explanation": "Sustainable means something can continue over time without causing serious damage, especially to the environment."
}
```

---

## 8. Achievement Seed

```json
[
  {
    "code": "FIRST_LESSON",
    "title": "First Lesson",
    "description": "Complete your first IELTS vocabulary lesson.",
    "icon": "BookOpen",
    "xpReward": 50
  },
  {
    "code": "SEVEN_DAY_STREAK",
    "title": "7 Day Streak",
    "description": "Study for 7 days in a row.",
    "icon": "Flame",
    "xpReward": 100
  },
  {
    "code": "HUNDRED_WORDS",
    "title": "100 Words Learned",
    "description": "Learn 100 IELTS vocabulary words.",
    "icon": "Trophy",
    "xpReward": 150
  },
  {
    "code": "EDUCATION_MASTER",
    "title": "Education Master",
    "description": "Complete all Education lessons.",
    "icon": "GraduationCap",
    "xpReward": 200
  }
]
```

---

## 9. Demo User

```json
{
  "email": "demo@example.com",
  "password": "password123",
  "name": "Demo User",
  "targetBand": 7.0
}
```

---

## 10. Go Seed Script Rules

Create:

```txt
backend/seeds/seed.go
```

The seed script should:

1. Connect to MySQL.
2. AutoMigrate models.
3. Create course.
4. Create band levels.
5. Create topics.
6. Create lessons.
7. Create vocabulary.
8. Attach vocabulary to lessons.
9. Create quiz questions and options.
10. Create achievements.
11. Create demo user.

For MVP, seed script may use `FirstOrCreate` to avoid duplicates.

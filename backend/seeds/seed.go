package seeds

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ielts-learning/backend/internal/models"
)

type topicSeed struct {
	Title       string
	Slug        string
	Description string
	Icon        string
	Emoji       string
	Color       string
	BandScore   float64
}

type vocabularySeed struct {
	Word            string
	IPA             string
	PartOfSpeech    string
	MeaningVI       string
	MeaningEN       string
	ExampleSentence string
	Difficulty      models.DifficultyLevel
	TargetBand      float64
	Synonyms        []string
	Collocations    []string
}

func Run(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		course, err := seedCourse(tx)
		if err != nil {
			return err
		}

		bands, err := seedBandLevels(tx, course.ID)
		if err != nil {
			return err
		}

		topics, err := seedTopics(tx, bands)
		if err != nil {
			return err
		}

		lessons, err := seedLessons(tx, topics)
		if err != nil {
			return err
		}

		vocabularies, err := seedVocabularies(tx, lessons)
		if err != nil {
			return err
		}

		if err := seedQuiz(tx, lessons, vocabularies); err != nil {
			return err
		}

		if err := seedAchievements(tx); err != nil {
			return err
		}

		if err := seedDemoUser(tx); err != nil {
			return err
		}

		return updateCourseTotals(tx, course.ID)
	})
}

func seedCourse(db *gorm.DB) (models.Course, error) {
	course := models.Course{
		Title:        "IELTS Vocabulary Roadmap",
		Description:  "A structured IELTS vocabulary path from foundation words to advanced academic language.",
		Slug:         "ielts-vocabulary-roadmap",
		IsPublished:  true,
		OrderIndex:   1,
		BandMin:      5.0,
		BandMax:      8.0,
		TotalWords:   50,
		TotalLessons: 5,
		TotalTopics:  5,
	}

	if err := db.Where("slug = ?", course.Slug).Assign(course).FirstOrCreate(&course).Error; err != nil {
		return models.Course{}, fmt.Errorf("seed course: %w", err)
	}

	return course, nil
}

func seedBandLevels(db *gorm.DB, courseID uint) (map[float64]models.BandLevel, error) {
	bandSeeds := []struct {
		BandScore   float64
		MinScore    float64
		MaxScore    float64
		Title       string
		Description string
		StatusLabel string
	}{
		{5.0, 5.0, 5.5, "Band 5.0 Foundation", "Build core IELTS topic vocabulary for clear everyday academic answers.", "Start here"},
		{6.0, 6.0, 6.5, "Band 6.0 Developing", "Expand range and accuracy with common IELTS academic themes.", "Recommended"},
		{7.0, 7.0, 7.5, "Band 7.0 Confident", "Use precise topic language for stronger writing and speaking responses.", "Popular"},
		{8.0, 8.0, 8.5, "Band 8.0 Advanced", "Master nuanced academic vocabulary and high-scoring collocations.", "Advanced"},
	}

	bands := make(map[float64]models.BandLevel, len(bandSeeds))
	for index, seed := range bandSeeds {
		minScore := seed.MinScore
		maxScore := seed.MaxScore
		statusLabel := seed.StatusLabel
		band := models.BandLevel{
			CourseID:    courseID,
			BandScore:   seed.BandScore,
			MinScore:    &minScore,
			MaxScore:    &maxScore,
			Title:       seed.Title,
			Description: seed.Description,
			StatusLabel: &statusLabel,
			OrderIndex:  index + 1,
		}

		if err := db.Where("course_id = ? AND band_score = ?", courseID, seed.BandScore).Assign(band).FirstOrCreate(&band).Error; err != nil {
			return nil, fmt.Errorf("seed band %.1f: %w", seed.BandScore, err)
		}
		bands[seed.BandScore] = band
	}

	return bands, nil
}

func seedTopics(db *gorm.DB, bands map[float64]models.BandLevel) (map[string]models.Topic, error) {
	topicSeeds := []topicSeed{
		{"Education", "education", "Vocabulary for schools, universities, learning methods, and academic success.", "graduation-cap", "ED", "indigo", 5.0},
		{"Technology", "technology", "Words for digital tools, innovation, automation, and online life.", "cpu", "TC", "blue", 6.0},
		{"Environment", "environment", "Language for climate, conservation, pollution, and sustainability.", "leaf", "EV", "emerald", 6.0},
		{"Health", "health", "Vocabulary for wellbeing, healthcare, lifestyle, and public health.", "heart-pulse", "HL", "rose", 7.0},
		{"Work", "work", "Terms for employment, productivity, career growth, and workplace change.", "briefcase", "WK", "amber", 7.0},
	}

	topics := make(map[string]models.Topic, len(topicSeeds))
	for index, seed := range topicSeeds {
		band := bands[seed.BandScore]
		topic := models.Topic{
			BandLevelID: band.ID,
			Title:       seed.Title,
			Slug:        seed.Slug,
			Description: seed.Description,
			Icon:        seed.Icon,
			Emoji:       seed.Emoji,
			Color:       seed.Color,
			OrderIndex:  index + 1,
		}

		if err := db.Where("band_level_id = ? AND slug = ?", band.ID, seed.Slug).Assign(topic).FirstOrCreate(&topic).Error; err != nil {
			return nil, fmt.Errorf("seed topic %s: %w", seed.Slug, err)
		}
		topics[seed.Slug] = topic
	}

	return topics, nil
}

func seedLessons(db *gorm.DB, topics map[string]models.Topic) (map[string]models.Lesson, error) {
	type lessonSeed struct {
		TopicSlug   string
		Title       string
		Slug        string
		Description string
		BandMin     float64
		BandMax     float64
	}

	lessonSeeds := []lessonSeed{
		{"education", "Education Foundations", "education-foundations", "Learn essential words for schools, study habits, and academic progress.", 5.0, 5.5},
		{"technology", "Digital Society", "digital-society", "Practice vocabulary about technology, innovation, and online behaviour.", 6.0, 6.5},
		{"environment", "Sustainability and Climate", "sustainability-and-climate", "Build language for environmental problems and practical solutions.", 6.0, 6.5},
		{"health", "Public Health and Lifestyle", "public-health-and-lifestyle", "Learn precise health vocabulary for IELTS speaking and writing.", 7.0, 7.5},
		{"work", "Careers and Productivity", "careers-and-productivity", "Use career and workplace vocabulary with confident collocations.", 7.0, 7.5},
	}

	lessons := make(map[string]models.Lesson, len(lessonSeeds))
	for index, seed := range lessonSeeds {
		topic := topics[seed.TopicSlug]
		bandMin := seed.BandMin
		bandMax := seed.BandMax
		timeLimit := 600
		lesson := models.Lesson{
			TopicID:              topic.ID,
			Title:                seed.Title,
			Slug:                 seed.Slug,
			Description:          seed.Description,
			RequiredScore:        80,
			EstimatedMinutes:     12,
			BandMin:              &bandMin,
			BandMax:              &bandMax,
			XPReward:             100,
			QuizTimeLimitSeconds: &timeLimit,
			OrderIndex:           index + 1,
			IsPublished:          true,
		}

		if err := db.Where("topic_id = ? AND slug = ?", topic.ID, seed.Slug).Assign(lesson).FirstOrCreate(&lesson).Error; err != nil {
			return nil, fmt.Errorf("seed lesson %s: %w", seed.Slug, err)
		}
		lessons[seed.Slug] = lesson
	}

	return lessons, nil
}

func seedVocabularies(db *gorm.DB, lessons map[string]models.Lesson) (map[string]models.Vocabulary, error) {
	vocabByLesson := map[string][]vocabularySeed{
		"education-foundations": {
			vocab("curriculum", "/kəˈrɪkjələm/", "noun", "chương trình học", "The subjects and content taught in a course.", "A balanced curriculum helps students develop practical and academic skills.", models.DifficultyIntermediate, 5.5, []string{"syllabus", "program"}, []string{"school curriculum", "national curriculum"}),
			vocab("literacy", "/ˈlɪtərəsi/", "noun", "khả năng đọc viết", "The ability to read and write effectively.", "Improving literacy should be a priority in rural education.", models.DifficultyBeginner, 5.0, []string{"reading ability"}, []string{"digital literacy", "literacy rate"}),
			vocab("discipline", "/ˈdɪsəplɪn/", "noun", "kỷ luật", "Controlled behaviour or a field of study.", "Classroom discipline creates a better learning environment.", models.DifficultyIntermediate, 5.5, []string{"order", "self-control"}, []string{"strict discipline", "academic discipline"}),
			vocab("assessment", "/əˈsesmənt/", "noun", "sự đánh giá", "A method of judging progress or ability.", "Continuous assessment can reduce exam pressure.", models.DifficultyIntermediate, 6.0, []string{"evaluation", "appraisal"}, []string{"formal assessment", "assessment criteria"}),
			vocab("scholarship", "/ˈskɒlərʃɪp/", "noun", "học bổng", "Financial support awarded for study.", "Scholarships make higher education more accessible.", models.DifficultyBeginner, 5.0, []string{"grant", "funding"}, []string{"win a scholarship", "full scholarship"}),
			vocab("tuition", "/tjuˈɪʃən/", "noun", "học phí", "Money paid for instruction.", "High tuition fees discourage some students from attending university.", models.DifficultyBeginner, 5.0, []string{"fees"}, []string{"tuition fees", "private tuition"}),
			vocab("interactive", "/ˌɪntərˈæktɪv/", "adjective", "có tính tương tác", "Involving active participation.", "Interactive lessons are more engaging than lectures.", models.DifficultyIntermediate, 5.5, []string{"participatory"}, []string{"interactive learning", "interactive platform"}),
			vocab("competence", "/ˈkɒmpɪtəns/", "noun", "năng lực", "The ability to do something well.", "Language competence improves with regular practice.", models.DifficultyIntermediate, 6.0, []string{"ability", "proficiency"}, []string{"professional competence", "communicative competence"}),
			vocab("mentor", "/ˈmentɔːr/", "noun", "người cố vấn", "An experienced person who guides someone.", "A mentor can help learners set realistic goals.", models.DifficultyBeginner, 5.0, []string{"coach", "advisor"}, []string{"career mentor", "experienced mentor"}),
			vocab("attendance", "/əˈtendəns/", "noun", "sự tham dự", "Being present at a class or event.", "Poor attendance often affects academic performance.", models.DifficultyBeginner, 5.0, []string{"presence"}, []string{"school attendance", "attendance record"}),
		},
		"digital-society": {
			vocab("innovation", "/ˌɪnəˈveɪʃən/", "noun", "sự đổi mới", "A new idea, method, or product.", "Technological innovation can improve public services.", models.DifficultyIntermediate, 6.0, []string{"invention", "creativity"}, []string{"drive innovation", "technological innovation"}),
			vocab("automation", "/ˌɔːtəˈmeɪʃən/", "noun", "tự động hóa", "The use of machines or software to do work.", "Automation may replace repetitive jobs.", models.DifficultyIntermediate, 6.5, []string{"mechanisation"}, []string{"industrial automation", "automation technology"}),
			vocab("algorithm", "/ˈælɡərɪðəm/", "noun", "thuật toán", "A set of rules used by a computer.", "Search engines use algorithms to rank information.", models.DifficultyAdvanced, 7.0, []string{"procedure"}, []string{"recommendation algorithm", "complex algorithm"}),
			vocab("privacy", "/ˈprɪvəsi/", "noun", "quyền riêng tư", "The right to keep personal information secret.", "Online privacy is a major concern for users.", models.DifficultyIntermediate, 6.0, []string{"confidentiality"}, []string{"protect privacy", "privacy policy"}),
			vocab("device", "/dɪˈvaɪs/", "noun", "thiết bị", "A machine or tool for a specific purpose.", "Mobile devices have changed how people learn.", models.DifficultyBeginner, 5.0, []string{"gadget", "tool"}, []string{"digital device", "mobile device"}),
			vocab("platform", "/ˈplætfɔːrm/", "noun", "nền tảng", "A digital service or system.", "Online platforms provide flexible learning opportunities.", models.DifficultyIntermediate, 6.0, []string{"system", "service"}, []string{"learning platform", "social platform"}),
			vocab("cybersecurity", "/ˌsaɪbəsɪˈkjʊərəti/", "noun", "an ninh mạng", "Protection of computer systems and data.", "Companies must invest in cybersecurity.", models.DifficultyAdvanced, 7.0, []string{"digital security"}, []string{"cybersecurity threat", "improve cybersecurity"}),
			vocab("accessible", "/əkˈsesəbl/", "adjective", "dễ tiếp cận", "Easy to reach, use, or understand.", "Technology makes education more accessible.", models.DifficultyIntermediate, 6.0, []string{"available", "reachable"}, []string{"accessible information", "accessible design"}),
			vocab("disruptive", "/dɪsˈrʌptɪv/", "adjective", "mang tính đột phá", "Causing major change to an industry or system.", "Disruptive technology can create new markets.", models.DifficultyAdvanced, 7.0, []string{"transformative"}, []string{"disruptive innovation", "disruptive change"}),
			vocab("data", "/ˈdeɪtə/", "noun", "dữ liệu", "Information collected for analysis.", "Reliable data helps governments make better decisions.", models.DifficultyBeginner, 5.0, []string{"information"}, []string{"collect data", "data analysis"}),
		},
		"sustainability-and-climate": {
			vocab("sustainable", "/səˈsteɪnəbl/", "adjective", "bền vững", "Able to continue without damaging resources.", "Cities need sustainable transport systems.", models.DifficultyIntermediate, 6.5, []string{"eco-friendly", "renewable"}, []string{"sustainable development", "sustainable energy"}),
			vocab("emission", "/ɪˈmɪʃən/", "noun", "khí thải", "Gas or substance released into the air.", "Vehicle emissions contribute to air pollution.", models.DifficultyIntermediate, 6.0, []string{"discharge", "release"}, []string{"carbon emissions", "reduce emissions"}),
			vocab("conservation", "/ˌkɒnsəˈveɪʃən/", "noun", "sự bảo tồn", "Protection of nature and resources.", "Wildlife conservation requires international cooperation.", models.DifficultyIntermediate, 6.5, []string{"protection", "preservation"}, []string{"nature conservation", "conservation effort"}),
			vocab("biodiversity", "/ˌbaɪəʊdaɪˈvɜːsəti/", "noun", "đa dạng sinh học", "The variety of living things in an area.", "Deforestation threatens biodiversity.", models.DifficultyAdvanced, 7.0, []string{"biological variety"}, []string{"protect biodiversity", "loss of biodiversity"}),
			vocab("renewable", "/rɪˈnjuːəbl/", "adjective", "có thể tái tạo", "Naturally replaced and not exhausted.", "Renewable energy can reduce dependence on fossil fuels.", models.DifficultyIntermediate, 6.0, []string{"reusable"}, []string{"renewable energy", "renewable resource"}),
			vocab("pollution", "/pəˈluːʃən/", "noun", "ô nhiễm", "Damage caused by harmful substances.", "Plastic pollution affects marine life.", models.DifficultyBeginner, 5.0, []string{"contamination"}, []string{"air pollution", "noise pollution"}),
			vocab("habitat", "/ˈhæbɪtæt/", "noun", "môi trường sống", "The natural home of an animal or plant.", "Urban expansion destroys natural habitats.", models.DifficultyIntermediate, 6.0, []string{"environment", "home"}, []string{"natural habitat", "habitat loss"}),
			vocab("scarce", "/skeəs/", "adjective", "khan hiếm", "Not enough for demand.", "Clean water is scarce in some regions.", models.DifficultyIntermediate, 6.0, []string{"limited", "rare"}, []string{"scarce resources", "increasingly scarce"}),
			vocab("mitigate", "/ˈmɪtɪɡeɪt/", "verb", "giảm nhẹ", "To make something less harmful.", "Governments should mitigate the effects of climate change.", models.DifficultyAdvanced, 7.0, []string{"reduce", "alleviate"}, []string{"mitigate risk", "mitigate impact"}),
			vocab("ecosystem", "/ˈiːkəʊsɪstəm/", "noun", "hệ sinh thái", "A community of organisms and their environment.", "A healthy ecosystem supports biodiversity.", models.DifficultyIntermediate, 6.0, []string{"ecological system"}, []string{"marine ecosystem", "fragile ecosystem"}),
		},
		"public-health-and-lifestyle": {
			vocab("prevention", "/prɪˈvenʃən/", "noun", "sự phòng ngừa", "Action taken to stop something happening.", "Disease prevention is cheaper than treatment.", models.DifficultyIntermediate, 6.0, []string{"avoidance"}, []string{"prevention strategy", "disease prevention"}),
			vocab("diagnosis", "/ˌdaɪəɡˈnəʊsɪs/", "noun", "sự chẩn đoán", "Identification of an illness or problem.", "Early diagnosis improves recovery rates.", models.DifficultyIntermediate, 6.5, []string{"identification"}, []string{"accurate diagnosis", "medical diagnosis"}),
			vocab("sedentary", "/ˈsedntri/", "adjective", "ít vận động", "Involving a lot of sitting and little exercise.", "A sedentary lifestyle increases health risks.", models.DifficultyAdvanced, 7.0, []string{"inactive"}, []string{"sedentary lifestyle", "sedentary job"}),
			vocab("nutrition", "/njuˈtrɪʃən/", "noun", "dinh dưỡng", "The food needed for health and growth.", "Good nutrition supports children's development.", models.DifficultyIntermediate, 6.0, []string{"diet", "nourishment"}, []string{"balanced nutrition", "poor nutrition"}),
			vocab("immunity", "/ɪˈmjuːnəti/", "noun", "khả năng miễn dịch", "Protection against disease.", "Regular sleep can support immunity.", models.DifficultyIntermediate, 6.0, []string{"resistance"}, []string{"natural immunity", "boost immunity"}),
			vocab("therapy", "/ˈθerəpi/", "noun", "liệu pháp", "Treatment for physical or mental illness.", "Therapy can help patients manage stress.", models.DifficultyBeginner, 5.5, []string{"treatment"}, []string{"physical therapy", "therapy session"}),
			vocab("chronic", "/ˈkrɒnɪk/", "adjective", "mãn tính", "Continuing for a long time.", "Chronic diseases require long-term care.", models.DifficultyAdvanced, 7.0, []string{"long-term", "persistent"}, []string{"chronic illness", "chronic pain"}),
			vocab("hygiene", "/ˈhaɪdʒiːn/", "noun", "vệ sinh", "Practices that keep people healthy and clean.", "Good hygiene reduces the spread of infection.", models.DifficultyBeginner, 5.0, []string{"cleanliness"}, []string{"personal hygiene", "hygiene standards"}),
			vocab("wellbeing", "/ˌwelˈbiːɪŋ/", "noun", "sức khỏe tinh thần và thể chất", "General health and happiness.", "Work-life balance improves wellbeing.", models.DifficultyIntermediate, 6.0, []string{"welfare", "health"}, []string{"mental wellbeing", "employee wellbeing"}),
			vocab("vaccination", "/ˌvæksɪˈneɪʃən/", "noun", "tiêm chủng", "Giving a vaccine to protect against disease.", "Vaccination programmes protect vulnerable groups.", models.DifficultyIntermediate, 6.5, []string{"immunisation"}, []string{"vaccination programme", "mass vaccination"}),
		},
		"careers-and-productivity": {
			vocab("collaboration", "/kəˌlæbəˈreɪʃən/", "noun", "sự hợp tác", "Working with others to achieve a goal.", "Collaboration improves problem-solving at work.", models.DifficultyIntermediate, 6.0, []string{"cooperation", "teamwork"}, []string{"close collaboration", "international collaboration"}),
			vocab("efficiency", "/ɪˈfɪʃənsi/", "noun", "hiệu quả", "Doing something with little waste.", "Flexible schedules can improve efficiency.", models.DifficultyIntermediate, 6.0, []string{"productivity"}, []string{"increase efficiency", "energy efficiency"}),
			vocab("entrepreneur", "/ˌɒntrəprəˈnɜːr/", "noun", "doanh nhân", "A person who starts and runs a business.", "Young entrepreneurs often bring fresh ideas.", models.DifficultyAdvanced, 7.0, []string{"business owner"}, []string{"successful entrepreneur", "social entrepreneur"}),
			vocab("incentive", "/ɪnˈsentɪv/", "noun", "động lực khuyến khích", "Something that encourages action.", "Financial incentives can motivate employees.", models.DifficultyIntermediate, 6.5, []string{"motivation", "reward"}, []string{"strong incentive", "tax incentive"}),
			vocab("resilience", "/rɪˈzɪliəns/", "noun", "khả năng phục hồi", "The ability to recover from difficulty.", "Resilience is important in a changing job market.", models.DifficultyAdvanced, 7.0, []string{"toughness", "adaptability"}, []string{"build resilience", "emotional resilience"}),
			vocab("deadline", "/ˈdedlaɪn/", "noun", "hạn chót", "The latest time by which work must be finished.", "Clear deadlines help teams manage tasks.", models.DifficultyBeginner, 5.0, []string{"time limit"}, []string{"meet a deadline", "tight deadline"}),
			vocab("workload", "/ˈwɜːkləʊd/", "noun", "khối lượng công việc", "The amount of work someone has.", "A heavy workload can cause stress.", models.DifficultyIntermediate, 6.0, []string{"tasks", "responsibilities"}, []string{"manage workload", "heavy workload"}),
			vocab("promotion", "/prəˈməʊʃən/", "noun", "sự thăng chức", "Moving to a higher position at work.", "Promotion opportunities can increase staff loyalty.", models.DifficultyBeginner, 5.0, []string{"advancement"}, []string{"get a promotion", "promotion prospects"}),
			vocab("adaptability", "/əˌdæptəˈbɪləti/", "noun", "khả năng thích nghi", "The ability to change for new conditions.", "Adaptability is essential in modern workplaces.", models.DifficultyAdvanced, 7.0, []string{"flexibility"}, []string{"show adaptability", "workplace adaptability"}),
			vocab("remuneration", "/rɪˌmjuːnəˈreɪʃən/", "noun", "thù lao", "Payment received for work.", "Fair remuneration helps companies retain skilled workers.", models.DifficultyAdvanced, 7.5, []string{"pay", "compensation"}, []string{"remuneration package", "fair remuneration"}),
		},
	}

	vocabularies := make(map[string]models.Vocabulary)
	for lessonSlug, seeds := range vocabByLesson {
		lesson := lessons[lessonSlug]
		for index, seed := range seeds {
			targetBand := seed.TargetBand
			vocabulary := models.Vocabulary{
				Word:             seed.Word,
				Slug:             fmt.Sprintf("%s-%s", lessonSlug, slugify(seed.Word)),
				IPA:              seed.IPA,
				PartOfSpeech:     seed.PartOfSpeech,
				MeaningVI:        seed.MeaningVI,
				MeaningEN:        seed.MeaningEN,
				ShortDefinition:  seed.MeaningEN,
				ExampleSentence:  seed.ExampleSentence,
				ExampleMeaningVI: "Câu ví dụ tiếng Việt sẽ được bổ sung trong nội dung bài học.",
				ExampleSource:    "Seed data",
				SynonymsJSON:     mustJSON(seed.Synonyms),
				AntonymsJSON:     mustJSON([]string{}),
				CollocationsJSON: mustJSON(seed.Collocations),
				IELTSUsage:       "Useful for IELTS Writing Task 2 and Speaking Part 3 topic responses.",
				Difficulty:       seed.Difficulty,
				TargetBand:       &targetBand,
			}

			if err := db.Where("slug = ?", vocabulary.Slug).Assign(vocabulary).FirstOrCreate(&vocabulary).Error; err != nil {
				return nil, fmt.Errorf("seed vocabulary %s: %w", seed.Word, err)
			}

			link := models.LessonVocabulary{
				LessonID:     lesson.ID,
				VocabularyID: vocabulary.ID,
				OrderIndex:   index + 1,
				IsRequired:   true,
			}
			if err := db.Where("lesson_id = ? AND vocabulary_id = ?", lesson.ID, vocabulary.ID).Assign(link).FirstOrCreate(&link).Error; err != nil {
				return nil, fmt.Errorf("seed lesson vocabulary %s: %w", seed.Word, err)
			}

			vocabularies[vocabulary.Slug] = vocabulary
		}
	}

	return vocabularies, nil
}

func seedQuiz(db *gorm.DB, lessons map[string]models.Lesson, vocabularies map[string]models.Vocabulary) error {
	for lessonSlug, lesson := range lessons {
		words, err := findLessonVocabularies(db, lesson.ID)
		if err != nil {
			return fmt.Errorf("load quiz words for %s: %w", lessonSlug, err)
		}
		if len(words) < 4 {
			continue
		}

		for questionIndex := 0; questionIndex < 2; questionIndex++ {
			word := words[questionIndex]
			topicID := lesson.TopicID
			vocabularyID := word.ID
			timeLimit := 45
			question := models.QuizQuestion{
				LessonID:         lesson.ID,
				TopicID:          &topicID,
				VocabularyID:     &vocabularyID,
				Type:             models.QuizQuestionMeaningChoice,
				Question:         fmt.Sprintf("What is the best meaning of \"%s\"?", word.Word),
				Explanation:      fmt.Sprintf("\"%s\" means: %s", word.Word, word.MeaningEN),
				Points:           20,
				TimeLimitSeconds: &timeLimit,
				OrderIndex:       questionIndex + 1,
			}

			if err := db.Where("lesson_id = ? AND order_index = ?", lesson.ID, question.OrderIndex).Assign(question).FirstOrCreate(&question).Error; err != nil {
				return fmt.Errorf("seed quiz question for %s: %w", word.Word, err)
			}

			options := []models.QuizOption{
				{QuestionID: question.ID, Label: "A", Content: word.MeaningEN, IsCorrect: true, OrderIndex: 1},
				{QuestionID: question.ID, Label: "B", Content: words[(questionIndex+2)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 2},
				{QuestionID: question.ID, Label: "C", Content: words[(questionIndex+4)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 3},
				{QuestionID: question.ID, Label: "D", Content: words[(questionIndex+6)%len(words)].MeaningEN, IsCorrect: false, OrderIndex: 4},
			}
			for _, option := range options {
				if err := db.Where("question_id = ? AND order_index = ?", question.ID, option.OrderIndex).Assign(option).FirstOrCreate(&option).Error; err != nil {
					return fmt.Errorf("seed quiz option for %s: %w", word.Word, err)
				}
			}
		}
	}

	_ = vocabularies
	return nil
}

func seedAchievements(db *gorm.DB) error {
	achievements := []models.Achievement{
		{Code: "FIRST_LESSON", Title: "First Lesson", Description: "Complete your first vocabulary lesson.", Icon: "book-open-check", Category: "lesson", RequirementType: "completed_lessons", RequirementValue: 1, XPReward: 50, SortOrder: 1, IsActive: true},
		{Code: "WORD_COLLECTOR_50", Title: "Word Collector", Description: "Learn 50 IELTS vocabulary items.", Icon: "layers", Category: "vocabulary", RequirementType: "learned_words", RequirementValue: 50, XPReward: 100, SortOrder: 2, IsActive: true},
		{Code: "QUIZ_STARTER", Title: "Quiz Starter", Description: "Pass your first lesson quiz.", Icon: "circle-check", Category: "quiz", RequirementType: "passed_quizzes", RequirementValue: 1, XPReward: 75, SortOrder: 3, IsActive: true},
		{Code: "STREAK_7", Title: "Seven Day Streak", Description: "Study for seven days in a row.", Icon: "flame", Category: "streak", RequirementType: "current_streak", RequirementValue: 7, XPReward: 120, SortOrder: 4, IsActive: true},
		{Code: "BAND_7_READY", Title: "Band 7 Ready", Description: "Master your first Band 7 vocabulary lesson.", Icon: "trophy", Category: "roadmap", RequirementType: "band_7_lessons", RequirementValue: 1, XPReward: 150, SortOrder: 5, IsActive: true},
	}

	for _, achievement := range achievements {
		if err := db.Where("code = ?", achievement.Code).Assign(achievement).FirstOrCreate(&achievement).Error; err != nil {
			return fmt.Errorf("seed achievement %s: %w", achievement.Code, err)
		}
	}

	return nil
}

func seedDemoUser(db *gorm.DB) error {
	username := "demo"
	currentBand := 5.5
	startingBand := 5.0
	recommendedBand := 6.0
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash demo password: %w", err)
	}

	user := models.User{
		Email:           "demo@example.com",
		Name:            "Demo Learner",
		Username:        &username,
		PasswordHash:    string(passwordHash),
		Role:            models.UserRoleUser,
		TargetBand:      7.0,
		CurrentBand:     &currentBand,
		StartingBand:    &startingBand,
		RecommendedBand: &recommendedBand,
		TotalXP:         0,
		Level:           1,
		CurrentStreak:   0,
		LongestStreak:   0,
		Timezone:        "Asia/Ho_Chi_Minh",
		Locale:          "en",
	}

	if err := db.Where("email = ?", user.Email).Assign(user).FirstOrCreate(&user).Error; err != nil {
		return fmt.Errorf("seed demo user: %w", err)
	}

	return nil
}

func updateCourseTotals(db *gorm.DB, courseID uint) error {
	var course models.Course
	if err := db.First(&course, courseID).Error; err != nil {
		return fmt.Errorf("load course totals: %w", err)
	}

	var totalTopics int64
	if err := db.Model(&models.Topic{}).
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalTopics).Error; err != nil {
		return fmt.Errorf("count topics: %w", err)
	}

	var totalLessons int64
	if err := db.Model(&models.Lesson{}).
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalLessons).Error; err != nil {
		return fmt.Errorf("count lessons: %w", err)
	}

	var totalWords int64
	if err := db.Model(&models.LessonVocabulary{}).
		Joins("JOIN lessons ON lessons.id = lesson_vocabularies.lesson_id").
		Joins("JOIN topics ON topics.id = lessons.topic_id").
		Joins("JOIN band_levels ON band_levels.id = topics.band_level_id").
		Where("band_levels.course_id = ?", courseID).
		Count(&totalWords).Error; err != nil {
		return fmt.Errorf("count words: %w", err)
	}

	course.TotalTopics = int(totalTopics)
	course.TotalLessons = int(totalLessons)
	course.TotalWords = int(totalWords)

	if err := db.Save(&course).Error; err != nil {
		return fmt.Errorf("update course totals: %w", err)
	}

	return nil
}

func findLessonVocabularies(db *gorm.DB, lessonID uint) ([]models.Vocabulary, error) {
	var words []models.Vocabulary
	err := db.Model(&models.Vocabulary{}).
		Joins("JOIN lesson_vocabularies ON lesson_vocabularies.vocabulary_id = vocabularies.id").
		Where("lesson_vocabularies.lesson_id = ?", lessonID).
		Order("lesson_vocabularies.order_index ASC").
		Find(&words).Error
	return words, err
}

func vocab(word string, ipa string, partOfSpeech string, meaningVI string, meaningEN string, example string, difficulty models.DifficultyLevel, targetBand float64, synonyms []string, collocations []string) vocabularySeed {
	return vocabularySeed{
		Word:            word,
		IPA:             ipa,
		PartOfSpeech:    partOfSpeech,
		MeaningVI:       meaningVI,
		MeaningEN:       meaningEN,
		ExampleSentence: example,
		Difficulty:      difficulty,
		TargetBand:      targetBand,
		Synonyms:        synonyms,
		Collocations:    collocations,
	}
}

func mustJSON(value []string) datatypes.JSON {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return datatypes.JSON(data)
}

func slugify(value string) string {
	slug := ""
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z':
			slug += string(char)
		case char >= 'A' && char <= 'Z':
			slug += string(char + 32)
		case char >= '0' && char <= '9':
			slug += string(char)
		case char == ' ' || char == '-' || char == '_':
			if len(slug) > 0 && slug[len(slug)-1] != '-' {
				slug += "-"
			}
		}
	}

	if len(slug) > 0 && slug[len(slug)-1] == '-' {
		return slug[:len(slug)-1]
	}

	return slug
}

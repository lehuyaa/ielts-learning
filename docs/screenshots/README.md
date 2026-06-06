# Screenshots - Approved Design References

The files in this folder are the visual source of truth for frontend UI implementation.

Latest approved Figma learning flow:

```txt
Roadmap
↓
Topic Detail
↓
Lesson Detail
↓
Flashcards
↓
Quiz
```

Canonical references:

```txt
roadmap.png
topic-detail.png
lesson-detail.png
vocabulary-detail.png
```

Current approved files in this repository:

```txt
Roadmap:
- roadmap1.png
- roadmap2.png

Topic Detail:
- topic-detail1.png
- topic-detail2.png
- topic-detail3.png

Lesson Detail:
- lesson-detail1.png
- lesson-detail2.png

Vocabulary Detail:
- vocabulary-detail.png
```

Navigation rules:

- Roadmap topic cards open Topic Detail.
- Topic Detail lesson cards open Lesson Detail.
- Lesson Detail launches Flashcards and Quiz.
- Lesson Detail "Full Detail" opens Vocabulary Detail.
- Vocabulary Detail at `/vocabulary/:vocabularyId` uses `vocabulary-detail.png` as the visual source of truth.

Implementation rules:

- Analyze the screenshot first.
- Do not redesign.
- Match layout closely.
- Match spacing hierarchy.
- Match typography hierarchy.
- Match component hierarchy.
- Match section order.
- Match card hierarchy and visual scale.
- Report visual differences before marking UI work complete.

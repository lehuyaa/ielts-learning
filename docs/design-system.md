# Design System - IELTS Vocabulary Platform

## 1. Design Direction

The UI should feel:

- Modern
- Premium
- Friendly
- Motivating
- Clean
- Educational
- Slightly gamified

Inspired by:

- Duolingo: streaks, progress, rewards
- Linear: clean SaaS layout
- Notion: calm typography and spacing
- Anki: focused flashcard learning

---

## 2. Frontend UI Stack

Use:

```txt
React.js
TypeScript
Vite
Tailwind CSS
shadcn/ui
Lucide React
```

---

## 3. Color Palette

### Primary

```txt
Indigo: #4F46E5
```

Use for:

- Primary buttons
- Links
- Active navigation
- Progress highlights
- Main CTA

### Secondary

```txt
Emerald: #10B981
```

Use for:

- Success states
- Completed lessons
- Correct answers
- Streak positive state

### Warning

```txt
Amber: #F59E0B
```

Use for:

- Hard difficulty
- Review due
- Warnings

### Error

```txt
Rose: #F43F5E
```

Use for:

- Wrong answers
- Form errors
- Again flashcard rating

### Neutral

```txt
Slate 50:  #F8FAFC
Slate 100: #F1F5F9
Slate 200: #E2E8F0
Slate 500: #64748B
Slate 700: #334155
Slate 900: #0F172A
```

---

## 4. Typography

Recommended font:

```txt
Inter
```

Fallback:

```css
font-family: Inter, system-ui, sans-serif;
```

### Type Scale

```txt
Display: 48px / 56px / 700
H1:      36px / 44px / 700
H2:      30px / 38px / 700
H3:      24px / 32px / 600
H4:      20px / 28px / 600
Body:    16px / 24px / 400
Small:   14px / 20px / 400
Tiny:    12px / 16px / 500
```

---

## 5. Border Radius

Use rounded UI.

```txt
Small: 8px
Medium: 12px
Large: 16px
XL: 24px
Full: 9999px
```

Default card radius:

```txt
16px
```

---

## 6. Spacing System

Use the shared spacing scale:

```txt
4px, 8px, 12px, 16px, 24px, 32px, 48px, 64px
```

Rules:

- Use values from the spacing scale.
- Do not use arbitrary spacing unless absolutely necessary.
- Label to input spacing: 8px.
- Input group to input group spacing: 24px.
- Page title/subtitle block to first form section: 48px.
- Card padding desktop: 32px.
- Card padding mobile: 24px.
- Section to section spacing: 64px.
- Button to previous form field spacing: 32px.

Common layout:

```txt
Page padding desktop: 32px
Page padding mobile: 16px
Card padding desktop: 32px
Card padding mobile: 24px
Section gap: 64px
Component gap: 16px
```

---

## 7. Components

### Button

Variants:

- Primary
- Secondary
- Outline
- Ghost
- Destructive
- Success

Primary button:

```txt
Background: Indigo
Text: White
Radius: 12px
Height: 44px
```

### Card

Default:

```txt
Background: White
Border: Slate 200
Radius: 16px
Padding: 24px
```

### Badge

Use for:

- Difficulty
- Status
- Band level
- Streak
- XP

Examples:

```txt
Easy: Emerald
Medium: Amber
Hard: Rose
Band 7.0: Indigo
Completed: Emerald
Locked: Slate
```

### Progress Bar

Use for:

- Band progress
- Lesson progress
- Quiz progress
- Flashcard progress

Default height:

```txt
8px
```

### Stat Card

Contains:

- Icon
- Label
- Value
- Optional trend

Example:

```txt
Words Learned Today
24
+8 from yesterday
```

### Empty State

Contains:

- Icon or illustration
- Title
- Description
- CTA

### Skeleton

Use for:

- Dashboard cards
- Vocabulary list
- Lesson cards
- Quiz loading

---

## 8. Page-Specific UI Guidelines

### Landing Page

Must feel conversion-oriented.

Hero:

- Big headline
- Short subheadline
- Two CTA buttons
- Product preview card

Use gradient background:

```txt
Indigo tint → white
```

### Dashboard

Should feel like a learning control center.

Important cards:

- Continue Learning
- Review Due
- Current Streak
- Target Band Progress

### Roadmap

Should feel gamified.

Visual style:

- Band sections
- Topic cards
- Lesson nodes
- Locked/completed states
- Path-like layout

### Flashcard Page

Should feel focused.

Rules:

- Minimal distractions
- Large centered card
- Clear progress at top
- Rating buttons at bottom
- Smooth flip animation optional

### Quiz Page

Should feel calm and exam-like.

Rules:

- One question at a time
- Clear options
- Progress indicator
- Submit button
- Result screen

---

## 9. Icons

Use Lucide React.

Recommended icons:

```txt
BookOpen
GraduationCap
Brain
Flame
Trophy
Target
Clock
CheckCircle
Lock
Star
Search
Bell
User
BarChart
Sparkles
```

---

## 10. Motion

Use subtle animations only.

Good:

- Card hover
- Progress bar transition
- Flashcard flip
- Badge pop on achievement
- Button tap scale

Avoid:

- Too many animations
- Slow transitions
- Distracting effects

Default transition:

```txt
150ms - 250ms
ease-out
```

---

## 11. Responsive Rules

### Desktop

- Sidebar layout
- Max content width around 1200px
- Grid cards

### Tablet

- Collapse some grids to 2 columns

### Mobile

- Bottom navigation or collapsible menu
- Single column layout
- Flashcard full width
- Quiz options stacked

---

## 12. Tone of Voice

The product copy should be:

- Encouraging
- Clear
- Short
- Motivational

Examples:

```txt
Welcome back, Alex.
You have 18 words to review today.
Great job! You completed Education Lesson 1.
Keep your streak alive.
```

Avoid:

```txt
You failed.
You are wrong.
```

Use:

```txt
Not quite. Review this word again.
```

---

## 13. Accessibility

1. Buttons must have clear labels.
2. Color should not be the only indicator.
3. Text contrast must be readable.
4. Quiz options should be keyboard accessible.
5. Use semantic HTML.
6. Add aria labels where needed.

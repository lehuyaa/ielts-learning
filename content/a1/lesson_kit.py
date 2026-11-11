"""Shared helpers for the A1 week builders (build_week1.py, build_week2.py…).

A week builder calls start_week(n, …), sets the current lesson with
lesson("l1") while it writes each lesson, and ends with finish(…), which adds
the AI role-plays, the optional writing tasks and the checkpoint pool, checks
vocabulary slugs and writes week-N.json and audio_manifest_week-N.tsv.
"""
import copy
import json
import random
import re
import sys

from writing_prompts import PROMPTS, writing_activity

AUDIO_ROOT = "/audio/a1"
WEEK = "w1"
WEEK_NUMBER = 1
k = ""             # key of the lesson being written ("l1"…, "rv")
audio_lines = {}   # audio key -> (text spoken, voice or "")
passages = []      # listening/reading texts (docs/foundation-spec.md §3.3)
DIALOGUE_VOICES = {}
READING_GENRES = {}

# Speaker voices: generate_lesson_audio.sh maps "female" / "male" to an
# installed macOS voice (FEMALE_VOICE / MALE_VOICE, or the best one found).
VOICE = {"F": "female", "M": "male"}


def start_week(number, dialogue_voices, reading_genres):
    """Resets the builder for week `number`."""
    global WEEK, WEEK_NUMBER
    WEEK, WEEK_NUMBER = f"w{number}", number
    audio_lines.clear()
    passages.clear()
    DIALOGUE_VOICES.clear()
    DIALOGUE_VOICES.update(dialogue_voices)
    READING_GENRES.clear()
    READING_GENRES.update(reading_genres)


def lesson(key):
    """Sets the lesson being written (used for audio and passage keys)."""
    global k
    k = key
    return key


def audio(key, text, voice=""):
    """Registers a sentence to synthesise and returns its public URL."""
    path = f"{WEEK}/{key}"
    audio_lines[path] = (text, voice)
    return f"{AUDIO_ROOT}/{path}.m4a"


# ---------------------------------------------------------------------------
# Exercise helpers (docs/foundation-spec.md §5)

def ex(type_, prompt, payload, answer, explanation="", tag=""):
    return {"type": type_, "prompt": prompt, "payload": payload, "answer": answer,
            "explanationVi": explanation, "skillTag": tag}


def mc(prompt, options, correct, explanation="", tag="", audio_url=None):
    ids = "abcd"
    payload = {"options": [{"id": ids[i], "text": text} for i, text in enumerate(options)]}
    if audio_url:
        payload["audioUrl"] = audio_url
    return ex("MULTIPLE_CHOICE", prompt, payload, {"correct": ids[correct]}, explanation, tag)


def minimal_pair(key, words, correct, explanation, tag):
    """words: [(word, ipa), ...]; the learner hears words[correct]."""
    ids = "abcd"
    url = audio(f"{key}-mp", words[correct][0])
    return ex("MINIMAL_PAIR", "Bạn nghe thấy từ nào?",
              {"audioUrl": url, "options": [{"id": ids[i], "word": w, "ipa": p} for i, (w, p) in enumerate(words)]},
              {"correct": ids[correct]}, explanation, tag)


def ipa_to_word(ipa, words, correct, explanation, tag):
    ids = "abcd"
    return ex("IPA_TO_WORD", "Phiên âm này là từ nào?",
              {"ipa": ipa, "options": [{"id": ids[i], "word": w} for i, w in enumerate(words)]},
              {"correct": ids[correct]}, explanation, tag)


def word_to_ipa(key, word, ipas, correct, explanation, tag):
    ids = "abcd"
    return ex("WORD_TO_IPA", f"Chọn phiên âm đúng của \"{word}\"",
              {"word": word, "audioUrl": audio(f"{key}-w2i", word),
               "options": [{"id": ids[i], "ipa": p} for i, p in enumerate(ipas)]},
              {"correct": ids[correct]}, explanation, tag)


def ending(key, word, correct, explanation, tag="sound:ending-s"):
    options = ["/s/", "/z/", "/ɪz/"] if correct in ("/s/", "/z/", "/ɪz/") else ["/t/", "/d/", "/ɪd/"]
    return ex("ENDING_SOUND", f"Đuôi trong \"{word}\" đọc là gì?",
              {"word": word, "audioUrl": audio(f"{key}-end", word), "options": options},
              {"correct": correct}, explanation, tag)


def stress(key, word, syllables, index, explanation, tag="stress"):
    return ex("WORD_STRESS", f"Trọng âm của \"{word}\" ở âm tiết nào?",
              {"word": word, "audioUrl": audio(f"{key}-stress", word), "syllables": syllables},
              {"stressIndex": index}, explanation, tag)


def fill(sentence, accepted, explanation, tag):
    return ex("FILL_BLANK", "Điền vào chỗ trống", {"sentence": sentence},
              {"accepted": accepted}, explanation, tag)


def reorder(sentence, explanation, tag, seed=0):
    tokens = re.findall(r"[\w'’-]+|[.,!?]", sentence)
    shuffled = tokens[:]
    rng = random.Random(sentence + str(seed))
    while shuffled == tokens:
        rng.shuffle(shuffled)
    return ex("REORDER", "Sắp xếp thành câu đúng", {"tokens": shuffled},
              {"accepted": [sentence]}, explanation, tag)


def error_correction(sentence, wrong_index, accepted, explanation, tag):
    tokens = re.findall(r"[\w'’-]+|[.,!?]", sentence)
    return ex("ERROR_CORRECTION", "Chạm vào từ sai và sửa lại", {"tokens": tokens},
              {"wrongIndex": wrong_index, "accepted": accepted}, explanation, tag)


def transform(sentence, instruction, accepted, explanation, tag):
    return ex("TRANSFORM", "Biến đổi câu", {"sentence": sentence, "instruction": instruction},
              {"accepted": accepted, "aiFallback": True}, explanation, tag)


def translate(vi, references, explanation, tag):
    return ex("TRANSLATE_VI_EN", "Dịch sang tiếng Anh", {"sourceVi": vi},
              {"reference": references, "aiGraded": True}, explanation, tag)


def true_false(statement, value, explanation="", tag="listening", audio_url=None):
    payload = {"statement": statement}
    if audio_url:
        payload["audioUrl"] = audio_url
    return ex("TRUE_FALSE", "Đúng hay sai?", payload, {"value": value}, explanation, tag)


def dictation(key, sentence, explanation="", tag="listening:dictation"):
    return ex("DICTATION", "Nghe và chép lại câu",
              {"audioUrl": audio(f"{key}-dict", sentence), "wordCount": len(sentence.split())},
              {"accepted": [sentence]}, explanation, tag)


def gap_fill(key, spoken, text, blanks, explanation="", tag="listening:gap-fill"):
    return ex("LISTEN_GAP_FILL", "Nghe và điền từ còn thiếu",
              {"audioUrl": audio(f"{key}-gap", spoken), "text": text},
              {"blanks": blanks}, explanation, tag)


def numbers(key, spoken, prompt, fmt, accepted, explanation="", tag="listening:numbers"):
    return ex("LISTEN_NUMBERS", prompt, {"audioUrl": audio(f"{key}-num", spoken), "format": fmt},
              {"accepted": accepted}, explanation, tag)


def speak(text, ipa=""):
    payload = {"text": text}
    if ipa:
        payload["ipa"] = ipa
    return ex("SPEAK", "Nghe mẫu rồi đọc to câu này", payload, {"minAccuracy": 70}, "", "speaking")


# ---------------------------------------------------------------------------
# Activity helpers

def pronunciation(title, instructions, phonemes, exercises):
    return {"slug": "sounds", "type": "PRONUNCIATION", "title": title, "instructionsVi": instructions,
            "content": {"phonemes": phonemes}, "exercises": exercises}


def vocabulary(title, slugs, exercises):
    return {"slug": "words", "type": "VOCABULARY", "title": title,
            "instructionsVi": "Chạm vào từng thẻ để xem nghĩa và nghe phát âm, sau đó làm bài kiểm tra nhanh.",
            "content": {"vocabularySlugs": slugs}, "exercises": exercises}


def grammar(title, instructions, point, exercises):
    return {"slug": "grammar", "type": "GRAMMAR", "title": title, "instructionsVi": instructions,
            "content": {"grammarPointSlug": point, "showTheory": True}, "exercises": exercises}


def listening(key, title, lines, exercises, extra_dictation=2):
    """lines: [(speaker, text), ...]. Each line is its own audio file (one
    passage segment), so the player can replay sentence by sentence. Adds
    dictation items from short dialogue lines."""
    genders = DIALOGUE_VOICES.get(key, "FM")
    segments = []
    for index, (speaker, text) in enumerate(lines):
        gender = genders[(ord(speaker) - ord("A")) % len(genders)]
        segments.append({"index": index, "speaker": speaker, "text": text,
                         "audioUrl": audio(f"{key}-s{index + 1}", text, VOICE[gender])})
    slug = f"a1-{WEEK}-{key}-dialogue"
    passages.append({"slug": slug, "genre": "DIALOGUE", "title": title,
                     "bodyText": "\n".join(f"{s}: {t}" for s, t in lines), "segments": segments,
                     "source": f"generated (A1 Week {WEEK_NUMBER})"})

    dictated = {e["answer"]["accepted"][0] for e in exercises if e["type"] == "DICTATION"}
    added = 0
    for segment in segments:
        words = segment["text"].split()
        if added >= extra_dictation or not 3 <= len(words) <= 9 or segment["text"] in dictated:
            continue
        exercises.append(ex("DICTATION", "Nghe và chép lại câu",
                            {"audioUrl": segment["audioUrl"], "wordCount": len(words)},
                            {"accepted": [segment["text"]]}, "", "listening:dictation"))
        added += 1

    return {"slug": "listening", "type": "LISTENING", "title": title,
            "instructionsVi": "Nghe hội thoại (có thể nghe chậm và nghe lại từng câu), rồi trả lời câu hỏi.",
            "passageSlug": slug, "content": {"speeds": [0.75, 1.0], "transcript": "AFTER_ANSWER"},
            "exercises": exercises}


def reading(title, text, exercises):
    key = k  # the lesson being written (set with lesson())
    slug = f"a1-{WEEK}-{key}-reading"
    passages.append({"slug": slug, "genre": READING_GENRES.get(key, "MESSAGE"), "title": title,
                     "bodyText": text, "source": f"generated (A1 Week {WEEK_NUMBER})"})
    return {"slug": "reading", "type": "READING", "title": title,
            "instructionsVi": "Đọc đoạn văn (chạm vào từ để xem nghĩa), rồi trả lời câu hỏi.", "isRequired": False,
            "passageSlug": slug, "content": {"tapToTranslate": True}, "exercises": exercises}


def speaking(sentences):
    return {"slug": "speaking", "type": "SPEAKING", "title": "Luyện nói",
            "instructionsVi": "Nghe câu mẫu, đọc to theo. Luyện nói với AI sẽ có ở bản sau.",
            "isRequired": False, "exercises": [speak(s) for s in sentences]}



# ---------------------------------------------------------------------------
# Week assembly

def finish(lessons, topic, role_plays, writing_keys=None):
    """Adds role-plays, writing tasks and the checkpoint pool, validates
    vocabulary slugs (--vocab rows.json) and writes the week's files.

    topic: the topic fields (slug, title, …) without "lessons".
    role_plays: lesson slug -> scenario slug.
    writing_keys: prompt key -> lesson slug (default: PROMPTS[key]["lesson"]).
    """
    # AI role-play (Tasks 67-70, D10): the read-aloud sentences become their
    # own activity and the lesson ends with a conversation with the AI.
    for item in lessons:
        scenario = role_plays.get(item["slug"])
        if not scenario:
            continue
        for activity in item["activities"]:
            if activity["type"] == "SPEAKING":
                activity["slug"] = "read-aloud"
                activity["title"] = "Đọc to theo mẫu"
                activity["instructionsVi"] = "Nghe câu mẫu rồi bấm nói và đọc to theo. AI chấm phát âm từng từ (cần độ chính xác từ 70)."
        item["activities"].append({
            "slug": "speaking", "type": "SPEAKING", "title": "Nói chuyện với AI",
            "instructionsVi": "Đóng vai trong tình huống của bài. AI nói chậm, câu ngắn; bí thì bấm Gợi ý. Kết thúc để xem báo cáo.",
            "scenarioSlug": scenario, "content": {"hintsAllowed": True},
            "isRequired": False, "passScore": 60, "xpReward": 20, "exercises": [],
        })

    # Optional writing tasks (Task D9), placed before the speaking activity.
    if writing_keys is None:
        writing_keys = {key: prompt.get("lesson") for key, prompt in PROMPTS.items()}
    for key, lesson_slug in writing_keys.items():
        for item in lessons:
            if item["slug"] == lesson_slug:
                activities = item["activities"]
                position = next((i for i, a in enumerate(activities) if a["type"] == "SPEAKING"), len(activities))
                activities.insert(position, writing_activity(key))

    # Checkpoint pool (Task 75): each attempt draws content.itemCount items
    # from the checkpoint's own exercises, so copy in the week's
    # self-contained items (grammar, dictation, numbers, sounds…). Reading /
    # listening comprehension questions need their passage on screen.
    checkpoint = next(a for item in lessons for a in item["activities"] if a["type"] == "CHECKPOINT")
    pool_types = set(checkpoint["content"]["fromActivityTypes"])
    seen_items = {json.dumps([e["type"], e["payload"]], sort_keys=True, ensure_ascii=False) for e in checkpoint["exercises"]}
    for item in lessons:
        for activity in item["activities"]:
            if activity is checkpoint or activity["type"] not in pool_types:
                continue
            for exercise in activity.get("exercises", []):
                if exercise.get("skillTag", "") in ("reading", "listening", ""):
                    continue
                key = json.dumps([exercise["type"], exercise["payload"]], sort_keys=True, ensure_ascii=False)
                if key in seen_items:
                    continue
                seen_items.add(key)
                checkpoint["exercises"].append(copy.deepcopy(exercise))

    bundle = {
        "version": 1,
        "courseSlug": "foundation-english",
        "cefrLevel": "A1",
        "passages": passages,
        "topics": [dict(topic, lessons=lessons)],
    }

    # Validate vocabulary slugs against the A1 vocabulary list when available.
    if "--vocab" in sys.argv:
        vocab_path = sys.argv[sys.argv.index("--vocab") + 1]
        known = {row["slug"] for row in json.load(open(vocab_path, encoding="utf-8"))}
        used = {slug for item in lessons for activity in item["activities"]
                for slug in activity.get("content", {}).get("vocabularySlugs", [])}
        used |= {slug for item in lessons for activity in item["activities"]
                 for slug in activity.get("content", {}).get("targetVocabulary", [])}
        missing = sorted(slug for slug in used if slug not in known)
        if missing:
            sys.exit("unknown vocabulary slugs: " + ", ".join(missing))

    with open(f"week-{WEEK_NUMBER}.json", "w", encoding="utf-8") as fh:
        json.dump(bundle, fh, ensure_ascii=False, indent=2)
    with open(f"audio_manifest_week-{WEEK_NUMBER}.tsv", "w", encoding="utf-8") as fh:
        for path, (text, voice) in sorted(audio_lines.items()):
            fh.write(f"{path}\t{text}\t{voice}\n" if voice else f"{path}\t{text}\n")

    count = sum(len(a.get("exercises", [])) for item in lessons for a in item["activities"])
    print(f"passages={len(passages)} lessons={len(lessons)} activities={sum(len(l['activities']) for l in lessons)} "
          f"exercises={count} audio={len(audio_lines)}")
    return bundle

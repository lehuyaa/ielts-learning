"""Build A1 vocabulary import file + seed structure from authored source files.

Usage: python3 -I build_a1.py <src_dir> <cmudict> <out_dir>
"""
import glob
import json
import os
import re
import sys

src_dir, cmu_path, out_dir = sys.argv[1:4]
MAX_PER_LESSON = 15

ARPA = {
    "AA": "ɑː", "AE": "æ", "AW": "aʊ", "AY": "aɪ", "EH": "e", "EY": "eɪ",
    "IH": "ɪ", "IY": "iː", "OW": "oʊ", "OY": "ɔɪ", "UH": "ʊ", "UW": "uː", "AO": "ɔː",
    "B": "b", "CH": "tʃ", "D": "d", "DH": "ð", "F": "f", "G": "ɡ", "HH": "h",
    "JH": "dʒ", "K": "k", "L": "l", "M": "m", "N": "n", "NG": "ŋ", "P": "p",
    "R": "r", "S": "s", "SH": "ʃ", "T": "t", "TH": "θ", "V": "v", "W": "w",
    "Y": "j", "Z": "z", "ZH": "ʒ",
}
VOWELS = {"AA", "AE", "AH", "AO", "AW", "AY", "EH", "ER", "EY", "IH", "IY", "OW", "OY", "UH", "UW"}

# Manual pronunciations (CMU keys) for heteronyms, abbreviations and missing words.
OVERRIDE_ARPA = {
    "live": "L IH1 V", "read": "R IY1 D", "close": "K L OW1 Z", "wind": "W IH1 N D",
    "use": "Y UW1 Z", "mr.": "M IH1 S T ER0", "mrs.": "M IH1 S IH0 Z", "a.m.": "EY1 EH1 M",
    "p.m.": "P IY1 EH1 M", "tv": "T IY1 V IY1", "t-shirt": "T IY1 SH ER2 T",
    "a": "AH0", "the": "DH AH0", "an": "AE1 N", "lead": "L IY1 D",
    "salesperson": "S EY1 L Z P ER2 S AH0 N", "motorbike": "M OW1 T ER0 B AY2 K",
    "chopsticks": "CH AA1 P S T IH0 K S", "badminton": "B AE1 D M IH2 N T AH0 N",
    "website": "W EH1 B S AY2 T", "online": "AA1 N L AY2 N", "email": "IY1 M EY2 L",
    "internet": "IH1 N T ER0 N EH2 T", "laptop": "L AE1 P T AA2 P",
    "hometown": "HH OW1 M T AW1 N", "yogurt": "Y OW1 G ER0 T", "o'clock": "AH0 K L AA1 K",
    "password": "P AE1 S W ER2 D", "volleyball": "V AA1 L IY0 B AO2 L",
    "colleague": "K AA1 L IY0 G", "developer": "D IH0 V EH1 L AH0 P ER0",
}

cmu = {}
with open(cmu_path, encoding="utf-8") as fh:
    for line in fh:
        parts = line.split("#")[0].strip().split()
        if not parts:
            continue
        key = parts[0].lower()
        if "(" in key:
            continue
        cmu.setdefault(key, " ".join(parts[1:]))


def to_ipa(arpa: str) -> str:
    phones = arpa.split()
    out = []
    n_vowels = sum(1 for p in phones if re.sub(r"\d", "", p) in VOWELS)
    primary = [k for k, p in enumerate(phones) if p.endswith("1")]
    last_primary = primary[-1] if primary else -1
    for i, ph in enumerate(phones):
        base = re.sub(r"\d", "", ph)
        stress = ph[-1] if ph[-1].isdigit() else ""
        if base in VOWELS:
            if base == "ER":
                sym = "ɜːr" if stress in "12" else "ər"
            elif base == "AH":
                sym = "ʌ" if stress in "12" else "ə"
            elif base == "IY" and stress == "0":
                sym = "i"
            else:
                sym = ARPA[base]
            # only primary stress is marked; secondary stress is noise for A1 learners
            if n_vowels > 1 and i == last_primary:
                mark = "ˈ"
                # place the mark before the consonant that starts this syllable
                j = len(out)
                if i > 0 and re.sub(r"\d", "", phones[i - 1]) not in VOWELS:
                    j -= 1
                    # keep common two-consonant onsets together (st, pl, tr, ...)
                    if i > 1 and re.sub(r"\d", "", phones[i - 2]) not in VOWELS and \
                            (phones[i - 2] in ("S",) or phones[i - 1] in ("L", "R", "W", "Y")):
                        j -= 1
                out.insert(j, mark)
            out.append(sym)
        else:
            out.append(ARPA[base])
    return "".join(out)


def ipa_for(word: str):
    tokens = word.lower().split()
    parts = []
    for tok in tokens:
        arpa = OVERRIDE_ARPA.get(tok) or cmu.get(tok) or cmu.get(tok.strip(".'"))
        if not arpa:
            return None
        parts.append(to_ipa(arpa))
    return "/" + " ".join(parts) + "/"


def slugify(value: str) -> str:
    return re.sub(r"[^a-z0-9]+", "-", value.lower()).strip("-")


topics, rows, problems = [], [], []
seen = {}
for path in sorted(glob.glob(os.path.join(src_dir, "*.txt"))):
    topic = lesson = None
    for raw in open(path, encoding="utf-8"):
        line = raw.strip()
        if not line:
            continue
        if line.startswith("## "):
            slug, title, desc, icon, color = [x.strip() for x in line[3:].split("|")]
            topic = {"slug": slug, "title": title, "description": desc, "icon": icon,
                     "color": color, "lessons": []}
            topics.append(topic)
            continue
        if line.startswith("### "):
            lesson = {"title": line[4:].strip(), "words": []}
            topic["lessons"].append(lesson)
            continue
        cols = line.split("|")
        if len(cols) != 6:
            problems.append(f"bad columns: {line}")
            continue
        word, pos, vi, en, ex, ex_vi = [c.strip() for c in cols]
        key = (word.lower(), pos)
        if key in seen:
            continue  # duplicate word+pos: keep the first occurrence only
        same_word = [k for k in seen if k[0] == word.lower()]
        slug = "a1-" + slugify(word) + (f"-{slugify(pos)}" if same_word else "")
        seen[key] = slug
        ipa = ipa_for(word)
        if not ipa:
            problems.append(f"no IPA: {word}")
        if word.lower().split()[0] not in ex.lower():
            problems.append(f"example may not contain word: {word} -> {ex}")
        lesson["words"].append({"word": word, "slug": slug, "pos": pos, "vi": vi, "en": en,
                                "ex": ex, "exVi": ex_vi, "ipa": ipa or ""})

EMOJI = {
    "a1-me-and-my-family": "👪", "a1-numbers-and-time": "🕒", "a1-grammar-words": "🧩",
    "a1-work-and-study": "💼", "a1-daily-life": "☀️", "a1-home": "🏠", "a1-food-and-drink": "🍜",
    "a1-shopping-and-money": "🛍️", "a1-city-and-travel": "🗺️", "a1-body-and-health": "🏥",
    "a1-describing-people-and-things": "✨", "a1-free-time": "🎵", "a1-nature-and-weather": "🌿",
    "a1-communication": "💬",
}

# split long lessons, assign slugs/order
structure = []
for t_index, topic in enumerate(topics, start=1):
    t_out = {k: topic[k] for k in ("slug", "title", "description", "icon", "color")}
    t_out["orderIndex"] = t_index
    t_out["emoji"] = EMOJI.get(topic["slug"], "📚")
    t_out["lessons"] = []
    order = 0
    for lesson in topic["lessons"]:
        words = lesson["words"]
        if not words:
            continue
        chunks = max(1, -(-len(words) // MAX_PER_LESSON))
        size = -(-len(words) // chunks)
        for c in range(chunks):
            part = words[c * size:(c + 1) * size]
            order += 1
            title = lesson["title"] + (f" ({c + 1})" if chunks > 1 else "")
            l_slug = "a1-" + slugify(title)
            t_out["lessons"].append({"slug": l_slug, "title": title, "orderIndex": order,
                                     "wordCount": len(part)})
            for w_index, w in enumerate(part, start=1):
                # Audio is generated by generate_audio.sh; until a file exists the
                # speaker button falls back to browser speech synthesis.
                audio = f"/audio/vocab/{w['slug']}.m4a"
                rows.append({
                    "word": w["word"], "slug": w["slug"], "meaningVi": w["vi"],
                    "topicSlug": topic["slug"], "lessonSlug": l_slug, "ipa": w["ipa"],
                    "audioUrl": audio, "partOfSpeech": w["pos"], "meaningEn": w["en"],
                    "shortDefinition": w["en"], "exampleSentence": w["ex"],
                    "exampleMeaningVi": w["exVi"], "exampleSource": "A1 Foundation",
                    "difficulty": "BEGINNER", "orderIndex": w_index, "isRequired": "true",
                })
    structure.append(t_out)

lesson_slugs = [l["slug"] for t in structure for l in t["lessons"]]
assert len(lesson_slugs) == len(set(lesson_slugs)), "duplicate lesson slug"

os.makedirs(out_dir, exist_ok=True)
with open(os.path.join(out_dir, "foundation_a1.json"), "w", encoding="utf-8") as fh:
    json.dump({"bandLevel": {"bandScore": 3.0, "title": "A1 Foundation"}, "topics": structure},
              fh, ensure_ascii=False, indent=2)
with open(os.path.join(out_dir, "a1_rows.json"), "w", encoding="utf-8") as fh:
    json.dump(rows, fh, ensure_ascii=False, indent=1)

from openpyxl import Workbook
from openpyxl.styles import Font, PatternFill

wb = Workbook()
ws = wb.active
ws.title = "A1 Vocabulary"
headers = list(rows[0].keys())
ws.append(headers)
for cell in ws[1]:
    cell.font = Font(bold=True, color="FFFFFF")
    cell.fill = PatternFill("solid", fgColor="4F46E5")
for r in rows:
    ws.append([r[h] for h in headers])
widths = {"word": 16, "meaningVi": 28, "topicSlug": 24, "lessonSlug": 30, "ipa": 18,
          "meaningEn": 40, "shortDefinition": 40, "exampleSentence": 40, "exampleMeaningVi": 40}
for i, h in enumerate(headers, start=1):
    ws.column_dimensions[ws.cell(1, i).column_letter].width = widths.get(h, 14)
ws.freeze_panes = "A2"
wb.save(os.path.join(out_dir, "a1_vocabulary_import.xlsx"))

with open(os.path.join(out_dir, "audio_manifest.tsv"), "w", encoding="utf-8") as fh:
    for r in rows:
        fh.write(f"{r['slug']}\t{r['word']}\n")

print(f"topics={len(structure)} lessons={len(lesson_slugs)} words={len(rows)}")
for t in structure:
    print(f"  {t['slug']}: {len(t['lessons'])} lessons, {sum(l['wordCount'] for l in t['lessons'])} words")
print(f"problems={len(problems)}")
for p in problems:
    print("  " + p)

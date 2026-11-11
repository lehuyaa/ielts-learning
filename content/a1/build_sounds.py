"""Authoring source for the A1 sound cards and the "Sound Lab" practice topic
(Tasks 57-58, D6).

Builds:
  a1-sounds.json             40 phonemes (IPA chart) + topic "a1-sound-lab"
                             (minimal pairs, ending -s / -ed, word stress)
  audio_manifest_sounds.tsv  audio for generate_lesson_audio.sh

IPA is generated from the CMU Pronouncing Dictionary and converted to the
Oxford NAmE style used by our vocabulary data (/ʃɪp/, /ˈtiːtʃər/, /kɑːr/).
Words missing from CMU, or whose CMU form is wrong for this style, are listed
in IPA_OVERRIDES.

Usage (from content/a1):
    python3 build_sounds.py --cmu path/to/cmudict.dict
CMU dict: https://github.com/cmusphinx/cmudict (cmudict.dict)
"""
import argparse
import json
import random
import re

AUDIO_ROOT = "/audio/a1"
audio_lines = {}


def audio(path, text):
    audio_lines[path] = text
    return f"{AUDIO_ROOT}/{path}.m4a"


def word_audio(word):
    """Public audio of one word (also used on the sound card page)."""
    return audio(f"words/{word.lower()}", word)


# ---------------------------------------------------------------------------
# CMU -> Oxford NAmE IPA

VOWELS = {
    "AA": "ɑː", "AE": "æ", "AH": "ʌ", "AO": "ɔː", "AW": "aʊ", "AY": "aɪ", "EH": "e",
    "ER": "ɜːr", "EY": "eɪ", "IH": "ɪ", "IY": "iː", "OW": "oʊ", "OY": "ɔɪ", "UH": "ʊ", "UW": "uː",
}
CONSONANTS = {
    "B": "b", "CH": "tʃ", "D": "d", "DH": "ð", "F": "f", "G": "ɡ", "HH": "h", "JH": "dʒ", "K": "k",
    "L": "l", "M": "m", "N": "n", "NG": "ŋ", "P": "p", "R": "r", "S": "s", "SH": "ʃ", "T": "t",
    "TH": "θ", "V": "v", "W": "w", "Y": "j", "Z": "z", "ZH": "ʒ",
}
ONSETS = {
    "pl", "pr", "bl", "br", "tr", "dr", "kl", "kr", "ɡl", "ɡr", "fl", "fr", "θr", "ʃr", "sp", "st",
    "sk", "sm", "sn", "sl", "sw", "spl", "spr", "str", "skr", "tw", "kw", "dw", "ɡw", "pj", "bj",
    "kj", "fj", "vj", "mj", "hj",
}

IPA_OVERRIDES = {
    "usually": "/ˈjuːʒuəli/",
    "television": "/ˈteləvɪʒn/",
    "caught": "/kɔːt/",
    "bought": "/bɔːt/",
    "law": "/lɔː/",
    "saw": "/sɔː/",
    "raw": "/rɔː/",
    "called": "/kɔːld/",
    "ball": "/bɔːl/",
    "than": "/ðæn/",
    "and": "/ænd/",
    "use": "/juːs/",
    "who'd": "/huːd/",
    "luke": "/luːk/",
    "yolk": "/joʊk/",
    "orange": "/ˈɔːrɪndʒ/",
    "oranges": "/ˈɔːrɪndʒɪz/",
    "studied": "/ˈstʌdid/",
    "hotel": "/hoʊˈtel/",
    "relax": "/rɪˈlæks/",
    "begin": "/bɪˈɡɪn/",
    "return": "/rɪˈtɜːrn/",
    "enjoy": "/ɪnˈdʒɔɪ/",
    "decided": "/dɪˈsaɪdɪd/",
    "between": "/bɪˈtwiːn/",
    "again": "/əˈɡen/",
    "today": "/təˈdeɪ/",
    "about": "/əˈbaʊt/",
    "july": "/dʒuˈlaɪ/",
    "collect": "/kəˈlekt/",
    "correct": "/kəˈrekt/",
    "lead": "/liːd/",
    "read": "/riːd/",
    "live": "/lɪv/",
}

cmu = {}


def load_cmu(path):
    with open(path, encoding="utf-8") as handle:
        for line in handle:
            parts = line.split("#")[0].split()
            if not parts:
                continue
            word = re.sub(r"\(\d+\)$", "", parts[0])
            cmu.setdefault(word, []).append(parts[1:])


def to_ipa_phones(arpabet):
    """Returns [(ipa, stress)] where stress is None for consonants."""
    phones = []
    for i, phone in enumerate(arpabet):
        base, stress = re.match(r"([A-Z]+)(\d?)", phone).groups()
        nxt = re.sub(r"\d", "", arpabet[i + 1]) if i + 1 < len(arpabet) else ""
        if base in VOWELS:
            stress = int(stress)
            ipa = VOWELS[base]
            if base == "AH" and stress == 0:
                ipa = "ə"
            elif base == "ER" and stress == 0:
                ipa = "ər"
            elif base == "IY" and stress == 0:
                ipa = "i"
            elif base == "UW" and stress == 0:
                ipa = "u"
            elif nxt == "R" and base == "IY":
                ipa = "ɪ"
            elif nxt == "R" and base in ("EH", "EY"):
                ipa = "e"
            elif nxt == "R" and base == "UW":
                ipa = "ʊ"
            phones.append((ipa, stress))
        else:
            phones.append((CONSONANTS[base], None))
    return phones


def render(phones):
    vowel_count = sum(1 for _, stress in phones if stress is not None)
    out = []
    for i, (ipa, stress) in enumerate(phones):
        if stress in (1, 2) and vowel_count > 1:
            # Move the stress mark before the longest valid onset.
            run = 0
            while run < len(out) and out[-1 - run][1] is None and i - 1 - run >= 0:
                run += 1
            onset = 0
            for size in range(run, 0, -1):
                cluster = "".join(p for p, _ in out[len(out) - size:])
                if size == 1 and cluster != "ŋ" or cluster in ONSETS:
                    onset = size
                    break
            mark = "ˈ" if stress == 1 else "ˌ"
            out.insert(len(out) - onset, (mark, "mark"))
        out.append((ipa, stress))
    return "/" + "".join(p for p, _ in out) + "/"


SYMBOLS = sorted(["iː", "ɪ", "e", "æ", "ɑː", "ɔː", "ʊ", "uː", "ʌ", "ɜː", "ə", "eɪ", "aɪ", "ɔɪ", "aʊ", "oʊ",
                  "p", "b", "t", "d", "k", "ɡ", "f", "v", "θ", "ð", "s", "z", "ʃ", "ʒ", "h", "tʃ", "dʒ",
                  "m", "n", "ŋ", "l", "r", "w", "j", "i", "u"], key=len, reverse=True)


def tokens(transcription):
    """Splits /ˈtiːtʃər/ into phonemes (ɜːr is ɜː + r here, stress marks dropped)."""
    text = transcription.strip("/").replace("ˈ", "").replace("ˌ", "").replace(" ", "")
    out = []
    while text:
        for symbol in SYMBOLS:
            if text.startswith(symbol):
                out.append(symbol)
                text = text[len(symbol):]
                break
        else:
            raise SystemExit(f"cannot split {transcription!r} at {text!r}")
    return out


def has_sound(transcription, symbol):
    return symbol.replace("ɜːr", "ɜː") in tokens(transcription)


def variants(word):
    key = word.lower()
    if key in IPA_OVERRIDES:
        return [IPA_OVERRIDES[key]]
    found = cmu.get(key)
    if not found:
        raise SystemExit(f"no IPA for {word!r}: add it to IPA_OVERRIDES")
    return [render(to_ipa_phones(v)) for v in found]


def ipa(word, must_contain=None):
    rendered = variants(word)
    if must_contain:
        for candidate in rendered:
            if has_sound(candidate, must_contain):
                return candidate
        raise SystemExit(f"{word!r} {rendered} has no /{must_contain}/")
    return rendered[0]


def contrast(word_a, word_b, a, b):
    """Picks transcriptions that differ only in sound a -> b (or drop a when b is None)."""
    a, b = a.replace("ɜːr", "ɜː"), b.replace("ɜːr", "ɜː") if b else None
    for ipa_a in variants(word_a):
        for ipa_b in variants(word_b):
            ta, tb = tokens(ipa_a), tokens(ipa_b)
            if b is None:
                ok = any(ta[i] == a and ta[:i] + ta[i + 1:] == tb for i in range(len(ta)))
            else:
                diff = [i for i in range(min(len(ta), len(tb))) if ta[i] != tb[i]]
                ok = len(ta) == len(tb) and len(diff) == 1 and ta[diff[0]] == a and tb[diff[0]] == b
            if ok:
                return ipa_a, ipa_b
    raise SystemExit(f"{word_a}/{word_b} is not a minimal pair for /{a}/–/{b}/: {variants(word_a)} {variants(word_b)}")


# ---------------------------------------------------------------------------
# Sound cards (spec §3.3: 11 vowels, 5 diphthongs, 24 consonants)
# (symbol, key, kind, voiced, nameVi, examples, mouthVi, mistakesVi)

V, D, C = "VOWEL", "DIPHTHONG", "CONSONANT"
PHONEMES = [
    ("iː", "iy", V, None, "i dài", ["see", "eat", "green"],
     "Môi kéo sang hai bên như đang cười, lưỡi nâng cao và đẩy về phía trước, gần chạm vòm miệng. Cơ miệng căng, kéo dài âm gấp đôi /ɪ/.",
     "Đọc ngắn như 'i' tiếng Việt nên 'sheep' nghe thành 'ship', 'leave' thành 'live'. Chú ý kéo dài cả khi có phụ âm cuối (eat, seat)."),
    ("ɪ", "ih", V, None, "i ngắn", ["sit", "big", "live"],
     "Môi thả lỏng, hơi mở; lưỡi cao nhưng thấp và lùi hơn /iː/. Âm rất ngắn, nằm giữa 'i' và 'ê'.",
     "Đọc thành /iː/ dài và căng (live → leave, it → eat). Thả lỏng cơ miệng và bật âm thật nhanh."),
    ("e", "eh", V, None, "e", ["bed", "ten", "red"],
     "Miệng mở vừa, môi hơi kéo sang hai bên, lưỡi ở giữa và hơi phía trước. Giống 'e' tiếng Việt nhưng ngắn.",
     "Nhầm với /æ/: 'men' và 'man' đọc giống nhau. /e/ mở miệng ít hơn /æ/."),
    ("æ", "ae", V, None, "a bẹt", ["cat", "bag", "apple"],
     "Hạ hàm xuống thấp, mở miệng rộng, môi kéo sang hai bên; lưỡi thấp và đẩy ra trước. Âm nằm giữa 'e' và 'a', hơi kéo dài.",
     "Đọc thành 'e' (bad → bed) hoặc 'a' tiếng Việt. Hãy đặt miệng như nói 'e' rồi hạ hàm thật thấp."),
    ("ɑː", "aa", V, None, "a dài", ["father", "hot", "car"],
     "Mở miệng rộng, hàm hạ thấp, lưỡi nằm thấp và lùi về sau, môi không tròn. Kéo dài như 'a' khi ngạc nhiên.",
     "Trong giọng Mỹ, nhiều từ viết 'o' đọc /ɑː/ (hot, stop, job) nhưng người học hay tròn môi thành 'o'. Đừng tròn môi."),
    ("ɔː", "ao", V, None, "o dài", ["saw", "ball", "talk"],
     "Môi tròn vừa (không chu nhiều), miệng mở khá rộng, lưỡi thấp và lùi sau. Âm dài, giữa 'o' và 'a'.",
     "Đọc thành /oʊ/ (law → low, caught → coat). Nhiều người Mỹ đọc âm này gần /ɑː/, nên nghe cả hai đều đúng."),
    ("ʊ", "uh", V, None, "u ngắn", ["book", "good", "put"],
     "Môi hơi tròn nhưng thả lỏng, không chu; lưỡi cao và lùi sau. Âm ngắn, giữa 'u' và 'ư'.",
     "Đọc thành /uː/ dài (full → fool, look → Luke). Giữ âm thật ngắn, môi thả lỏng."),
    ("uː", "uw", V, None, "u dài", ["food", "blue", "school"],
     "Môi tròn và chu ra phía trước, lưỡi cao và lùi sau; kéo dài âm.",
     "Đọc quá ngắn nên 'pool' nghe như 'pull'. Hãy chu môi và kéo dài."),
    ("ʌ", "ah", V, None, "ă / ơ ngắn", ["cup", "sun", "love"],
     "Miệng mở vừa, môi thả lỏng, lưỡi ở giữa. Âm ngắn, giống 'ă' hoặc 'ơ' bật nhanh.",
     "Đọc thành /ɑː/ dài (cut → cot) hoặc 'a' tiếng Việt. Nhiều từ viết 'u' và 'o' có âm này: bus, mother, money."),
    ("ɜːr", "er", V, None, "ơ cong lưỡi", ["bird", "work", "girl"],
     "Môi thả lỏng, hơi mở; cong đầu lưỡi lên và lùi về sau (không chạm vòm miệng), giữ nguyên trong suốt âm. Âm dài.",
     "Bỏ phần cong lưỡi nên 'work' nghe như 'woke' hay 'oắc'. Phải cong lưỡi ngay từ đầu âm."),
    ("ə", "ax", V, None, "ơ nhẹ (schwa)", ["about", "sofa", "teacher"],
     "Âm 'ơ' rất ngắn và nhẹ, miệng và lưỡi hoàn toàn thả lỏng. Chỉ xuất hiện ở âm tiết KHÔNG nhấn — là âm hay gặp nhất trong tiếng Anh.",
     "Đọc rõ từng chữ cái (banana thành 'ba-na-na', teacher thành 'ti-chơ' rõ ràng). Âm tiết không nhấn thường bị 'nuốt' thành /ə/."),
    ("eɪ", "ey", D, None, "ây", ["day", "name", "rain"],
     "Bắt đầu từ /e/ rồi trượt nhanh lên /ɪ/: miệng khép dần, môi kéo sang hai bên. Giống 'ây' tiếng Việt.",
     "Đọc thành /e/ đơn (late → let, wait → wet), nhất là khi có phụ âm cuối. Phải trượt đủ lên /ɪ/."),
    ("aɪ", "ay", D, None, "ai", ["five", "my", "time"],
     "Bắt đầu với miệng mở rộng như 'a', rồi trượt lên /ɪ/, hàm nâng dần. Giống 'ai' tiếng Việt.",
     "Bỏ phụ âm cuối sau âm này (time → 'tai', like → 'lai', five → 'phai')."),
    ("ɔɪ", "oy", D, None, "oi", ["boy", "toy", "coin"],
     "Bắt đầu với môi tròn như /ɔː/ rồi trượt sang /ɪ/, môi chuyển từ tròn sang dẹt. Giống 'oi' tiếng Việt.",
     "Ít bị sai; chú ý đừng bỏ phụ âm cuối (coin, noise, enjoyed)."),
    ("aʊ", "aw", D, None, "ao", ["now", "house", "down"],
     "Bắt đầu mở rộng miệng như 'a', rồi trượt về /ʊ/: môi tròn dần lại. Giống 'ao' tiếng Việt.",
     "Đọc thành /oʊ/ (town → tone) hoặc bỏ phụ âm cuối (house → 'hao', down → 'đao')."),
    ("oʊ", "ow", D, None, "ô-u", ["go", "home", "phone"],
     "Bắt đầu với môi tròn vừa như 'ô', rồi tròn và khép hơn về /ʊ/. Giống 'âu' nhưng môi tròn.",
     "Đọc thành 'ô' đơn tiếng Việt (no → 'nô'), thiếu phần trượt cuối. Hay nhầm với /ɔː/ (coat – caught)."),

    ("p", "p", C, False, "p bật hơi", ["pen", "happy", "stop"],
     "Mím hai môi, giữ hơi rồi bật mạnh ra. Đặt tờ giấy trước miệng: tờ giấy phải rung. Dây thanh không rung.",
     "Không bật hơi nên nghe như 'b' (pig → big). Ở cuối từ (stop, help) hay bị bỏ hẳn — hãy mím môi rõ ở cuối."),
    ("b", "b", C, True, "b", ["bag", "baby", "job"],
     "Mím hai môi như /p/ nhưng dây thanh rung và hơi bật ra nhẹ hơn.",
     "Âm 'b' cuối từ (job, cab) thường bị bỏ hoặc đọc thành /p/. Giữ rung cổ họng đến hết âm."),
    ("t", "t", C, False, "t bật hơi", ["ten", "time", "cat"],
     "Đầu lưỡi chạm lợi trên (sau răng cửa trên), chặn hơi rồi bật mạnh. Không rung dây thanh.",
     "Đọc như 't' tiếng Việt không bật hơi, hoặc bỏ /t/ ở cuối (eight, what, cat). Giọng Mỹ đọc /t/ giữa hai nguyên âm nhẹ như 'đ' (water, better)."),
    ("d", "d", C, True, "d", ["dog", "day", "red"],
     "Vị trí như /t/ (đầu lưỡi chạm lợi trên) nhưng dây thanh rung.",
     "Đọc thành 'đ' tiếng Việt (đầu lưỡi chạm răng). Bỏ /d/ cuối (red, good, played)."),
    ("k", "k", C, False, "k bật hơi", ["key", "cake", "book"],
     "Cuống lưỡi nâng lên chạm vòm mềm phía sau, chặn hơi rồi bật mạnh. Không rung.",
     "Bỏ /k/ cuối (book → 'bu', like → 'lai'). Cụm /ks/ (six, books) hay bị đọc thiếu."),
    ("ɡ", "g", C, True, "g", ["go", "bag", "egg"],
     "Vị trí như /k/ nhưng dây thanh rung.",
     "Đọc thành 'g' tiếng Việt hơi xát (ga). Bỏ /ɡ/ cuối (bag, dog)."),
    ("f", "f", C, False, "f", ["fish", "coffee", "life"],
     "Răng cửa trên chạm nhẹ môi dưới, thổi hơi qua khe hẹp. Không rung.",
     "Bỏ /f/ cuối (life, laugh, half). Đọc 'f' thành /p/ (fan → pan)."),
    ("v", "v", C, True, "v", ["very", "five", "love"],
     "Như /f/ (răng trên chạm môi dưới) nhưng dây thanh rung — sờ cổ sẽ thấy rung.",
     "Đọc thành /f/ ở cuối từ (five → 'phai', love → 'lớp') hoặc thành /w/ (very → 'wery'). Cần thấy môi rung."),
    ("θ", "th", C, False, "th vô thanh", ["think", "three", "mouth"],
     "Đặt đầu lưỡi giữa hai hàm răng (hoặc chạm mặt sau răng trên) và thổi hơi qua. Không rung.",
     "Đọc thành 'th' tiếng Việt (bật như /t/) hoặc /s/: think → 'sink', three → 'tree'. Phải thấy đầu lưỡi giữa hai răng."),
    ("ð", "dh", C, True, "th hữu thanh", ["this", "mother", "they"],
     "Như /θ/ (đầu lưỡi giữa hai răng) nhưng dây thanh rung.",
     "Đọc thành 'đ' hoặc 'd': they → 'day', this → 'đít'. Rất phổ biến vì gặp trong the, this, that, there."),
    ("s", "s", C, False, "s", ["see", "bus", "class"],
     "Đầu lưỡi gần lợi trên (không chạm), răng gần khép, thổi hơi thành tiếng 'xì'. Không rung.",
     "Bỏ /s/ cuối (bus, books, likes) — lỗi rất hay gặp vì tiếng Việt không có phụ âm cuối /s/. Hay đọc thành /ʃ/ ('sh')."),
    ("z", "z", C, True, "z", ["zoo", "zero", "nose"],
     "Như /s/ nhưng dây thanh rung, nghe như tiếng ong 'zzz'.",
     "Đọc thành /s/ hoặc 'd' (zoo → 'du'). Đuôi -s sau âm hữu thanh đọc /z/ (dogs, plays) nhưng hay bị đọc /s/ hoặc bỏ."),
    ("ʃ", "sh", C, False, "sh", ["she", "shop", "fish"],
     "Môi hơi tròn và chu ra, đầu lưỡi lùi sau hơn /s/, thổi hơi như 'suỵt'. Không rung.",
     "Đọc thành /s/ (she → see, ship → sip). Phải chu môi."),
    ("ʒ", "zh", C, True, "zh", ["vision", "measure", "television"],
     "Như /ʃ/ (môi chu, lưỡi lùi) nhưng dây thanh rung.",
     "Đọc thành 'd' hoặc /z/. Âm này hiếm ở đầu từ; thường gặp trong -sion, -sure (television, measure, usually)."),
    ("h", "hh", C, False, "h", ["hello", "house", "have"],
     "Miệng mở sẵn theo nguyên âm phía sau, đẩy hơi nhẹ từ cổ họng như khi hà hơi lên kính.",
     "Bỏ /h/ (hair → air, hold → old) hoặc đọc quá nặng như 'kh'."),
    ("tʃ", "ch", C, False, "ch", ["chair", "teacher", "watch"],
     "Bắt đầu như /t/ (lưỡi chạm lợi) rồi nhả ra thành /ʃ/, môi hơi chu. Không rung.",
     "Đọc thành 'ch' tiếng Việt không bật hơi hoặc /ʃ/ (chair → share). Bỏ ở cuối từ (watch → 'wát')."),
    ("dʒ", "jh", C, True, "j", ["job", "July", "page"],
     "Như /tʃ/ nhưng dây thanh rung.",
     "Đọc thành 'gi'/'d' tiếng Việt (job → 'dóp') hoặc /j/ (jet → yet). Bỏ ở cuối từ (page, orange)."),
    ("m", "m", C, True, "m", ["man", "summer", "name"],
     "Mím hai môi, hơi ra đằng mũi, dây thanh rung.",
     "Ở cuối từ (name, time, room) môi hay không khép lại nên âm bị mất."),
    ("n", "n", C, True, "n", ["no", "nine", "ten"],
     "Đầu lưỡi chạm lợi trên, hơi ra đằng mũi, dây thanh rung.",
     "Nhầm /n/ cuối với /ŋ/ (thin → thing) và nhầm /n/ với /l/ ở đầu từ (night → light)."),
    ("ŋ", "ng", C, True, "ng", ["sing", "long", "English"],
     "Cuống lưỡi nâng chạm vòm mềm (như /k/), hơi ra đằng mũi. Không bật /ɡ/ ở cuối.",
     "Đọc thành /n/ (sing → sin) hoặc thêm /k/, /ɡ/ phía sau. Đuôi -ing phải kết thúc bằng /ŋ/."),
    ("l", "l", C, True, "l", ["light", "yellow", "feel"],
     "Đầu lưỡi chạm lợi trên, hơi đi ra hai bên lưỡi. Ở cuối từ (feel, school) lưỡi lùi về sau thành 'l tối'.",
     "Nhầm với /n/ (light → night) và /r/. Bỏ /l/ cuối (school → 'sku', milk → 'mia')."),
    ("r", "r", C, True, "r", ["red", "very", "car"],
     "Cong đầu lưỡi lên và lùi về sau, KHÔNG chạm vòm miệng, môi hơi tròn. Giọng Mỹ vẫn đọc /r/ sau nguyên âm (car, water).",
     "Rung lưỡi như 'r' tiếng Việt, hoặc đọc thành /z/ (rice → 'zai') hay /l/ (right → light)."),
    ("w", "w", C, True, "w", ["we", "water", "one"],
     "Chu tròn môi như 'u' rồi mở nhanh sang nguyên âm phía sau. Răng không chạm môi.",
     "Đọc thành /v/ (wine → vine) hoặc bỏ hẳn (would → 'ud')."),
    ("j", "y", C, True, "y", ["yes", "you", "music"],
     "Lưỡi nâng cao gần vòm miệng như 'i' rồi trượt nhanh sang nguyên âm sau. Giống 'd' miền Nam ('dạ') nhẹ.",
     "Đọc thành /dʒ/ (yes → 'jes') hoặc bỏ /j/ trước 'u' (music → 'mu-sic', use → 'ús')."),
]

# ---------------------------------------------------------------------------
# Minimal pairs: (sound A, sound B or None, [(word with A, word with B), ...])

PAIRS = [
    ("iː", "ɪ", [("sheep", "ship"), ("leave", "live"), ("eat", "it"), ("seat", "sit"), ("feel", "fill"),
                 ("heat", "hit"), ("beat", "bit"), ("cheap", "chip"), ("sleep", "slip"), ("deep", "dip"),
                 ("peel", "pill"), ("read", "rid")]),
    ("uː", "ʊ", [("pool", "pull"), ("fool", "full"), ("Luke", "look"), ("who'd", "hood")]),
    ("e", "æ", [("bed", "bad"), ("men", "man"), ("pen", "pan"), ("said", "sad"), ("head", "had"),
                ("set", "sat"), ("beg", "bag"), ("then", "than"), ("ten", "tan"), ("met", "mat")]),
    ("ʌ", "ɑː", [("cut", "cot"), ("luck", "lock"), ("duck", "dock"), ("hut", "hot"), ("cup", "cop"),
                 ("nut", "not"), ("shut", "shot"), ("rub", "rob")]),
    ("ɜːr", "ɑː", [("heard", "hard"), ("firm", "farm"), ("burn", "barn"), ("hurt", "heart"), ("fur", "far"),
                   ("stir", "star")]),
    ("ɔː", "oʊ", [("law", "low"), ("caught", "coat"), ("bought", "boat"), ("saw", "so"), ("called", "cold"),
                  ("ball", "bowl"), ("raw", "row")]),
    ("e", "eɪ", [("wet", "wait"), ("pen", "pain"), ("get", "gate"), ("men", "main"), ("sell", "sale"),
                 ("tell", "tail"), ("let", "late"), ("west", "waste")]),
    ("aɪ", "ɔɪ", [("buy", "boy"), ("tie", "toy"), ("file", "foil"), ("aisle", "oil")]),
    ("aʊ", "oʊ", [("now", "know"), ("town", "tone"), ("loud", "load"), ("found", "phoned")]),
    ("p", "b", [("pig", "big"), ("pear", "bear"), ("pie", "buy"), ("pack", "back"), ("path", "bath"),
                ("cap", "cab"), ("rope", "robe"), ("pull", "bull"), ("pin", "bin"), ("pet", "bet")]),
    ("t", "d", [("ten", "den"), ("time", "dime"), ("tie", "die"), ("town", "down"), ("try", "dry"),
                ("write", "ride"), ("bat", "bad"), ("coat", "code"), ("seat", "seed"), ("cart", "card")]),
    ("k", "ɡ", [("coat", "goat"), ("cold", "gold"), ("come", "gum"), ("back", "bag"), ("pick", "pig"),
                ("card", "guard"), ("class", "glass"), ("duck", "dug")]),
    ("θ", "s", [("think", "sink"), ("thing", "sing"), ("thick", "sick"), ("path", "pass"), ("mouth", "mouse"),
                ("worth", "worse")]),
    ("θ", "t", [("three", "tree"), ("thin", "tin"), ("thank", "tank"), ("both", "boat"), ("theme", "team"),
                ("thick", "tick")]),
    ("ð", "d", [("they", "day"), ("then", "den"), ("those", "doze"), ("there", "dare"), ("though", "dough")]),
    ("s", "z", [("sip", "zip"), ("sue", "zoo"), ("price", "prize"), ("bus", "buzz"), ("place", "plays"),
                ("rice", "rise")]),
    ("ʃ", "s", [("ship", "sip"), ("she", "see"), ("shell", "sell"), ("shoe", "sue"), ("sheet", "seat"),
                ("short", "sort")]),
    ("ʃ", "tʃ", [("ship", "chip"), ("shop", "chop"), ("share", "chair"), ("wash", "watch"), ("cash", "catch"),
                 ("sheep", "cheap")]),
    ("tʃ", "dʒ", [("cheap", "jeep"), ("choke", "joke"), ("chin", "gin"), ("rich", "ridge")]),
    ("dʒ", "j", [("jet", "yet"), ("joke", "yolk"), ("jam", "yam"), ("juice", "use")]),
    ("f", "v", [("fan", "van"), ("fine", "vine"), ("ferry", "very"), ("leaf", "leave"), ("safe", "save"),
                ("fast", "vast"), ("few", "view")]),
    ("f", "p", [("fan", "pan"), ("fast", "past"), ("fork", "pork"), ("coffee", "copy"), ("wife", "wipe"),
                ("feel", "peel")]),
    ("w", "v", [("wet", "vet"), ("west", "vest"), ("wine", "vine"), ("while", "vile")]),
    ("l", "r", [("light", "right"), ("long", "wrong"), ("lead", "read"), ("play", "pray"), ("glass", "grass"),
                ("collect", "correct"), ("fly", "fry"), ("lock", "rock")]),
    ("l", "n", [("light", "night"), ("line", "nine"), ("low", "no"), ("lot", "not")]),
    ("n", "ŋ", [("thin", "thing"), ("sin", "sing"), ("ran", "rang"), ("win", "wing"), ("ban", "bang"),
                ("kin", "king")]),
    ("m", "n", [("mail", "nail"), ("some", "sun"), ("map", "nap"), ("mine", "nine")]),
    ("h", None, [("hair", "air"), ("heat", "eat"), ("hill", "ill"), ("hand", "and"), ("hold", "old")]),
]

ENDING_S = [
    ("books", "/s/"), ("cats", "/s/"), ("stops", "/s/"), ("cakes", "/s/"), ("laughs", "/s/"),
    ("dogs", "/z/"), ("plays", "/z/"), ("rooms", "/z/"), ("cars", "/z/"), ("lives", "/z/"),
    ("buses", "/ɪz/"), ("watches", "/ɪz/"), ("boxes", "/ɪz/"), ("oranges", "/ɪz/"), ("washes", "/ɪz/"),
]
ENDING_ED = [
    ("worked", "/t/"), ("watched", "/t/"), ("stopped", "/t/"), ("washed", "/t/"), ("helped", "/t/"),
    ("played", "/d/"), ("called", "/d/"), ("lived", "/d/"), ("opened", "/d/"), ("cleaned", "/d/"),
    ("wanted", "/ɪd/"), ("needed", "/ɪd/"), ("visited", "/ɪd/"), ("waited", "/ɪd/"), ("decided", "/ɪd/"),
]
# (word, syllables, stressed index)
STRESS = [
    ("doctor", ["doc", "tor"], 0), ("window", ["win", "dow"], 0), ("happy", ["hap", "py"], 0),
    ("morning", ["mor", "ning"], 0), ("water", ["wa", "ter"], 0), ("sister", ["sis", "ter"], 0),
    ("number", ["num", "ber"], 0), ("yellow", ["yel", "low"], 0), ("table", ["ta", "ble"], 0),
    ("coffee", ["cof", "fee"], 0),
    ("begin", ["be", "gin"], 1), ("relax", ["re", "lax"], 1), ("hotel", ["ho", "tel"], 1),
    ("today", ["to", "day"], 1), ("between", ["be", "tween"], 1), ("again", ["a", "gain"], 1),
    ("July", ["Ju", "ly"], 1), ("enjoy", ["en", "joy"], 1), ("return", ["re", "turn"], 1),
    ("about", ["a", "bout"], 1),
]

ENDING_S_NOTE = {
    "/s/": "Âm cuối trước -s là âm vô thanh (p, t, k, f) nên đọc /s/.",
    "/z/": "Âm cuối trước -s là âm hữu thanh hoặc nguyên âm nên đọc /z/.",
    "/ɪz/": "Từ tận cùng bằng s, z, sh, ch, x, ge thì -es đọc /ɪz/.",
}
ENDING_ED_NOTE = {
    "/t/": "Âm cuối là âm vô thanh (p, k, f, s, sh, ch) nên -ed đọc /t/.",
    "/d/": "Âm cuối là âm hữu thanh hoặc nguyên âm nên -ed đọc /d/.",
    "/ɪd/": "Từ tận cùng bằng t hoặc d thì -ed đọc /ɪd/ (thêm một âm tiết).",
}


def pair_label(a, b):
    return f"/{a}/ – /{b}/" if b else f"/{a}/ – không có /{a}/"


def build(cmu_path):
    load_cmu(cmu_path)
    rng = random.Random(57)

    # Minimal pairs with IPA, checked: the A word must contain sound A, the B word sound B.
    pair_sets = []
    for a, b, words in PAIRS:
        items = []
        for word_a, word_b in words:
            ipa_a, ipa_b = contrast(word_a, word_b, a, b)
            items.append((word_a, ipa_a, word_b, ipa_b))
        pair_sets.append((a, b, items))

    # Sound cards
    phonemes = []
    for order, (symbol, key, kind, voiced, name_vi, examples, mouth, mistakes) in enumerate(PHONEMES, start=1):
        example_words = [{"word": w, "ipa": ipa(w, symbol),
                          "audioUrl": word_audio(w)} for w in examples]
        pairs = []
        for a, b, items in pair_sets:
            if symbol not in (a, b):
                continue
            for word_a, ipa_a, word_b, ipa_b in items:
                mine, other = ((word_a, ipa_a), (word_b, ipa_b)) if symbol == a else ((word_b, ipa_b), (word_a, ipa_a))
                pairs.append({"a": mine[0], "aIpa": mine[1], "b": other[0], "bIpa": other[1],
                              "contrast": (b if symbol == a else a) or "",
                              "aAudioUrl": word_audio(mine[0]), "bAudioUrl": word_audio(other[0])})
        entry = {
            "symbol": symbol, "kind": kind, "nameVi": name_vi,
            "exampleWords": example_words, "minimalPairs": pairs,
            "mouthVi": mouth, "mistakesVi": mistakes,
            "audioUrl": audio(f"sounds/ph-{key}", ", ".join(examples)),
            "orderIndex": order,
        }
        if kind == C:
            entry["voiced"] = voiced
        phonemes.append(entry)

    # Practice exercises
    ids = "ab"

    def minimal_pair_exercises(key, a, b, items):
        exercises = []
        for n, (word_a, ipa_a, word_b, ipa_b) in enumerate(items, start=1):
            options = [(word_a, ipa_a), (word_b, ipa_b)]
            rng.shuffle(options)
            correct = rng.randrange(2)
            heard = options[correct][0]
            exercises.append({
                "type": "MINIMAL_PAIR", "prompt": "Bạn nghe thấy từ nào?",
                "payload": {"audioUrl": audio(f"sounds/{key}-{n}", heard),
                            "options": [{"id": ids[i], "word": w, "ipa": p} for i, (w, p) in enumerate(options)]},
                "answer": {"correct": ids[correct]},
                "explanationVi": f"Bạn vừa nghe \"{heard}\" {dict(options)[heard]}.",
                "skillTag": f"sound:{a}" if b is None else f"sound:{a}-{b}",
            })
        # Two "read the IPA" items from the same set.
        for word_a, ipa_a, word_b, ipa_b in rng.sample(items, min(2, len(items))):
            target = rng.randrange(2)
            shown = [ipa_a, ipa_b][target]
            options = [word_a, word_b]
            rng.shuffle(options)
            exercises.append({
                "type": "IPA_TO_WORD", "prompt": "Phiên âm này là từ nào?",
                "payload": {"ipa": shown, "options": [{"id": ids[i], "word": w} for i, w in enumerate(options)]},
                "answer": {"correct": ids[options.index([word_a, word_b][target])]},
                "explanationVi": f"{shown} là \"{[word_a, word_b][target]}\".",
                "skillTag": f"sound:{a}" if b is None else f"sound:{a}-{b}",
            })
        return exercises

    def contrast_activity(a, b, items):
        key = f"{a}-{b or 'none'}"
        slug_key = "-".join(PHONEME_KEYS[s] for s in (a, b) if s)
        return {
            "slug": f"pairs-{slug_key}" + ("-h-drop" if b is None else ""),
            "type": "PRONUNCIATION",
            "title": pair_label(a, b),
            "instructionsVi": "Nghe và chọn từ đúng. Xem thẻ âm nếu chưa phân biệt được.",
            "content": {"phonemes": [s for s in (a, b) if s]},
            "exercises": minimal_pair_exercises(f"mp-{slug_key}" + ("-0" if b is None else ""), a, b, items),
        }

    pair_index = {(a, b): items for a, b, items in pair_sets}
    groups = [
        ("long-short-vowels", "Nguyên âm dài và ngắn", "Phân biệt i dài – i ngắn, u dài – u ngắn.",
         [("iː", "ɪ"), ("uː", "ʊ")]),
        ("open-vowels", "Nguyên âm e, æ, ʌ, ɑː", "Các nguyên âm mở hay bị đọc giống nhau.",
         [("e", "æ"), ("ʌ", "ɑː")]),
        ("r-vowels-diphthongs", "Âm /ɜːr/ và nguyên âm đôi", "Âm cong lưỡi và các nguyên âm đôi.",
         [("ɜːr", "ɑː"), ("ɔː", "oʊ"), ("e", "eɪ"), ("aɪ", "ɔɪ"), ("aʊ", "oʊ")]),
        ("stops", "Phụ âm p–b, t–d, k–g", "Âm bật hơi và âm hữu thanh, nhất là ở cuối từ.",
         [("p", "b"), ("t", "d"), ("k", "ɡ")]),
        ("th-sounds", "Âm th /θ/ và /ð/", "Đặt lưỡi giữa hai răng — không đọc thành s, t, d.",
         [("θ", "s"), ("θ", "t"), ("ð", "d")]),
        ("hissing", "Âm s, z, sh, ch, j", "Các âm gió và âm tắc xát.",
         [("s", "z"), ("ʃ", "s"), ("ʃ", "tʃ"), ("tʃ", "dʒ"), ("dʒ", "j")]),
        ("lips-teeth", "Âm f, v, w", "Răng trên chạm môi dưới, hay môi tròn?",
         [("f", "v"), ("f", "p"), ("w", "v")]),
        ("nasals-liquids", "Âm l, r, n, ng, m, h", "Âm mũi, l và r, và âm h đầu từ.",
         [("l", "r"), ("l", "n"), ("n", "ŋ"), ("m", "n"), ("h", None)]),
    ]

    lessons = []
    for order, (slug, title, description, contrasts) in enumerate(groups, start=1):
        activities = [contrast_activity(a, b, pair_index[(a, b)]) for a, b in contrasts]
        lessons.append({"slug": f"sounds-{slug}", "title": title, "description": description,
                        "estimatedMinutes": 4 * len(activities), "xpReward": 10 * len(activities),
                        "orderIndex": order, "activities": activities})

    def ending_activity(slug, title, instructions, items, notes, tag, options):
        exercises = []
        for word, correct in items:
            exercises.append({
                "type": "ENDING_SOUND", "prompt": f"Đuôi trong \"{word}\" đọc là gì?",
                "payload": {"word": word, "audioUrl": word_audio(word), "options": options},
                "answer": {"correct": correct}, "explanationVi": notes[correct], "skillTag": tag,
            })
        rng.shuffle(exercises)
        return {"slug": slug, "type": "PRONUNCIATION", "title": title, "instructionsVi": instructions,
                "content": {"phonemes": [], "rule": tag.split(":")[1]}, "exercises": exercises}

    lessons.append({
        "slug": "sounds-ending-s", "title": "Đọc đuôi -s / -es", "description": "/s/, /z/ hay /ɪz/?",
        "estimatedMinutes": 5, "xpReward": 10, "orderIndex": len(lessons) + 1,
        "activities": [ending_activity(
            "ending-s", "Đuôi -s: /s/, /z/, /ɪz/",
            "Sau âm vô thanh (p, t, k, f) đọc /s/; sau âm hữu thanh và nguyên âm đọc /z/; sau s, z, sh, ch, x, ge đọc /ɪz/.",
            ENDING_S, ENDING_S_NOTE, "sound:ending-s", ["/s/", "/z/", "/ɪz/"])],
    })
    lessons.append({
        "slug": "sounds-ending-ed", "title": "Đọc đuôi -ed", "description": "/t/, /d/ hay /ɪd/?",
        "estimatedMinutes": 5, "xpReward": 10, "orderIndex": len(lessons) + 1,
        "activities": [ending_activity(
            "ending-ed", "Đuôi -ed: /t/, /d/, /ɪd/",
            "Sau âm vô thanh đọc /t/; sau âm hữu thanh và nguyên âm đọc /d/; sau t, d đọc /ɪd/ (thêm một âm tiết).",
            ENDING_ED, ENDING_ED_NOTE, "sound:ending-ed", ["/t/", "/d/", "/ɪd/"])],
    })

    stress_exercises = []
    for word, syllables, index in STRESS:
        stress_exercises.append({
            "type": "WORD_STRESS", "prompt": f"Trọng âm của \"{word}\" ở âm tiết nào?",
            "payload": {"word": word, "audioUrl": word_audio(word), "syllables": syllables},
            "answer": {"stressIndex": index},
            "explanationVi": f"{'-'.join(s.upper() if i == index else s for i, s in enumerate(syllables))} {ipa(word)}.",
            "skillTag": "stress:2-syllable",
        })
    rng.shuffle(stress_exercises)
    lessons.append({
        "slug": "sounds-word-stress", "title": "Trọng âm từ hai âm tiết",
        "description": "Danh từ, tính từ thường nhấn âm đầu; nhiều động từ nhấn âm sau.",
        "estimatedMinutes": 6, "xpReward": 10, "orderIndex": len(lessons) + 1,
        "activities": [{
            "slug": "stress-2", "type": "PRONUNCIATION", "title": "Nhấn âm tiết nào?",
            "instructionsVi": "Nghe từ và chọn âm tiết được nhấn (đọc to và dài hơn). Danh từ/tính từ hai âm tiết thường nhấn âm đầu (DOC-tor), nhiều động từ và từ có tiền tố nhấn âm sau (be-GIN).",
            "content": {"phonemes": [], "rule": "stress-2"}, "exercises": stress_exercises,
        }],
    })

    bundle = {
        "version": 1, "courseSlug": "foundation-english", "cefrLevel": "A1",
        "phonemes": phonemes,
        "topics": [{
            "slug": "a1-sound-lab", "title": "Sound Lab · Luyện phát âm", "track": "VOCABULARY",
            "description": "Luyện nghe phân biệt âm, đuôi -s/-ed và trọng âm. Mở bất cứ lúc nào.",
            "emoji": "🔊", "orderIndex": 90, "lessons": lessons,
        }],
    }
    return bundle


PHONEME_KEYS = {symbol: key for symbol, key, *_ in PHONEMES}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--cmu", default="cmudict.dict")
    args = parser.parse_args()
    bundle = build(args.cmu)
    with open("a1-sounds.json", "w", encoding="utf-8") as handle:
        json.dump(bundle, handle, ensure_ascii=False, indent=2)
        handle.write("\n")
    with open("audio_manifest_sounds.tsv", "w", encoding="utf-8") as handle:
        for path in sorted(audio_lines):
            handle.write(f"{path}\t{audio_lines[path]}\n")

    pairs = sum(len(items) for _, _, items in PAIRS)
    exercises = sum(len(a["exercises"]) for t in bundle["topics"] for l in t["lessons"] for a in l["activities"])
    print(f"{len(bundle['phonemes'])} sound cards, {pairs} minimal pairs, "
          f"{exercises} exercises, {len(audio_lines)} audio lines")


if __name__ == "__main__":
    main()

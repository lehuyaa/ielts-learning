package pronunciation

import "strings"

// soundMapping turns one phoneme as Azure writes it (en-US IPA, or SAPI when
// IPA is not returned) into our Oxford NAmE notation: display is what the
// learner sees, symbols are the sound cards it counts for (spec §3.3).
type soundMapping struct {
	display string
	symbols []string
}

func m(display string, symbols ...string) soundMapping {
	return soundMapping{display: display, symbols: symbols}
}

var azureSounds = map[string]soundMapping{
	// Vowels (Azure IPA).
	"i": m("iː", "iː"), "ɪ": m("ɪ", "ɪ"), "ɛ": m("e", "e"), "e": m("e", "e"), "æ": m("æ", "æ"),
	"ɑ": m("ɑː", "ɑː"), "ɔ": m("ɔː", "ɔː"), "ʊ": m("ʊ", "ʊ"), "u": m("uː", "uː"), "ʌ": m("ʌ", "ʌ"),
	"ə": m("ə", "ə"), "ɝ": m("ɜːr", "ɜːr"), "ɚ": m("ər", "ə", "r"),
	"eɪ": m("eɪ", "eɪ"), "aɪ": m("aɪ", "aɪ"), "ɔɪ": m("ɔɪ", "ɔɪ"), "aʊ": m("aʊ", "aʊ"), "oʊ": m("oʊ", "oʊ"),
	"ju": m("juː", "j", "uː"),
	"ɪɹ": m("ɪr", "ɪ", "r"), "ɛɹ": m("er", "e", "r"), "ʊɹ": m("ʊr", "ʊ", "r"), "ɑɹ": m("ɑːr", "ɑː", "r"),
	"ɔɹ": m("ɔːr", "ɔː", "r"), "aɪɹ": m("aɪr", "aɪ", "r"), "aʊɹ": m("aʊr", "aʊ", "r"),
	// Consonants (Azure IPA).
	"p": m("p", "p"), "b": m("b", "b"), "t": m("t", "t"), "d": m("d", "d"), "k": m("k", "k"),
	"g": m("ɡ", "ɡ"), "ɡ": m("ɡ", "ɡ"), "f": m("f", "f"), "v": m("v", "v"), "θ": m("θ", "θ"), "ð": m("ð", "ð"),
	"s": m("s", "s"), "z": m("z", "z"), "ʃ": m("ʃ", "ʃ"), "ʒ": m("ʒ", "ʒ"), "h": m("h", "h"),
	"tʃ": m("tʃ", "tʃ"), "dʒ": m("dʒ", "dʒ"), "m": m("m", "m"), "n": m("n", "n"), "ŋ": m("ŋ", "ŋ"),
	"l": m("l", "l"), "ɹ": m("r", "r"), "r": m("r", "r"), "w": m("w", "w"), "j": m("j", "j"),
	// SAPI fallback.
	"iy": m("iː", "iː"), "ih": m("ɪ", "ɪ"), "eh": m("e", "e"), "ae": m("æ", "æ"), "aa": m("ɑː", "ɑː"),
	"ao": m("ɔː", "ɔː"), "uh": m("ʊ", "ʊ"), "uw": m("uː", "uː"), "ah": m("ʌ", "ʌ"), "ax": m("ə", "ə"),
	"er": m("ɜːr", "ɜːr"), "ey": m("eɪ", "eɪ"), "ay": m("aɪ", "aɪ"), "oy": m("ɔɪ", "ɔɪ"), "aw": m("aʊ", "aʊ"),
	"ow": m("oʊ", "oʊ"), "th": m("θ", "θ"), "dh": m("ð", "ð"), "sh": m("ʃ", "ʃ"), "zh": m("ʒ", "ʒ"),
	"hh": m("h", "h"), "ch": m("tʃ", "tʃ"), "jh": m("dʒ", "dʒ"), "ng": m("ŋ", "ŋ"), "y": m("j", "j"),
}

// mapSound returns the learner-facing notation and sound cards of an Azure phoneme.
func mapSound(azure string) soundMapping {
	key := strings.TrimSpace(azure)
	key = strings.NewReplacer("ˈ", "", "ˌ", "", "ː", "", ".", "").Replace(key)
	if mapping, ok := azureSounds[key]; ok {
		return mapping
	}
	if mapping, ok := azureSounds[strings.ToLower(key)]; ok {
		return mapping
	}
	return soundMapping{display: key}
}

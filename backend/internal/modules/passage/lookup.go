package passage

import "strings"

// wordCandidates returns the forms to look up for a word as it appears in a
// text: the word itself, then simple base forms (plural, past, -ing).
func wordCandidates(raw string) []string {
	word := strings.ToLower(strings.Trim(raw, " .,!?;:\"'“”‘’()[]"))
	word = strings.NewReplacer("’", "'", "‘", "'").Replace(word)
	word = strings.TrimSuffix(word, "'s")
	if word == "" {
		return nil
	}

	seen := map[string]bool{}
	var out []string
	add := func(candidate string) {
		if seen[candidate] || (len(candidate) < 2 && candidate != word) {
			return
		}
		seen[candidate] = true
		out = append(out, candidate)
	}
	add(word)

	switch {
	case strings.HasSuffix(word, "ies"):
		add(strings.TrimSuffix(word, "ies") + "y")
	case strings.HasSuffix(word, "es"):
		add(strings.TrimSuffix(word, "es"))
		add(strings.TrimSuffix(word, "s"))
	case strings.HasSuffix(word, "s") && !strings.HasSuffix(word, "ss"):
		add(strings.TrimSuffix(word, "s"))
	}
	switch {
	case strings.HasSuffix(word, "ied"):
		add(strings.TrimSuffix(word, "ied") + "y")
	case strings.HasSuffix(word, "ed"):
		stem := strings.TrimSuffix(word, "ed")
		add(stem)
		add(stem + "e")
		add(undouble(stem))
	}
	if strings.HasSuffix(word, "ing") {
		stem := strings.TrimSuffix(word, "ing")
		add(stem)
		add(stem + "e")
		add(undouble(stem))
	}
	return out
}

// undouble turns "stopp" into "stop" (stopped, running).
func undouble(stem string) string {
	n := len(stem)
	if n >= 3 && stem[n-1] == stem[n-2] && !strings.ContainsRune("aeiou", rune(stem[n-1])) {
		return stem[:n-1]
	}
	return stem
}

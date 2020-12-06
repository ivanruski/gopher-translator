package translation

import (
	"errors"
	"io"
	"regexp"
	"strings"
)

type WordTranslator interface {
	TranslateWord(word string) string
}

type WordSkipper interface {
	SkipWord(word string) error
}

type Context struct {
	Translators  []WordTranslator
	WordSkippers []WordSkipper
}

var (
	DefaultContext = &Context{
		Translators: []WordTranslator{
			VowelFirstLetterTranslator{},
			ConsonantFirstLettersXRTranslator{},
			ConsonantSoundFirstLetterTranslator{},
		},
		WordSkippers: []WordSkipper{
			ShortFormWordSkipper{},
			InvalidWordSkipper{},
		},
	}

	vowels     = "aeiou"
	consonants = "bcdfghjklmnpqrstvwxzy"
)

func TranslateWord(context *Context, word string) (string, error) {
	word = strings.TrimSpace(word)
	gopherWord, err := translateWord(context, word)

	if err == nil {
		go WordsHistory.Cache(word, gopherWord)
	}

	return gopherWord, err
}

func TranslateSentence(context *Context, sentence string) (string, error) {
	words := strings.Split(sentence, " ")
	builder := strings.Builder{}

	for idx, word := range words {

		// If the space between the words is more than one ignore it
		if word == "" {
			continue
		}

		if idx == len(words)-1 {
			// exclude the sentence's last punctuation mark
			word = word[:len(word)-1]
		}

		gopherWord, err := translateWord(context, word)
		if err != nil {
			return "", errors.New("Could not translate sentence: " + err.Error())
		}

		builder.WriteString(gopherWord)
		if idx == len(words)-1 {
			// Finish the translation with the same punctuation mark
			builder.WriteString(sentence[len(sentence)-1:])
		} else {
			builder.WriteString(" ")
		}
	}

	gopherSentence := builder.String()
	go SentencesHistory.Cache(sentence, gopherSentence)

	return gopherSentence, nil
}

func translateWord(context *Context, word string) (string, error) {
	for _, skipper := range context.WordSkippers {
		if err := skipper.SkipWord(word); err != nil {
			return "", err
		}
	}

	for _, translator := range context.Translators {
		if gopherWord := translator.TranslateWord(word); gopherWord != "" {
			return gopherWord, nil
		}
	}

	return "", errors.New(word + ": was not matched by any of the translators")
}

// Skip words in short form
// e.g. don't, I'm
type ShortFormWordSkipper struct{}
type ShortFormSkipperError struct{ Word string }

func (err ShortFormSkipperError) Error() string {
	return err.Word + ": is in short form"
}

func (skipper ShortFormWordSkipper) SkipWord(word string) error {
	matched, err := regexp.MatchString(`[a-zA-Z]'[a-z]`, word)
	if err != nil {
		return err
	} else if matched {
		return &ShortFormSkipperError{Word: word}
	} else {
		return nil
	}
}

// Skip wrods which doesn't contain only letters
type InvalidWordSkipper struct{}
type InvalidWordError struct{ Word string }

func (err InvalidWordError) Error() string {
	return err.Word + ": must contain only letters a-z or A-Z"
}

func (skipper InvalidWordSkipper) SkipWord(word string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z]{1,}$`, word)
	if err != nil {
		return err
	} else if matched == false {
		return &InvalidWordError{Word: word}
	} else {
		return nil
	}
}

// If a word starts with a vowel letter, add prefix "g" to the word
// e.g. apple => gapple
type VowelFirstLetterTranslator struct{}

func (translator VowelFirstLetterTranslator) TranslateWord(word string) string {
	wordToLower := strings.ToLower(word)

	if i := strings.IndexAny(wordToLower, vowels); i == 0 {
		return "g" + wordToLower
	}

	return ""
}

// If a word starts with the consonant letters "xr", add the prefix "ge"
// to the begging of the word
// e.g. xray => gexray
type ConsonantFirstLettersXRTranslator struct{}

func (translator ConsonantFirstLettersXRTranslator) TranslateWord(word string) string {
	wordToLower := strings.ToLower(word)

	if strings.HasPrefix(wordToLower, "xr") {
		return "ge" + wordToLower
	}

	return ""
}

// If a word starts with a consonant sound, move it to the end of the
// word and then add "ogo" suffix to the word
// e.g. chair => airchogo
// or
// If a word starts with a consonant sound followed by "qu", move it
// to the end of the word and then add "ogo" suffix to the word
// e.g. square => aresquogo
//
// The second rule has precedence over the first
// e.g. square is a valid case for the first rule, but it will be
// matched by the second
type ConsonantSoundFirstLetterTranslator struct {}

func (translator ConsonantSoundFirstLetterTranslator) TranslateWord(word string) string {
	var (
		wordToLower = strings.ToLower(word)
		gopherWord string
		prefix = extractConsonantSoundPrefix(wordToLower)
		prefixQU = prefix + "qu"
	)

	if prefix == "" {
		return ""
	} else if strings.HasPrefix(wordToLower, prefixQU) {
		gopherWord = strings.Replace(wordToLower, prefixQU, "", 1)
		gopherWord = gopherWord + prefixQU + "ogo"
	} else {
		gopherWord = strings.Replace(wordToLower, prefix, "", 1)
		gopherWord = gopherWord + prefix + "ogo"
	}

	return gopherWord
}

func extractConsonantSoundPrefix(word string) string {
	var (
		reader = strings.NewReader(word)
		prefix = strings.Builder{}
		q      = 'q'
		u      = 'u'
	)

	for ch, size, _ := reader.ReadRune(); size != 0; ch, size, _ = reader.ReadRune() {
		if ch == q {
			ch, size, _ = reader.ReadRune()

			if size == 0 || // q is the last letter
				ch == u && prefix.Len() == 0 { // u is the second letter
				prefix.WriteRune(q)
				break
			} else if ch == u {
				break
			} else {
				// if not qu write q to prefix
				prefix.WriteRune(q)

				// re-read the current rune again
				// in order to cover cases like 'qqu'
				reader.Seek(-1, io.SeekCurrent)
				continue
			}
		}

		if strings.ContainsRune(consonants, ch) {
			prefix.WriteRune(ch)
		} else {
			break
		}
	}

	return prefix.String()
}

package translation

import (
	"testing"
)

func TestWordTranslate_Apple(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "apple")
	if tw != "gapple" {
		t.Errorf("Expected: apple => gapple. Got: %s", tw)
	}
}

func TestWordTranslate_Ear(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "ear")
	if tw != "gear" {
		t.Errorf("Expected: ear => gear. Got: %s", tw)
	}
}

func TestWordTranslate_Context(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "context")
	if tw != "ontextcogo" {
		t.Errorf("Expected: context => ontextcogo. Got: %s", tw)
	}
}

func TestWordTranslate_Whitespaces(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "     squared     ")
	if tw != "aredsquogo" {
		t.Errorf("Expected: squared => aredsquogo. Got: %s", tw)
	}
}

func TestWordTranslate_XR(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "xriphone")
	if tw != "gexriphone" {
		t.Errorf("Expected: xriphone => gexriphone. Got: %s", tw)
	}
}

func TestWordTranslaate_CH(t *testing.T) {
	tw, _ := TranslateWord(DefaultContext, "chewbacca")
	if tw != "ewbaccachogo" {
		t.Errorf("Expected: chewbacca => ewbaccachogo. Got: %s", tw)
	}
}

func Test_SingleLetter(t *testing.T) {
	var tw string

	tw, _ = TranslateWord(DefaultContext, "q")
	if tw != "qogo" {
		t.Errorf("Expected: q => qogo. Got: %s", tw)
	}

	tw, _ = TranslateWord(DefaultContext, "")
	if tw != "" {
		t.Errorf("Expected: Empty result. Got: %s", tw)
	}
}

func Test_QU(t *testing.T) {
	var tw string

	tw, _ = TranslateWord(DefaultContext, "quake")
	if tw != "uakeqogo" {
		t.Errorf("Expected: quake => uakeqogo. Got: %s", tw)
	}

	tw, _ = TranslateWord(DefaultContext, "qu")
	if tw != "uqogo" {
		t.Errorf("Expected: qu => uqogo. Got: %s", tw)
	}

	tw, _ = TranslateWord(DefaultContext, "qqu")
	if tw != "qquogo" {
		t.Errorf("Expected: qqu => qquogo. Got: %s", tw)
	}

	tw, _ = TranslateWord(DefaultContext, "sssqqquuu")
	if tw != "uusssqqquogo" {
		t.Errorf("Expected: sssqqquuu => uusssqqquogo. Got: %s", tw)
	}

	tw, _ = TranslateWord(DefaultContext, "sqquare")
	if tw != "aresqquogo" {
		t.Errorf("Expected: sqquare => aresqquogo. Got: %s", tw)
	}
}

func TestInvalidInputs(t *testing.T) {
	var err error

	_, err = TranslateWord(DefaultContext, "I'm")
	errSkipper1 := &ShortFormSkipperError{Word: "I'm"}
	if err == nil || err.(*ShortFormSkipperError).Error() != errSkipper1.Error() {
		t.Errorf("Expected: %s", errSkipper1.Error())
	}

	_, err = TranslateWord(DefaultContext, "woo   oord")
	errSkipper2 := &InvalidWordError{Word: "woo   oord"}
	if err == nil || err.(*InvalidWordError).Error() != errSkipper2.Error() {
		t.Errorf("Expected: %s", errSkipper2.Error())
	}

	_, err = TranslateWord(DefaultContext, "google.com")
	errSkipper3 := &InvalidWordError{Word: "google.com"}
	if err == nil || err.(*InvalidWordError).Error() != errSkipper3.Error() {
		t.Errorf("Expected: %s", errSkipper3.Error())
	}
}

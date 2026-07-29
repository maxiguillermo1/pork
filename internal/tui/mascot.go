package tui

import (
	"strings"
	"time"
)

// pork's mascot is a pig that sniffs out torrents for you. Its mood follows what
// the app is doing - asleep when idle, curious while hunting, content while
// seeding - a small bit of warmth in the corners of the UI.

type mood int

const (
	moodSleeping mood = iota
	moodCurious
	moodBusy
	moodHappy
)

// smallPig is a compact three-line pig whose face matches the mood.
func smallPig(m mood) []string {
	face := "-.-"
	switch m {
	case moodCurious:
		face = "o.o"
	case moodBusy:
		face = ">.<"
	case moodHappy:
		face = "^.^"
	}
	return []string{
		`  ^~~^`,
		`( ` + face + ` )`,
		` > ᴥ <`,
	}
}

// sleepingPig is a curled-up pig for cozy empty states.
var sleepingPig = []string{
	"    (~~)",
	"   (    )",
	"   |____|",
}

// renderPig styles a pig block, padded to a rectangle so it stays aligned
// when centered. A sleeping pig gets a "z Z" beside its head.
func renderPig(lines []string, m mood) string {
	pigLines := append([]string(nil), lines...)
	if m == moodSleeping && len(pigLines) > 0 {
		pigLines[0] += "  ~ z Z"
	}
	w := 0
	for _, l := range pigLines {
		if len(l) > w {
			w = len(l)
		}
	}
	for i, l := range pigLines {
		pigLines[i] = l + strings.Repeat(" ", w-len(l))
	}
	return styleDim.Render(strings.Join(pigLines, "\n"))
}

func cozyGreeting() string {
	switch h := time.Now().Hour(); {
	case h < 5:
		return "still up?"
	case h < 12:
		return "good morning"
	case h < 18:
		return "good afternoon"
	case h < 22:
		return "good evening"
	default:
		return "good night"
	}
}

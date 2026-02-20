// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"fmt"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/tview"
)

// FKeyBar renders a persistent F-key navigation legend at the bottom of the TUI.
type FKeyBar struct {
	*tview.TextView
	styles *config.Styles
}

// NewFKeyBar returns a new F-key bar.
func NewFKeyBar(styles *config.Styles) *FKeyBar {
	f := &FKeyBar{
		TextView: tview.NewTextView(),
		styles:   styles,
	}
	f.SetBackgroundColor(styles.BgColor())
	f.SetDynamicColors(true)
	f.SetTextAlign(tview.AlignCenter)
	f.SetBorderPadding(0, 0, 0, 0)
	f.refresh()
	styles.AddListener(f)

	return f
}

// StylesChanged notifies skin changed.
func (f *FKeyBar) StylesChanged(s *config.Styles) {
	f.styles = s
	f.SetBackgroundColor(s.BgColor())
	f.refresh()
}

func (f *FKeyBar) refresh() {
	f.Clear()
	keyColor := "[green::b]"
	sepColor := "[gray::-]"
	descColor := "[white::-]"
	reset := "[-::-]"

	legend := fmt.Sprintf(
		"%sF1%s%sHome%s %s│%s %sF2%s%sRancher%s %s│%s %sF3%s%sDistro%s %s│%s %sF4%s%setcd%s %s│%s %sF5%s%sNodes%s %s│%s %sF6%s%sFleet%s %s│%s %sF7%s%sLH%s %s│%s %sF8%s%sVMs%s %s│%s %sF9%s%sInfo%s %s│%s %sF10%s%sCtx%s",
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset, sepColor, reset,
		keyColor, reset, descColor, reset,
	)
	fmt.Fprint(f, legend)
}

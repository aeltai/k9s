// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui

import (
	"fmt"
	"strings"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/tview"
)

// LogoSmall rk9s small logo — shown in the top-left header while running.
var LogoSmall = []string{
	`rK9s`,
	`────────────────────`,
	`SUSE Rancher · K8s  `,
	`Opinionated TUI     `,
	``,
	``,
}

// LogoBig rk9s big logo for splash page.
var LogoBig = []string{
	``,
	`  rK9s`,
	`  ─────────────────────────────`,
	`  SUSE Rancher · Kubernetes TUI`,
	`  Opinionated multi-cluster ops`,
	``,
}

// Splash represents a splash screen.
type Splash struct {
	*tview.Flex
}

// NewSplash instantiates a new splash screen with product and company info.
func NewSplash(styles *config.Styles, version string) *Splash {
	s := Splash{Flex: tview.NewFlex()}
	s.SetBackgroundColor(styles.BgColor())

	logo := tview.NewTextView()
	logo.SetDynamicColors(true)
	logo.SetTextAlign(tview.AlignCenter)
	s.layoutLogo(logo, styles)

	vers := tview.NewTextView()
	vers.SetDynamicColors(true)
	vers.SetTextAlign(tview.AlignCenter)
	s.layoutRev(vers, version, styles)

	s.SetDirection(tview.FlexRow)
	s.AddItem(logo, 10, 1, false)
	s.AddItem(vers, 1, 1, false)

	return &s
}

func (*Splash) layoutLogo(t *tview.TextView, styles *config.Styles) {
	c := styles.Body().LogoColor
	_, _ = fmt.Fprintf(t, "%s", strings.Repeat("\n", 2))
	for i, line := range LogoBig {
		if i == 1 && len(line) > 3 {
			// "  rK9s" — color 'r' green, 'K9s' in logo color
			_, _ = fmt.Fprintf(t, "[%s::b]%s[green::b]r[%s::b]%s", c, line[:2], c, line[3:])
		} else {
			_, _ = fmt.Fprintf(t, "[%s::b]%s", c, line)
		}
		if i+1 < len(LogoBig) {
			_, _ = fmt.Fprintf(t, "\n")
		}
	}
	_, _ = fmt.Fprintf(t, "\n")
}

func (*Splash) layoutRev(t *tview.TextView, rev string, styles *config.Styles) {
	_, _ = fmt.Fprintf(t, "[%s::b]Revision [red::b]%s", styles.Body().FgColor, rev)
}

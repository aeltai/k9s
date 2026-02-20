// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package ui_test

import (
	"testing"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/k9s/internal/ui"
	"github.com/stretchr/testify/assert"
)

func TestNewLogoView(t *testing.T) {
	v := ui.NewLogo(config.NewStyles())
	v.Reset()

	const elogo = "[green::b] ____[#ffa500::b]  _ __  __ ________     \n[green::b]| _ \\[#ffa500::b]| |\\ \\/ /   __   \\______\n[green::b]|   /[#ffa500::b]| | \\  /\\____    /  ___/\n[green::b]|_|_\\[#ffa500::b]|_|  \\  \\  /    /\\___  \\\n[#ffa500::b]       |_|  /_/____//____  /\n[#ffa500::b]            \\/           \\/  \n"
	assert.Equal(t, elogo, v.Logo().GetText(false))
	assert.Empty(t, v.Status().GetText(false))
}

func TestLogoStatus(t *testing.T) {
	uu := map[string]struct {
		logo, msg, e string
	}{
		"info": {
			"[green::b] ____[#008000::b]  _ __  __ ________     \n[green::b]| _ \\[#008000::b]| |\\ \\/ /   __   \\______\n[green::b]|   /[#008000::b]| | \\  /\\____    /  ___/\n[green::b]|_|_\\[#008000::b]|_|  \\  \\  /    /\\___  \\\n[#008000::b]       |_|  /_/____//____  /\n[#008000::b]            \\/           \\/  \n",
			"blee",
			"[#ffffff::b]blee\n",
		},
		"warn": {
			"[green::b] ____[#c71585::b]  _ __  __ ________     \n[green::b]| _ \\[#c71585::b]| |\\ \\/ /   __   \\______\n[green::b]|   /[#c71585::b]| | \\  /\\____    /  ___/\n[green::b]|_|_\\[#c71585::b]|_|  \\  \\  /    /\\___  \\\n[#c71585::b]       |_|  /_/____//____  /\n[#c71585::b]            \\/           \\/  \n",
			"blee",
			"[#ffffff::b]blee\n",
		},
		"err": {
			"[green::b] ____[#ff0000::b]  _ __  __ ________     \n[green::b]| _ \\[#ff0000::b]| |\\ \\/ /   __   \\______\n[green::b]|   /[#ff0000::b]| | \\  /\\____    /  ___/\n[green::b]|_|_\\[#ff0000::b]|_|  \\  \\  /    /\\___  \\\n[#ff0000::b]       |_|  /_/____//____  /\n[#ff0000::b]            \\/           \\/  \n",
			"blee",
			"[#ffffff::b]blee\n",
		},
	}

	v := ui.NewLogo(config.NewStyles())
	for n := range uu {
		k, u := n, uu[n]
		t.Run(k, func(t *testing.T) {
			switch k {
			case "info":
				v.Info(u.msg)
			case "warn":
				v.Warn(u.msg)
			case "err":
				v.Err(u.msg)
			}
			assert.Equal(t, u.logo, v.Logo().GetText(false))
			assert.Equal(t, u.e, v.Status().GetText(false))
		})
	}
}

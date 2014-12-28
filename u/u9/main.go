package u9

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

// AddTabSupport is a helper that modifies a <textarea>, so that pressing tab key will insert tabs.
func AddTabSupport(textArea *dom.HTMLTextAreaElement) {
	textArea.AddEventListener("keydown", false, func(event dom.Event) {
		switch ke := event.(*dom.KeyboardEvent); {
		case ke.KeyCode == 9 && !ke.CtrlKey && !ke.AltKey && !ke.MetaKey && !ke.ShiftKey: // Tab.
			value, start, end := textArea.Value, textArea.SelectionStart, textArea.SelectionEnd

			textArea.Value = value[:start] + "\t" + value[end:]

			textArea.SelectionStart, textArea.SelectionEnd = start+1, start+1

			event.PreventDefault()

			// Trigger "input" event listeners.
			inputEvent := js.Global.Get("CustomEvent").New("input")
			textArea.Underlying().Call("dispatchEvent", inputEvent)
		}
	})
}

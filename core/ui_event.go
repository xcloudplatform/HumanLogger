package core

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"time"
)

type UIEvent struct {
	Timestamp time.Time `json:"timestamp"`

	hook.Event
}

func NewUIEvent(event hook.Event) UIEvent {
	return UIEvent{
		Timestamp: time.Now(),
		Event:     event,
	}
}

// String return formatted hook kind string
func (e UIEvent) String() string {
	switch e.Kind {
	case hook.HookEnabled:
		return fmt.Sprintf("{Kind: HookEnabled}")
	case hook.HookDisabled:
		return fmt.Sprintf("{Kind: HookDisabled}")
	case hook.KeyUp:
		return fmt.Sprintf("{Kind: KeyUp, Rawcode: %v, Keychar: %v}", e.Rawcode, e.Keychar)
	case hook.KeyHold:
		return fmt.Sprintf(
			"{Kind: KeyHold, Rawcode: %v, Keychar: %v}", e.Rawcode, e.Keychar)
	case hook.KeyDown:
		return fmt.Sprintf(
			"{Kind: KeyDown, Rawcode: %v, Keychar: %v}",
			e.Rawcode, e.Keychar)
	case hook.MouseUp:
		return fmt.Sprintf(
			"{Kind: MouseUp, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.Button, e.X, e.Y, e.Clicks)
	case hook.MouseHold:
		return fmt.Sprintf(
			"{Kind: MouseHold, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.Button, e.X, e.Y, e.Clicks)
	case hook.MouseDown:
		return fmt.Sprintf(
			"{Kind: MouseDown, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.Button, e.X, e.Y, e.Clicks)
	case hook.MouseMove:
		return fmt.Sprintf(
			"{Kind: MouseMove, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.Button, e.X, e.Y, e.Clicks)
	case hook.MouseDrag:
		return fmt.Sprintf(
			"{Kind: MouseDrag, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.Button, e.X, e.Y, e.Clicks)
	case hook.MouseWheel:
		return fmt.Sprintf(
			"{Kind: MouseWheel, Amount: %v, Rotation: %v, Direction: %v}",
			e.Amount, e.Rotation, e.Direction)
	case hook.FakeEvent:
		return fmt.Sprintf("{Kind: FakeEvent}")
	}

	return "Unknown event, contact the mantainers."
}

package core

import (
	"encoding/json"
	"fmt"
	hook "github.com/robotn/gohook"
	"strconv"
	"time"
)

type UIEvent struct {
	hook.Event
}

func NewUIEvent(event hook.Event) UIEvent {
	return UIEvent{
		Event: event,
	}
}

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

func (e UIEvent) MarshalJSON() ([]byte, error) {
	type Alias UIEvent

	switch e.Kind {
	case hook.KeyUp, hook.KeyHold, hook.KeyDown:
		// If it's a Keyboard event the relevant fields are:
		// Mask, Keycode, Rawcode, and Keychar,

		return json.Marshal(&struct {
			When    int64  `json:"timestamp"`
			Kind    string `json:"kind"`
			Mask    uint16 `json:"mask"`
			Keycode uint16 `json:"keycode"`
			Rawcode uint16 `json:"rawcode"`
			//Keychar rune   `json:"keychar"`
		}{
			When:    e.When.UnixNano() / int64(time.Millisecond),
			Kind:    e.kindToStr(),
			Mask:    e.Mask,
			Keycode: e.Keycode,
			Rawcode: e.Rawcode,
			//Keychar: e.Keychar,
		})

	case hook.MouseHold, hook.MouseUp, hook.MouseMove, hook.MouseDrag, hook.MouseDown, hook.MouseWheel:
		// If it's a Mouse event the relevant fields are:
		// Button, Clicks, X, Y, Amount, Rotation and Direction

		return json.Marshal(&struct {
			When      int64  `json:"timestamp"`
			Kind      string `json:"kind"`
			Button    uint16 `json:"button"`
			Clicks    uint16 `json:"clicks"`
			X         int16  `json:"x"`
			Y         int16  `json:"y"`
			Amount    uint16 `json:"amount"`
			Rotation  int32  `json:"rotation"`
			Direction uint8  `json:"direction"`
		}{
			When:      e.When.UnixNano() / int64(time.Millisecond),
			Kind:      e.kindToStr(),
			Button:    e.Button,
			Clicks:    e.Clicks,
			X:         e.X,
			Y:         e.Y,
			Amount:    e.Amount,
			Rotation:  e.Rotation,
			Direction: e.Direction,
		})
	default:
		return json.Marshal(&struct {
			When int64  `json:"timestamp"`
			Kind string `json:"kind"`
		}{
			When: e.When.UnixNano() / int64(time.Millisecond),
			Kind: e.kindToStr(),
		})

	}
}

func (e UIEvent) kindToStr() string {
	switch e.Kind {
	case hook.HookEnabled:
		return "HookEnabled"
	case hook.HookDisabled:
		return "HookDisabled"
	case hook.KeyDown:
		return "KeyDown"
	case hook.KeyHold:
		return "KeyHold"
	case hook.KeyUp:
		return "KeyUp"
	case hook.MouseUp:
		return "MouseUp"
	case hook.MouseHold:
		return "MouseHold"
	case hook.MouseDown:
		return "MouseDown"
	case hook.MouseMove:
		return "MouseMove"
	case hook.MouseDrag:
		return "MouseDrag"
	case hook.MouseWheel:
		return "MouseWheel"
	case hook.FakeEvent:
		return "FakeEvent"
	default:
		return "Unknown Kind " + strconv.Itoa(int(e.Kind))
	}
}

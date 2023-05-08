package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-vgo/robotgo"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

type Command struct {
	Name   string      `json:"name"`
	Param1 interface{} `json:"param1"`
	Param2 interface{} `json:"param2"`
	Param3 interface{} `json:"param3"`
}

/*
func websocket_handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

*/

func websocket_handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Read JSON message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "failed to read message from client",
			})
		}

		var cmd Command
		err = json.Unmarshal(msg, &cmd)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to parse command JSON",
			})
		}

		// Handle command
		err = handleCommand(cmd)
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		// Return executed command
		resp := map[string]interface{}{
			"executedCommand": cmd,
		}

		err = ws.WriteJSON(resp)
		if err != nil {
			c.Logger().Error(err)
		}

	}
}

func handleCommand(cmd Command) error {
	// Execute robotgo command based on received command and parameter
	switch cmd.Name {
	case "scroll":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				robotgo.Scroll(int(x), int(y))
				return nil
			}
		}
		return fmt.Errorf("invalid parameters for 'scroll' command")
	case "scrollSmooth":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				robotgo.ScrollSmooth(int(x), int(y))
				return nil
			}
		}
		return fmt.Errorf("invalid parameters for 'scrollSmooth' command")
	case "milliSleep":
		if x, ok := cmd.Param1.(float64); ok {
			robotgo.MilliSleep(int(x))
			return nil
		}
		return fmt.Errorf("invalid parameters for 'milliSleep' command")
	case "move":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				robotgo.Move(int(x), int(y))
				return nil
			}
		}
		return fmt.Errorf("invalid parameters for 'move' command")
	case "moveRelative":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				robotgo.MoveRelative(int(x), int(y))
				return nil
			}
		}
		return fmt.Errorf("invalid parameters for 'moveRelative' command")
	case "dragSmooth":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				robotgo.DragSmooth(int(x), int(y))
				return nil
			}
		}
		return fmt.Errorf("invalid parameters for 'dragSmooth' command")
	case "click":
		if btn, ok := cmd.Param1.(string); ok {
			if args, ok := cmd.Param2.([]interface{}); ok {
				if len(args) == 0 {
					robotgo.Click(btn)
					return nil
				} else if len(args) == 1 {
					if b, ok := args[0].(bool); ok {
						robotgo.Click(btn, b)
						return nil
					}
				}
			}
		}
		return fmt.Errorf("invalid parameters for 'click' command")
	case "moveSmooth":
		if x, ok := cmd.Param1.(float64); ok {
			if y, ok := cmd.Param2.(float64); ok {
				if args, ok := cmd.Param3.([]interface{}); ok {
					if len(args) == 2 {
						if dur, ok := args[0].(float64); ok {
							if accel, ok := args[1].(float64); ok {
								robotgo.MoveSmooth(int(x), int(y), dur, accel)
								return nil
							}
						}
					}
				}
			}
		}
		return fmt.Errorf("invalid parameters for 'moveSmooth' command")
	case "toggle":
		if btn, ok := cmd.Param1.(string); ok {
			if args, ok := cmd.Param2.([]interface{}); ok {
				if len(args) == 0 {
					robotgo.Toggle(btn)
					return nil
				} else if len(args) == 1 {
					if direction, ok := args[0].(string); ok {
						robotgo.Toggle(btn, direction)
						return nil
					}
				}
			}
		}
		return fmt.Errorf("invalid parameters for 'toggle' command")
	default:
		return fmt.Errorf("unknown command '%s'", cmd.Name)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "./public")
	e.GET("/ws", websocket_handler)
	e.Logger.Fatal(e.Start(":1323"))
}

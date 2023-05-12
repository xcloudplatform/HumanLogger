package ports

import "github.com/ClickerAI/ClickerAI/core"

type LoggingSession interface {
	GetID() string
}

type Logger interface {
	StartLogging() (LoggingSession, error)
	StopLogging(session LoggingSession) error
	LogScreenshot(session LoggingSession, screenshot *core.Screenshot) error
	LogUIEvent(session LoggingSession, event *core.UIEvent) error
}

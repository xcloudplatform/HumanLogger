package logger

import (
	"fmt"
	"time"

	"humanlogger/core"
	"humanlogger/core/ports"
)

type consoleLogger struct{}

type consoleLoggerSession struct {
	ID        string
	StartTime time.Time
}

func (session *consoleLoggerSession) IsNeedsRotation() bool {
	return false
}

func (session *consoleLoggerSession) GetDirPath() string {
	panic("you cant rotate console")
}

func (session *consoleLoggerSession) GetID() string {
	return session.ID
}

func NewConsoleLogger() ports.Logger {
	return &consoleLogger{}
}

func (l *consoleLogger) StartLogging() (ports.LoggingSession, error) {
	session := &consoleLoggerSession{
		ID: generateUUID(),
	}

	fmt.Printf("[%s] Logging started: %s\n", session.GetID()[:4], session.StartTime.Format(time.RFC3339))

	return session, nil
}

func (l *consoleLogger) StopLogging(session ports.LoggingSession) error {
	fmt.Printf("[%s] Logging stopped\n", session.GetID()[:4])
	return nil
}

func (l *consoleLogger) LogScreenshot(session ports.LoggingSession, screenshot *core.Screenshot) error {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] %s: %v\n", session.GetID()[:4], timestamp, screenshot)
	return nil
}

func (l *consoleLogger) LogUIEvent(session ports.LoggingSession, event *core.UIEvent) error {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] %s: %v\n", session.GetID()[:4], timestamp, event)
	return nil
}

package logger

import (
	"fmt"
	"time"

	"github.com/ClickerAI/ClickerAI/core"
	"github.com/ClickerAI/ClickerAI/core/ports"
	"github.com/google/uuid"
)

type consoleLogger struct{}

type consoleLoggerSession struct {
	ID        string
	StartTime time.Time
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
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] %s: %v\n", session.GetID()[:4], timestamp, screenshot)
	return nil
}

func (l *consoleLogger) LogUIEvent(session ports.LoggingSession, event *core.UIEvent) error {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] %s: %v\n", session.GetID()[:4], timestamp, event)
	return nil
}

func generateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return u.String()
}

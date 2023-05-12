package logger

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/ClickerAI/ClickerAI/core"
	"github.com/ClickerAI/ClickerAI/core/ports"
)

type LocalFileLogger struct {
	sessionDir string
}

type LocalFileLoggerSession struct {
	ID        string
	StartTime time.Time

	sessionDir string
}

func (session *LocalFileLoggerSession) GetID() string {
	return session.ID
}

func (session *LocalFileLoggerSession) GetDirPath() string {
	return filepath.Join("logs", session.ID)

}

func (session *LocalFileLoggerSession) GetFilePathForScreenshot(screenshot *core.Screenshot) string {
	filename := fmt.Sprintf("%d_%s_%d.png", screenshot.Timestamp.UnixMilli(), session.GetID(), screenshot.DisplayID)
	return filepath.Join(session.GetDirPath(), filename)
}

func (session *LocalFileLoggerSession) GetFilePathForEvents() string {
	filename := fmt.Sprintf("%d_%s.json", session.StartTime.UnixMilli(), session.GetID())
	return filepath.Join(session.GetDirPath(), filename)

}

func NewLocalFileLogger() ports.Logger {
	return &LocalFileLogger{}
}

func (l *LocalFileLogger) StartLogging() (ports.LoggingSession, error) {
	session := &LocalFileLoggerSession{
		ID:        generateUUID(),
		StartTime: time.Now(),
	}

	l.sessionDir = session.GetDirPath()
	if err := os.MkdirAll(l.sessionDir, 0755); err != nil {
		return nil, err
	}

	logFile := session.GetFilePathForEvents()
	if _, err := os.Create(logFile); err != nil {
		return nil, err
	}

	return session, nil
}

func (l *LocalFileLogger) StopLogging(session ports.LoggingSession) error {
	//todo close
	return nil
}

func (l *LocalFileLogger) LogScreenshot(session ports.LoggingSession, screenshot *core.Screenshot) error {
	localSession, ok := session.(*LocalFileLoggerSession)
	if !ok {
		return fmt.Errorf("session is not of type LocalFileLoggerSession")
	}

	path := localSession.GetFilePathForScreenshot(screenshot)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := png.Encode(file, screenshot.Image); err != nil {
		return err
	}
	return nil
}
func (l *LocalFileLogger) LogUIEvent(session ports.LoggingSession, event *core.UIEvent) error {
	localSession, ok := session.(*LocalFileLoggerSession)
	if !ok {
		return fmt.Errorf("session is not of type LocalFileLoggerSession")
	}

	path := localSession.GetFilePathForEvents()

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}
	eventJSON = append(eventJSON, '\n')
	if _, err := file.Write(eventJSON); err != nil {
		return err
	}
	return nil
}

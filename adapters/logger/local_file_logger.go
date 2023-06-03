package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/tiff"

	"humanlogger/core"
	"humanlogger/core/ports"
)

type LocalFileLogger struct {
	sessionDir string
}

type LocalFileLoggerSession struct {
	ID        string
	StartTime time.Time

	sessionDir string
}

const LocalSessionSizeRotation = 30 * 1024 * 1024

func (session *LocalFileLoggerSession) IsNeedsRotation() bool {
	dirPath := session.GetDirPath()

	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return false
	}

	// Walk through all the files in the directory and sum up their sizes
	var totalSize int64
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return false
	}

	// Check if the total size of the directory is bigger than 5MB (5*1024*1024 bytes)
	return totalSize > LocalSessionSizeRotation
}

func (session *LocalFileLoggerSession) GetID() string {
	return session.ID
}

func (session *LocalFileLoggerSession) GetDirPath() string {
	return filepath.Join("logs", session.ID)

}

func (session *LocalFileLoggerSession) GetFilePathForScreenshot(screenshot *core.Screenshot) string {
	filename := fmt.Sprintf("%d_%s_%d.tiff", screenshot.Timestamp.UnixMilli(), session.GetID(), screenshot.DisplayID)
	return filepath.Join(session.GetDirPath(), filename)
}

func (session *LocalFileLoggerSession) GetFilePathForEvents() string {
	filename := fmt.Sprintf("%d_%s.json", session.StartTime.UnixMilli(), session.GetID())
	return filepath.Join(session.GetDirPath(), filename)

}

func NewLocalFileLogger() ports.Logger {

	//ocr := NewOcr()
	//go ocr.ProcessQueue()

	return &LocalFileLogger{
		//ocr: *ocr,
	}
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

	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] %s: started logging\n", session.GetID()[:4], timestamp)
	return session, nil
}

func (l *LocalFileLogger) StopLogging(session ports.LoggingSession) error {
	//todo close
	return nil
}
func (l *LocalFileLogger) LogScreenshot(session ports.LoggingSession, screenshot *core.Screenshot) error {
	timestamp := time.Now().Format("15:04:05.000")

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

	// Save the screenshot image as a TIFF file with compression and predictor options
	if err := tiff.Encode(file, screenshot.Image, &tiff.Options{Compression: tiff.Deflate, Predictor: true}); err != nil {
		return err
	}
	fmt.Printf("[%s] %s: %v %s\n", session.GetID()[:4], timestamp, screenshot, path)

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

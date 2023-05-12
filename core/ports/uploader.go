package ports

type LogUploader interface {
	Upload(session *LoggingSession) error
}

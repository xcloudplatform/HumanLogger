package ports

type LogUploader interface {
	Upload(session *LoggingSession, packedPath string) error
}

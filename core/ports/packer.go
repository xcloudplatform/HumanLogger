package ports

type LogPacker interface {
	Pack(session *LoggingSession) (string, error)
}

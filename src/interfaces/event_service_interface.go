package interfaces

type EventServiceInterface interface {
	Log(level, service, message string)
	Startup(logLevel string)
}

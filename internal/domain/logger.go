package domain

type Logger interface {
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
	Fatal(format string, v ...any)
	Stop()
}

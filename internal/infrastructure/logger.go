package infrastructure

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	def "github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

var _ def.Logger = (*asyncLogger)(nil)

type asyncLogger struct {
	entries chan string
	done    chan struct{}
	file    *os.File
}

const (
	logsDir  = ".logs"
	logsFile = "app.log"
)

func ensureLogDirectory() error {
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}
	return nil
}

func NewLogger(bufSize int) *asyncLogger {
	if err := ensureLogDirectory(); err != nil {
		log.Fatalf("could not create logs directory: %v", err)
	}

	logFilePath := fmt.Sprintf("%s%s%s", logsDir, "/", logsFile)
	const logFileFlags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	const logFilePerm os.FileMode = 0o644

	logFile, err := os.OpenFile(logFilePath, logFileFlags, logFilePerm)
	if err != nil {
		log.Fatalf("could not open log file: %v", err)
	}
	multi := io.MultiWriter(os.Stdout, logFile)

	l := &asyncLogger{
		entries: make(chan string, bufSize),
		done:    make(chan struct{}),
		file:    logFile,
	}

	log.SetOutput(multi)

	go l.run()
	return l
}

func (l *asyncLogger) run() {
	for msg := range l.entries {
		timestamp := time.Now().Format(time.RFC3339)
		logLine := fmt.Sprintf("%s %s\n", timestamp, msg)
		fmt.Print(logLine)
		_, err := l.file.WriteString(logLine)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing to log file: %v\n", err)
		}
	}

	l.file.Close()
	close(l.done)
}

func (l *asyncLogger) Infof(format string, v ...any) {
	l.send("[INFO] "+format, v...)
}

func (l *asyncLogger) Errorf(format string, v ...any) {
	l.send("[ERROR] "+format, v...)
}

func (l *asyncLogger) Fatal(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	full := fmt.Sprintf("[FATAL] %s", msg)

	l.entries <- full
	l.Stop()

	log.Fatal(msg)
}

func (l *asyncLogger) send(format string, v ...any) {
	select {
	case l.entries <- fmt.Sprintf(format, v...):
	default:
	}
}

func (l *asyncLogger) Stop() {
	close(l.entries)
	<-l.done
}

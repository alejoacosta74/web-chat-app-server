package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/alejoacosta74/chatserver/server"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app            = kingpin.New("chat-server", "A real time multithreading chat application server.")
	logFile        = app.Flag("log-file", "write logs to a file").Envar("LOG_FILE").Default("").String()
	devMode        = app.Flag("dev", "[Insecure] Developer mode").Envar("DEV").Default("false").Bool()
	serverIP       = app.Flag("ip address", "Server address.").Default("127.0.0.1").IP()
	serverPort     = app.Flag("port", "Server port number.").Default("8000").Int()
	singleThreaded = app.Flag("singleThreaded", "[Non-production] Process http request in a single thread").Envar("SINGLE_THREADED").Default("false").Bool()
)

func action(pc *kingpin.ParseContext) error {
	addr := fmt.Sprintf("%s:%d", *serverIP, *serverPort)
	writers := []io.Writer{os.Stdout}

	// If 'logfile' was provided, create a file for loggin
	if logFile != nil && (*logFile) != "" {
		_, err := os.Stat(*logFile)
		if os.IsNotExist(err) {
			newLogFile, err := os.Create(*logFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to create log file %s", *logFile)
			} else {
				writers = append(writers, newLogFile)
			}
		} else {
			existingLogFile, err := os.Open(*logFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to open log file %s", *logFile)
			} else {
				writers = append(writers, existingLogFile)
			}
		}
	}

	logWriter := io.MultiWriter(writers...)
	logger := log.NewLogfmtLogger(logWriter)

	if !*devMode {
		logger = level.NewFilter(logger, level.AllowWarn())
	}

	s, err := server.New(
		addr,
		server.SetLogWriter(logWriter),
		server.SetLogger(logger),
		server.SetDebug(*devMode),
		server.SetSingleThreaded(*singleThreaded),
	)
	if err != nil {
		return errors.Wrap(err, "server#New")
	}

	return s.Start()
}

func Run() {
	app.Version(VersionWithGitSha)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func init() {
	app.Action(action)
}

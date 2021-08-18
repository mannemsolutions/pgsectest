package internal

import (
	"os"
	"strings"

	"github.com/mannemsolutions/pgsectest/pkg/pg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.SugaredLogger
	atom zap.AtomicLevel
)

func Initialize() {
	atom = zap.NewAtomicLevel()
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	log = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).Sugar()

	pg.Initialize(log)
}

func Handle() {
	var scores float64
	var maxScores float64
	configs, err := GetConfigs()
	if err != nil {
		log.Errorf("could not parse all configs: %s", err.Error())
		os.Exit(125)
	}
	for _, config := range configs {
		name := config.Name()
		log.Debugf(strings.Repeat("=", 19+len(name)))
		log.Debugf("Running tests from %s", name)
		log.Debugf(strings.Repeat("=", 19+len(name)))
		if config.Debug {
			atom.SetLevel(zapcore.DebugLevel)
		} else {
			atom.SetLevel(zapcore.InfoLevel)
		}
		conn := pg.NewConn(config.DSN, config.Retries, config.Delay)
		for i, test := range config.Tests {
			flawLess := test.Score.Flawless()
			maxScores += flawLess
			if err = test.Validate(); err != nil {
				log.Errorf("Test %d (%s): Invalid test: %s", i, test.Name, err.Error())
			} else if result, err := conn.RunQuery(test.Query); err != nil {
				log.Errorf("Test %d (%s): error occurred while running query : %s", i, test.Name, err.Error())
			} else {
				if value, err := result.OneField(); err != nil {
					log.Errorf("Test %d (%s): error occurred while retrieving field: %s", i, test.Name, err.Error())
				} else if f, err := value.AsFloat(); err != nil {
					log.Errorf("Test %d (%s): field is no valid numeric: %s", i, test.Name, err.Error())
				} else {
					score := test.Score.FromResult(f)
					log.Debugf("Score for test %d: %f out of %f", i, score, flawLess)
					scores += score

				}
			}
		}
	}
	log.Infof("Score: %.2f%% (%.2f out of %.2f)", 100*scores/maxScores, scores, maxScores)
}

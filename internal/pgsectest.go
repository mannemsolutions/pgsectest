package internal

import (
	"os"
	"strconv"
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

func getResult(query string, conn *pg.Conn, defValue float64) (float64, error) {
	if query == "" {
		return defValue, nil
	} else if f, err := strconv.ParseFloat(query, 64); err == nil {
		return f, nil
	} else if result, err := conn.RunQuery(query); err != nil {
		return 0, err
	} else if value, err := result.OneField(); err != nil {
		return 0, err
	} else if f, err = value.AsFloat(); err != nil {
		return 0, err
	} else {
		return f, nil
	}
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
			} else if dividend, err := getResult(test.Dividend, conn, float64(test.Score.Min)); err != nil {
				log.Errorf("Test %d (%s): error occurred while running dividend query : %s", i, test.Name, err.Error())
			} else if divisor, err := getResult(test.Divisor, conn, 1); err != nil {
				log.Errorf("Test %d (%s): error occurred while running dividend query : %s", i, test.Name, err.Error())
			} else {
				score := test.Score.FromResult(dividend, divisor)
				if config.Verbosity > 2 || (config.Verbosity > 0 && score < flawLess) {
					log.Infof("Score for test %d (%s): %.2f out of %.2f", i, test.Name, score, flawLess)
				}
				if config.Verbosity > 1 && score < flawLess {
					if test.Advice != "" {
						log.Infof("  | You can improve your score by:")
						for _, line := range strings.Split(test.Advice, "\n") {
							log.Infof("  | %s", line)
						}
					}
					if test.Url != "" {
						log.Infof("  |")
						log.Infof("  | For more info, please see <%s>.", test.Url)
					}
				}
				scores += score
			}
		}
	}
	log.Infof("Score: %.2f%% (%.2f out of %.2f)", 100*scores/maxScores, scores, maxScores)
}

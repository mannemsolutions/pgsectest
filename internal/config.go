package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/pflag"

	"github.com/mannemsolutions/pgsectest/pkg/pg"
	"gopkg.in/yaml.v3"
)

/*
 * This module reads the config file and returns a config object with all entries from the config yaml file.
 */

type Configs []Config

type Config struct {
	path      string
	index     int
	Debug     bool  `yaml:"debug"`
	Verbosity int   `yaml:"verbosity"`
	Tests     Tests `yaml:"tests"`

	Delay   time.Duration `yaml:"delay"`
	Retries uint          `yaml:"retries"`
	DSN     pg.Dsn        `yaml:"dsn"`
}

func (c Config) Name() (name string) {
	return fmt.Sprintf("%s (%d)", c.path, c.index)
}

func NewConfigsFromReader(reader io.Reader, name string) (configs Configs, err error) {
	var i int
	decoder := yaml.NewDecoder(reader)
	for {
		// create new spec here
		config := new(Config)
		// pass a reference to spec reference
		err := decoder.Decode(&config)
		// check it was parsed
		if config == nil {
			continue
		}
		// break the loop in case of EOF
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return Configs{}, err
		}
		config.path = name
		config.index = i
		i += 1
		configs = append(configs, *config)
	}
	return configs, nil
}
func NewConfigsFromFile(path string) (c Configs, err error) {
	// This only parsed as yaml, nothing else
	// #nosec
	reader, err := os.Open(path)
	if err != nil {
		return c, err
	}
	return NewConfigsFromReader(reader, path)
}

func NewConfigsFromStdin() (configs Configs, err error) {
	reader := bufio.NewReader(os.Stdin)
	return NewConfigsFromReader(reader, "(stdin)")
}

// ReadFromFileOrDir returns an array of Configs parsed from all yaml files, found while recursively walking
// through a directory, while following symlinks as needed.
func ReadFromFileOrDir(path string) (configs Configs, err error) {
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return Configs{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Configs{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return Configs{}, err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		// file is a directory
		entries, err := file.ReadDir(0)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
		// I want the entries sorted, so adding them to a list of strings
		var entryNames []string
		for _, entry := range entries {
			entryNames = append(entryNames, entry.Name())
		}
		sort.Strings(entryNames)
		for _, entryName := range entryNames {
			newConfigs, err := ReadFromFileOrDir(filepath.Join(path, entryName))
			if err != nil {
				_ = file.Close()
				return Configs{}, err
			}
			configs = append(configs, newConfigs...)
		}
	} else {
		// file is not a directory
		configs, err = NewConfigsFromFile(path)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
	}
	return configs, file.Close()
}

func GetConfigs() (configs Configs, err error) {
	var debug *bool
	var version *bool
	debug = pflag.BoolP("debug", "d", false, "Add debugging output")
	version = pflag.BoolP("version", "V", false, "Show version information")
	verbosity := pflag.CountP("verbose", "v", "Make output more verbose")

	pflag.Parse()
	if *version {
		fmt.Println(appVersion)
		os.Exit(0)
	}
	paths := pflag.Args()
	if len(paths) == 0 {
		log.Info("Reading tests from stdin")
		return NewConfigsFromStdin()
	}
	for _, path := range paths {
		newConfigs, err := ReadFromFileOrDir(path)
		if err != nil {
			return Configs{}, nil
		}
		configs = append(configs, newConfigs...)
	}

	for i := range configs {
		configs[i].Debug = configs[i].Debug || *debug
		if configs[i].Verbosity == 0 {
			configs[i].Verbosity = *verbosity
		}
	}

	return configs, err
}

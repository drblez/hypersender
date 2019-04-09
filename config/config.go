package config

import (
	"github.com/jessevdk/go-flags"
	"github.com/joomcode/errorx"
)

var (
	Errors       = errorx.NewNamespace("config")
	CommonErrors = Errors.NewType("common_error")
	UnknownFlag  = Errors.NewType("unknown_flags")
)

type Config struct {
	Debug          bool   `long:"debug" description:"Debug level logging" env:"DEBUG"`
	Console        bool   `long:"console" description:"Output to console" env:"CONSOLE"`
	Path           string `short:"p" long:"path" default:"." description:"Path to scan"`
	URL            string `short:"u" long:"url" description:"URL to send" required:"true"`
	PathSubst      bool   `short:"s" long:"path-substitution" description:"Substitute file name in place of %f in URL"`
	SubstString    string `short:"q" long:"substitute-sequence" default:"%f" description:"Change default sequence '%f' to user sequence"`
	LogPath        string `long:"log-path" default:"." description:"Path to save log"`
	FSParallelism  int    `short:"f" long:"fs-parallelism" default:"10"`
	NetParallelism int    `short:"n" long:"net-parallelism" default:"10"`
	ContentType    string `short:"t" long:"content-type" default:"application/json"`
}

func Init() (*Config, error) {
	config := &Config{}
	f := flags.NewParser(config, flags.Default)
	_, err := f.Parse()
	if err != nil {
		switch err := err.(type) {
		case *flags.Error:
			switch err.Type {
			case flags.ErrUnknownFlag:
				return nil, UnknownFlag.New(err.Message)
			}
		}
		return nil, CommonErrors.Wrap(err, "Config parse error")
	}
	return config, nil
}
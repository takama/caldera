package commands

import (
	"strings"

	"{{[ .Project ]}}/pkg/config"
	"{{[ .Project ]}}/pkg/helper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type bind struct {
	src  string
	dest string
}

func bindFlags(src, dest string, flags []bind, cmd *cobra.Command) {
	for _, flag := range flags {
		helper.LogF("Flag error",
			viper.BindPFlag(src+"."+flag.src, cmd.PersistentFlags().Lookup(dest+"-"+flag.dest)),
		)
	}
}

func bindEnvs(env string, flags []bind) {
	for _, flag := range flags {
		helper.LogF("Env error", viper.BindEnv(env+"."+flag.src))
	}
}

func bindCustomEnvs(src, dest string, flags []bind) {
	for _, flag := range flags {
		helper.LogF("Env error",
			viper.BindEnv(src+"."+flag.src,
				strings.ToUpper(config.ServiceName+"."+dest+"."+flag.src)),
		)
	}
}

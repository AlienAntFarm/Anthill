package utils

import (
	"flag"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strconv"
)

const configName = "anthive.toml"

type Configuration struct {
	Debug    bool
	Dev      bool
	Host     string
	Port     int
	Url      string
	Database struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     int
	}
	Minio struct {
		User     string
		Password string
		Host     string
		Port     int
	}
	Assets struct {
		Images string
	}
}

func PreRun(cmd *cobra.Command, args []string) {
	// reinit args for glog
	os.Args = os.Args[:1]

	// load configuration
	err := viper.ReadInConfig()
	if viper.GetBool("Dev") {
		viper.Set("Assets.Images", path.Join(".", "static", "images"))
		if err := os.MkdirAll(viper.GetString("Assets.Images"), 0755); err != nil {
			glog.Fatalf("%s", err)
		}
	}
	if err != nil {
		glog.Fatalf("when reading config file: %s", err)
	}
	err = viper.Unmarshal(Config)
	if err != nil {
		glog.Fatalf("when unmarshalling the json: %s", err)
	}
	if Config.Debug { // debug is same as -vvvvv
		verbosity = 5
	}
	flag.Set("v", strconv.Itoa(verbosity))
	flag.Set("logtostderr", "true")
	flag.Parse()
	glog.V(1).Infoln("debug mode enabled")
}

var (
	verbosity int
	sep       = string(os.PathSeparator)
	Config    = &Configuration{}
	Command   = &cobra.Command{
		Use:              "anthive",
		Short:            "Start anthive server",
		PersistentPreRun: PreRun,
	}
	OCICommand = &cobra.Command{
		Use:   "aif [docker tag]",
		Short: "Convert a docker image to an antling compatible container",
	}
)

func init() {
	debugMsg := "trigger debug logs, same as -vvvvv, take precedence over verbose flag"
	devMsg := "dev mode, instead of getting path from the system use those at $PWD"
	verboseMsg := "verbose output, can be stacked to increase verbosity"

	Command.PersistentFlags().Bool("debug", false, debugMsg)
	Command.PersistentFlags().Bool("dev", false, devMsg)
	Command.PersistentFlags().CountVarP(&verbosity, "verbose", "v", verboseMsg)

	Command.AddCommand(OCICommand)

	viper.BindPFlag("Debug", Command.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("Dev", Command.PersistentFlags().Lookup("dev"))

	viper.Set("PROJECT", "github.com/alienantfarm/anthive")
	// set some paths
	viper.Set("Assets.Images", path.Join(sep, "var", "lib", "antling", "images"))

	viper.SetConfigName(configName[:len(configName)-5])

	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$ANTHIVE_CONFIG" + sep)
}

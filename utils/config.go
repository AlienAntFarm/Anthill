package utils

import (
	"flag"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

type Configuration struct {
	Debug    bool   `json:"debug"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Url      string `json:"url"`
	Database struct {
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"database"`
	Assets struct {
		Static    string `json:"static"`
		Templates string `json:"templates"`
		Images    string `json:"-"`
	} `json:"assets"`
}

func PreRun(cmd *cobra.Command, args []string) {
	// reinit args for glog
	os.Args = os.Args[:1]

	// load configuration
	err := viper.ReadInConfig()
	if err != nil {
		glog.Fatalf("when reading config file: %s", err)
	}
	// set Images path
	viper.Set("Assets.Images", viper.GetString("Assets.Static")+sep+"images")
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
	verboseMsg := "verbose output, can be stacked to increase verbosity"

	Command.PersistentFlags().Bool("debug", false, debugMsg)
	Command.PersistentFlags().CountVarP(&verbosity, "verbose", "v", verboseMsg)

	Command.AddCommand(OCICommand)

	viper.BindPFlag("Debug", Command.PersistentFlags().Lookup("debug"))
	viper.Set("PROJECT", "github.com/alienantfarm/anthive")
	viper.Set("CONFIG", "ANTHIVE_CONFIG")

	viper.SetConfigName("config")
	viper.AddConfigPath("$" + viper.GetString("CONFIG") + sep)
	viper.AddConfigPath(".")
}

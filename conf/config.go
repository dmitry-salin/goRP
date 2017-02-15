package conf

import (
	"log"

	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

//Registry represents type of used service discovery server
type Registry string

const (
	//Consul service discovery
	Consul Registry = "consul"
	//Eureka service discovery
	Eureka Registry = "eureka"
)

//ServerConfig represents Main service configuration
type ServerConfig struct {
	Hostname string
	Port     int
}

//EurekaConfig represents Eureka Discovery service configuration
type EurekaConfig struct {
	URL          string
	PollInterval int
}

//ConsulConfig represents Consul Discovery service configuration
type ConsulConfig struct {
	Address      string
	Scheme       string
	Token        string
	PollInterval int
	Tags         []string
}

//RpConfig represents Composite of all app configs
type RpConfig struct {
	AppName  string
	Registry Registry
	Server   ServerConfig
	Eureka   EurekaConfig
	Consul   ConsulConfig

	rawConfig *viper.Viper
}

//Param reads parameter/property value from config (env,file,defaults)
func (cfg *RpConfig) Param(key string) interface{} {
	return cfg.rawConfig.Get(key)
}

//LoadConfig loads configuration from provided file and serializes it into RpConfig struct
func LoadConfig(file string, defaults map[string]interface{}) *RpConfig {
	var vpr = viper.New()

	if "" != file {
		vpr.SetConfigType(strings.TrimLeft(filepath.Ext(file), "."))
		vpr.SetConfigFile(file)
	}

	vpr.SetEnvPrefix("RP")
	vpr.AutomaticEnv()

	applyDefaults(vpr)
	if nil != defaults {
		for k, v := range defaults {
			vpr.SetDefault(k, v)
		}
	}

	err := vpr.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}

	var rpConf RpConfig
	vpr.Unmarshal(&rpConf)
	rpConf.rawConfig = vpr

	//vpr.Debug()
	return &rpConf
}

func applyDefaults(vpr *viper.Viper) {
	vpr.SetDefault("appname", "goRP")
	vpr.SetDefault("AuthServerURL", "http://localhost:9998/sso/me")
	vpr.SetDefault("registry", Consul)

	vpr.SetDefault("server.port", 9999)
	vpr.SetDefault("server.hostname", nil)

	vpr.SetDefault("eureka.url", "http://localhost:8761/eureka")
	vpr.SetDefault("eureka.appname", "goRP")
	vpr.SetDefault("eureka.pollInterval", 5)

	vpr.SetDefault("consul.address", "localhost:8500")
	vpr.SetDefault("consul.scheme", "http")
	vpr.SetDefault("consul.pollInterval", 5)
	vpr.SetDefault("consul.tags", nil)

}

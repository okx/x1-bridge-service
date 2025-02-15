package config

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"

	"github.com/0xPolygonHermez/zkevm-bridge-service/bridgectrl"
	"github.com/0xPolygonHermez/zkevm-bridge-service/claimtxman"
	"github.com/0xPolygonHermez/zkevm-bridge-service/coinmiddleware"
	"github.com/0xPolygonHermez/zkevm-bridge-service/config/apolloconfig"
	"github.com/0xPolygonHermez/zkevm-bridge-service/config/businessconfig"
	"github.com/0xPolygonHermez/zkevm-bridge-service/db"
	"github.com/0xPolygonHermez/zkevm-bridge-service/etherman"
	"github.com/0xPolygonHermez/zkevm-bridge-service/log"
	"github.com/0xPolygonHermez/zkevm-bridge-service/messagepush"
	"github.com/0xPolygonHermez/zkevm-bridge-service/metrics"
	"github.com/0xPolygonHermez/zkevm-bridge-service/nacos"
	"github.com/0xPolygonHermez/zkevm-bridge-service/server"
	"github.com/0xPolygonHermez/zkevm-bridge-service/server/iprestriction"
	"github.com/0xPolygonHermez/zkevm-bridge-service/server/tokenlogoinfo"
	"github.com/0xPolygonHermez/zkevm-bridge-service/synchronizer"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	Log                    log.Config
	Apollo                 apolloconfig.Config
	SyncDB                 db.Config             `apollo:"SyncDB"`
	ClaimTxManager         claimtxman.Config     `apollo:"ClaimTxManager"`
	Etherman               etherman.Config       `apollo:"Etherman"`
	Synchronizer           synchronizer.Config   `apollo:"Synchronizer"`
	BridgeController       bridgectrl.Config     `apollo:"BridgeController"`
	BridgeServer           server.Config         `apollo:"BridgeServer"`
	CoinKafkaConsumer      coinmiddleware.Config `apollo:"CoinKafkaConsumer"`
	MessagePushProducer    messagepush.Config    `apollo:"MessagePushProducer"`
	NetworkConfig          `apollo:"NetworkConfig"`
	NacosConfig            nacos.Config
	BusinessConfig         businessconfig.Config `apollo:"BusinessConfig"`
	Metrics                metrics.Config        `apollo:"Metrics"`
	IPRestriction          iprestriction.Config  `apollo:"IPRestriction"`
	TokenLogoServiceConfig tokenlogoinfo.Config  `apollo:"TokenLogoServiceConfig"`
}

// Load loads the configuration
func Load(configFilePath string, network string) (*Config, error) {
	var cfg Config
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer([]byte(DefaultValues)))
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}
	if configFilePath != "" {
		dirName, fileName := filepath.Split(configFilePath)

		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")
		fileNameWithoutExtension := strings.TrimSuffix(fileName, "."+fileExtension)

		viper.AddConfigPath(dirName)
		viper.SetConfigName(fileNameWithoutExtension)
		viper.SetConfigType(fileExtension)
	}
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("ZKEVM_BRIDGE")
	err = viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			log.Infof("config file not found")
		} else {
			log.Infof("error reading config file: ", err)
			return nil, err
		}
	}

	err = viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	if err != nil {
		return nil, err
	}

	if viper.IsSet("NetworkConfig") && network != "" {
		return nil, errors.New("Network details are provided in the config file (the [NetworkConfig] section) and as a flag (the --network or -n). Configure it only once and try again please.")
	}
	if !viper.IsSet("NetworkConfig") && network == "" {
		return nil, errors.New("Network details are not provided. Please configure the [NetworkConfig] section in your config file, or provide a --network flag.")
	}
	if !viper.IsSet("NetworkConfig") && network != "" {
		cfg.loadNetworkConfig(network)
	}

	if cfg.Apollo.Enabled {
		err = apolloconfig.Init(cfg.Apollo)
		if err != nil {
			return nil, err
		}
		err = apolloconfig.Load(&cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

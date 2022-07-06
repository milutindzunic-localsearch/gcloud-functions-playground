package conf

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	OnlimApiURL         string `mapstructure:"ONLIM_API_URL"`
	OnlimApiKey         string `mapstructure:"ONLIM_API_KEY"`
	AcceptedCategoryIDs string `mapstructure:"ACCEPTED_CATEGORY_IDS"`
}

func LoadConfigFromEnvVariables(c *Config) error {
	viper.AutomaticEnv()

	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(c, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			return bindErr
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		return err
	}

	return nil
}

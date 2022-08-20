package cfg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var NotionDBCOnfig string

type Config struct {
	DBConn          string
	TGConfig        string
	NotionSecretKey string
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath("config")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	connString := fmt.Sprintf("%s:%s@/%s",
		v.Get("config.db_conn_config.user"),
		v.Get("config.db_conn_config.password"),
		v.Get("config.db_conn_config.db_name"))

	var cfg Config

	NotionDBCOnfig = v.GetString("config.notion_db_config.config")
	cfg.TGConfig = v.GetString("config.tg_config.token")
	cfg.NotionSecretKey = v.GetString("config.notion_db_config.secret_key")
	cfg.DBConn = connString

	return &cfg, nil
}

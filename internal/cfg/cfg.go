package cfg

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var NotionDBCOnfig string

type Config struct {
	PGConfig *pgxpool.Config `yaml:"pg_config"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath("config")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config")
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s",
		v.Get("config.db_conn_config.user"),
		v.Get("config.db_conn_config.password"),
		v.Get("config.db_conn_config.host"),
		v.Get("config.db_conn_config.port"),
		v.Get("config.db_conn_config.db_name"),
		v.Get("config.db_conn_config.pool_max_conns"))

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	var cfg Config

	NotionDBCOnfig = v.GetString("config.notion_db_config.config")
	cfg.PGConfig = config

	return &cfg, nil
}

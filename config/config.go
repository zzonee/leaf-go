package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	Mode_Snowflake = 1
	Mode_Segment   = 2
)

const (
	DB_Type_MySQL = 1
	DB_Type_Redis = 2
)

type Config struct {
	Mode int

	Segment   Segment
	Snowflake Snowflake
	DB        DBConfig

	Http HttpConfig
}

type Snowflake struct {
	WorkerId int64

	WorkerIdGetter func() int64
}

type Segment struct {
	CacheDir string
}

type DBConfig struct {
	Type       int
	DataSource []string
}

type HttpConfig struct {
	Addr        string
	RequestPath string
	Query       string
}

var Global Config

func Init() error {
	v := viper.New()
	v.AddConfigPath("./")
	v.SetConfigFile("leaf.yaml")
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	Global.Mode = v.GetInt("mode")
	if Global.Mode == Mode_Snowflake {
		if err := v.UnmarshalKey("snowflake", &Global.Snowflake); err != nil {
			return err
		} else if Global.Snowflake.WorkerId >= 0 { // 否则 需要自己设置 Snowflake.WorkerIdGetter
			Global.Snowflake.WorkerIdGetter = func() int64 {
				return Global.Snowflake.WorkerId
			}
		}
	} else if Global.Mode == Mode_Segment {
		if err := v.UnmarshalKey("segment", &Global.Segment); err != nil {
			return err
		}
		if err := v.UnmarshalKey("db", &Global.DB); err != nil {
			return err
		}
	}

	if err := v.UnmarshalKey("http", &Global.Http); err != nil {
		return err
	} else if Global.Http.RequestPath == "" || Global.Http.Query == "" {
		err = fmt.Errorf("http config not valid")
		return err
	}
	return nil
}

package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-sql-driver/mysql"
	"github.com/lixiangzhong/cast"
	"github.com/lixiangzhong/viper"
	"path/filepath"
	"strings"
	"time"
)

var (
	Viper = viper.New()
)

func MustLoad(filename string) {
	if err := Load(filename); err != nil {
		panic(err)
	}
}

// 支持 JSON, TOML, YAML, INI后缀
func Load(filename string) error {
	dir := filepath.Dir(filename)
	file := filepath.Base(filename)
	Viper.AddConfigPath(dir)
	Viper.SetConfigName(strings.TrimSuffix(file, filepath.Ext(filename)))
	return Viper.ReadInConfig()
}

func WatchChange(funcs ...func()) {
	Viper.WatchConfig()
	Viper.OnConfigChange(func(fsnotify.Event) {
		for _, f := range funcs {
			f()
		}
	})
}

func SetDefault(key string, value interface{}) {
	Viper.SetDefault(key, value)
}

func Set(key string, value interface{}) {
	Viper.Set(key, value)
}

func Get(key string) interface{} {
	return Viper.Get(key)
}

func String(key string) string {
	return Viper.GetString(key)
}

func Bool(key string) bool {
	return Viper.GetBool(key)
}

func Int(key string) int {
	return Viper.GetInt(key)
}

func Int32(key string) int32 {
	return Viper.GetInt32(key)
}

func Int64(key string) int64 {
	return Viper.GetInt64(key)
}

func Uint32(key string) uint32 {
	return cast.ToUint32(Get(key))
}

func Uint64(key string) uint64 {
	return cast.ToUint64(Get(key))
}

func Float64(key string) float64 {
	return Viper.GetFloat64(key)
}

func Duration(key string) time.Duration {
	return Viper.GetDuration(key)
}

func Time(key string) time.Time {
	return cast.ToTime(Get(key))
}

func IntSlice(key string) []int {
	return cast.ToIntSlice(Get(key))
}

func StringSlice(key string) []string {
	return Viper.GetStringSlice(key)
}

func StringMapString(key string) map[string]string {
	return Viper.GetStringMapString(key)
}

func StringMapStringSlice(key string) map[string][]string {
	return Viper.GetStringMapStringSlice(key)
}

func MySQLConfig(section string) *mysql.Config {
	var cfg = mysql.NewConfig()
	fields := []string{
		"host", "addr",
		"user", "username",
		"password", "passwd",
		"db", "dbname", "database",
	}
	for _, field := range fields {
		key := fmt.Sprintf("%v.%v", section, field)
		if Viper.IsSet(key) == false {
			continue
		}
		switch field {
		case "host", "addr":
			cfg.Addr = String(key)
		case "user", "username":
			cfg.User = String(key)
		case "passwd", "password":
			cfg.Passwd = String(key)
		case "db", "dbname", "database":
			cfg.DBName = String(key)
		}
	}
	cfg.Loc = time.Local
	cfg.Net = "tcp"
	cfg.ParseTime = true
	return cfg
}

package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-sql-driver/mysql"
	"github.com/lixiangzhong/cast"
	"github.com/lixiangzhong/viper"
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
	ext := filepath.Ext(filename)
	sortSupportExtsPriority(ext)
	Viper.AddConfigPath(dir)
	Viper.SetConfigName(strings.TrimSuffix(file, ext))
	return Viper.ReadInConfig()
}

func sortSupportExtsPriority(first string) {
	first = strings.Trim(first, ".")
	if first == "" {
		return
	}
	var exts = []string{first}
	for _, ext := range viper.SupportedExts {
		switch ext {
		case first:
		default:
			exts = append(exts, ext)
		}
	}
	viper.SupportedExts = exts
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

func Size(key string) int64 {
	return parseSize(Viper.GetString(key))
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

var validUnits = []struct {
	symbol     string
	multiplier int64
}{
	{"K", 1024},
	{"M", 1024 * 1024},
	{"G", 1024 * 1024 * 1024},
	{"B", 1},
	{"", 1}, // defaulting to "B"
}

// parseSize parses the given string as size limit
// Size are positive numbers followed by a unit (case insensitive)
// Allowed units: "B" (bytes), "KB" (kilo), "MB" (mega), "GB" (giga)
// If the unit is omitted, "b" is assumed
// Returns the parsed size in bytes, or -1 if cannot parse
func parseSize(sizeStr string) int64 {
	sizeStr = strings.ToUpper(sizeStr)
	sizeStr = strings.TrimSuffix(sizeStr, "B")
	for _, unit := range validUnits {
		if strings.HasSuffix(sizeStr, unit.symbol) {
			size, err := strconv.ParseInt(sizeStr[0:len(sizeStr)-len(unit.symbol)], 10, 64)
			if err != nil {
				return -1
			}
			return size * unit.multiplier
		}
	}

	// Unreachable code
	return -1
}

package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
)

type EnvConfiguration struct {
	Environment           string   `def:"debug" desc:"Environment level" enums:"debug,dev,prod"`
	Host                  string   `def:"localhost:8000" desc:"Host for api server"`
	Port                  string   `def:"8000" desc:"Port for api server"`
	SentryEndpoint        string   `desc:"Sentry endpoint for stacktrace's'"`
	DbUsername            string   `isReq:"true" desc:"Postgres username"`
	DbHost                string   `isReq:"true" desc:"Postgres host"`
	DbPort                string   `isReq:"true" desc:"Postgres port"`
	DbPassword            string   `isReq:"true" desc:"Postgres password"`
	DbName                string   `isReq:"true" desc:"Name of Postgres database"`
	DirMigrations         string   `def:"migrations" desc:"Path to migrations directory"`
	RedisHost             string   `isReq:"true" desc:"Redis host"`
	RedisDb               int      `isReq:"true" desc:"Redis DB number"`
	RedisUser             string   `isReq:"true" desc:"Redis username"`
	RedisPassword         string   `isReq:"true" desc:"Redis password"`
	AccessTokenLifespan   int      `isReq:"true" desc:"Lifespan of Access-token (in hours)"`
	RefreshTokenLifespan  int      `isReq:"true" desc:"Lifespan of Refresh-token (in hours)"`
	TokenSecret           string   `isReq:"true" desc:"Secret salt of JWT-generation"`
	RootAdminId           int      `isReq:"true" desc:"Id of root admin in Postgres"`
	KafkaBootstrapServers []string `def:"localhost:9092" desc:"Kafka bootstrap servers"`
	KafkaTopic            string   `def:"eventryx_raw_messages" desc:"Kafka topic"`
}

var Config EnvConfiguration

func InitConfig() error {
	_ = godotenv.Load(".env")

	cfgValue := reflect.ValueOf(&Config).Elem()
	cfgType := cfgValue.Type()

	var errs []string

	for i := 0; i < cfgType.NumField(); i++ {
		fieldType := cfgType.Field(i)
		fieldValue := cfgValue.Field(i)

		meta := parseTags(fieldType)

		envVal := resolveEnvValue(meta)

		if meta.isRequired && envVal == "" {
			errs = append(errs,
				fmt.Sprintf("environment variable %s (%s) required but not set",
					meta.envName, meta.description))
			continue
		}

		if err := assignValue(fieldValue, envVal, meta); err != nil {
			errs = append(errs, fmt.Sprintf("%s(%s): %v", meta.envName, meta.description, err))
		}
	}

	if len(errs) > 0 {
		log.Printf("config initialization failed:\n%s", strings.Join(errs, "\n"))
		return fmt.Errorf("config initialization failed")
	}

	return nil
}

type tagMeta struct {
	envName     string
	isRequired  bool
	defaultVal  string
	description string
	enums       []string
	separator   string
}

func parseTags(f reflect.StructField) tagMeta {
	envName := f.Tag.Get("env")
	if envName == "" {
		envName = toMacroCase(f.Name)
	}

	separator := f.Tag.Get("sep")
	if separator == "" {
		separator = ","
	}

	var enums []string
	if rawEnums := f.Tag.Get("enums"); rawEnums != "" {
		enums = strings.Split(rawEnums, ",")
	}

	return tagMeta{
		envName:     envName,
		isRequired:  f.Tag.Get("isReq") == "true",
		defaultVal:  f.Tag.Get("def"),
		description: f.Tag.Get("desc"),
		enums:       enums,
		separator:   separator,
	}
}

func resolveEnvValue(meta tagMeta) string {
	val := os.Getenv(meta.envName)
	if val == "" && meta.defaultVal != "" {
		return meta.defaultVal
	}
	return val
}

func assignValue(field reflect.Value, raw string, meta tagMeta) error {
	if raw == "" {
		return nil
	}

	switch field.Kind() {

	case reflect.String:
		if len(meta.enums) > 0 && !slices.Contains(meta.enums, raw) {
			return fmt.Errorf("value must be one of [%s]",
				strings.Join(meta.enums, ", "))
		}
		field.SetString(raw)

	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid value for boolean: %w", err)
		}
		field.SetBool(v)

	case reflect.Int:
		v, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid value for integer: %w", err)
		}
		field.SetInt(int64(v))

	case reflect.Slice:
		values := strings.Split(raw, meta.separator)
		field.Set(reflect.ValueOf(values))

	default:
		return fmt.Errorf("unsupported value type: %s", field.Kind())
	}

	return nil
}

func toMacroCase(input string) string {
	var result []rune

	for i, r := range input {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToUpper(r))
	}

	return string(result)
}

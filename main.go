package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/m-kuzmin/simple-rest-api/api"
	"github.com/m-kuzmin/simple-rest-api/api/swaggerui"
	"github.com/m-kuzmin/simple-rest-api/db"
	"github.com/m-kuzmin/simple-rest-api/logging"
	"github.com/spf13/viper"
)

const (
	dbDriver         = "postgres"
	migrationsSource = "file:///migrations"
)

type config struct {
	Database databaseConfig `mapstructure:"database"`
	RestAPI  restAPIConfig  `mapstructure:"rest_api"`
}

type databaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}
type restAPIConfig struct {
	Swagger            string       `mapstructure:"swagger"`
	Address            string       `mapstructure:"address"`
	DatabaseConnection dbConnConfig `mapstructure:"database_connection"`
	ReadTimeoutSecs    uint         `mapstructure:"readTimeoutSecs"`
}
type dbConnConfig struct {
	Retries  uint `mapstructure:"retries"`
	Interval uint `mapstructure:"interval"`
}

func main() {
	logging.GlobalLogger = logging.StdLogger{}

	config, err := loadConfig()
	if err != nil {
		logging.Warnf("Viper config error: %s", err)
	}

	logging.Debugf("App config: %v", config)

	postgres, router := ServerFromConfig(config)

	httpServer := StartServer(router, config.RestAPI.Address, time.Second*time.Duration(config.RestAPI.ReadTimeoutSecs))

	logging.Infof("Server started")

	WaitForCtrcC()
	logging.Infof("Shutting down the server")

	err = httpServer.Shutdown(context.Background())
	if err != nil {
		logging.Errorf("[1/2] Error during shutdown: %s", err)
	} else {
		logging.Infof("[1/2] HTTP handler stopped")
	}

	if err = postgres.Close(); err != nil {
		logging.Errorf("[2/2] Error closing PostgreSQL connection: %s", err)
	} else {
		logging.Infof("[2/2] SQL connection closed")
	}

	logging.Infof("Server stopped.")
}

func loadConfig() (*config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/simple-rest-api.d")
	viper.AddConfigPath("$HOME/.config/simple-rest-api.d")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	//nolint:gomnd // Default values for config file
	conf := &config{
		Database: databaseConfig{
			User:     "root",
			Password: "",
			Database: "users",
			Host:     "postgres",
			Port:     "5432",
		},
		RestAPI: restAPIConfig{
			Swagger: "swaggerui",
			Address: ":8000",
			DatabaseConnection: dbConnConfig{
				Retries:  5,
				Interval: 7,
			},
			ReadTimeoutSecs: 60,
		},
	}

	var err error
	if err = viper.Unmarshal(conf); err != nil {
		err = fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	switch conf.RestAPI.Swagger {
	case "true", "1": // In yaml "true" => "true" and true => "1". Note the quotes around the value.
		conf.RestAPI.Swagger = "swaggerui"
	case "false", "0":
		conf.RestAPI.Swagger = ""
	}

	return conf, err
}

func ServerFromConfig(conf *config) (*db.Postgres, *gin.Engine) {
	logging.Infof("Connecting to Postgres")

	postgres := MustSetupPostgres(conf.Database.User, conf.Database.Password,
		net.JoinHostPort(conf.Database.Host, conf.Database.Port), conf.Database.Database,
		conf.RestAPI.DatabaseConnection.Retries, conf.RestAPI.DatabaseConnection.Interval)

	logging.Infof("Connected to Postgres")
	gin.SetMode(gin.ReleaseMode)

	server := api.NewServer(postgres)
	router := api.NewGinRouter(server)

	if conf.RestAPI.Swagger != "" {
		swaggerui.Mount(router, conf.RestAPI.Swagger)
	}

	return postgres, router
}

func StartServer(engine http.Handler, address string, readTimeout time.Duration) *http.Server {
	server := &http.Server{
		Addr:        address,
		Handler:     engine,
		ReadTimeout: readTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logging.Errorf("HTTP server error: %s", err)
		}

		logging.Infof("HTTP server shutdown")
	}()

	return server
}

func WaitForCtrcC() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
}

func MustSetupPostgres(user, password, address, database string, retries, intervalSecs uint) *db.Postgres {
	dbAddress := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		user, password, address, database)

	conn, err := db.ConnectToDBWithRetry(dbDriver, dbAddress,
		retries, time.Duration(intervalSecs)*time.Second)
	if err != nil {
		logging.Fatalf("failed to connect to PostgreSQL: %s", err)
	}

	if err = db.PostgresMigrateUp(conn, migrationsSource, database); err != nil {
		logging.Fatalf("failed to migrate PostgreSQL to latest version: %s", err)
	}

	postgres := db.NewPostgres(conn)

	return postgres
}

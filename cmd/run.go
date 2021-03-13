package main

//go:generate swag init -g ./web/router.go -o ./docs

import (
	"fmt"
	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/chiahsoon/go_scaffold/internal/models/users"
	"log"
	"os"
	"strings"

	"github.com/chiahsoon/go_scaffold/web"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Logging
	err := initLogger()
	if err != nil {
		log.Fatal("failed to initialise logger", err)
	}

	// Configuration
	env := "dev"
	if envVar := os.Getenv("ENV"); envVar != "" {
		env = strings.ToLower(envVar)
	}

	err = initConfigs(env)
	if err != nil {
		log.Fatal("failed to initialise configurations file: \n", err)
	}

	// Database
	err = initDB()
	if err != nil {
		log.Fatal("failed to initialise database: \n", err)
	}

	// fmt.Println("FLAG: ", os.Args[1:])
	web.Run()
}

func initLogger() error {
	logDir := "log"
	infoPath := fmt.Sprintf("%s/info.log", logDir)
	errorPath := fmt.Sprintf("%s/info.log", logDir)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err = os.Mkdir(logDir, os.ModePerm); err != nil {
			log.Fatal("unable to create /log directory")
		}
	}

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	prodEncoder := zap.NewProductionEncoderConfig()
	prodEncoder.EncodeTime = zapcore.ISO8601TimeEncoder

	lowWriteSyncer, lowClose, err := zap.Open(infoPath)
	if err != nil {
		lowClose()
		return err
	}

	highWriteSyncer, highClose, err := zap.Open(errorPath)
	if err != nil {
		highClose()
		return err
	}

	highCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), highWriteSyncer, highPriority)
	lowCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder), lowWriteSyncer, lowPriority)

	logger := zap.New(zapcore.NewTee(highCore, lowCore), zap.AddCaller())
	zap.ReplaceGlobals(logger)

	// Here are some examples of logging
	// logger.Debug("i am debug",zap.String("key","debug"))
	// logger.Info("i am info",zap.String("key","info"))
	// logger.Error("i am error",zap.String("key","error"))

	return nil
}

func initConfigs(env string) error {
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs/")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// viper.SetEnvPrefix("") // Require env variables to be prepended with $ENV_

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return err
}

func initDB() error {
	var dbUser, dbPassword, dbProtocol, dbAddress, dbName string

	databases := viper.Get("databases")
	for databaseName, databaseConfig := range databases.(map[string]interface{}) {
		dbConfig := databaseConfig.(map[string]interface{})
		dbName = databaseName
		dbProtocol = dbConfig["protocol"].(string)
		dbUser = dbConfig["user"].(string)
		dbPassword = dbConfig["password"].(string)
		dbAddress = dbConfig["address"].(string)
	}

	zap.L().Info("database_config",
		zap.String("name", dbName),
		zap.String("user", dbUser),
		zap.String("password", dbPassword),
		zap.String("address", dbAddress),
	)

	dsn := fmt.Sprintf("%v:%v@%v(%v)/%v?charset=utf8mb4&parseTime=True",
		dbUser, dbPassword, dbProtocol, dbAddress, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	models.DB = db
	err = models.DB.AutoMigrate(&users.User{})

	return err
}

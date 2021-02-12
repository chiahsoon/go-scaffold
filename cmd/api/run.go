package main

//go:generate swag init -g ./web/router.go -o ./docs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chiahsoon/go_scaffold/internal/model"
	"github.com/chiahsoon/go_scaffold/web"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitLogger(infoPath, errorPath string) error {
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

func InitConfigs(env string) error {
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

func InitDB() error {
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

	model.DB = db
	err = model.DB.AutoMigrate(&model.User{})

	return err
}

func main() {
	// Logging
	infoLogFilepath := "logs/info.log"
	errorLogFilepath := "logs/error.log"
	err := InitLogger(infoLogFilepath, errorLogFilepath)
	if err != nil {
		log.Fatal("failed to initialise logger", err)
	}

	// Configuration
	env := "dev"
	if envVar := os.Getenv("ENV"); envVar != "" {
		env = strings.ToLower(envVar)
	}

	err = InitConfigs(env)
	if err != nil {
		log.Fatal("failed to initialise configurations file: \n", err)
	}

	// Database
	err = InitDB()
	if err != nil {
		log.Fatal("failed to initialise database: \n", err)
	}

	// fmt.Println("FLAG: ", os.Args[1:])
	web.Run()
}

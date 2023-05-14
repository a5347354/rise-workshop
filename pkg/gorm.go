package pkg

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/viper"
	gormotel "github.com/wei840222/gorm-otel"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/prometheus"
)

// NewPostgresClient return new *gorm.DB with mysql configs with fx.Lifecycle
func NewPostgresClient(lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := NewPostgresClientWithoutLC()
	if err != nil {
		return nil, err
	}

	if err := db.Use(gormotel.New(gormotel.WithLogResult(true), gormotel.WithSqlParameters(true))); err != nil {
		return nil, err
	}

	if err := db.Use(prometheus.New(prometheus.Config{})); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if viper.GetString("postgre.readonly.host") != "" {
		if err := db.Use(
			dbresolver.Register(dbresolver.Config{
				Replicas: []gorm.Dialector{mysql.Open(fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&loc=%s&time_zone=%s",
					viper.GetString("postgres.readonly.username"),
					viper.GetString("postgres.readonly.password"),
					viper.GetString("postgres.readonly.host"),
					viper.GetString("postgres.database"),
					url.QueryEscape("Asia/Taipei"),
					url.QueryEscape("'+8:00'"),
				))},
				Policy: dbresolver.RandomPolicy{},
			}).
				SetConnMaxIdleTime(30 * time.Second).
				SetConnMaxLifetime(60 * time.Second).
				SetMaxIdleConns(3).
				SetMaxOpenConns(30),
		); err != nil {
			return nil, err
		}
	} else {
		sqlDB.SetMaxIdleConns(3)
		sqlDB.SetMaxOpenConns(30)
		sqlDB.SetConnMaxIdleTime(30 * time.Second)
		sqlDB.SetConnMaxLifetime(60 * time.Second)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return sqlDB.Close()
		},
	})
	return db, nil
}

// NewPostgresClientWithoutLC return new *gorm.DB with mysql configs
func NewPostgresClientWithoutLC() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		viper.GetString("postgres.host"),
		viper.GetString("postgres.username"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.password"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 100 * time.Millisecond,
			LogLevel: func() logger.LogLevel {
				if viper.GetBool("debug") {
					return logger.Info
				}
				return logger.Warn
			}(),
			Colorful: viper.GetBool("debug") || isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
		}),
	})
}

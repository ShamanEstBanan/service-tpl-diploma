package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress          string `env:"RUN_ADDRESS"`
	PostgresDSN         string `env:"DATABASE_URI"`
	AccrualSystemAddres string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() *Config {
	cfg := Config{}

	// заполнение конфига из значений аргументов командной строки
	flag.StringVar(
		&cfg.RunAddress,
		"a",
		"localhost:8080",
		"адрес и порт запуска сервиса",
	)
	flag.StringVar(
		&cfg.PostgresDSN,
		"d",
		"postgresql://localhost:5432/service-tpl",
		"адрес подключения к базе данных, дефолтного значения нет",
	)
	flag.StringVar(
		&cfg.AccrualSystemAddres,
		"r",
		"",
		"адрес системы расчёта начислений",
	)
	flag.Parse()

	// переопределяем значения конфига переменными ENV, eсли они определены в ОС
	if err := env.Parse(&cfg); err != nil {
		log.Printf("err while parsing env-values: %v\n", err)
	}
	return &cfg
}

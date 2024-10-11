package config

import (
	"flag"
	"time"
)

var (
	FlagDir           string
	FlagAll           bool
	FlagCheckout      bool
	FlagFetch         bool
	FlagPull          bool
	FlagResetHard     bool
	FlagDefaultBranch string
	FlagVerbose       bool

	FlagExecuteTimeout *time.Duration

	FlagVersion bool
)

func init() {
	flag.StringVar(&FlagDir, "dir", "../", "Путь к каталогу с проектами")
	flag.StringVar(&FlagDefaultBranch, "branch", "master", "Наименование ветки: master/main")

	flag.BoolVar(&FlagAll, "all", false, "Выполнить все шаги checkout master, fetch, pull")
	flag.BoolVar(&FlagCheckout, "checkout", false, "Выполнить checkout master")
	flag.BoolVar(&FlagFetch, "fetch", false, "Выполнить fetch")
	flag.BoolVar(&FlagPull, "pull", false, "Выполнить pull")
	flag.BoolVar(&FlagResetHard, "reset-hard", false, "При необходимости перед git checkout выполнить git reset --hard")
	flag.BoolVar(&FlagVerbose, "verbose", false, "Выводить результат выполнения команд")

	FlagExecuteTimeout = flag.Duration("execute_timeout", time.Second*30, "Максимальное время обработки одного каталога")

	flag.BoolVar(&FlagVersion, "version", false, "Показать версию приложения")

	flag.Parse()
}

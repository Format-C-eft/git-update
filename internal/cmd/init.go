package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git_update",
	Short: "Обновление всех репозиториев до актуального состояния",
	RunE:  cmdRun,
}

func init() {
	rootCmd.Version = "v0.0.1"

	rootCmd.Flags().StringVarP(&flagDir, "dir", "d", "../", "Путь к каталогу с проектами")
	rootCmd.Flags().StringVarP(&flagDefaultBranch, "branch", "b", "master", "Наименование ветки: master/main")
	rootCmd.Flags().BoolVarP(&flagAll, "all", "a", false, "Выполнить все шаги checkout master, fetch, pull")
	rootCmd.Flags().BoolVarP(&flagCheckout, "checkout", "c", false, "Выполнить checkout master")
	rootCmd.Flags().BoolVarP(&flagFetch, "fetch", "f", false, "Выполнить fetch")
	rootCmd.Flags().BoolVarP(&flagPull, "pull", "p", false, "Выполнить pull")
	rootCmd.Flags().BoolVarP(&flagResetHard, "reset-hard", "r", false, "При необходимости перед git checkout выполнить git reset --hard")

	flagExecuteTimeout = rootCmd.Flags().DurationP("execute_timeout", "e", time.Second*30, "Максимальное время обработки одного каталога")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

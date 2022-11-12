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
	rootCmd.Version = "v0.0.1-alpha1"

	rootCmd.Flags().StringVarP(&flagDir, "dir", "d", "../", "Путь к каталогу с проектами")
	rootCmd.Flags().StringVarP(&flagDefaultBranch, "branch", "b", "master", "Наименование ветки: master/main")
	rootCmd.Flags().BoolVarP(&flagAll, "all", "a", false, "Выполнить все шаги checkout master, fetch, pull")
	rootCmd.Flags().BoolVarP(&flagAll, "checkout", "c", false, "Выполнить checkout master")
	rootCmd.Flags().BoolVarP(&flagAll, "fetch", "f", false, "Выполнить fetch")
	rootCmd.Flags().BoolVarP(&flagAll, "pull", "p", false, "Выполнить pull")

	flagExecuteTimeout = rootCmd.Flags().DurationP("execute_timeout", "e", time.Second*30, "Максимальное время обработки одного каталога")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/*
Copyright © 2021 Christian Lerrahn <github@penpal4u.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package sub

import (
	"fmt"
	"github.com/jsfan/t3migrate/internal/config"
	"github.com/jsfan/t3migrate/internal/storage"
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

var cfgFile string
var dryRun *bool
var srcDbConfig, dstDbConfig *config.MySQLConfig
var srcStore, dstStore *storage.Store

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "t3migrate",
	Short: "Selectively migrates TYPO3 data from one database to another",
	Long: `This utility is used to copy important data from an existing TYPO3
database to a new (empty) TYPO3 database. It can be used for rebuildign a TYPO3
from scratch on a different version.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./t3migrate.yaml)")
	dryRun = rootCmd.Flags().BoolP("dryRun", "n", false, "DRY RUN")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("t3migrate")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		srcDbPort := 3306
		if viper.IsSet("databases.source.port") {
			srcDbPort = viper.GetInt("databases.source.port")
		}
		srcDbConfig = &config.MySQLConfig{
			User:     viper.GetString("databases.source.user"),
			Password: viper.GetString("databases.source.password"),
			Host:     viper.GetString("databases.source.host"),
			Port:     srcDbPort,
			Database: viper.GetString("databases.source.database"),
		}
		dstDbPort := 3306
		if viper.IsSet("databases.destination.port") {
			dstDbPort = viper.GetInt("databases.destination.port")
		}
		dstDbConfig = &config.MySQLConfig{
			User:     viper.GetString("databases.destination.user"),
			Password: viper.GetString("databases.destination.password"),
			Host:     viper.GetString("databases.destination.host"),
			Port:     dstDbPort,
			Database: viper.GetString("databases.destination.database"),
		}
		srcStore = &storage.Store{}
		if err := srcStore.Connect(srcDbConfig, *dryRun); err != nil {
			cobra.CheckErr(fmt.Errorf("could not connect to source database: %w", err))
		}
		dstStore = &storage.Store{}
		if err := dstStore.Connect(dstDbConfig, *dryRun); err != nil {
			cobra.CheckErr(fmt.Errorf("could not connect to destination database: %w", err))
		}
	}
}

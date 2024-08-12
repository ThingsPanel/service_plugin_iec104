package main

import (
	"github.com/spf13/cobra"
	iec104_slave "iec104-slave"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file (default is ./config.yaml)")
}

var rootCmd = &cobra.Command{
	Use:   "iec104-master",
	Short: "iec104-master service",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		if err := iec104_slave.Init(cfgFile); err != nil {
			log.Fatal(err)
		}

		go func() {
			err := iec104_slave.Run()
			if err != nil {
				log.Fatal(err)
			}
		}()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		<-sigCh

		iec104_slave.Close()
	},
}

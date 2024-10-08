/*
Copyright © 2024 WajahatAliAbid
*/
package cmd

import (
	"fmt"
	"os"

	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/cobra"
)

type AwsConfig struct {
}

func PreExecute(cmd *cobra.Command, _ []string) {
	profile := cmd.Flag("profile")
	region := cmd.Flag("region")

	if profile.Value.String() == "" {
		profile.Value.Set(
			os.Getenv("AWS_PROFILE"),
		)
	}
	if region.Value.String() == "" {
		region.Value.Set(
			os.Getenv("AWS_REGION"),
		)
	}
	var profileName *string
	var regionName *string

	if profile.Value.String() != "" {
		profileName = new(string)
		*profileName = profile.Value.String()
	}
	if region.Value.String() != "" {
		regionName = new(string)
		*regionName = region.Value.String()
	}

	if profileName == nil {
		cfg, err := config.LoadDefaultConfig(
			cmd.Context(),
			config.WithRegion(*regionName),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmd.SetContext(
			context.WithValue(
				cmd.Context(),
				AwsConfig{},
				cfg,
			),
		)
		return
	} else {
		cfg, err := config.LoadDefaultConfig(
			cmd.Context(),
			config.WithRegion(*regionName),
			config.WithSharedConfigProfile(*profileName),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmd.SetContext(
			context.WithValue(
				cmd.Context(),
				AwsConfig{},
				cfg,
			),
		)
	}

}

var rootCmd = &cobra.Command{
	Use:       "zqs [sub]",
	Short:     "Send messages to the queue",
	Long:      `Send messages to the queue, from either a single file or from a directory containing json files`,
	PreRun:    PreExecute,
	Run:       RunSendMessages,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"queue_url"},
}

func Execute() {
	rootCmd.PersistentFlags().StringP(
		"profile",
		"p",
		"",
		"Profile to use from ~/.aws/credentials",
	)

	rootCmd.PersistentFlags().StringP(
		"region",
		"r",
		"",
		"Region to use",
	)

	rootCmd.Flags().StringP(
		"file",
		"f",
		"",
		"File to send",
	)

	rootCmd.Flags().StringP(
		"directory",
		"d",
		"",
		"Directory to send",
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

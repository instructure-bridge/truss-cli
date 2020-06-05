package cmd

import (
	"os"

	"github.com/instructure-bridge/truss-cli/truss"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getKubeconfigCmd represents the getKubeconfig command
var getKubeconfigCmd = &cobra.Command{
	Use:   "get-kubeconfig",
	Short: "Get Kubeconfigs from source",
	Long: `
kubeconfigfiles:
  s3:
    bucket: my-aws-bucket
    region: us-east-1
`,
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.GetStringMap("kubeconfigfiles")

		dest, err := getKubeDir()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}

		s3 := config["s3"]
		if s3 != nil {
			s3Config, ok := s3.(map[string]interface{})
			if !ok {
				log.Errorln("invalid s3 config")
				os.Exit(1)
			}

			awsrole := s3Config["awsrole"].(string)
			bucket, ok := s3Config["bucket"].(string)
			if !ok {
				log.Errorln("s3 config must have bucket name")
				os.Exit(1)
			}
			region, ok := s3Config["region"].(string)
			if !ok {
				log.Errorln("s3 config must have region")
				os.Exit(1)
			}
			if err := truss.GetKubeconfigS3(awsrole, bucket, dest, region).Fetch(); err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getKubeconfigCmd)
}

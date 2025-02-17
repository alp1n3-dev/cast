package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

var (
    body    string
    headers []string
)

func init() {
  //rootCmd.AddCommand(sendHTTPCmd)
  sendHTTPCmd.Flags().StringVarP(&body, "body", "b", "", "Request body")
  sendHTTPCmd.Flags().StringArrayVarP(&headers, "header", "h", []string{}, "Request headers")
}

var sendHTTPCmd = &cobra.Command{
  Use:   "get || post || put || delete || patch || trace || head || options || connect",
  Short: "send an HTTP request.",
  Long:  `Send a fully customizable HTTP request and retrieve the response.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("reached sendhttpcmd")
  },
}

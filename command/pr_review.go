package command

import (
	"github.com/spf13/cobra"
)

func init() {
	prCmd.AddCommand(prReviewCmd)

	prCmd.Flags().StringP("approve", "a", "", "Approve pull request")
	prCmd.Flags().StringP("request-changes", "r", "", "Request changes on a pull request")
	prCmd.Flags().StringP("comment", "c", "", "Comment on a pull request")
}

var prReviewCmd = &cobra.Command{
	Use:   "TODO",
	Short: "TODO",
	Long:  "TODO",
	RunE:  prReview,
}

func prReview(cmd *cobra.Command, args []string) error {

	// TODO process flags, make some decisions

	return nil
}

package command

import (
	"errors"

	"github.com/spf13/cobra"
)

// This probably never-passed-by-a-user string is used to signal the absence of a flag being passed
// at all (in order to distinguish from an empty string being passed).
const noOptSigil = "!HIDEOUS SIGIL!"

func init() {
	prCmd.AddCommand(prReviewCmd)

	prCmd.Flags().StringP("approve", "a", "", "Approve pull request")
	prCmd.Flags().StringP("request-changes", "r", "", "Request changes on a pull request")
	prCmd.Flags().StringP("comment", "c", "", "Comment on a pull request")

	prCmd.Flags().Lookup("approve").NoOptDefVal = noOptSigil
	prCmd.Flags().Lookup("request-changes").NoOptDefVal = noOptSigil
	prCmd.Flags().Lookup("comment").NoOptDefVal = noOptSigil
}

var prReviewCmd = &cobra.Command{
	Use:   "TODO",
	Short: "TODO",
	Long:  "TODO",
	RunE:  prReview,
}

func stringsEqual(target string, ss ...string) bool {
	for _, s := range ss {
		if s != target {
			return false
		}
	}
	return true
}

func oneEqual(target string, ss ...string) bool {
	count := 0
	for _, s := range ss {
		if s == target {
			count += 1
		}
	}

	return count == 1
}

func prReview(cmd *cobra.Command, args []string) error {
	approveVal, err := cmd.Flags().GetString("approve")
	if err != nil {
		return err
	}
	requestChangesVal, err := cmd.Flags().GetString("request-changes")
	if err != nil {
		return err
	}
	commentVal, err := cmd.Flags().GetString("comment")
	if err != nil {
		return err
	}

	validationErr := errors.New("need exactly on of approve, request-changes, or comment")
	if stringsEqual(noOptSigil, approveVal, requestChangesVal, commentVal) {
		return validationErr
	}

	if !oneEqual(noOptSigil, approveVal, requestChangesVal, commentVal) {
		return validationErr
	}

	// TODO process flags, make some decisions

	return nil
}

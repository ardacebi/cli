package command

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// This probably never-passed-by-a-user string is used to signal the absence of a flag being passed
// at all (in order to distinguish from an empty string being passed).
const noOptSigil = "!HIDEOUS SIGIL!"

func init() {
	prCmd.AddCommand(prReviewCmd)

	prReviewCmd.Flags().StringP("approve", "a", "", "Approve pull request")
	prReviewCmd.Flags().StringP("request-changes", "r", "", "Request changes on a pull request")
	prReviewCmd.Flags().StringP("comment", "c", "", "Comment on a pull request")

	prReviewCmd.Flags().Lookup("approve").NoOptDefVal = noOptSigil
	prReviewCmd.Flags().Lookup("request-changes").NoOptDefVal = noOptSigil
	prReviewCmd.Flags().Lookup("comment").NoOptDefVal = noOptSigil
}

var prReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "TODO",
	Args:  cobra.MaximumNArgs(1),
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
	ctx := contextForCommand(cmd)
	baseRepo, err := determineBaseRepo(cmd, ctx)
	if err != nil {
		return fmt.Errorf("could not determine base repo: %w", err)
	}

	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		return err
	}

	prArg := ""

	// determine PR
	if len(args) == 0 {
		prNum, _, err := prSelectorForCurrentBranch(ctx, baseRepo)
		if err != nil && err.Error() != "git: not on any branch" {
			return fmt.Errorf("could not query for pull request for current branch: %w", err)
		}
		prArg = fmt.Sprintf("%d", prNum)
	} else {
		if prNum, repo := prFromURL(args[0]); repo != nil {
			prArg = prNum
			baseRepo = repo
		}
	}

	pr, err := prFromArg(apiClient, baseRepo, prArg)
	if err != nil {
		return fmt.Errorf("could not find pull request %d: %w", pr.Number, err)
	}

	// process PR action
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

	validationErr := errors.New("need exactly one of approve, request-changes, or comment")
	if stringsEqual(noOptSigil, approveVal, requestChangesVal, commentVal) {
		return validationErr
	}

	if !oneEqual(noOptSigil, approveVal, requestChangesVal, commentVal) {
		return validationErr
	}

	// TODO process flags, make some decisions

	return nil
}

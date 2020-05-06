package command

import (
	"errors"
	"fmt"

	"github.com/cli/cli/api"
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

func processReviewOpt(cmd *cobra.Command) (*api.PullRequestReviewInput, error) {
	approveVal, err := cmd.Flags().GetString("approve")
	if err != nil {
		return nil, err
	}
	changesVal, err := cmd.Flags().GetString("request-changes")
	if err != nil {
		return nil, err
	}
	commentVal, err := cmd.Flags().GetString("comment")
	if err != nil {
		return nil, err
	}

	found := 0
	var body string
	var state api.PullRequestReviewState

	if approveVal != noOptSigil {
		found++
		body = approveVal
		state = api.ReviewApprove
	} else if changesVal != noOptSigil {
		found++
		body = changesVal
		state = api.ReviewRequestChanges
	} else if commentVal != noOptSigil {
		found++
		body = commentVal
		state = api.ReviewComment
	}

	if found != 1 {
		return nil, errors.New("need exactly one of approve, request-changes, or comment")
	}

	return &api.PullRequestReviewInput{
		Body:  body,
		State: state,
	}, nil
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

	input, err := processReviewOpt(cmd)
	if err != nil {
		return fmt.Errorf("did not understand desired review action: %w", err)
	}

	err = api.AddReview(apiClient, pr, input)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

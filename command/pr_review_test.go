package command

import (
	"testing"
)

func TestPRReview_validation(t *testing.T) {
	for _, cmd := range []string{
		`pr review`,
		`pr review --approve --comment`,
		`pr review --approve="cool" --comment="rad"`,
	} {
		_, err := RunCommand(prReviewCmd, cmd)
		eq(t, err.Error(), "need exactly one of approve, request-changes, or comment")
	}
}

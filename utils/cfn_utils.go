package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func statusIsValid(status string) bool {
	switch status {
	case "CREATE_IN_PROGRESS",
		"CREATE_FAILED",
		"CREATE_COMPLETE",
		"ROLLBACK_IN_PROGRESS",
		"ROLLBACK_FAILED",
		"ROLLBACK_COMPLETE",
		"DELETE_IN_PROGRESS",
		"DELETE_FAILED",
		"DELETE_COMPLETE",
		"UPDATE_IN_PROGRESS",
		"UPDATE_COMPLETE_CLEANUP_IN_PROGRESS",
		"UPDATE_COMPLETE",
		"UPDATE_ROLLBACK_IN_PROGRESS",
		"UPDATE_ROLLBACK_FAILED",
		"UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS",
		"UPDATE_ROLLBACK_COMPLETE",
		"REVIEW_IN_PROGRESS",
		"IMPORT_IN_PROGRESS",
		"IMPORT_COMPLETE",
		"IMPORT_ROLLBACK_IN_PROGRESS",
		"IMPORT_ROLLBACK_FAILED",
		"IMPORT_ROLLBACK_COMPLETE":
		return true
	}

	return false
}

func GetStackSummaries(sess *session.Session, status string) (*cloudformation.ListStacksOutput, error) {
	svc := cloudformation.New(sess)
	var filter []*string

	if status != "all" {
		filter = append(filter, aws.String(status))
	}

	input := &cloudformation.ListStacksInput{StackStatusFilter: filter}

	resp, err := svc.ListStacks(input)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

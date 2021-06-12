package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func GetStackOutput(cf *cloudformation.CloudFormation, stackName string, outputKey string) (string, error) {
	describeStacksInput := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	describeStacksOutput, err := cf.DescribeStacks(describeStacksInput)

	if err != nil {
		return "", err
	}

	stacks := describeStacksOutput.Stacks

	for _, stack := range stacks {
		outputs := stack.Outputs
		for _, output := range outputs {
			if strings.Compare(*output.OutputKey, outputKey) == 0 {
				return *output.OutputValue, nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("outputKey: %s not found in stack %s", outputKey, stackName))
}

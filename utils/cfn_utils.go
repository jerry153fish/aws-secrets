package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/patrickmn/go-cache"
)

func GetStackOutput(cf *cloudformation.CloudFormation, stackName string, outputKey string, c *cache.Cache) (string, error) {
	// TODO: need to refactor to cache all outputs and query only for stack once very 5 min
	cacheKey := stackName + "-" + outputKey

	if x, found := c.Get(cacheKey); found {
		return x.(string), nil
	}

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
				c.Set(cacheKey, *output.OutputValue, cache.DefaultExpiration)
				return *output.OutputValue, nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("outputKey: %s not found in stack %s", outputKey, stackName))
}

func GetCfnClient() *cloudformation.CloudFormation {

	if _, ok := os.LookupEnv("TEST_ENV"); ok {
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-west-2"),
		}))

		return cloudformation.New(sess, aws.NewConfig().WithEndpoint("http://localhost:4566"))
	}

	return cloudformation.New(session.Must(session.NewSession()))
}

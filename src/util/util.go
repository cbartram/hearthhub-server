package util

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"strings"
)

// GetUserAttribute Retrieves and parses a user attribute from Cognito into an array of strings. Most
// attributes are CSV strings. Examples include: purchased plugins, plugin expiration dates, plugin purchase dates etc...
func GetUserAttribute(attributes []types.AttributeType, attributeName string) []string {
	for _, attribute := range attributes {
		if aws.ToString(attribute.Name) == attributeName {
			return strings.Split(aws.ToString(attribute.Value), ",")
		}
	}

	return make([]string, 0)
}

func GetUserAttributeString(attributes []types.AttributeType, attributeName string) string {
	for _, attribute := range attributes {
		if aws.ToString(attribute.Name) == attributeName {
			return aws.ToString(attribute.Value)
		}
	}

	return ""
}

func MakeAttribute(key, value string) types.AttributeType {
	attr := types.AttributeType{
		Name:  &key,
		Value: &value,
	}
	return attr
}

func Map[T any, O any](things []T, mapper func(thing T) O) []O {
	result := make([]O, 0, len(things))
	for _, thing := range things {
		result = append(result, mapper(thing))
	}
	return result
}

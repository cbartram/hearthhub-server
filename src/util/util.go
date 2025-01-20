package util

import (
	"crypto/rand"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"regexp"
	"strings"
)

const (
	// Define the format: 5 chars + hyphen + 5 chars + hyphen + 5 chars + hyphen + 5 chars
	keyLength    = 23 // Total length including hyphens
	sectionSize  = 5  // Number of characters in each section
	numSections  = 4  // Number of sections
	validPattern = "^[A-Z]{5}-[A-Z]{5}-[A-Z]{5}-[A-Z]{5}$"
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

// GenerateLicenseKey generates a random license key in the format "XXXXX-XXXXX-XXXXX-XXXXX"
func GenerateLicenseKey() (string, error) {
	totalChars := sectionSize * numSections
	bytes := make([]byte, totalChars)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	// Convert random bytes to uppercase letters
	result := make([]string, numSections)
	for i := 0; i < totalChars; i++ {
		// Convert to character in range A-Z (0-25 + 65 for ASCII 'A')
		bytes[i] = byte(bytes[i]%26 + 65)
		result[i/sectionSize] += string(bytes[i])
	}

	return strings.Join(result, "-"), nil
}

func ValidateLicenseKey(key string) bool {
	if len(key) != keyLength {
		return false
	}

	// Check if the key matches the pattern using regex
	match, err := regexp.MatchString(validPattern, key)
	if err != nil {
		return false
	}

	return match
}

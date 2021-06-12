package utils

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecretType(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudformation Utils Suite")
}

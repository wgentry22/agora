package agora_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAgora(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agora Suite")
}

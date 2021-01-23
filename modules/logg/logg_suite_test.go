package logg_test

import (
	"github.com/wgentry22/agora/types/config"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	conf config.Logging
)

func TestLogg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logg Suite")
}

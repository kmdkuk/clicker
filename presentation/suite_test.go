package presentation

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPresentation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Presentation Suite")
}

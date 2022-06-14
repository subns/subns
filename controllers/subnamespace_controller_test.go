package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/subns/subns/api/v1alpha1"
)

var _ = Describe("Subnamespace CR", func() {
	Context("Defaults", func() {
		var underTest *v1alpha1.Subnamespace

		err := k8sClient.Create(context.Background(), underTest)
		Expect(err).NotTo(HaveOccurred())
	})
})

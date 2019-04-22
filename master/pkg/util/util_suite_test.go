package util_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/oshosanya/go-dm/pkg/util"
)

var _ = Describe("Util", func() {
	var (
		testURL           string
		incompleteTestURL string
	)

	BeforeEach(func() {
		testURL = "http://31.210.87.4/ringtones_new/fullmp3low/t/Timaya_feat_Phyno_feat_Olamide_Telli_Person.mp3"
		incompleteTestURL = "www.wapdam.com"
	})

	Describe("Get file name from url", func() {
		Context("With web url", func() {
			It("return a file name and file extensions", func() {
				Expect(util.GetFileNameFromURL(testURL)).To(Equal("Timaya_feat_Phyno_feat_Olamide_Telli_Person.mp3"))
			})
		})
	})

	Describe("Validate URL", func() {
		Context("With web url", func() {
			It("Return true if a url is a valid http or https url", func() {
				Expect(util.ValidateURL(testURL)).To(Equal(true))
			})

			It("Return false if url is invalid", func() {
				Expect(util.ValidateURL("http://terra")).To(Equal(false))
			})
		})
	})

	Describe("Build URL", func() {
		Context("With web url", func() {
			It("Appends http to a url without the protocol", func() {
				Expect(util.BuildURL(incompleteTestURL)).To(Equal("http://www.wapdam.com"))
			})
		})
	})
})

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Suite")
}

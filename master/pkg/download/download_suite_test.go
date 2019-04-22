package download_test

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
})

func TestDownload(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Download Suite")
}

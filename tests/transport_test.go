package tests

import (
	"fmt"
	"strconv"
	"time"

	"net/url"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"

	"kubevirt.io/containerized-data-importer/pkg/common"
	"kubevirt.io/containerized-data-importer/pkg/controller"
	"kubevirt.io/containerized-data-importer/tests/framework"
	"kubevirt.io/containerized-data-importer/tests/utils"
)

var _ = Describe("Transport Tests", func() {

	const (
		secretPrefix = "transport-e2e-sec"
		targetFile   = "tinyCore.iso"
		sizeCheckPod = "size-checker"
	)

	f := framework.NewFrameworkOrDie("transport", framework.Config{SkipNamespaceCreation: false})
	c := f.K8sClient

	hostURL, err := url.Parse(f.RestConfig.Host)
	handelError(err)
	ip := hostURL.Hostname()

	fileHostService := utils.GetServiceInNamespaceOrDie(c, utils.FileHostNs, utils.FileHostName)
	httpAuthPort, err := utils.GetServiceNodePortByName(fileHostService, utils.HttpAuthPortName)
	handelError(err)
	httpNoAuthPort, err := utils.GetServiceNodePortByName(fileHostService, utils.HttpNoAuthPortName)
	handelError(err)

	httpAuthEp := fmt.Sprintf("http://%s:%d", ip, httpAuthPort)
	httpNoAuthEp := fmt.Sprintf("http://%s:%d", ip, httpNoAuthPort)

	fc, err := utils.GetCatalog(httpNoAuthEp, "", "")
	handelError(err)
	targetSize, err := fc.Size(targetFile)
	handelError(err)

	var ns string
	BeforeEach(func() {
		ns = f.Namespace.Name
		By(fmt.Sprintf("Waiting for all \"%s/%s\" deployment replicas to be Ready", utils.FileHostNs, utils.FileHostName))
		utils.WaitForDeploymentReplicasReadyOrDie(c, utils.FileHostNs, utils.FileHostName)
	})

	// it() is the body of the test and is executed once per Entry() by DescribeTable()
	it := func(ip, file, accessKey, secretKey string) {
		var (
			err error // prevent shadowing
			sec *v1.Secret
		)

		pvcAnn := map[string]string{
			controller.AnnEndpoint: ip + "/" + file,
			controller.AnnSecret:   "",
		}

		if accessKey != "" || secretKey != "" {
			By(fmt.Sprintf("Creating secret for endpoint %s", ip))
			stringData := map[string]string{
				common.KEY_ACCESS: utils.AccessKeyValue,
				common.KEY_SECRET: utils.SecretKeyValue,
			}
			sec, err = utils.CreateSecretFromDefinition(c, utils.NewSecretDefinition(nil, stringData, nil, ns, secretPrefix))
			Expect(err).NotTo(HaveOccurred(), "Error creating test secret")
			pvcAnn[controller.AnnSecret] = sec.Name
		}

		By(fmt.Sprintf("Creating PVC with endpoint annotation %q", pvcAnn[controller.AnnEndpoint]))
		pvc, err := utils.CreatePVCFromDefinition(c, ns, utils.NewPVCDefinition("transport-e2e", "20M", pvcAnn, nil))
		Expect(err).NotTo(HaveOccurred(), "Error creating PVC")

		importer, err := utils.FindPodByPrefix(c, ns, common.IMPORTER_PODNAME, common.CDI_LABEL_SELECTOR)
		Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Unable to get importer pod %q", ns+"/"+common.IMPORTER_PODNAME))

		Expect(utils.WaitTimeoutForPodStatus(c, importer.Name, importer.Namespace, v1.PodSucceeded, utils.PodWaitForTime)).To(Succeed())

		By("Verifying PVC is not empty")
		Expect(framework.VerifyPVCIsEmpty(f, pvc)).To(BeFalse(), fmt.Sprintf("Found 0 imported files on PVC %q", pvc.Namespace+"/"+pvc.Name))

		By(fmt.Sprintf("Verifying imported file size matches %q size of %d bytes", targetFile, targetSize))
		pod, err := utils.CreateExecutorPodWithPVC(c, sizeCheckPod, ns, pvc)
		Expect(err).NotTo(HaveOccurred())
		Expect(utils.WaitTimeoutForPodReady(c, sizeCheckPod, ns, 20*time.Second)).To(Succeed())

		stdout := f.ExecShellInPod(pod.Name, ns, "wc -c < /pvc/disk.img")
		Expect(err).NotTo(HaveOccurred(), "Error getting size of imported file")

		importSize, err := strconv.Atoi(stdout)
		Expect(err).NotTo(HaveOccurred())
		Expect(importSize).To(Equal(targetSize), "Expect imported file size to match remote file size")
	}

	DescribeTable("Transport Test Table", it,
		Entry("should connect to http endpoint without credentials", httpNoAuthEp, targetFile, "", ""),
		Entry("should connect to http endpoint with credentials", httpAuthEp, targetFile, utils.AccessKeyValue, utils.SecretKeyValue))
})

// handelError is intended for use outside It(), BeforeEach() and AfterEach() blocks where Expect() cannot be called.
func handelError(e error) {
	if e != nil {
		glog.Errorf("%#v", e)
		Fail(fmt.Sprintf("Encountered error: %v", e), 0)
	}
}

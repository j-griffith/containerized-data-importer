package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"kubevirt.io/containerized-data-importer/tests/framework"
	"kubevirt.io/containerized-data-importer/tests/utils"

	"k8s.io/api/core/v1"

	"kubevirt.io/containerized-data-importer/pkg/common"
	"kubevirt.io/containerized-data-importer/pkg/controller"
	"github.com/golang/glog"
	"strings"
)

var _ = Describe("Transport Tests", func() {

	const (
		secretPrefix = "transport-e2e-sec"
		targetFile   = "tinyCore.iso"
	)

	f, err := framework.NewFramework("", framework.Config{SkipNamespaceCreation: false})
	handelError(errors.Wrap(err, "error creating test framework"))

	c := f.K8sClient
	handelError(errors.Wrap(err, "error creating k8s client"))

	// Use http and trim api port from url
	host := strings.Replace(f.RestConfig.Host, "https", "http", 1)
	if strings.Count(host, ":") > 1 {
		host = host[:strings.LastIndex(host, ":")]
	}

	fileHostService := utils.GetServiceInNamespaceOrDie(c, utils.FileHostNs, utils.FileHostName)

	httpAuthPort, err := utils.GetServiceNodePortByName(fileHostService, utils.HttpAuthPortName)
	handelError(err)
	httpNoAuthPort, err := utils.GetServiceNodePortByName(fileHostService, utils.HttpNoAuthPortName)
	handelError(err)

	httpAuthEp := fmt.Sprintf("%s:%d", host, httpAuthPort)
	httpNoAuthEp := fmt.Sprintf("%s:%d", host, httpNoAuthPort)

	var ns string
	BeforeEach(func() {
		ns = f.Namespace.Name
		By(fmt.Sprintf("Waiting for all \"%s/%s\" deployment replicas to be Ready", utils.FileHostNs, utils.FileHostName))
		utils.WaitForDeploymentReplicasReadyOrDie(c, utils.FileHostNs, utils.FileHostName)
	})

	// it() is the body of the test and is executed once per Entry() by DescribeTable()
	// closes over c and ns
	it := func(ip, file, accessKey, secretKey string) {
		var (
			err error // prevent shadowing
			sec *v1.Secret
		)

		fc, err := utils.GetFileHostCatalog(ip, accessKey, secretKey)
		Expect(err).NotTo(HaveOccurred())
		glog.Infof("Catalog: %#v", fc)


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

		By(fmt.Sprintf("Creating PVC with endpoint annotation %q", ip))
		pvc, err := utils.CreatePVCFromDefinition(c, ns, utils.NewPVCDefinition("transport-e2e", "20M", pvcAnn, nil))
		Expect(err).NotTo(HaveOccurred(), "Error creating PVC")

		err = utils.WaitForPersistentVolumeClaimPhase(c, ns, v1.ClaimBound, pvc.Name)
		Expect(err).NotTo(HaveOccurred(), "Error waiting for claim phase Bound")

		By("Verifying PVC is not empty")
		Expect(framework.VerifyPVCIsEmpty(f, pvc)).To(BeFalse())
	}

	DescribeTable("Transport Test Table", it,
		Entry("should connect to http endpoint without credentials", httpNoAuthEp, targetFile, "", ""),
		Entry("should connect to http endpoint with credentials", httpAuthEp, targetFile, utils.AccessKeyValue, utils.SecretKeyValue))
})

// handelError is intended for use outside It(), BeforeEach() and AfterEach() blocks where Expect() cannot be called.
func handelError(e error) {
	if e != nil {
		Fail(fmt.Sprintf("Encountered error: %v", e), 2)
	}
}

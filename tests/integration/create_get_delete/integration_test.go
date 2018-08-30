// +build integration

package create_get_delete

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/weaveworks/eksctl/pkg/testutils/aws"
	. "github.com/weaveworks/eksctl/pkg/testutils/matchers"
	"github.com/weaveworks/eksctl/tests/integration"
	//"k8s.io/client-go/tools/clientcmd"
)

const (
	createTimeout = 20
	region        = "us-west-2"
)

var (
	eksctlPath string

	// Flags to help with the development of the integration tests
	clusterName    string
	doCreate       bool
	doDelete       bool
	kubeconfigPath string

	kubeconfigTemp bool
)

func TestCreateIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration - Create Suite")
}

type tInterface interface {
	GinkgoTInterface
	Helper()
}

type tHelper struct{ GinkgoTInterface }

func (t *tHelper) Helper()      { return }
func (t *tHelper) Name() string { return "eksctl-test" }

func newKubeTest(h *harness.Harness) *harness.Test {
	return h.NewTest(&tHelper{GinkgoT()})
}

var _ = Describe("Create (Integration)", func() {

	BeforeSuite(func() {
		kubeconfigTemp = false
		if kubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
			kubeconfigPath = f.Name()
			kubeconfigTemp = true
		}
	})

	AfterSuite(func() {
		gexec.KillAndWait()
		if kubeconfigTemp {
			os.Remove(kubeconfigPath)
		}
		if doCreate && doDelete {
			integration.CleanupAws(clusterName, region)
		}
	})

	Describe("when creating a cluster with 1 node", func() {
		var (
			err     error
			session *gexec.Session
		)

		It("should not return an error", func() {
			if !doCreate {
				fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", clusterName)
				return
			}

			fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", kubeconfigPath)

			if clusterName == "" {
				clusterName = utils.ClusterName("", "")
			}

			args := []string{"create", "cluster", "--name", clusterName, "--node-type", "t2.medium", "--nodes", "1", "--region", region, "--kubeconfig", kubeconfigPath}

			command := exec.Command(eksctlPath, args...)
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)

			if err != nil {
				Fail(fmt.Sprintf("error starting process: %v", err), 1)
			}

			session.Wait(createTimeout * time.Minute)
			Expect(session.ExitCode()).Should(Equal(0))
		})

		It("should have created an EKS cluster", func() {
			session := aws.NewSession(region)
			Expect(session).To(HaveEksCluster(clusterName, awseks.ClusterStatusActive, "1.10"))
		})

		It("should have the required cloudformation stacks", func() {
			session := aws.NewSession(region)

			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-VPC", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-ControlPlane", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-ServiceRole", clusterName)))
			Expect(session).To(HaveCfnStack(fmt.Sprintf("EKS-%s-DefaultNodeGroup", clusterName)))
		})

		/*It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(pathKubeCtlConf)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
		})*/

		Context("and we create a deployment using kubectl", func() {
			var test *harness.Test

			h := harness.New(harness.Options{
				Kubeconfig: kubeconfigPath,
				Logger:     &logger.TestLogger{},
			})

			BeforeEach(func() {
				err := h.Setup()
				Expect(err).ShouldNot(HaveOccurred())
				test = newKubeTest(h)
				test.CreateNamespace(test.Namespace)
			})

			AfterEach(func() {
				test.Close()
			})

			It("should deploy the service to the cluster", func() {
				d := test.CreateDeploymentFromFile(test.Namespace, "podinfo.yaml")
				test.WaitForDeploymentReady(d, 1*time.Minute)

				pods := test.ListPodsFromDeployment(d)
				Expect(len(pods.Items)).To(Equal(2))

				// For each pod of the Deployment, check we receive a sensible response to a
				// GET request on /version.
				for _, pod := range pods.Items {
					Expect(pod.Namespace).To(Equal(test.Namespace))

					req := test.PodProxyGet(&pod, "", "/version")
					fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

					var js interface{}
					test.PodProxyGetJSON(&pod, "", "/version", &js)

					Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.0.1"))
				}
			})
		})
	})
})

func init() {
	flag.StringVar(&eksctlPath, "eksctl.path", "../../../eksctl", "Path to eksctl")

	// Flags to help with the development of the integration tests
	flag.StringVar(&clusterName, "eksctl.cluster", "", "Cluster name (default: generate one)")
	flag.BoolVar(&doCreate, "eksctl.create", true, "Skip the creation tests. Useful for debugging the tests")
	flag.BoolVar(&doDelete, "eksctl.delete", true, "Skip the cleanup after the tests have run")
	flag.StringVar(&kubeconfigPath, "eksctl.kubeconfig", "", "Path to kubeconfig (default: create it a temporary file)")
}

package domains

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cloudfoundry/cf-performance-tests/helpers"
	"github.com/cloudfoundry/cf-test-helpers/v2/workflowhelpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
)

var _ = Describe("domains", func() {
	Describe("GET /v3/domains", func() {

		It("gets /v3/domains as admin efficiently", func() {
			experiment := gmeasure.NewExperiment("GET /v3/domains as admin")
			AddReportEntry(experiment.Name, experiment) // #TODO include if using built-in Ginkgo reporter.

			experiment.Sample(func(idx int) {
				experiment.MeasureDuration("GET /v3/domains as admin", func() {
					workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
						helpers.V2TimeCFCurl(testConfig.BasicTimeout, "/v3/domains")
					})
				})
			}, gmeasure.SamplingConfig{N: testConfig.Samples, Duration: time.Duration(testConfig.SampleLength)})
		})

		It("gets /v3/domains as a regular user efficiently", func() {
			experiment := gmeasure.NewExperiment("GET /v3/domains as user")
			AddReportEntry(experiment.Name, experiment) // #TODO include if using built-in Ginkgo reporter.

			experiment.Sample(func(idx int) {
				experiment.MeasureDuration("GET /v3/domains as user", func() {
					workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
						helpers.V2TimeCFCurl(testConfig.BasicTimeout, "/v3/domains")
					})
				})
			}, gmeasure.SamplingConfig{N: testConfig.Samples, Duration: time.Duration(testConfig.SampleLength)})
		})

		It(fmt.Sprintf("gets /v3/domains as admin with page size %d efficiently", testConfig.LargePageSize), func() {
			experiment := gmeasure.NewExperiment(fmt.Sprintf("GET /v3/domains as admin with page size %d", testConfig.LargePageSize))
			AddReportEntry(experiment.Name, experiment) // #TODO include if using built-in Ginkgo reporter.

			experiment.Sample(func(idx int) {
				experiment.MeasureDuration(fmt.Sprintf("GET /v3/domains as admin with page size %d", testConfig.LargePageSize), func() {
					workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.LongTimeout, func() {
						helpers.V2TimeCFCurl(testConfig.LongTimeout, fmt.Sprintf("/v3/domains?per_page=%d", testConfig.LargePageSize))
					})
				})
			}, gmeasure.SamplingConfig{N: testConfig.Samples, Duration: time.Duration(testConfig.SampleLength)})
		})
	})

	Describe("GET /v3/organizations/:guid/domains", func() {
		Measure("as admin", func(b Benchmarker) {
			orgGUIDs := helpers.GetGUIDs(testSetup.AdminUserContext(), testConfig, "/v3/organizations")
			Expect(orgGUIDs).NotTo(BeNil())
			orgGUID := orgGUIDs[rand.Intn(len(orgGUIDs))]
			workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
				helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/organizations/%s/domains", orgGUID))
			})
		}, testConfig.Samples)

		Measure("as regular user", func(b Benchmarker) {
			orgGUIDs := helpers.GetGUIDs(testSetup.RegularUserContext(), testConfig, "/v3/organizations")
			Expect(orgGUIDs).NotTo(BeNil())
			orgGUID := orgGUIDs[rand.Intn(len(orgGUIDs))]
			workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
				helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/organizations/%s/domains", orgGUID))
			})
		}, testConfig.Samples)
	})

	Describe("individually", func() {
		Describe("as admin", func() {
			var domainGUID string
			BeforeEach(func() {
				domainGUIDs := helpers.GetGUIDs(testSetup.AdminUserContext(), testConfig, "/v3/domains")
				Expect(domainGUIDs).NotTo(BeNil())
				domainGUID = domainGUIDs[rand.Intn(len(domainGUIDs))]
			})

			Measure("GET /v3/domains/:guid", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/domains/%s", domainGUID))
				})
			}, testConfig.Samples)

			Measure("PATCH /v3/domains/:guid", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					data := `{ "metadata": { "annotations": { "test": "PATCH /v3/domains/:guid" } } }`
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, "-X", "PATCH", "-d", data, fmt.Sprintf("/v3/domains/%s", domainGUID))
				})
			}, testConfig.Samples)

			Measure("DELETE /v3/domains/:guid", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, "-X", "DELETE", fmt.Sprintf("/v3/domains/%s", domainGUID))

					// Wait until "GET /v3/domains/:guid" fails.
					helpers.WaitToFail(testSetup.AdminUserContext(), testConfig, fmt.Sprintf("/v3/domains/%s", domainGUID))
				})
			}, testConfig.Samples)
		})

		Describe("as regular user", func() {
			Measure("GET /v3/domains/:guid", func(b Benchmarker) {
				domainGUIDs := helpers.GetGUIDs(testSetup.RegularUserContext(), testConfig, "/v3/domains")
				Expect(domainGUIDs).NotTo(BeNil())
				domainGUID := domainGUIDs[rand.Intn(len(domainGUIDs))]
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/domains/%s", domainGUID))
				})
			}, testConfig.Samples)
		})
	})
})

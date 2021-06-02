package main

import (
	"testing"

	"github.com/cloudfoundry-incubator/cf-performance-tests/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var testConfig helpers.Config = helpers.NewConfig()
var testSetup *workflowhelpers.ReproducibleTestSuiteSetup

var _ = BeforeSuite(func() {
	testSetup = workflowhelpers.NewTestSuiteSetup(&testConfig)
	testSetup.Setup()
})

func TestCfPerformanceTests(t *testing.T) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.cf-performance-tests")
	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("error loading config: %s", err.Error())
	}

	err = viper.Unmarshal(&testConfig)
	if err != nil {
		t.Fatalf("error parsing config: %s", err.Error())
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "CfPerformanceTests Suite")
}
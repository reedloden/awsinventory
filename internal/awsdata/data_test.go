package awsdata_test

import (
	"testing"

	. "github.com/manywho/awsinventory/internal/awsdata"
	logrustest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestLoadExitsEarlyWhenRegionsIsEmptyAndRegionalServicesAreIncluded(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{})

	d.Load([]string{}, []string{ServiceEC2})

	require.Contains(t, hook.LastEntry().Message, ErrNoRegions.Error())
}

func TestLoadExitsEarlyWhenServicesIsEmpty(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{})

	d.Load([]string{ValidRegions[0]}, []string{})

	require.Contains(t, hook.LastEntry().Message, ErrNoServices.Error())
}

func TestLoadCatchesInvalidRegion(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{})

	d.Load([]string{"test-region"}, []string{ValidServices[0]})

	require.Contains(t, hook.LastEntry().Message, "invalid region: test-region")
}

func TestLoadCatchesInvalidService(t *testing.T) {
	logger, hook := logrustest.NewNullLogger()

	d := New(logger, TestClients{})

	d.Load([]string{ValidRegions[0]}, []string{"invalid-service"})

	require.Contains(t, hook.LastEntry().Message, "invalid service: invalid-service")
}
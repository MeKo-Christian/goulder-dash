package game

import (
	"github.com/cucumber/godog"
	"testing"
)

func TestFeatures(t *testing.T) {
	opts := godog.Options{
		Format:   "pretty",
		Paths:    []string{"../features"},
		TestingT: t, // Use the *testing.T provided by the test function
	}

	suite := godog.TestSuite{
		Name:                "godog",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	InitializeSteps(ctx)
}
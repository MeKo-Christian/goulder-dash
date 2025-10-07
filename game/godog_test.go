package game_test

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	t.Parallel()

	opts := godog.Options{
		Format:              "pretty",
		Paths:               []string{"../features"},
		TestingT:            t, // Use the *testing.T provided by the test function
		ShowStepDefinitions: false,
		Randomize:           0,
		StopOnFailure:       false,
		Strict:              false,
		NoColors:            false,
		Tags:                "",
		Concurrency:         0,
		Output:              nil,
		DefaultContext:      nil,
		FeatureContents:     nil,
		ShowHelp:            false,
	}

	suite := godog.TestSuite{
		Name:                 "godog",
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
		TestSuiteInitializer: nil,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	InitializeSteps(ctx)
}

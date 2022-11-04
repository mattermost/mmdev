package utils

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/mattermost/mmdev/model"
)

func IsCurrentVersionIncluded(currentVersion string, targetVersion model.VersionConfig) bool {
	if currentVersion == "" {
		return false
	}

	currentSemVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		panic(fmt.Errorf("cannot parse current version: %w", err))
	}

	if targetVersion.MinVersion != "" {
		c, err := semver.NewConstraint(">= " + targetVersion.MinVersion)
		if err != nil {
			panic(fmt.Errorf("cannot parse minimum version: %w", err))
		}
		if !c.Check(currentSemVer) {
			return false
		}
	}

	if targetVersion.MaxVersion != "" {
		c, err := semver.NewConstraint("<= " + targetVersion.MaxVersion)
		if err != nil {
			panic(fmt.Errorf("cannot parse maximum version: %w", err))
		}
		if !c.Check(currentSemVer) {
			return false
		}
	}

	return true
}

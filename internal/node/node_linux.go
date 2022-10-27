package node

import (
	"fmt"
)

func InstallNodeIfNeeded(version string) error {
	fmt.Printf("Check & install NodeJS version %s on Linux\n", version)

	return nil
}

package node

import (
	"fmt"
)

func InstallNodeIfNeeded(version string) error {
	fmt.Printf("Check & install NodeJS version %s on macOS\n", version)

	return nil
}

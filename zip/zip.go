package zip

import (
	"fmt"
	"os/exec"
)

func Zip(srcDir string, destFile string) error {
	fmt.Println("Zipping: ", srcDir, " -> ", destFile)
	cmd := exec.Command("zip", "-r", destFile, ".")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error zipping: ", err)
		return err
	}
	return nil
}

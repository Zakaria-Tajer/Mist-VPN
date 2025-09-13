package helpers

import (
	"log"
	"os"
	"os/exec"
)




func RunBin(bin string, args ...string) {
    cmd := exec.Command(bin, args...)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin

    err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
}
package main

import "fmt"

var (
	gitBuildCommit = "N/A"
	gitBuildTag    = "N/A"
	buildDateTime  = "N/A"
)

func main() {

	fmt.Printf("Client. Build commit: %s, build date and time: %s, build tag: %s \n",
		gitBuildCommit,
		buildDateTime,
		gitBuildTag,
	)
}

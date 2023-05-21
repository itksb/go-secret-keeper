package main

import "fmt"

var (
	gitBuildCommit string = "N/A"
	gitBuildTag    string = "N/A"
	buildDateTime  string = "N/A"
)

func main() {

	fmt.Printf("Client. Build commit: %s, build date and time: %s, build tag: %s \n",
		gitBuildCommit,
		buildDateTime,
		gitBuildTag,
	)

}

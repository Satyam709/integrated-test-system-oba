package docker_utils

// func StartServer(composeFile string)  {
// 	// check if docker containers are running
// 	cmd := "docker ps -q --filter 'name=oba_server'"
// 	output, err := RunCommand(cmd)
// 	if err != nil {
// 		panic("Failed to check if OBA server is running: " + err.Error())
// 	}
// 	if output == "" {
// 		fmt.Printf("OBA server is not running, starting it now...")
// 		SpinContainers(composeFile)
// 	} else {
// 		fmt.Printf("OBA conatiner running with container ID: %s\n", output)
// 	}

// 	cmd = "docker exec "
// }
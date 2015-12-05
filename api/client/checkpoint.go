// +build experimental

package client

import (
	"fmt"

	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/runconfig"
)

// CmdCheckpoint checkpoints the process running in a container
//
// Usage: docker checkpoint CONTAINER
func (cli *DockerCli) CmdCheckpoint(args ...string) error {
	cmd := Cli.Subcmd("checkpoint", []string{"CONTAINER [CONTAINER...]"}, "Checkpoint one or more running containers", true)
	cmd.Require(flag.Min, 1)

	var (
		flImgDir       = cmd.String([]string{"-image-dir"}, "", "directory for storing checkpoint image files")
		flWorkDir      = cmd.String([]string{"-work-dir"}, "", "directory for storing log file")
		flLeaveRunning = cmd.Bool([]string{"-leave-running"}, false, "leave the container running after checkpoint")
	)

	if err := cmd.ParseFlags(args, true); err != nil {
		return err
	}

	if cmd.NArg() < 1 {
		cmd.Usage()
		return nil
	}

	criuOpts := &runconfig.CriuConfig{
		ImagesDirectory:         *flImgDir,
		WorkDirectory:           *flWorkDir,
		LeaveRunning:            *flLeaveRunning,
		TCPEstablished:          true,
		ExternalUnixConnections: true,
		FileLocks:               true,
	}

	var encounteredError error
	for _, name := range cmd.Args() {
		_, _, err := readBody(cli.call("POST", "/containers/"+name+"/checkpoint", criuOpts, nil))
		if err != nil {
			fmt.Fprintf(cli.err, "%s\n", err)
			encounteredError = fmt.Errorf("Error: failed to checkpoint one or more containers")
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	return encounteredError
}

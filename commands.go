package flexvolume

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
)

// Commands is a list of actions supported by this flex volume
func Commands(fv FlexVolume) []cli.Command {
	var commands []cli.Command
	commands = append(commands, cli.Command{
		Name:  "init",
		Usage: "Initialize the driver",
		Action: func(c *cli.Context) error {
			var resp Response = fv.Init()
			var caps Capabilities = fv.Capabilities()
			// Format the output as JSON.
			output, err := json.Marshal(struct {
				Response
				Capabilities Capabilities `json:"capabilities,omitempty"`
			}{resp, caps})

			if err != nil {
				return handle(WrapError(errors.Wrap(err, "Unable to create json")))
			}
			fmt.Println(string(output))
			return nil
		},
	})

	if fv.Capabilities().Attach {
		commands = append(commands, cli.Command{
			Name:  "attach",
			Usage: "Attach the volume",
			Action: func(c *cli.Context) error {
				var opts map[string]string
				if err := json.Unmarshal([]byte(c.Args().Get(0)), &opts); err != nil {
					return handle(WrapError(errors.Wrapf(err, "Unable to process json: %s", c.Args().Get(0))))
				}

				return handle(fv.Attach(opts))
			},
		})
	}
	if fv.Capabilities().Detach {
		commands = append(commands, cli.Command{
			Name:  "detach",
			Usage: "Detach the volume",
			Action: func(c *cli.Context) error {
				return handle(fv.Detach(c.Args().Get(0)))
			},
		})
	}
	commands = append(commands, cli.Command{
		Name:  "mount",
		Usage: "Mount the volume",
		Action: func(c *cli.Context) error {
			var opts map[string]string

			if err := json.Unmarshal([]byte(c.Args().Get(2)), &opts); err != nil {
				return handle(WrapError(errors.Wrapf(err, "Unable to process json: |%s|", c.Args().Get(2))))
			}

			return handle(fv.Mount(c.Args().Get(0), c.Args().Get(1), opts))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "unmount",
		Usage: "Mount the volume",
		Action: func(c *cli.Context) error {
			return handle(fv.Unmount(c.Args().Get(0)))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "getvolumename",
		Usage: "Returns the unique name of the volume",
		Action: func(c *cli.Context) error {
			var opts map[string]string

			if err := json.Unmarshal([]byte(c.Args().Get(1)), &opts); err != nil {
				return err
			}

			return handle(fv.GetVolumeName(opts))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "waitforattach",
		Usage: "Waits until a volume is fully attached to a node and its device emerges",
		Action: func(c *cli.Context) error {
			var opts map[string]string

			if err := json.Unmarshal([]byte(c.Args().Get(2)), &opts); err != nil {
				return err
			}

			return handle(fv.WaitForAttach(c.Args().Get(1), opts))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "isattached",
		Usage: "Checks that a volume is attached to a node",
		Action: func(c *cli.Context) error {
			var opts map[string]string

			if err := json.Unmarshal([]byte(c.Args().Get(1)), &opts); err != nil {
				return err
			}

			return handle(fv.Detach(opts, c.Args().Get(2)))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "mountdevice",
		Usage: "Mounts a volume’s device to a directory",
		Action: func(c *cli.Context) error {
			var opts map[string]string

			if err := json.Unmarshal([]byte(c.Args().Get(3)), &opts); err != nil {
				return err
			}

			return handle(fv.MountDevice(c.Args().Get(1), c.Args().Get(2), opts))
		},
	})
	commands = append(commands, cli.Command{
		Name:  "unmountdevice",
		Usage: "Unmounts a volume’s device from a directory",
		Action: func(c *cli.Context) error {
			return handle(fv.UnmountDevice(c.Args().Get(1)))
		},
	})
	return commands
}

// The following handles:
//   * Output of the Response object.
//   * Sets an error so we can bubble up an error code.
func handle(resp Response) error {
	// Format the output as JSON.
	output, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	os.Stderr.Write(output)
	if resp.Status != StatusSuccess {
		os.Exit(1)
	}
	return nil
}

func WrapError(err error) Response {
	return Response{
		Status:  StatusFailure,
		Message: fmt.Sprintf("%s", err),
	}
}

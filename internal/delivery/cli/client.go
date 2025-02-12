package cli

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Shayan-Ghani/GOAD/config"
	"github.com/Shayan-Ghani/GOAD/internal/delivery/command"
	"github.com/Shayan-Ghani/GOAD/pkg/validation"
)

type Client struct {
	printer Printer
}

func NewClient() *Client {
	return &Client{
		printer: NewPrinter(),
	}
}

func (c *Client) Run(args []string) error {
	cmd, err := command.NewCommand(args)
	if err != nil {
		if t, isHelp := err.(validation.Help); isHelp {
			c.printer.PrintUsage(t.Message)
			return nil
		}
		return err
	}

	req := &CliRequest{
		Flags:    cmd.GetFlags(),
		Resource: args[0],
		Command:  args[1],
	}

	response, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	return c.handleResponse(req, response)
}

func (c *Client) sendRequest(req *CliRequest) (*CliResponse, error) {
	conn, err := net.Dial("tcp", config.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	if err = json.NewEncoder(conn).Encode(req); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	var res CliResponse
	if err = json.NewDecoder(conn).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &res, nil
}

func (c *Client) handleResponse(req *CliRequest, res *CliResponse) error {
	if res.Err != "" {
		return fmt.Errorf(res.Err)
	}
	if res.Items != nil{
		c.printer.PrintResponse(req.Flags.Format, res.Items)
	}

	return nil
}

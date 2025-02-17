package cli

import (
	"fmt"

	"github.com/Shayan-Ghani/GOAD/internal/delivery/cli/requesthandler"
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

func (c *Client) Run(args []string, itemSvcURL string, tagSvcURL string) error {
	cmd, err := command.NewCommand(args)
	if err != nil {
		if t, isHelp := err.(validation.Help); isHelp {
			c.printer.PrintUsage(t.Message)
			return nil
		}
		return err
	}

	req := requesthandler.CliRequest{
		Flags:    cmd.GetFlags(),
		Resource: args[0],
		Command:  args[1],
	}

	response, err := c.sendRequest(req, itemSvcURL, tagSvcURL)
	if err != nil {
		return err
	}

	if response != nil {
		return c.handleResponse(req.Flags.Format, response)
	}

	return nil
}

func (c *Client) sendRequest(req requesthandler.CliRequest, itemSvcURL string, tagSvcURL string) (*requesthandler.CliResponse, error) {

	h := requesthandler.NewHandler(itemSvcURL, tagSvcURL)

	res, err := h.HandleRequest(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) handleResponse(format string, res *requesthandler.CliResponse) error {
	if res.Err != "" {
		return fmt.Errorf(res.Err)
	}
	if res.Items != nil {
		c.printer.PrintResponse(format, res.Items)
	}

	return nil
}

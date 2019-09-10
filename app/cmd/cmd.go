package cmd

// Revision sets from main
var Revision = "unknown"

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	SetCommon(commonOpts CommonOpts)
	Execute(args []string) error
}

// CommonOpts sets externally from main, shared across all commands
type CommonOpts struct {
	SharedSecret string
}

// SetCommon satisfies CommonOptionsCommander interface and sets common option fields
func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.SharedSecret = commonOpts.SharedSecret
}

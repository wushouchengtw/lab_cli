package cmd

import (
	"fmt"
	"strings"
)

// HelpItem represents a single command or option line in the help message.
type HelpItem struct {
	Name string
	Desc string
}

// HelpBuilder constructs a formatted CLI help message with automatic column alignment.
type HelpBuilder struct {
	usage        string
	commands     []HelpItem
	options      []HelpItem
	optionsTitle string
	hint         string
}

// NewHelpBuilder creates a new HelpBuilder.
func NewHelpBuilder() *HelpBuilder {
	return &HelpBuilder{
		optionsTitle: "Options:",
	}
}

// SetUsage sets the usage line (e.g. "lab_cli build <command> [options]").
func (b *HelpBuilder) SetUsage(usage string) *HelpBuilder {
	b.usage = usage
	return b
}

// AddCommand adds a command with its description.
func (b *HelpBuilder) AddCommand(name, desc string) *HelpBuilder {
	b.commands = append(b.commands, HelpItem{Name: name, Desc: desc})
	return b
}

// AddOption adds an option/flag with its description.
func (b *HelpBuilder) AddOption(name, desc string) *HelpBuilder {
	b.options = append(b.options, HelpItem{Name: name, Desc: desc})
	return b
}

// SetOptionsTitle sets a custom section title for options (default: "Options:").
func (b *HelpBuilder) SetOptionsTitle(title string) *HelpBuilder {
	b.optionsTitle = title
	return b
}

// SetHintCommand sets the hint line referencing a command (e.g. "lab_cli build <command>").
func (b *HelpBuilder) SetHintCommand(cmd string) *HelpBuilder {
	b.hint = fmt.Sprintf("Use %q for more information about a subcommand.", cmd+" -h")
	return b
}

// SetHint sets a raw custom hint text at the end of the help message.
func (b *HelpBuilder) SetHint(hint string) *HelpBuilder {
	b.hint = hint
	return b
}

// String formats and returns the complete help message.
func (b *HelpBuilder) String() string {
	var sb strings.Builder

	// Usage section
	if b.usage != "" {
		sb.WriteString(Header("Usage:"))
		sb.WriteString("\n  ")
		sb.WriteString(b.usage)
		sb.WriteString("\n\n")
	}

	// Commands section
	if len(b.commands) > 0 {
		sb.WriteString(Header("Available Commands:"))
		sb.WriteByte('\n')
		maxLen := 0
		for _, c := range b.commands {
			if len(c.Name) > maxLen {
				maxLen = len(c.Name)
			}
		}
		for _, c := range b.commands {
			fmt.Fprintf(&sb, "  %-*s   %s\n", maxLen, c.Name, c.Desc)
		}
		sb.WriteString("\n")
	}

	// Options section
	if len(b.options) > 0 {
		title := b.optionsTitle
		if title == "" {
			title = "Options:"
		}
		sb.WriteString(Header(title))
		sb.WriteByte('\n')
		maxLen := 0
		for _, o := range b.options {
			if len(o.Name) > maxLen {
				maxLen = len(o.Name)
			}
		}
		for _, o := range b.options {
			fmt.Fprintf(&sb, "  %-*s   %s\n", maxLen, o.Name, o.Desc)
		}
		sb.WriteString("\n")
	}

	// Hint section
	if b.hint != "" {
		sb.WriteString(b.hint)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// Print outputs the generated help message to stdout.
func (b *HelpBuilder) Print() {
	fmt.Print(b.String())
}

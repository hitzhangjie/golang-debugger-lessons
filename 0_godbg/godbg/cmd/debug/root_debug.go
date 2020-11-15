package debug

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stromland/cobra-prompt"
)

const (
	commandGroup = "group_of_debug_commands"
	prefix       = "godbg> "
	description  = "interactive debugging commands"
)

// NewDebugShell 创建一个debug专用的交互管理器
func NewDebugShell() *cobraprompt.CobraPrompt {
	debugRootCmd := &cobra.Command{
		Use:   "",
		Short: description,
	}
	fn := func() func(cmd *cobra.Command) error {
		return func(cmd *cobra.Command) error {
			usage := groupDebugCommands(cmd)
			fmt.Println(usage)
			return nil
		}
	}
	debugRootCmd.SetUsageFunc(fn())
	debugRootCmd.AddCommand(breakCmd, clearCmd, exitCmd)

	return &cobraprompt.CobraPrompt{
		RootCmd:                debugRootCmd,
		//DynamicSuggestionsFunc: dynamicSuggestions,
		ResetFlagsFlag:         true,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle(description),
			prompt.OptionPrefix(prefix),
			prompt.OptionSuggestionBGColor(prompt.DarkBlue),
			prompt.OptionDescriptionBGColor(prompt.DarkBlue),
			prompt.OptionSelectedSuggestionBGColor(prompt.Red),
			prompt.OptionSelectedDescriptionBGColor(prompt.Red),
			prompt.OptionMaxSuggestion(5),
			prompt.OptionCompletionOnDown(),
		},
	}
}

func groupDebugCommands(cmd *cobra.Command) string {

	// key:group, val:sorted commands in same group
	groups := map[string][]string{}
	for _, c := range cmd.Commands() {
		// 如果没有指定命令分组，放入other组
		var groupName string
		v, ok := c.Annotations[commandGroup]
		if !ok {
			groupName = "other"
		} else {
			groupName = v
		}

		groupCmds, ok := groups[groupName]
		groupCmds = append(groupCmds, fmt.Sprintf("%-16s:\t%s", c.Use, c.Short))
		sort.Strings(groupCmds)

		groups[groupName] = groupCmds
	}

	// 按照group分组，并对组内命令进行排序
	buf := bytes.Buffer{}
	for grp, commands := range groups {
		buf.WriteString(fmt.Sprintf("[%s]\n", grp))
		for _, cmd := range commands {
			buf.WriteString(fmt.Sprintf("%s\n", cmd))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func dynamicSuggestions(annotation string, _ prompt.Document) []prompt.Suggest {
	switch annotation {
	case "GetFood":
		return GetFood()
	default:
		return []prompt.Suggest{}
	}
}

func GetFood() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "apple", Description: "Green apple"},
		{Text: "tomato", Description: "Red tomato"},
	}
}

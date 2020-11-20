package debug

import (
	"bytes"
	"fmt"
	"sort"

	"godbg/target"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stromland/cobra-prompt"
)

const (
	cmdGroupKey         = "cmd_group_key"
	cmdGroupBreakpoints = "breakpoint"
	cmdGroupSource      = "code"
	cmdGroupCtrlFlow    = "ctrlflow"
	cmdGroupInfo        = "information"
	cmdGroupOthers      = "other"

	prefix      = "godbg> "
	description = "interactive debugging commands"
)

const (
	suggestionListSourceFiles = "ListSourceFiles"
)

var (
	TraceePID   int
	breakpoints = map[uintptr]*target.Breakpoint{}
)

var debugRootCmd = &cobra.Command{
	Use:   "",
	Short: description,
}

// NewDebugShell 创建一个debug专用的交互管理器
func NewDebugShell() *cobraprompt.CobraPrompt {

	fn := func() func(cmd *cobra.Command) error {
		return func(cmd *cobra.Command) error {
			usage := groupDebugCommands(cmd)
			fmt.Println(usage)
			return nil
		}
	}
	debugRootCmd.SetUsageFunc(fn())

	return &cobraprompt.CobraPrompt{
		RootCmd:                debugRootCmd,
		DynamicSuggestionsFunc: dynamicSuggestions,
		ResetFlagsFlag:         true,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle(description),
			prompt.OptionPrefix(prefix),
			prompt.OptionSuggestionBGColor(prompt.DarkBlue),
			prompt.OptionDescriptionBGColor(prompt.DarkBlue),
			prompt.OptionSelectedSuggestionBGColor(prompt.Red),
			prompt.OptionSelectedDescriptionBGColor(prompt.Red),
			// here, hide prompt dropdown list
			// TODO do we have a better way to show/hide the prompt dropdown list?
			prompt.OptionMaxSuggestion(10),
			prompt.OptionShowCompletionAtStart(),
			prompt.OptionCompletionOnDown(),
		},
		EnableSilentPrompt: true,
		EnableShowAtStart:  true,
	}
}

// groupDebugCommands 将各个命令按照分组归类，再展示帮助信息
func groupDebugCommands(cmd *cobra.Command) string {

	// key:group, val:sorted commands in same group
	groups := map[string][]string{}
	for _, c := range cmd.Commands() {
		// 如果没有指定命令分组，放入other组
		var groupName string
		v, ok := c.Annotations[cmdGroupKey]
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

	// 按照分组名进行排序
	groupNames := []string{}
	for k, _ := range groups {
		groupNames = append(groupNames, k)
	}
	sort.Strings(groupNames)

	// 按照group分组，并对组内命令进行排序
	buf := bytes.Buffer{}
	for _, groupName := range groupNames {
		commands, _ := groups[groupName]
		buf.WriteString(fmt.Sprintf("[%s]\n", groupName))
		for _, cmd := range commands {
			buf.WriteString(fmt.Sprintf("%s\n", cmd))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func dynamicSuggestions(annotation string, _ prompt.Document) []prompt.Suggest {
	switch annotation {
	case suggestionListSourceFiles:
		return GetSourceFiles()
	default:
		return []prompt.Suggest{}
	}
}

func GetSourceFiles() []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "main.go", Description: "main.go"},
		{Text: "helloworld.go", Description: "helloworld.go"},
	}
}

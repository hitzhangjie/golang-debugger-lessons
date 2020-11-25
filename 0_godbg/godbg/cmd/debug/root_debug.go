package debug

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"godbg/target"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/stromland/cobra-prompt"
)

const (
	cmdGroupAnnotation = "cmd_group_annotation"

	cmdGroupBreakpoints = "1-breaks"
	cmdGroupSource      = "2-source"
	cmdGroupCtrlFlow    = "3-execute"
	cmdGroupInfo        = "4-info"
	cmdGroupOthers      = "5-other"
	cmdGroupCobra       = "other"

	cmdGroupDelimiter = "-"

	prefix    = "godbg> "
	descShort = "godbg interactive debugging commands"
)

const (
	suggestionListSourceFiles = "ListSourceFiles"
)

var (
	TraceePID   int
	breakpoints = map[uintptr]*target.Breakpoint{}
)

var debugRootCmd = &cobra.Command{
	Use:   "help [command]",
	Short: descShort,
}

// NewDebugShell 创建一个debug专用的交互管理器
func NewDebugShell() *cobraprompt.CobraPrompt {

	fn := func(cmd *cobra.Command, args []string) {
		// 描述信息
		fmt.Println(cmd.Short)
		fmt.Println()

		// 使用信息
		fmt.Println(cmd.Use)
		fmt.Println(cmd.Flags().FlagUsages())

		// 命令分组
		usage := helpMessageByGroups(cmd)
		fmt.Println(usage)
	}
	debugRootCmd.SetHelpFunc(fn)

	return &cobraprompt.CobraPrompt{
		RootCmd:                debugRootCmd,
		DynamicSuggestionsFunc: dynamicSuggestions,
		ResetFlagsFlag:         false,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle(descShort),
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

// helpMessageByGroups 将各个命令按照分组归类，再展示帮助信息
func helpMessageByGroups(cmd *cobra.Command) string {

	// key:group, val:sorted commands in same group
	groups := map[string][]string{}
	for _, c := range cmd.Commands() {
		// 如果没有指定命令分组，放入other组
		var groupName string
		v, ok := c.Annotations[cmdGroupAnnotation]
		if !ok {
			groupName = "other"
		} else {
			groupName = v
		}

		groupCmds, ok := groups[groupName]
		groupCmds = append(groupCmds, fmt.Sprintf("  %-16s:%s", c.Name(), c.Short))
		sort.Strings(groupCmds)

		groups[groupName] = groupCmds
	}

	if len(groups[cmdGroupCobra]) != 0 {
		groups[cmdGroupOthers] = append(groups[cmdGroupOthers], groups[cmdGroupCobra]...)
	}
	delete(groups, cmdGroupCobra)

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

		group := strings.Split(groupName, cmdGroupDelimiter)[1]
		buf.WriteString(fmt.Sprintf("- [%s]\n", group))

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

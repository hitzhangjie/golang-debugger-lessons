package debug

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var pregsCmd = &cobra.Command{
	Use:   "pregs",
	Short: "打印寄存器数据",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		regsOut := syscall.PtraceRegs{}
		err := syscall.PtraceGetRegs(TraceePID, &regsOut)
		if err != nil {
			return fmt.Errorf("get regs error: %v", err)
		}
		prettyPrintRegs(regsOut)
		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(pregsCmd)
}

func prettyPrintRegs(regs syscall.PtraceRegs) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	rt := reflect.TypeOf(regs)
	rv := reflect.ValueOf(regs)
	for i := 0; i < rv.NumField(); i++ {
		fmt.Fprintf(w, "Register\t%s\t%#x\t\n", rt.Field(i).Name, rv.Field(i).Uint())
	}
	w.Flush()
}

package cmd

import (
	"fmt"
	"github.com/dyammarcano/cobra-cli-m/licenses"
	"os"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	packageName string
	parentName  string

	addCmd = &cobra.Command{
		Use:     "add [command name]",
		Aliases: []string{"command"},
		Short:   "Add a command to a Cobra Application",
		Long: `Add (cobra-cli add) will create a new command, with a license and
the appropriate structure for a Cobra-based CLI application,
and register it to its parent (default rootCmd).

If you want your command to be public, pass in the command name
with an initial uppercase letter.

Example: cobra-cli add server -> resulting in a new cmd/server.go`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var comps []string
			if len(args) == 0 {
				comps = cobra.AppendActiveHelp(comps, "Please specify the name for the new command")
			} else if len(args) == 1 {
				comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
			} else {
				comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
			}
			return comps, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cobra.CheckErr(fmt.Errorf("add needs a name for the command"))
			}

			wd, err := os.Getwd()
			cobra.CheckErr(err)

			commandName := validateCmdName(args[0])
			command := &Command{
				CmdName:   commandName,
				CmdParent: parentName,
				Project: &Project{
					AbsolutePath: wd,
					Legal:        licenses.GetLicense(userLicense),
					Copyright:    licenses.CopyrightLine(),
				},
			}

			cobra.CheckErr(command.Create())

			fmt.Printf("%s created at %s\n", command.CmdName, command.AbsolutePath)
		},
	}
)

func init() {
	addCmd.Flags().StringVarP(&packageName, "package", "t", "", "target package name (e.g. github.com/spf13/hugo)")
	addCmd.Flags().StringVarP(&parentName, "parent", "p", "rootCmd", "variable name of parent command for this command")
	cobra.CheckErr(addCmd.Flags().MarkDeprecated("package", "this operation has been removed."))
}

// validateCmdName returns source without any dashes and underscore.
// If there will be dash or underscore, next letter will be uppered.
// It supports only ASCII (1-byte character) strings.
// https://github.com/spf13/cobra/issues/269
func validateCmdName(source string) string {
	output := convertDashesAndUnderscoresToUppercase(source)

	if output == "" {
		return source // source is initially valid name.
	}
	return output
}

func convertDashesAndUnderscoresToUppercase(source string) string {
	var output string
	for i := 0; i < len(source); i++ {
		if isDashOrUnderscore(source[i]) {
			if output == "" {
				output = source[:i]
			}

			if isLastRune(i, source) {
				break
			}

			if isDashOrUnderscore(source[i+1]) {
				i++
				continue
			}

			output += uppercaseNextLetter(source, i)
			i += 2
			continue
		}

		if output != "" {
			output += string(source[i])
		}
	}
	return output
}

func isDashOrUnderscore(char byte) bool {
	return char == '-' || char == '_'
}

func isLastRune(index int, source string) bool {
	return index == len(source)-1
}

func uppercaseNextLetter(source string, index int) string {
	return string(unicode.ToUpper(rune(source[index+1])))
}

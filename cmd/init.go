package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dyammarcano/cobra-cli-m/licenses"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use:     "init [path]",
		Aliases: []string{"initialize", "initialise", "create"},
		Short:   "Initialize a Cobra Application",
		Long: `Initialize (cobra-cli init) will create a new application, with a license
and the appropriate structure for a Cobra-based CLI application.

Cobra init must be run inside of a go module (please run "go mod init <MODNAME>" first)
`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var comps []string
			var directive cobra.ShellCompDirective
			if len(args) == 0 {
				comps = cobra.AppendActiveHelp(comps, "Optionally specify the path of the go module to initialize")
				directive = cobra.ShellCompDirectiveDefault
			} else if len(args) == 1 {
				comps = cobra.AppendActiveHelp(comps, "This command does not take any more arguments (but may accept flags)")
				directive = cobra.ShellCompDirectiveNoFileComp
			} else {
				comps = cobra.AppendActiveHelp(comps, "ERROR: Too many arguments specified")
				directive = cobra.ShellCompDirectiveNoFileComp
			}
			return comps, directive
		},
		Run: func(_ *cobra.Command, args []string) {
			projectPath, err := initializeProject(args)
			cobra.CheckErr(err)
			cobra.CheckErr(goGet("github.com/spf13/cobra"))
			if viper.GetBool("useViper") {
				cobra.CheckErr(goGet("github.com/spf13/viper"))
			}
			fmt.Printf("Your Cobra application is ready at\n%s\n", projectPath)
		},
	}
)

func initializeProject(args []string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if args[0] != "." {
			wd = fmt.Sprintf("%s/%s", wd, args[0])
		}
	}

	modName := getModImportPath()

	project := &Project{
		AbsolutePath: wd,
		PkgName:      modName,
		Legal:        licenses.GetLicense(userLicense),
		Copyright:    licenses.CopyrightLine(),
		Viper:        viper.GetBool("useViper"),
		AppName:      path.Base(modName),
	}

	if err := project.Create(); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

func getModImportPath() string {
	mod, cd := parseModInfo()
	return path.Join(mod.Path, fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)))
}

func fileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))
	return path.Join(i...)
}

func parseModInfo() (Mod, CurDir) {
	var mod Mod
	var dir CurDir

	m := modInfoJSON("-m")
	cobra.CheckErr(json.Unmarshal(m, &mod))

	// Unsure why, but if no module is present Path is set to this string.
	if mod.Path == "command-line-arguments" {
		cobra.CheckErr("Please run `go mod init <MODNAME>` before `cobra-cli init`")
	}

	e := modInfoJSON("-e")
	cobra.CheckErr(json.Unmarshal(e, &dir))

	return mod, dir
}

type Mod struct {
	Path, Dir, GoMod string
}

type CurDir struct {
	Dir string
}

func goGet(mod string) error {
	return exec.Command("go", "get", mod).Run()
}

func modInfoJSON(args ...string) []byte {
	cmdArgs := append([]string{"list", "-json"}, args...)
	out, err := exec.Command("go", cmdArgs...).Output()
	cobra.CheckErr(err)

	return out
}

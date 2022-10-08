package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose        bool
	targetBoard    string
	installDestDir string

	rootCmd = &cobra.Command{
		Use:   "ccps",
		Short: "ccps is a development kit for CPS-1 board",
	}
	buildCmd = &cobra.Command{
		Use:     "build",
		Short:   "builds the target board",
		Long:    "builds the target board",
		Run:     build,
		Example: "ccps build -v -b sf2",
	}
	installCmd = &cobra.Command{
		Use:     "install",
		Short:   "installs into target directory",
		Long:    "installs into target directory",
		Run:     install,
		Example: "ccps install -v -d outputDir",
	}
	cleanCmd = &cobra.Command{
		Use:     "clean",
		Short:   "deletes the .tmp and out directories",
		Long:    "deletes the .tmp and out directories",
		Run:     clean,
		Example: "ccps clean",
	}
	dumpGFXCmd = &cobra.Command{
		Use:     "dumpgfx",
		Short:   "dumps the gfx from the board",
		Long:    "dumps the gfx from the board",
		Run:     dumpGFX,
		Example: "ccps dumpgfx -v -b sf2",
	}
	dumpSFXCmd = &cobra.Command{
		Use:     "dumpsfx",
		Short:   "dumps the sfx from the board",
		Long:    "dumps the sfx from the board",
		Run:     dumpSFX,
		Example: "ccps dumpsfx -v -b sf2",
	}
	postCmd = &cobra.Command{
		Use:     "post",
		Short:   "",
		Long:    "",
		Run:     post,
		Example: "ccps post -v -b sf2",
	}
	helloWorldCmd = &cobra.Command{
		Use:     "hw",
		Short:   "Hello World example",
		Long:    "Hello World example",
		Run:     helloWorld,
		Example: "ccps hw -v -b sf2",
	}
)

type FlagBuilder func(*cobra.Command)

func StringFlagBuilder(dst *string, long, short, defaultValue, help string) FlagBuilder {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringVarP(dst, long, short, defaultValue, help)
	}
}

func BoolFlagBuilder(dst *bool, long, short string, defaultValue bool, help string) FlagBuilder {
	return func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolVarP(dst, long, short, defaultValue, help)
	}
}

func init() {
	verboseFlag := BoolFlagBuilder(&verbose, "verbose", "v", false, "enables the verbose mode")
	targetBoardFlag := StringFlagBuilder(&targetBoard, "board", "b", "", "target board")
	verboseFlag(buildCmd)
	targetBoardFlag(buildCmd)
	buildCmd.MarkFlagRequired("board")

	installDestDirFlag := StringFlagBuilder(&installDestDir, "destination", "d", "", "destination directory")
	verboseFlag(installCmd)
	installDestDirFlag(installCmd)
	installCmd.MarkFlagRequired("destination")

	verboseFlag(dumpGFXCmd)
	targetBoardFlag(dumpGFXCmd)
	dumpGFXCmd.MarkFlagRequired("board")

	verboseFlag(dumpSFXCmd)
	targetBoardFlag(dumpSFXCmd)
	dumpSFXCmd.MarkFlagRequired("board")

	verboseFlag(postCmd)
	targetBoardFlag(postCmd)
	postCmd.MarkFlagRequired("board")

	verboseFlag(helloWorldCmd)
	targetBoardFlag(helloWorldCmd)
	helloWorldCmd.MarkFlagRequired("board")

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(dumpGFXCmd)
	rootCmd.AddCommand(dumpSFXCmd)
	rootCmd.AddCommand(postCmd)
	rootCmd.AddCommand(helloWorldCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

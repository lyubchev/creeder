package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/impzero/creeder/pkg/file"

	"github.com/spf13/cobra"
)

var (
	filter string
	ignore string
)

var rootCmd = &cobra.Command{
	Use:   "creeder [path]",
	Short: "Creeder is a tool that prints directory tree and file contents",
	Long: `Creeder is a tool that prints directory tree and file contents.
	
The tool accepts a path to a source code directory, reads all files, 
and prints the files tree of that directory together with the content of the files.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if file.ShouldIgnorePath(path, ignore) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if info.IsDir() {
				fmt.Printf("%s/\n", path)
				return nil
			}

			if file.ShouldIncludeFile(path, filter) {
				fmt.Printf("%s\n", path)
				content, err := file.ReadFile(path)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", content)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to scan directory: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "comma-separated list of file extensions to include in the output")
	rootCmd.Flags().StringVarP(&ignore, "ignore", "i", "", "comma-separated list of directories or file names to ignore")

	rootCmd.MarkFlagRequired("filter")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

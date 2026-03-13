// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mendixlabs/mxcli/cmd/mxcli/playwright"
	"github.com/spf13/cobra"
)

var playwrightCmd = &cobra.Command{
	Use:   "playwright",
	Short: "Browser testing with playwright-cli",
	Long:  `Commands for running browser-based verification tests using playwright-cli.`,
}

var playwrightVerifyCmd = &cobra.Command{
	Use:   "verify <file|dir> [file|dir...]",
	Short: "Run playwright-cli test scripts against a running Mendix app",
	Long: `Run .test.sh scripts that use playwright-cli to verify a Mendix application.

Test scripts are plain bash files using playwright-cli commands. Each script
runs sequentially, and a non-zero exit code marks the script as failed.
On failure, a screenshot is automatically captured for debugging.

Script naming convention: tests/verify-<name>.test.sh

Examples:
  # Run all test scripts in a directory
  mxcli playwright verify tests/ -p app.mpr

  # Run a specific script
  mxcli playwright verify tests/verify-customers.test.sh

  # Output JUnit XML for CI
  mxcli playwright verify tests/ -p app.mpr --junit results.xml

  # List scripts without executing
  mxcli playwright verify tests/ --list

  # Verbose output (show script stdout/stderr)
  mxcli playwright verify tests/ -p app.mpr --verbose

  # Custom app URL
  mxcli playwright verify tests/ --base-url http://localhost:9090
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		list, _ := cmd.Flags().GetBool("list")
		junitOutput, _ := cmd.Flags().GetString("junit")
		verbose, _ := cmd.Flags().GetBool("verbose")
		color, _ := cmd.Flags().GetBool("color")
		baseURL, _ := cmd.Flags().GetString("base-url")
		skipHealth, _ := cmd.Flags().GetBool("skip-health-check")
		timeoutStr, _ := cmd.Flags().GetString("timeout")
		projectPath, _ := cmd.Flags().GetString("project")

		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid timeout: %v\n", err)
			os.Exit(1)
		}

		if list {
			if err := playwright.ListScripts(args, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		opts := playwright.VerifyOptions{
			ProjectPath:     projectPath,
			TestFiles:       args,
			BaseURL:         baseURL,
			Timeout:         timeout,
			JUnitOutput:     junitOutput,
			Color:           color,
			Verbose:         verbose,
			SkipHealthCheck: skipHealth,
			Stdout:          os.Stdout,
			Stderr:          os.Stderr,
		}

		result, err := playwright.Verify(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !result.AllPassed() {
			os.Exit(1)
		}
	},
}

func init() {
	playwrightVerifyCmd.Flags().BoolP("list", "l", false, "List test scripts without executing")
	playwrightVerifyCmd.Flags().StringP("junit", "j", "", "Write JUnit XML results to file")
	playwrightVerifyCmd.Flags().BoolP("verbose", "v", false, "Show script stdout/stderr")
	playwrightVerifyCmd.Flags().BoolP("color", "", false, "Use colored output")
	playwrightVerifyCmd.Flags().StringP("timeout", "t", "2m", "Timeout per script execution")
	playwrightVerifyCmd.Flags().StringP("base-url", "", "http://localhost:8080", "Mendix app base URL")
	playwrightVerifyCmd.Flags().BoolP("skip-health-check", "", false, "Skip app reachability check")

	playwrightCmd.AddCommand(playwrightVerifyCmd)
}

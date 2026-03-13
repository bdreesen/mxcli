// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mendixlabs/mxcli/cmd/mxcli/docker"
	"github.com/mendixlabs/mxcli/sdk/mpr"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup development tools",
	Long: `Download and configure tools required for Mendix development.

Subcommands:
  mxbuild    Download MxBuild for the project's Mendix version
  mxruntime  Download the Mendix runtime for the project's Mendix version

Examples:
  mxcli setup mxbuild -p app.mpr
  mxcli setup mxbuild --version 11.6.3
  mxcli setup mxruntime -p app.mpr
  mxcli setup mxruntime --version 11.6.3
`,
}

var setupMxBuildCmd = &cobra.Command{
	Use:   "mxbuild",
	Short: "Download MxBuild from the Mendix CDN",
	Long: `Download and cache MxBuild for a specific Mendix version.

The version is detected from the project file (--project) or specified
explicitly (--version). The binary is cached at ~/.mxcli/mxbuild/{version}/
and automatically found by 'mxcli docker build' and 'mxcli docker check'.

Examples:
  mxcli setup mxbuild -p app.mpr
  mxcli setup mxbuild --version 11.6.3
  mxcli setup mxbuild -p app.mpr --dry-run
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("project")
		versionStr, _ := cmd.Flags().GetString("version")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Determine version
		if versionStr == "" && projectPath == "" {
			fmt.Fprintln(os.Stderr, "Error: specify --project (-p) or --version")
			os.Exit(1)
		}

		if versionStr == "" {
			// Detect from project
			reader, err := mpr.Open(projectPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening project: %v\n", err)
				os.Exit(1)
			}
			pv := reader.ProjectVersion()
			reader.Close()
			versionStr = pv.ProductVersion
			fmt.Fprintf(os.Stdout, "Detected Mendix version: %s\n", versionStr)
		}

		if dryRun {
			url := docker.MxBuildCDNURL(versionStr, runtime.GOARCH)
			cacheDir, _ := docker.MxBuildCacheDir(versionStr)
			fmt.Fprintf(os.Stdout, "Dry run:\n")
			fmt.Fprintf(os.Stdout, "  Version:      %s\n", versionStr)
			fmt.Fprintf(os.Stdout, "  Architecture: %s\n", runtime.GOARCH)
			fmt.Fprintf(os.Stdout, "  URL:          %s\n", url)
			fmt.Fprintf(os.Stdout, "  Cache dir:    %s\n", cacheDir)

			if cached := docker.CachedMxBuildPath(versionStr); cached != "" {
				fmt.Fprintf(os.Stdout, "  Status:       already cached at %s\n", cached)
			} else {
				fmt.Fprintf(os.Stdout, "  Status:       not cached, would download\n")
			}
			return
		}

		path, err := docker.DownloadMxBuild(versionStr, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "\nMxBuild ready: %s\n", path)
	},
}

var setupMxRuntimeCmd = &cobra.Command{
	Use:   "mxruntime",
	Short: "Download the Mendix runtime from the Mendix CDN",
	Long: `Download and cache the Mendix runtime for a specific Mendix version.

The version is detected from the project file (--project) or specified
explicitly (--version). The runtime is cached at ~/.mxcli/runtime/{version}/
and automatically used by 'mxcli docker build' when the PAD output does not
include the runtime (MxBuild 11.6.3+).

Examples:
  mxcli setup mxruntime -p app.mpr
  mxcli setup mxruntime --version 11.6.3
  mxcli setup mxruntime -p app.mpr --dry-run
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectPath, _ := cmd.Flags().GetString("project")
		versionStr, _ := cmd.Flags().GetString("version")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Determine version
		if versionStr == "" && projectPath == "" {
			fmt.Fprintln(os.Stderr, "Error: specify --project (-p) or --version")
			os.Exit(1)
		}

		if versionStr == "" {
			// Detect from project
			reader, err := mpr.Open(projectPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening project: %v\n", err)
				os.Exit(1)
			}
			pv := reader.ProjectVersion()
			reader.Close()
			versionStr = pv.ProductVersion
			fmt.Fprintf(os.Stdout, "Detected Mendix version: %s\n", versionStr)
		}

		if dryRun {
			url := docker.RuntimeCDNURL(versionStr)
			cacheDir, _ := docker.RuntimeCacheDir(versionStr)
			fmt.Fprintf(os.Stdout, "Dry run:\n")
			fmt.Fprintf(os.Stdout, "  Version:   %s\n", versionStr)
			fmt.Fprintf(os.Stdout, "  URL:       %s\n", url)
			fmt.Fprintf(os.Stdout, "  Cache dir: %s\n", cacheDir)

			if cached := docker.CachedRuntimePath(versionStr); cached != "" {
				fmt.Fprintf(os.Stdout, "  Status:    already cached at %s\n", cached)
			} else {
				fmt.Fprintf(os.Stdout, "  Status:    not cached, would download\n")
			}
			return
		}

		path, err := docker.DownloadRuntime(versionStr, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "\nMendix runtime ready: %s\n", path)
	},
}

func init() {
	setupMxBuildCmd.Flags().String("version", "", "Mendix version to download (e.g., 11.6.3)")
	setupMxBuildCmd.Flags().Bool("dry-run", false, "Show what would be downloaded without downloading")

	setupMxRuntimeCmd.Flags().String("version", "", "Mendix version to download (e.g., 11.6.3)")
	setupMxRuntimeCmd.Flags().Bool("dry-run", false, "Show what would be downloaded without downloading")

	setupCmd.AddCommand(setupMxBuildCmd)
	setupCmd.AddCommand(setupMxRuntimeCmd)
	rootCmd.AddCommand(setupCmd)
}

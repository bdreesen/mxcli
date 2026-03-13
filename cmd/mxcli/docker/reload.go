// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ReloadOptions configures the docker reload command.
type ReloadOptions struct {
	// ProjectPath is the path to the .mpr file.
	ProjectPath string

	// MxBuildPath is an explicit path to the mxbuild executable.
	MxBuildPath string

	// SkipCheck skips mx check before build.
	SkipCheck bool

	// SkipBuild skips mxbuild, just calls reload_model.
	SkipBuild bool

	// CSSOnly calls update_styling only (no build, no reload).
	CSSOnly bool

	// Host is the M2EE admin API host.
	Host string

	// Port is the M2EE admin API port.
	Port int

	// Token is the M2EE admin password.
	Token string

	// Direct bypasses docker exec for admin API calls.
	Direct bool

	// Stdout for output messages.
	Stdout io.Writer

	// Stderr for error output.
	Stderr io.Writer
}

// Reload performs a hot reload of the Mendix application.
//
// Modes:
//   - CSSOnly: calls update_styling (instant, no build or model reload)
//   - SkipBuild: calls reload_model without rebuilding first
//   - Default: runs Build() then reload_model
func Reload(opts ReloadOptions) error {
	w := opts.Stdout
	if w == nil {
		w = os.Stdout
	}

	m2eeOpts := M2EEOptions{
		Host:        opts.Host,
		Port:        opts.Port,
		Token:       opts.Token,
		ProjectPath: opts.ProjectPath,
		Direct:      opts.Direct,
		Timeout:     30 * time.Second,
	}

	// CSS-only mode: just update styling
	if opts.CSSOnly {
		resp, err := CallM2EE(m2eeOpts, "update_styling", nil)
		if err != nil {
			return fmt.Errorf("update_styling failed: %w", err)
		}
		if errMsg := resp.M2EEError(); errMsg != "" {
			return fmt.Errorf("update_styling failed: %s", errMsg)
		}
		fmt.Fprintln(w, "Styling updated.")
		return nil
	}

	// Build step (unless --model-only)
	if !opts.SkipBuild {
		buildOpts := BuildOptions{
			ProjectPath: opts.ProjectPath,
			MxBuildPath: opts.MxBuildPath,
			SkipCheck:   opts.SkipCheck,
			Stdout:      w,
		}
		if err := Build(buildOpts); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		fmt.Fprintln(w, "")
	}

	// Reload model
	fmt.Fprintln(w, "Reloading model...")
	resp, err := CallM2EE(m2eeOpts, "reload_model", nil)
	if err != nil {
		return fmt.Errorf("reload_model failed: %w", err)
	}
	if errMsg := resp.M2EEError(); errMsg != "" {
		return fmt.Errorf("reload failed: %s", errMsg)
	}

	// Try to extract duration from feedback.startup_metrics.duration
	if durationStr := extractReloadDuration(resp.Feedback()); durationStr != "" {
		fmt.Fprintf(w, "Model reloaded (%s).\n", durationStr)
	} else {
		fmt.Fprintln(w, "Model reloaded.")
	}

	return nil
}

// extractReloadDuration extracts the duration from feedback.startup_metrics.duration.
func extractReloadDuration(feedback map[string]any) string {
	if feedback == nil {
		return ""
	}

	metrics, ok := feedback["startup_metrics"]
	if !ok {
		return ""
	}

	metricsMap, ok := metrics.(map[string]any)
	if !ok {
		return ""
	}

	duration, ok := metricsMap["duration"]
	if !ok {
		return ""
	}

	// Duration may be a float64 (JSON number) representing milliseconds
	switch d := duration.(type) {
	case float64:
		if d < 1000 {
			return fmt.Sprintf("%.0fms", d)
		}
		return fmt.Sprintf("%.1fs", d/1000)
	case string:
		return d
	default:
		return fmt.Sprintf("%v", d)
	}
}

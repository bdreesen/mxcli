// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveMxBuild_ExplicitPath(t *testing.T) {
	dir := t.TempDir()
	fakeBin := filepath.Join(dir, "mxbuild")
	if runtime.GOOS == "windows" {
		fakeBin += ".exe"
	}
	os.WriteFile(fakeBin, []byte("fake"), 0755)

	result, err := resolveMxBuild(fakeBin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != fakeBin {
		t.Errorf("expected %s, got %s", fakeBin, result)
	}
}

func TestResolveMxBuild_ExplicitDir_FindsBinaryInRoot(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, mxbuildBinaryName())
	os.WriteFile(bin, []byte("fake"), 0755)

	result, err := resolveMxBuild(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != bin {
		t.Errorf("expected %s, got %s", bin, result)
	}
}

func TestResolveMxBuild_ExplicitDir_FindsBinaryInModeler(t *testing.T) {
	dir := t.TempDir()
	modelerDir := filepath.Join(dir, "modeler")
	os.MkdirAll(modelerDir, 0755)
	bin := filepath.Join(modelerDir, mxbuildBinaryName())
	os.WriteFile(bin, []byte("fake"), 0755)

	result, err := resolveMxBuild(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != bin {
		t.Errorf("expected %s, got %s", bin, result)
	}
}

func TestResolveMxBuild_ExplicitDir_NoBinaryInside(t *testing.T) {
	dir := t.TempDir()
	_, err := resolveMxBuild(dir)
	if err == nil {
		t.Error("expected error for directory without mxbuild binary")
	}
}

func TestResolveMxBuild_ExplicitPathNotFound(t *testing.T) {
	_, err := resolveMxBuild("/nonexistent/mxbuild")
	if err == nil {
		t.Error("expected error for nonexistent explicit path")
	}
}

func TestResolveMxBuild_NoExplicitPath_FallsThrough(t *testing.T) {
	// Without mxbuild in PATH or known locations, this should error
	_, err := resolveMxBuild("")
	if err == nil {
		// It's possible mxbuild is actually installed; skip in that case
		t.Skip("mxbuild found on system")
	}
}

func TestIsJDK21_InvalidPath(t *testing.T) {
	if isJDK21("/nonexistent/jdk") {
		t.Error("expected false for nonexistent path")
	}
}

func TestMxbuildSearchPaths_NonEmpty(t *testing.T) {
	paths := mxbuildSearchPaths()
	if len(paths) == 0 {
		t.Error("expected non-empty search paths")
	}
}

func TestJdkSearchPaths_NonEmpty(t *testing.T) {
	paths := jdkSearchPaths()
	if len(paths) == 0 {
		t.Error("expected non-empty search paths")
	}
}

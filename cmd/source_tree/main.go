// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// TTY detection
	if fi, err := os.Stdout.Stat(); err == nil {
		useColor = fi.Mode()&os.ModeCharDevice != 0
	}

	// Parse arguments
	coverFlag := false
	dupFlag := false
	fnFlag := false
	churnFlag := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--cover":
			coverFlag = true
		case "--dup":
			dupFlag = true
		case "--fn":
			fnFlag = true
		case "--churn":
			churnFlag = true
		case "--all":
			dupFlag = true
			fnFlag = true
			churnFlag = true
		}
	}

	// Find project root (directory containing go.mod, searching upward)
	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Get module path
	modPath := runCmd(root, "go", "list", "-m")
	modPath = strings.TrimSpace(modPath)

	// Phase 1: Collect files
	files := collectFiles(root)

	// Phase 1.5: Function metrics
	hasFn := false
	if fnFlag {
		computeFnMetrics(root, files)
		hasFn = true
	}

	// Phase 1.6: Duplication detection
	hasDup := false
	if dupFlag {
		computeDuplication(root, files)
		hasDup = true
	}

	// Phase 1.7: Churn (commit frequency)
	hasChurn := false
	if churnFlag {
		computeChurn(root, files)
		hasChurn = true
	}

	// Phase 1.6: Coverage
	if coverFlag {
		fmt.Fprintf(os.Stderr, "%sRunning tests with coverage (this may take a while)...%s\n", c(codeDim), c(codeReset))
		cmd := exec.Command("go", "test", "-coverprofile=coverage.out", "./...")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		_ = cmd.Run() // ignore errors (some tests may fail)
	}

	var coverage map[string]covData
	hasCoverage := false
	covFile := filepath.Join(root, "coverage.out")
	if _, err := os.Stat(covFile); err == nil {
		coverage = parseCoverage(covFile, modPath)
		hasCoverage = true
	}

	// Phase 1.6: Commit data
	commits, commitSubjects, commitData := collectCommits(root)

	// Filter out commits that have no changes in any displayed file
	fileSet := map[string]bool{}
	for _, f := range files {
		fileSet[f.path] = true
	}
	var visibleCommits []string
	for _, ch := range commits {
		if chMap, ok := commitData[ch]; ok {
			for fp := range chMap {
				if fileSet[fp] {
					visibleCommits = append(visibleCommits, ch)
					break
				}
			}
		}
	}
	commits = visibleCommits
	hasCommits := len(commits) > 0

	// Phase 2: Dependency depth
	depths := computeDepths(root, modPath)

	// Add TypeScript directories at max+1
	maxDepth := 0
	for _, d := range depths {
		if d > maxDepth {
			maxDepth = d
		}
	}
	tsDirs := map[string]bool{}
	for _, f := range files {
		if strings.HasSuffix(f.path, ".ts") {
			dir := filepath.Dir(f.path)
			tsDirs[dir] = true
		}
	}
	for dir := range tsDirs {
		depths[dir] = maxDepth + 1
	}

	// Phase 3: Find max lines for bar scaling
	maxLines := 0
	for _, f := range files {
		if f.lines > maxLines {
			maxLines = f.lines
		}
	}
	if maxLines == 0 {
		maxLines = 1
	}

	// Phase 4: Render
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// Header
	fmt.Fprintf(w, "%sSource Code Tree%s %s(sorted by dependency depth — leaf packages last)%s\n",
		c(codeBold), c(codeReset), c(codeDim), c(codeReset))
	fmt.Fprintf(w, "  %sLines:%s    %s██%s <500  %s██%s 500-999  %s██%s 1000-1499  %s██%s 1500+\n",
		c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
		c(codeOrange), c(codeReset), c(codeRed), c(codeReset))
	if hasFn {
		fmt.Fprintf(w, "  %sMaxFn:%s   %s██%s <30   %s██%s 30-59   %s██%s 60-99   %s██%s 100+  %s(longest function in lines)%s\n",
			c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
			c(codeOrange), c(codeReset), c(codeRed), c(codeReset),
			c(codeDim), c(codeReset))
	}
	if hasDup {
		fmt.Fprintf(w, "  %sDup:  %s    %s██%s <10%%  %s██%s 10-24%%  %s██%s 25-39%%  %s██%s 40%%+\n",
			c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
			c(codeOrange), c(codeReset), c(codeRed), c(codeReset))
	}
	if hasChurn {
		fmt.Fprintf(w, "  %sChurn:%s   %s██%s 1-10  %s██%s 11-25  %s██%s 26-50  %s██%s 51+   %s(commits touching file)%s\n",
			c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
			c(codeOrange), c(codeReset), c(codeRed), c(codeReset),
			c(codeDim), c(codeReset))
	}
	if hasCoverage {
		fmt.Fprintf(w, "  %sCoverage:%s %s██%s 80%%+  %s██%s 60-79%%  %s██%s 30-59%%  %s██%s <30%%\n",
			c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
			c(codeOrange), c(codeReset), c(codeRed), c(codeReset))
	}
	if hasCommits {
		fmt.Fprintf(w, "  %sCommits:%s  %s██%s <50   %s██%s 50-199  %s██%s 200-499  %s██%s 500+\n",
			c(codeDim), c(codeReset), c(codeGreen), c(codeReset), c(codeYellow), c(codeReset),
			c(codeOrange), c(codeReset), c(codeRed), c(codeReset))
	}
	// Print column label row when extra columns are active
	if hasFn || hasDup || hasChurn || hasCoverage || hasCommits {
		basePad := 2 + 4 + 36 + 5 + 2 + barMax
		fmt.Fprintf(w, "%*s", basePad, "")
		if hasFn {
			fmt.Fprintf(w, "  %s#fn  max%s", c(codeDim), c(codeReset))
		}
		if hasDup {
			fmt.Fprintf(w, "  %s dup%*s%s", c(codeDim), 1+dupBarMax, "", c(codeReset))
		}
		if hasChurn {
			fmt.Fprintf(w, "  %schurn%s", c(codeDim), c(codeReset))
		}
		if hasCoverage {
			fmt.Fprintf(w, "  %s cov%*s%s", c(codeDim), 1+covBarMax, "", c(codeReset))
		}
		if hasCommits {
			fmt.Fprintf(w, "  ")
			for _, ch := range commits {
				short := ch
				if len(short) > 4 {
					short = short[:4]
				}
				fmt.Fprintf(w, "%s%6s%s", c(codeDim), short, c(codeReset))
			}
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)

	// Sort directories by tier, then name
	type dirTier struct {
		dir  string
		tier int
	}
	var sortedDirs []dirTier
	for dir, tier := range depths {
		sortedDirs = append(sortedDirs, dirTier{dir, tier})
	}
	sort.Slice(sortedDirs, func(i, j int) bool {
		if sortedDirs[i].tier != sortedDirs[j].tier {
			return sortedDirs[i].tier > sortedDirs[j].tier
		}
		return sortedDirs[i].dir < sortedDirs[j].dir
	})

	prevTier := -1
	for _, dt := range sortedDirs {
		dirFiles := filesInDir(files, dt.dir)
		if len(dirFiles) == 0 {
			continue
		}

		// Tier separator
		if dt.tier != prevTier {
			if prevTier != -1 {
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, "%s──── tier %d ", c(codeDim), dt.tier)
			for range 60 {
				w.WriteString("─")
			}
			fmt.Fprintf(w, "%s\n", c(codeReset))
			prevTier = dt.tier
		}

		// Directory totals
		dirTotal := 0
		for _, f := range dirFiles {
			dirTotal += f.lines
		}

		displayDir := dt.dir + "/"
		if dt.dir == "." {
			displayDir = "(root)"
		}

		fmt.Fprintf(w, "%s%s%-48s%s %s%6d%s lines  %s(%d files)%s\n",
			c(codeBold), c(codeCyan), displayDir, c(codeReset),
			c(codeYellow), dirTotal, c(codeReset),
			c(codeDim), len(dirFiles), c(codeReset))

		// Print files with tree connectors and bars
		for i, f := range dirFiles {
			filename := filepath.Base(f.path)
			connector := "├── "
			if i == len(dirFiles)-1 {
				connector = "└── "
			}

			bar := makeBar(f.lines, maxLines)
			bColor := barColor(f.lines)

			// Base columns
			fmt.Fprintf(w, "  %s%s%s%-36s%s%5d%s  %s%s%s",
				c(codeDim), connector, c(codeReset),
				filename,
				c(codeWhite), f.lines, c(codeReset),
				c(bColor), bar, c(codeReset))

			// Function metrics columns
			if hasFn {
				if f.fnData != nil && f.fnData.count > 0 {
					mClr := maxFnColor(f.fnData.maxLen)
					fmt.Fprintf(w, "  %s%3d%s %s%4d%s",
						c(codeDim), f.fnData.count, c(codeReset),
						c(mClr), f.fnData.maxLen, c(codeReset))
				} else {
					fmt.Fprintf(w, "  %*s", 3+1+4, "")
				}
			}

			// Duplication bar
			if hasDup {
				dpct := f.dupPct
				dClr := dupBarColor(dpct)
				filledLen := dpct * dupBarMax / 100
				dupFilled := strings.Repeat("█", filledLen)
				dupEmpty := strings.Repeat("░", dupBarMax-filledLen)
				fmt.Fprintf(w, "  %s%3d%% %s%s%s%s%s",
					c(dClr), dpct, dupFilled, c(codeReset), c(codeDim), dupEmpty, c(codeReset))
			}

			// Churn column
			if hasChurn {
				if f.churn > 0 {
					chClr := churnColor(f.churn)
					fmt.Fprintf(w, "  %s%5d%s", c(chClr), f.churn, c(codeReset))
				} else {
					fmt.Fprintf(w, "  %s%5s%s", c(codeDim), "·", c(codeReset))
				}
			}

			// Coverage bar
			if hasCoverage {
				if cd, ok := coverage[f.path]; ok {
					covClr := covBarColor(cd.pct)
					filledLen := cd.pct * covBarMax / 100
					covFilled := strings.Repeat("█", filledLen)
					covEmpty := strings.Repeat("░", covBarMax-filledLen)
					fmt.Fprintf(w, "  %s%3d%% %s%s%s%s%s",
						c(covClr), cd.pct, covFilled, c(codeReset), c(codeDim), covEmpty, c(codeReset))
				} else {
					// Pad to keep commit columns aligned
					fmt.Fprintf(w, "  %*s", 4+1+covBarMax, "")
				}
			}

			// Commit columns
			if hasCommits {
				fmt.Fprintf(w, "  ")
				for _, ch := range commits {
					if chMap, ok := commitData[ch]; ok {
						if count, ok := chMap[f.path]; ok && count > 0 {
							chClr := commitColor(count)
							fmt.Fprintf(w, "%s%6d%s", c(chClr), count, c(codeReset))
						} else {
							fmt.Fprintf(w, "%s%6s%s", c(codeDim), "·", c(codeReset))
						}
					} else {
						fmt.Fprintf(w, "%s%6s%s", c(codeDim), "·", c(codeReset))
					}
				}
			}

			fmt.Fprintln(w)
		}
	}

	// Phase 5: Summary
	fmt.Fprintln(w)
	fmt.Fprintln(w)

	totalLines := 0
	totalFiles := len(files)
	goLines := 0
	goCount := 0
	tsLines := 0
	tsCount := 0
	for _, f := range files {
		totalLines += f.lines
		if strings.HasSuffix(f.path, ".go") {
			goLines += f.lines
			goCount++
		} else if strings.HasSuffix(f.path, ".ts") {
			tsLines += f.lines
			tsCount++
		}
	}

	// Count unique directories
	dirSet := map[string]bool{}
	for _, f := range files {
		dirSet[filepath.Dir(f.path)] = true
	}
	dirCount := len(dirSet)

	fmt.Fprintf(w, "%sSummary%s\n", c(codeBold), c(codeReset))
	fmt.Fprintf(w, "  %-14s %s%6d%s lines   %3d files   %2d packages\n",
		"Go", c(codeYellow), goLines, c(codeReset), goCount, dirCount)
	fmt.Fprintf(w, "  %-14s %s%6d%s lines   %3d files\n",
		"TypeScript", c(codeYellow), tsLines, c(codeReset), tsCount)
	fmt.Fprintf(w, "  %s──────────────────────────────────────────%s\n", c(codeDim), c(codeReset))
	fmt.Fprintf(w, "  %s%-14s %s%6d%s lines   %3d files%s\n",
		c(codeBold), "Total", c(codeYellow), totalLines, c(codeReset), totalFiles, c(codeReset))

	if hasFn {
		totalFns := 0
		maxFnLen := 0
		bigFnCount := 0
		for _, f := range files {
			if f.fnData != nil {
				totalFns += f.fnData.count
				if f.fnData.maxLen > maxFnLen {
					maxFnLen = f.fnData.maxLen
				}
				if f.fnData.maxLen >= 60 {
					bigFnCount++
				}
			}
		}
		mClr := maxFnColor(maxFnLen)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  %sFunctions%s     %s%4d%s total    %3d files with maxFn >=60   longest: %s%d%s\n",
			c(codeBold), c(codeReset), c(codeYellow), totalFns, c(codeReset),
			bigFnCount, c(mClr), maxFnLen, c(codeReset))
	}

	if hasChurn {
		maxChurn := 0
		maxChurnFile := ""
		hotCount := 0
		for _, f := range files {
			if f.churn > maxChurn {
				maxChurn = f.churn
				maxChurnFile = filepath.Base(f.path)
			}
			if f.churn >= 26 {
				hotCount++
			}
		}
		chClr := churnColor(maxChurn)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  %sChurn%s          %s%4d%s max      %3d files with >=26 commits   hottest: %s%s%s\n",
			c(codeBold), c(codeReset), c(chClr), maxChurn, c(codeReset),
			hotCount, c(chClr), maxChurnFile, c(codeReset))
	}

	if hasDup {
		dupTotal := 0
		dupCount := 0
		highDupCount := 0
		for _, f := range files {
			if f.dupPct > 0 {
				dupTotal += f.dupPct
				dupCount++
			}
			if f.dupPct >= 25 {
				highDupCount++
			}
		}
		avgDup := 0
		if dupCount > 0 {
			avgDup = dupTotal / dupCount
		}
		dupClr := dupBarColor(avgDup)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  %sDuplication%s    %s%3d%%%s avg      %3d files with dup >=25%%\n",
			c(codeBold), c(codeReset), c(dupClr), avgDup, c(codeReset), highDupCount)
	}

	if hasCoverage {
		covTotal := 0
		covCovered := 0
		covFiles := 0
		for _, cd := range coverage {
			covTotal += cd.total
			covCovered += cd.covered
			covFiles++
		}
		covPct := 0
		if covTotal > 0 {
			covPct = covCovered * 100 / covTotal
		}
		covClr := covBarColor(covPct)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  %sCoverage%s       %s%3d%%%s          %3d files with data\n",
			c(codeBold), c(codeReset), c(covClr), covPct, c(codeReset), covFiles)
	}

	if hasCommits {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "%sRecent commits%s\n", c(codeBold), c(codeReset))
		for _, ch := range commits {
			short := ch
			if len(short) > 7 {
				short = short[:7]
			}
			subject := commitSubjects[ch]
			fmt.Fprintf(w, "  %s%s%s %s\n", c(codeDim), short, c(codeReset), subject)
		}
	}
}

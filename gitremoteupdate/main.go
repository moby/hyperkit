// Package gitremoteupdate provides a parser for output of `git remote update`.
package gitremoteupdate

import (
	"fmt"
	"strings"
)

// Operation that happened on a branch.
type Operation uint8

const (
	// New is a newly created branch.
	New Operation = iota

	// Updated is a branch that was updated.
	Updated

	// Deleted is a branch that was deleted.
	Deleted
)

// Change is a single entry in the result.
type Change struct {
	Op     Operation
	Branch string
}

// Result is the result of parsing output of `git remote update`.
type Result struct {
	Changes []Change
}

// Parse parses stderr output from running `git remote update`.
func Parse(stderr []byte) (Result, error) {
	var result Result

	lines := strings.Split(string(stderr), "\n")
	for _, line := range lines[1 : len(lines)-1] {
		change, err := parseLine(line)
		if err != nil {
			return result, err
		}
		result.Changes = append(result.Changes, change)
	}

	return result, nil
}

// parseLine parses a line like `   e8569f7..de0ad17  master     -> master`.
func parseLine(line string) (Change, error) {
	var change Change

	// Shortest valid input is len(`   d6d0813..e8569f7  m -> m`) = 27 characters.
	if len(line) < 27 {
		return change, fmt.Errorf("line too short")
	}

	// Parse operation.
	switch line[:3] {
	case " * ":
		change.Op = New
	case "   ":
		change.Op = Updated
	case " x ":
		change.Op = Deleted
	default:
		return change, fmt.Errorf("unsupported format")
	}

	// Parse branch name.
	branch, err := parseBranchArrowBranch(line[21:])
	if err != nil {
		return change, fmt.Errorf("failed to parse branch name")
	}
	change.Branch = branch

	return change, nil
}

// parseBranchArrowBranch parses a `master     -> master` segment to extract
// relevant branch name.
func parseBranchArrowBranch(bab string) (branch string, err error) {
	branches := strings.SplitN(bab, " -> ", 2)
	if len(branches) != 2 {
		return "", fmt.Errorf("failed to parse `branch -> branch` segment")
	}
	// Note, if we wanted to use branches[0], we should trim whitespace on its right.
	// Return second branch name since it's always valid.
	return branches[1], nil
}

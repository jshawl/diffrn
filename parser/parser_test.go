package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	expected := "file"
	actual, _ := parseFile("+++ b/file")
	if expected != actual.Name {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual.Name))
	}
}

func TestParseFileWithExtension(t *testing.T) {
	expected := "file.txt"
	actual, _ := parseFile("+++ b/file.txt")
	if expected != actual.Name {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual.Name))
	}
}

func TestParseFileWithSlashes(t *testing.T) {
	expected := "folder/file.txt"
	actual, _ := parseFile("+++ b/folder/file.txt")
	if expected != actual.Name {
		t.Fatalf(fmt.Sprintf("Expected: %s Actual: %s", expected, actual.Name))
	}
}

func TestParseFileOnNonFile(t *testing.T) {
	_, err := parseFile("not a line with a filename")
	if err == nil {
		t.Fatalf("Expected error, got match")
	}
}

func TestParseBlockWithAddedLines(t *testing.T) {
	actual, _ := parseBlock("@@ -0,0 +1 @@")
	if actual.OldStart != 0 {
		t.Fatalf("Expected OldStart to be 1")
	}
	if actual.OldEnd != 0 {
		t.Fatalf("Expected OldEnd to be 0")
	}
	if actual.NewStart != 1 {
		t.Fatalf("Expected NewStart to be 1")
	}
	if actual.NewEnd != 1 {
		t.Fatalf("Expected NewEnd to be 1")
	}
}

func TestParseBlockWithChangedLines(t *testing.T) {
	actual, _ := parseBlock("@@ -1,4 +1,4 @@")
	if actual.OldStart != 1 {
		t.Fatalf("Expected OldStart to be 1")
	}
	if actual.OldEnd != 5 {
		t.Fatalf("Expected OldEnd to be 5")
	}
	if actual.NewStart != 1 {
		t.Fatalf("Expected NewStart to be 1")
	}
	if actual.NewEnd != 5 {
		t.Fatalf("Expected NewEnd to be 5")
	}
}

func TestParseLineDiffPreamble(t *testing.T) {
	actual, _ := parseLine("diff --git a/file b/file")
	if actual != "" {
		t.Fatalf("Expected diff line to be ignored")
	}
}

func TestParseLineIndexPreamble(t *testing.T) {
	actual, _ := parseLine("index e69de29..d00491f 100644")
	if actual != "" {
		t.Fatalf("Expected index line to be ignored")
	}
}

func TestParseLineOldPreamble(t *testing.T) {
	actual, _ := parseLine("--- a/file")
	if actual != "" {
		t.Fatalf("Expected old line to be ignored")
	}
}

func TestParseLineBlockPreamble(t *testing.T) {
	actual, _ := parseLine("@@ -0,0 +1 @@")
	if actual != "" {
		t.Fatalf("Expected block line to be ignored")
	}
}

func TestParseLineDiff(t *testing.T) {
	actual, _ := parseLine("- removed this line")
	if actual != "- removed this line" {
		t.Fatalf("Expected line to be parsed")
	}
}

func TestParseLineDiffEmptyLine(t *testing.T) {
	_, err := parseLine("")
	if err != nil {
		t.Fatalf("Expected line to be parsed")
	}
}

func TestParseLineNewPreamble(t *testing.T) {
	actual, _ := parseLine("+++ a/file")
	if actual != "" {
		t.Fatalf("Expected new line to be ignored")
	}
}

func TestParseDiff(t *testing.T) {
	contents, _ := os.ReadFile("./test/one-file-one-block.diff")
	actual, _ := ParseDiff(string(contents))

	if len(actual.Files) != 1 {
		t.Fatalf("Expected 1 file")
	}
	if actual.Files[0].Name != "file" {
		t.Fatalf("Expected 1 file name to be 'file'")
	}
	if len(actual.Files[0].Blocks[0].Lines) != 1 {
		t.Fatalf("Expected 1 File, 1 Block, 1 Line")
	}
}

func TestParseDiffOneFileTwoBlocks(t *testing.T) {
	contents, _ := os.ReadFile("./test/one-file-two-blocks.diff")
	actual, _ := ParseDiff(string(contents))
	if len(actual.Files) != 1 {
		t.Fatalf("Expected 1 file")
	}
	if len(actual.Files[0].Blocks) != 2 {
		t.Fatalf("Expected 1 File, 2 Blocks")
	}
}

func TestParseDiffTwoFilesTwoBlocks(t *testing.T) {
	contents, _ := os.ReadFile("./test/two-files-two-blocks.diff")
	actual, _ := ParseDiff(string(contents))
	if len(actual.Files) != 2 {
		t.Fatalf("Expected 2 files")
	}
	if len(actual.Files[0].Blocks) != 2 {
		t.Fatalf("Expected 2 Files, 2 Blocks")
	}
}

func TestParseDiffNewFile(t *testing.T) {
	contents, _ := os.ReadFile("./test/new-file.diff")
	actual, _ := ParseDiff(string(contents))
	if len(actual.Files) != 1 {
		t.Fatalf("Expected 1 file")
	}
	if actual.Files[0].Name != ".gitignore" {
		t.Fatalf("Expected filename .gitignore")
	}
}

func TestParseDiffDeletedFile(t *testing.T) {
	contents, _ := os.ReadFile("./test/deleted-file.diff")
	actual, _ := ParseDiff(string(contents))
	if len(actual.Files) != 1 {
		t.Fatalf("Expected 1 file")
	}
}

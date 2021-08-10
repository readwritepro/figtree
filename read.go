//=============================================================================
// File:     reader.go
// Contents: ReadConfig scans a user config file, and merges it with any baseline
//            file referenced by a !baseline pragma.
//           ReadFigtree scans a file that contains figtree syntax, and any files
//            embedded via an include pragma.
//           ParseBranch recursively scans figtree syntax creating an in-memory
//            tree of branches and items.
//=============================================================================

package figtree

import (
	"bufio"
	"os"
	"path"
	"strings"

	eh "github.com/readwritepro/error-handler"
)

// The ReadConfig function reads a user's configuration file into memory, honoring any baseline pragma it may contain.
//
// Returns the root branch of the tree created by merging the user's file with any baseline file it may point to.
//
// May return ErrEndOfBranch if the parser prematurely stopped, before parsing
// the entire file, due to a misconfigured closing brace.
func ReadConfig(inFilename string) (*Branch, error) {

	gBaselineTree = nil // reset the global baselineTree

	userTree, err := ReadFigtree(inFilename, UserFile)
	if err != nil {
		return nil, err
	}

	// Reading the user's file may have triggered the creation of a baseline tree via the !baseline pragma
	// Now that the user's tree and the baseline tree are both fully parsed and in memory, merge them.
	mergedBranch := mergeBaselineWithUser(gBaselineTree, userTree)
	return mergedBranch, nil
}

// The ReadFigtree function opens, parses, and closes the given file.
// This is a public function, and may be called to read any file containing
// figtree syntax, but it is rarely used publicly. See the ReadConfig function for that.
//
// The returned Branch is the root of the configuration tree which is used
// in subsequent calls to access and alter the tree's inner branches and items.
//
// Returns ErrEndOfBranch if the parser prematurely stopped, before parsing
// the entire file, due to a misconfigured closing brace.
func ReadFigtree(inFilename string, fileOrigin FileOrigin) (*Branch, error) {
	inFile, err := os.Open(inFilename)
	if eh.Invalid(err) {
		return nil, err
	}
	defer inFile.Close()

	// create a scanner that uses the "ScanLines" splitter
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	root := NewBranch()
	srcLine := 0
	err = root.ParseBranch(scanner, inFilename, &srcLine, fileOrigin)
	if err == ErrEndOfBranch {
		return nil, err
	}
	if err != ErrEOF {
		return nil, err
	}

	return root, nil
}

// Recursive function to read lines via a bufio scanner, adding
// key/value pairs and inner branches to the current branch.
// This function is typically only called by the ReadFigtree function,
// but it may safely be called in userland in order to graft one branch onto another.
//
// Returns the ErrEndOfBranch sentinal when finished parsing each inner branch.
// Return ErrEOF to the outermost caller.
func (branch *Branch) ParseBranch(scanner *bufio.Scanner, srcFile string, srcLine *int, srcOrigin FileOrigin) error {

	blockComments := make([]string, 0) // block comment accumulator

	for scanner.Scan() { // advance the scanner to the end of line
		var leftSide, rightSide string
		var key, val, terminalWhitespace, terminalComment string

		// copy the runes into a string and remove leading and trailing whitespace
		line := scanner.Text()
		line = strings.Trim(line, " \t")
		*srcLine++

		// send blank lines and comment lines to the block comment accumulator
		if len(line) == 0 || line[0] == '#' {
			blockComments = append(blockComments, line)
			continue
		}

		// split into two halves based on first whitespace or opening-brace
		whitespace := strings.IndexAny(line, " \t{")
		if whitespace == -1 {
			leftSide = line
			rightSide = ""
			key = leftSide
		} else {
			leftSide = line[:whitespace]
			rightSide = line[whitespace:]
			key = leftSide
		}

		// split right side into value and possible comment
		// hashtags that are not preceded by whitespace are treated as part
		// of the value, to allow for things like URLs with bookmarks
		hash := strings.Index(rightSide, "\t#")
		if hash == -1 {
			hash = strings.Index(rightSide, " #")
		}
		if hash == -1 {
			val = strings.Trim(rightSide, " \t")
			terminalWhitespace = ""
			terminalComment = ""
		} else {
			val = strings.TrimLeft(rightSide[:hash+1], " \t")         // keep the trailing whitespace for next step
			terminalComment = strings.Trim(rightSide[hash+2:], " \t") // drop the whitespace and hash

			// extract the whitespace between value and hash
			pos := -1
			for i := len(val) - 1; i >= 0; i-- {
				if val[i] != ' ' && val[i] != '\t' {
					pos = i
					break
				}
			}
			if pos != -1 {
				terminalWhitespace = val[pos+1:]
				val = val[:pos+1]
			}
		}

		// if the right-hand side is "{" create a branch and recurse
		if len(val) == 1 && val[0] == '{' {
			// begin branch
			err := branch.handleBranch(scanner, key, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
			if err != ErrEndOfBranch {
				return err
			}
		} else if len(leftSide) > 0 && leftSide[0] == '}' {
			// end branch
			return ErrEndOfBranch
		} else {
			// typical key/value
			err := branch.handleKeyValuePair(key, val, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
			if err != nil {
				return err
			}
		}

		// reset the block comment accumulator
		blockComments = make([]string, 0)
	}
	return ErrEOF
}

// Helper function used by ParseBranch to handle the beginning of a branch
// by recursively calling ParseBranch.
//
// The normal return is the sentinal ErrEndOfBranch, anything else should halt further processing
func (branch *Branch) handleBranch(scanner *bufio.Scanner, key string, blockComments []string, terminalWhitespace string, terminalComment string, srcFile string, srcLine *int, srcOrigin FileOrigin) error {
	innerBranch := NewBranch()
	branch.appendItem(key, innerBranch, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
	return innerBranch.ParseBranch(scanner, srcFile, srcLine, srcOrigin)
}

// Helper function used by ParseBranch to handle typical key/value pairs
// with special detection for the !include, !baseline, and !dtd pragmas.
func (branch *Branch) handleKeyValuePair(key string, value string, blockComments []string, terminalWhitespace string, terminalComment string, srcFile string, srcLine *int, srcOrigin FileOrigin) error {

	if strings.Index(key, "!include") == 0 {
		branch.appendItem("!include", value, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
		err := branch.readIncludeFile(value)
		if err != nil {
			return err
		}
	} else if strings.Index(key, "!baseline") == 0 {
		branch.appendItem("!baseline", value, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
		err := branch.readBaselineFile(value)
		if err != nil {
			return err
		}
	} else if strings.Index(key, "!dtd") == 0 {
		branch.appendItem("!dtd", value, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
		dtdRootBranch, err := branch.readDtdFile(value)
		if err != nil {
			return err
		}
		todo(dtdRootBranch)
	} else {
		branch.appendItem(key, value, blockComments, terminalWhitespace, terminalComment, srcFile, srcLine, srcOrigin)
	}
	return nil
}

// Special processing for including key/values from another file.
// When the filename is not an absolute path, prepend the current working directory.
func (branch *Branch) readIncludeFile(localFilename string) error {
	if len(localFilename) > 0 && localFilename[0] != '/' {
		cwd, _ := os.Getwd()
		localFilename = path.Join(cwd, localFilename)
	}
	includeBranch, err := ReadFigtree(localFilename, IncludeFile)
	if err != nil {
		return err
	}
	branch.Items = append(branch.Items, includeBranch.Items...)
	return nil
}

// Special processing for adding a default set of fallback key/values from a baseline file.
// When the filename is not an absolute path, prepend the current working directory.
func (branch *Branch) readBaselineFile(localFilename string) error {
	if len(localFilename) > 0 && localFilename[0] != '/' {
		cwd, _ := os.Getwd()
		localFilename = path.Join(cwd, localFilename)
	}
	var err error
	gBaselineTree, err = ReadFigtree(localFilename, BaselineFile)
	if err != nil {
		return err
	}
	//branch.Items = append(branch.Items, includeBranch.Items...)
	return nil
}

// Special processing for parsing a declared document type definition file,
// which may be used for validation.
// When the filename is not an absolute path, prepend the current working directory.
//
// Returns the dtd root branch, which should not become part of the user's actual figtree
func (branch *Branch) readDtdFile(localFilename string) (*Branch, error) {
	if len(localFilename) > 0 && localFilename[0] != '/' {
		cwd, _ := os.Getwd()
		localFilename = path.Join(cwd, localFilename)
	}
	dtdRootBranch, err := ReadFigtree(localFilename, DtdFile)
	if err != nil {
		return nil, err
	}
	return dtdRootBranch, nil
}

//=============================================================================
// File:     writers.go
// Contents: WriteToFile
//           WriteToBuffer
//           serializeConfig
//             WriteFigtree
//             WriteInternal
//             WriteJson
//             WriteYaml
//=============================================================================

package figtree

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	eh "github.com/readwritepro/error-handler"
)

// The WriteConfig interface is any type that implements serializeConfig allowing
// a figtree to be sent to either a file or a string.
type WriteConfig interface {
	serializeConfig(branch *Branch, w *bufio.Writer, depth int) error
}

// The WriteFigtree type is used with WriteToFile and WriteToBuffer to
// serialize a configuration using native figtree syntax.
type WriteFigtree struct{}

// The WriteInternal type is used with WriteToFile and WriteToBuffer to
// serialize a configuration with internal parsing and debugging information.
type WriteInternal struct{}

// The WriteJson type is used with WriteToFile and WriteToBuffer to
// serialize a configuration using JSON syntax.
type WriteJson struct{}

// The WriteYaml type is used with WriteToFile and WriteToBuffer to
// serialize a configuration using YAML syntax.
type WriteYaml struct{}

//-----------------------------------------------------------------------------
// Public write functions
//-----------------------------------------------------------------------------

// Write the contents of a configuration tree to the specified file.
// This function is typically only called on the root branch;
// nevertheless, it may safely  be called on any branch, regardless of
// whether it is a root, thus allowing a pruned portion of a tree to be serialized.
// The path to the output directory must already exist and have write permission.
// The WriteConfig interface is any type that implements serializeConfig.
//
// Example:
//  root, _ := ReadConfig("testdata/fixtures/sample")
//  wf := WriteFigtree{}
//  root.WriteToFile(wf, "testdata/actual/sample")
//
func (branch *Branch) WriteToFile(wc WriteConfig, outFilename string) error {
	fOut, err := os.Create(outFilename)
	if eh.Invalid(err, outFilename) {
		return err
	}
	defer fOut.Close()

	writer := bufio.NewWriter(fOut)
	err = wc.serializeConfig(branch, writer, 0)
	if eh.Invalid(err) {
		return err
	}
	return writer.Flush()
}

// Write the contents of a configuration tree to a buffer and return the buffer's string value.
//
// The WriteConfig interface is any type that implements serializeConfig.
func (branch *Branch) WriteToBuffer(wc WriteConfig) (string, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	err := wc.serializeConfig(branch, writer, 0)
	if err != nil {
		return "", err
	}

	writer.Flush()
	return buf.String(), nil
}

//-----------------------------------------------------------------------------
// Write Figtree
//-----------------------------------------------------------------------------

// Recursive function to write the current branch, in figtree syntax, to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
// This function is typically only called by WriteToFile or WriteToBuffer.
func (wf WriteFigtree) serializeConfig(branch *Branch, w *bufio.Writer, depth int) error {
	prefix := strings.Repeat("\t", depth)
	var err error

	for _, item := range branch.Items {
		key := item.key

		// write any blank lines or block comments
		for _, bc := range item.blockComments {
			_, err = fmt.Fprintln(w, prefix+bc)
			if err != nil {
				return err
			}
		}

		var wsComment string
		if item.terminalComment != "" {
			wsComment = fmt.Sprintf("%s# %s", item.terminalWhitespace, item.terminalComment)
		}

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s%s %s%s\n", prefix, key, value, wsComment)
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			_, err = fmt.Fprintf(w, "%s%s {%s\n", prefix, key, wsComment)
			if err != nil {
				return err
			}
			_ = wf.serializeConfig(value, w, depth+1)
			_, err = fmt.Fprintf(w, "%s}\n", prefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
// Write Internal
//-----------------------------------------------------------------------------

// Recursive function to write the current branch with internal debug info to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
func (wi WriteInternal) serializeConfig(branch *Branch, w *bufio.Writer, depth int) error {
	prefix := strings.Repeat("\t", depth)
	var err error

	for _, item := range branch.Items {
		key := item.key

		srcContext := fmt.Sprintf("(%v)[%s:%d]", item.srcOrigin, filepath.Base(item.srcFile), item.srcLine)
		srcContext = fmt.Sprintf("%-32s%s ", srcContext, prefix)

		// write any blank lines or block comments
		for _, bc := range item.blockComments {
			_, err = fmt.Fprintf(w, "%s%s\n", srcContext, bc)
			if err != nil {
				return err
			}
		}

		var wsComment string
		if item.terminalComment != "" {
			wsComment = fmt.Sprintf("%s# %s", item.terminalWhitespace, item.terminalComment)
		}

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s%s %s%s\n", srcContext, key, value, wsComment)
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			_, err = fmt.Fprintf(w, "%s%s {%s\n", srcContext, key, wsComment)
			if err != nil {
				return err
			}

			// fmt.Fprintf(w, "%s(branchPath: %s)\n", prefix, value.branchPath)

			_ = wi.serializeConfig(value, w, depth+1)
			_, err = fmt.Fprintf(w, "%s}\n", srcContext)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//-----------------------------------------------------------------------------
// Write JSON
//-----------------------------------------------------------------------------

// Function to write the current branch, in JSON syntax, to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
// This function is typically only called by WriteToFile or WriteToBuffer.
func (wj WriteJson) serializeConfig(branch *Branch, w *bufio.Writer, depth int) error {
	fmt.Fprintf(w, "{")
	err := wj.serializeBranch(branch, w, depth+1)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\n}")
	return nil
}

// Private recursive function to write the current branch, in JSON syntax, to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
// This function is typically only called by WriteToFile or WriteToBuffer.
func (wj WriteJson) serializeBranch(branch *Branch, w *bufio.Writer, depth int) error {
	prefix := strings.Repeat("\t", depth)
	var err error

	// keep track of all keys which have been determined to be arrays
	arrayKeys := map[string]bool{}

	commaLF := "\n" // begin first item without a comma

	for _, item := range branch.Items {
		key := item.key

		// has this keyName already been handled as an array
		if _, exists := arrayKeys[key]; exists {
			continue
		}

		// check to see if the keyname ends in [], if so, treat it as an array even if it is empty or has only one entry
		bracketPos := strings.Index(key, "[]")
		bIsArray := false
		if bracketPos != -1 && bracketPos == len(key)-2 {
			bIsArray = true
		}

		// check to see if this key occurs more than once, if so, treat it as an array
		allItems := branch.QueryAll(key)
		if bIsArray || len(allItems) > 1 {
			fmt.Fprintf(w, "%s", commaLF)
			err = wj.serializeArray(key, allItems, branch, w, depth)
			if err != nil {
				return err
			}
			commaLF = ",\n"
			arrayKeys[key] = true
			continue
		}

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s%s\"%s\": %s", commaLF, prefix, escapeJsonKey(key), escapeJsonValue(value))
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			_, err = fmt.Fprintf(w, "%s%s\"%s\": {", commaLF, prefix, escapeJsonKey(key))
			if err != nil {
				return err
			}
			_ = wj.serializeBranch(value, w, depth+1)
			_, err = fmt.Fprintf(w, "\n%s}", prefix)
			if err != nil {
				return err
			}
		}
		commaLF = ",\n" // second and subsequent lines
	}

	return nil
}

// keyName may end in trailing brackets []
// allItems may be a collection of 0, 1 or more items
func (wj WriteJson) serializeArray(keyName string, allItems []Item, branch *Branch, w *bufio.Writer, depth int) error {

	prefix0 := strings.Repeat("\t", depth)
	prefix1 := strings.Repeat("\t", depth+1)
	var err error

	// strip off the JSON-only "trailing bracket hack", leaving only the name
	bracketPos := strings.Index(keyName, "[]")
	if bracketPos != -1 && bracketPos == len(keyName)-2 {
		keyName = keyName[:bracketPos]
	}

	// write the opening bracket "["
	_, err = fmt.Fprintf(w, "%s\"%s\": [\n", prefix0, escapeJsonKey(keyName))
	if err != nil {
		return err
	}

	// an empty JSON array will have one item with no value
	// remove the zombie item from the slice
	if len(allItems) == 1 {
		if firstItemValue, _ := allItems[0].Value(); firstItemValue == "" {
			allItems = []Item{}
		}
	}

	commaLF := "" // begin first item without a comma

	// print each item's value on a separate line
	for _, item := range allItems {

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s%s%s", commaLF, prefix1, escapeJsonValue(value))
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			_, err = fmt.Fprintf(w, "%s%s{", commaLF, prefix1)
			if err != nil {
				return err
			}
			_ = wj.serializeBranch(value, w, depth+2)
			_, err = fmt.Fprintf(w, "\n%s}", prefix1)
			if err != nil {
				return err
			}
		}
		commaLF = ",\n" // second and subsequent lines
	}

	// write the closing bracket "]"
	_, err = fmt.Fprintf(w, "\n%s]", prefix0)
	if err != nil {
		return err
	}
	return nil
}

// \b   U+0008 backspace
// \f   U+000C form feed
// \n   U+000A line feed or newline
// \r   U+000D carriage return
// \t   U+0009 horizontal tab
// \v   U+000B vertical tab
// \"   U+0022 double quote
// \\   U+005C reverse solidus
func escapeJsonKey(unescaped string) string {
	r := strings.NewReplacer(
		"\u0008", "\b",
		"\u000C", "\f",
		"\u000A", "\n",
		"\u000D", "\r",
		"\u0009", "\t",
		"\u000B", "\v",
		"\u0022", "\u005C\u0022",
		"\u005C", "\u005C\u005C")
	return r.Replace(unescaped)
}

func escapeJsonValue(unescaped string) string {
	// special recognition and treatment for nulls
	if unescaped == "" || unescaped == "null" {
		return "null"
	}
	// special recognition and treatment for booleans
	if unescaped == "true" || unescaped == "false" {
		return unescaped
	}
	// make sure numbers don't get quote delimiters
	if _, err := strconv.ParseFloat(unescaped, 64); err == nil {
		return unescaped
	}
	// anything else follows the normal escaping rules for strings
	return "\u0022" + escapeJsonKey(unescaped) + "\u0022"
}

//-----------------------------------------------------------------------------
// Write YAML
//-----------------------------------------------------------------------------

// Function to write the current branch, in YAML syntax, to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
// This function is typically only called by WriteToFile or WriteToBuffer.
func (wy WriteYaml) serializeConfig(branch *Branch, w *bufio.Writer, depth int) error {
	fmt.Fprintf(w, "---\n")
	err := wy.serializeBranch(branch, w, depth)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\n")
	return nil
}

// Private recursive function to write the current branch, in JSON syntax, to the specified
// bufio writer.
// The depth parameter specifies how many tab characters to indent each line.
// This function is typically only called by WriteToFile or WriteToBuffer.
func (wy WriteYaml) serializeBranch(branch *Branch, w *bufio.Writer, depth int) error {
	prefix := strings.Repeat("  ", depth)
	var err error

	// keep track of all keys which have been determined to be arrays
	arrayKeys := map[string]bool{}

	for _, item := range branch.Items {
		key := item.key

		// has this keyName already been handled as an array
		if _, exists := arrayKeys[key]; exists {
			continue
		}

		// write any blank lines or block comments
		for _, bc := range item.blockComments {
			_, err = fmt.Fprintln(w, prefix+bc)
			if err != nil {
				return err
			}
		}

		var wsComment string
		if item.terminalComment != "" {
			wsComment = fmt.Sprintf("%s# %s", item.terminalWhitespace, item.terminalComment)
		}

		// check to see if the keyname ends in [], if so, treat it as an array even if it is empty or has only one entry
		bracketPos := strings.Index(key, "[]")
		bIsArray := false
		if bracketPos != -1 && bracketPos == len(key)-2 {
			bIsArray = true
		}

		// check to see if this key occurs more than once, if so, treat it as an array
		allItems := branch.QueryAll(key)
		if bIsArray || len(allItems) > 1 {
			err = wy.serializeArray(key, allItems, branch, w, depth)
			if err != nil {
				return err
			}
			arrayKeys[key] = true
			continue
		}

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s%s: %s%s\n", prefix, escapeYaml(key), escapeYaml(value), wsComment)
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			_, err = fmt.Fprintf(w, "%s%s:%s\n", prefix, escapeYaml(key), wsComment)
			if err != nil {
				return err
			}
			_ = wy.serializeBranch(value, w, depth+1)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// keyName may end in trailing brackets []
// allItems may be a collection of 0, 1 or more items
func (wy WriteYaml) serializeArray(keyName string, allItems []Item, branch *Branch, w *bufio.Writer, depth int) error {

	prefix0 := strings.Repeat("  ", depth)
	prefix1 := strings.Repeat("  ", depth+1)
	var err error

	// strip off the JSON-only "trailing bracket hack", leaving only the name
	bracketPos := strings.Index(keyName, "[]")
	if bracketPos != -1 && bracketPos == len(keyName)-2 {
		keyName = keyName[:bracketPos]
	}

	// an empty JSON array will have one item with no value
	if len(allItems) == 1 {
		if firstItemValue, _ := allItems[0].Value(); firstItemValue == "" {
			var wsComment string
			if allItems[0].terminalComment != "" {
				wsComment = fmt.Sprintf(" %s# %s", allItems[0].terminalWhitespace, allItems[0].terminalComment)
			}
			fmt.Fprintf(w, "%s%s: []%s\n", prefix0, escapeYaml(keyName), wsComment)
			return nil
		}
	}

	// write the opening keyName
	_, err = fmt.Fprintf(w, "%s%s:\n", prefix0, escapeYaml(keyName))
	if err != nil {
		return err
	}

	// print each item's value on a separate line
	for _, item := range allItems {

		var wsComment string
		if item.terminalComment != "" {
			wsComment = fmt.Sprintf("%s# %s", item.terminalWhitespace, item.terminalComment)
		}

		switch value := item.value.(type) {
		// simple key/value pair
		case string:
			_, err = fmt.Fprintf(w, "%s- %s%s\n", prefix1, escapeYaml(value), wsComment)
			if err != nil {
				return err
			}
		// nested branch
		case *Branch:
			fmt.Fprintf(w, "%s-\n", prefix1)
			_ = wy.serializeBranch(value, w, depth+2)
		}
	}

	return nil
}

// \b   U+0008 backspace
// \f   U+000C form feed
// \n   U+000A line feed or newline
// \r   U+000D carriage return
// \t   U+0009 horizontal tab
// \v   U+000B vertical tab
// \"   U+0022 double quote
// \'   U+0027 apostrophe
// \\   U+005C reverse solidus
func escapeYaml(unescaped string) string {
	if len(unescaped) == 0 {
		return "null "
	}

	// make sure numbers don't get quote delimiters
	if _, err := strconv.ParseFloat(unescaped, 64); err == nil {
		return unescaped
	}

	bNeedsDelimiter := false
	if strings.ContainsAny(unescaped, "-?:,[]{}#&*!|>`") {
		bNeedsDelimiter = true
	}

	r := strings.NewReplacer(
		"\u0008", "\b",
		"\u000C", "\f",
		"\u000A", "\n",
		"\u000D", "\r",
		"\u0009", "\t",
		"\u000B", "\v",
		"\u0022", "\u005C\u0022",
		"\u0027", "\u005C\u0027",
		"\u005C", "\u005C\u005C")
	escaped := r.Replace(unescaped)

	if unescaped != escaped || bNeedsDelimiter {
		escaped = "\u0022" + escaped + "\u0022"
	}

	return escaped
}

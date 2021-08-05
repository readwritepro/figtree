//=============================================================================
// File:     file-origin.go
// Contents: FileOrigin enum declaration
//=============================================================================

package figtree

// The FileOrigin type is used during parsing to distinguish whether the
// file being processed is: the user's configuration file, or a file being included,
// or a file containing baseline default fallback values, or a document type definition.
type FileOrigin int

const (
	UserFile     FileOrigin = iota // the user's configuration file specified directly in a call to Read
	IncludeFile                    // file containing a portion of the configuration tree, referenced with a !include pragma
	BaselineFile                   // file containing fallback defaults, referenced with a !baseline pragma
	DtdFile                        // a document type definition file, refereneced with a !dtd pragma
)

func (fileOrigin FileOrigin) String() string {
	return [...]string{"User", "Incl", "Base", "Dtd"}[fileOrigin]
}

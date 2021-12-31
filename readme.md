Package figtree provides a multi-paradigm SDK for sophisticated configuration
file access.

# Motivation

Figtree syntax is based on classic key/value pairs that may be organized
into a nested hierarchy of named sections.

Many of the design goals for the figtree syntax come from its predecessors
including XML, JSON, YAML, TOML, win.ini, and Apache.  The deficiencies that
figtree syntax addresses comprise:

 1. XML is verbose, requiring opening and closing tags that obscure the data.
 2. JSON does not allow comments, forcing usage documentation to be written in a separate place.
 3. YAML uses indentation to indicate hierarchical structure, which is not as appealing to some users as curly braces.
 4. TOML, JSON and YAML require excessive use of quotation-marks, colons, equal-signs, hyphens, or commas purely for the sake of the parsing algorithm.
 5. TOML, JSON and YAML require special notation to declare multi-value arrays.
 6. Some configuration file editors arbitrarily re-sort keys, loosing the original author's intended ordering.
 7. Some configuration file editors arbitrarily drop comments from file.

The features that figtree syntax incorporates comprise:
 1. Like TOML and YAML, keynames do not need to be delimited by quotation-marks.
 2. Like YAML, values do not need to be delimited by quotation-marks.
 3. Like XML's XPath, the configuration hierarchy can be accessed using keyPath notation.
 4. Like XML, the configuration can be validated against a document type definition.
 5. Like win.ini, and unlike the Windows system registry, the configuration is held in plain text files.
 6. Like Apache configuration files, portions of the configuration can be included from separate files.
 7. Like many Linux-based tools, default baseline settings can be kept separate from user overrides.

# Figtree syntax

There are only a few essential constructs to understand. First, simple key/value
pairs appear on a single line. Unlike other configuration syntax styles,
key/value pairs do not use colons or equal-signs as assignment operators.
Instead, the presence of whitespace between the key-name and the beginning
of the value is enough to parse the line into left-hand and right-hand halves.

Furthermore, values are not delimited by quotation marks. Instead, all leading
and trailing whitespace characters are stripped from the right-hand side of
the key/value pair to produce the value. Examples:


     hostname    figtree
     url         https:www.figtree.io
     version     1.0
     note        A fig tree should not be confused with a "figtree"

The second construct to understand are named sections, which are multi-line
collections of key/value pairs. Named sections have a key name followed by
a K&R-style pair of curly braces.

The opening brace must be the last non-whitespace character of a line. The
closing brace must be the first non-whitespace character of a line. The lines
between the opening and closing braces may contain both simple key/value pairs
and other named sections. Sections may be nested arbitrarily deep. Example:


     ip-settings {
         ip4 {
             inet       179.100.102.215
             netmask    255.255.240.0
             broadcast  167.111.255
             gateway    167.99.96.0
         }
         ip6 {
             inet       1fe80::e457:b8ff:fe48:8ead
         }
     }

Duplicate, non-unique keys are used to declare array-like values. Any key that
appears more than once per section is considered to be an array. The ItemIsArray
function can be used to test whether or not a key has multiple values. The QueryAll
function can be used to obtain multiple items with the same key. Example:


     name-servers {
         ns ns1.figtree.net
         ns ns2.figtree.net
         ns ns3.figtree.net
     }

Block comments are written using hashtags as the first non-whitespace character
of a line. Example:

     # All ip-settings are required

Terminal comments may appear on the same line together with key/value pairs when
whitespace and a hash tag follow the value. Example:

 

     hostname    figtree                   # the short device name
     url         https:www.figtree.io      # protocol, optional port, and DNS name
     version     1.0                       # minimum required version
     reference   allthedocs.org#figtree    # a hashtag embedded within the value (with no preceding space)
    
# Reading and writing

Configuration files written in figtree syntax can be read into memory using the ReadConfig function. Example:

 

    configFilename := "testdata/fixtures/sample"
     root, err := ReadConfig(configFilename)
     if err != nil {
         return
     }
     # root is a pointer to the in-memory hierarchical tree
    
An in-memory tree can be saved using any type that implements the SerializeBranch function,
which is called by WriteToFile and WriteToBuffer. Example:



     writer := FigtreeWriter{}
     configFilename := "testdata/fixtures/sample"
     err := root.WriteToFile(writer, configFilename)
     if err != nil {
         t.Errorf(err.Error())
     }


# Accessing figtree items

Programs may access individual items using either simpleKeyNames or keyPaths.
A simpleKeyName is a key as it appears in the figtree file. A keyPath is a
slash-separated sequence of key names where each segment matches successively
deeper sections of the hierarchy. Examples:



     var item *Item
     var err error
     item, err = QueryOne("hostname")               simpleKeyName
     item, err = QueryOne("ip-settings/ip4/inet")   keyPath

Arrays can be accessed similarly. Example:



     var items []Items
     items = QueryAll("name-servers/ns")
    
An item's value is an anonymous interface which may be a string value or a Branch pointer.
When a program knows the type of item that it is looking for it may be simpler to use
either the GetValue or GetBranch function to directly access it. When unsure, a program
may use GetItem at the cost of having to explicitly work around the ambiguity.

Testing for the existence of a simpleKeyName is done with ItemExists. Testing for the
existence of a keyPath is done with PathExists. Checking to see if a key has multiple values
is done with ItemIsArray.

# Manipulating the figtree

Figtree items and branches can be programmatically manipulated. Adding new keys
is done with AppendItem, PrependItem, InsertBeforeItem, or InsertAfterItem. Removing keys is done with RemoveItem.

package main

import (
	"fmt"

	"github.com/pedrogao/plib/pkg/jj"
)

/*
 simple json parse library implement by JIT
*/

const source = `{
  "glossary": {
    "title": "example glossary",
    "GlossDiv": {
      "title": "S",
      "GlossList": {
        "GlossEntry": {
          "ID": "SGML",
          "SortAs": "SGML",
          "GlossTerm": "Standard Generalized Markup Language",
          "Acronym": "SGML",
          "Abbrev": "ISO 8879:1986",
          "GlossDef": {
            "para": "A meta-markup language, used to create markup languages such as DocBook.",
            "GlossSeeAlso": [
              "GML",
              "XML"
            ]
          },
          "GlossSee": "markup"
        }
      }
    }
  }
}
`

func main() {
	node, err := jj.UnmarshalAst(source)
	if err != nil {
		panic(err)
	}

	ret := jj.Marshal(node, "")

	fmt.Println(ret)
}

# markdown
Markdown library

## Markdown to html

```
package main

import (
        "fmt"

        markdown "github.com/mokelab-go/markdown/html"
)

const src = `
# Hello markdown

This library outputs html from
markdown.

 * u1
 * u2
 * [u3](https://mokelab.com)
`


func main() {
        m := markdown.NewMarkdown()
        out, err := m.Compile(src)
        if err != nil {
                fmt.Errorf("Error :%s", err)
                return
        }
        fmt.Printf(out)
}
```

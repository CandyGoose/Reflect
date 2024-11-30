package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailru/easyjson"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	seenBrowsers := make(map[string]bool, 128)

	foundUsers := strings.Builder{}
	foundUsers.Grow(8192)

	user := &User{}

	reader := bufio.NewReader(file)

LOOP:
	for i := 0; ; i++ {
		line, prefix, err := reader.ReadLine()
		if prefix {
			panic(errors.New("failed to read line, line too long"))
		}
		if err != nil {
			switch err {
			case io.EOF:
				break LOOP
			default:
				panic(err)
			}
		}

		err = easyjson.Unmarshal(line, user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if containsAndroid := strings.Contains(browser, "Android"); containsAndroid {
				isAndroid = true
				seenBrowsers[browser] = true
			}
			if containsMSIE := strings.Contains(browser, "MSIE"); containsMSIE {
				isMSIE = true
				seenBrowsers[browser] = true
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))

		user.Reset()
	}

	fmt.Fprintf(out, "found users:\n%s\n", foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

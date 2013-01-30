HttpGet "https://gist.github.com/4668739.git" c5129558666c24a32c44298a986c3f3e945b56c4

import (
	"strings"
	//"fmt"
)

func GoKeywords() []string {
	//var go_spec = "/usr/local/go/doc/go_spec.html"
	//go_spec string
	//b, err := ioutil.ReadFile(go_spec)
	/*b, err := exec.Command("curl", "-s", "http://golang.org/ref/spec").Output()
	if err != nil {
		panic(err)
	}
	s := string(b)*/
	s := HttpGet("http://golang.org/ref/spec")
	//fmt.Println(s)
	f := strings.Index(s, "following keywords are reserved and may not be used as identifiers")
	s = s[f:]
	//fmt.Printf("%v", s)
	start := "<pre class=\"grammar\">"
	f = strings.Index(s, start)
	s = s[f+len(start)+0:]
	//fmt.Printf("%v", s)
	e := strings.Index(s, "</pre>")
	s = s[:e]
	//fmt.Printf(">%v<\n---\n", s)
	o := strings.Fields(s)
	//fmt.Printf("%v\n", o)
	//fmt.Printf("%v", strings.Join(o, ", "))
	return o
}
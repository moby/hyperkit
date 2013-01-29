import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) string {
	r, err := http.Get(url)
	if nil != err {
		panic(err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if nil != err {
		panic(err)
	}
	return string(b)
}
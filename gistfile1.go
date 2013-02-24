package gist5022726

import "encoding/json"
import . "gist.github.com/4668739.git"

func GetWeeklyIncome(username string) string {
	url := "https://www.gittip.com/" + username + "/public.json"

	b := HttpGetB(url)

	type Response struct {
		Receiving string
	}

	var animals Response
	err := json.Unmarshal(b, &animals)
	if err != nil {
		// TODO: Make this better
		println("error: ", err)
		return "Error."
	}
	return animals.Receiving
}

func main() {
	print("$" + GetWeeklyIncome("shurcooL"))
}
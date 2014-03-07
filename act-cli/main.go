package main

import "flag"

var (
	flagInput = flag.String("input", "../output/weibo.med", "Medium file to cook with.")
	flagModel = flag.String("model", "../output/weibo.actnet", "Model file to use.")
	flagCook  = flag.Bool("cook", false, "Generate the model from .med files.")
	flagShow  = flag.Bool("show", true, "Show the model")
)

func main() {
	flag.Parse()

	if *flagCook {
		cook()
	}
	if *flagShow {
		show()
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

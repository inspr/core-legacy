package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/examples/primes/yamls"
	"inspr.dev/inspr/pkg/controller/client"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/rest"
)

func main() {
	dapp := meta.App{}
	yaml.Unmarshal([]byte(yamls.PrimesYAML), &dapp)
	dapp.Meta.Name = "controllerprimes"
	config, err := client.GetInClusterConfigs()
	if err != nil {
		panic(err)
	}
	c := client.NewControllerClient(*config)
	mux := http.NewServeMux()

	mux.HandleFunc("/update", func(rw http.ResponseWriter, r *http.Request) {
		temp := dapp
		temp.Meta.Annotations = make(map[string]string)
		temp.Meta.Annotations["hahahahaha"] = "hehehehehe"
		diff, err := c.Apps().
			Update(context.Background(), temp.Meta.Name, &temp, false)
		diff.Print(os.Stdout)
		if err != nil {
			rest.ERROR(rw, err)
		}
	})
	mux.HandleFunc("/create", func(rw http.ResponseWriter, r *http.Request) {
		diff, err := c.Apps().Create(context.Background(), "", &dapp, false)
		diff.Print(os.Stdout)
		if err != nil {
			rest.ERROR(rw, err)
		}
	})
	mux.HandleFunc("/delete", func(rw http.ResponseWriter, r *http.Request) {
		diff, err := c.Apps().
			Delete(context.Background(), "controllerprimes", false)
		diff.Print(os.Stdout)
		if err != nil {
			rest.ERROR(rw, err)
		}
	})

	log.Fatalln(http.ListenAndServe(":8000", mux))
}

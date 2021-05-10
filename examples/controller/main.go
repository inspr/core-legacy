package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/inspr/inspr/examples/primes/yamls"
	"github.com/inspr/inspr/pkg/controller/client"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/rest"
	"gopkg.in/yaml.v2"
)

func main() {
	dapp := meta.App{}
	yaml.Unmarshal([]byte(yamls.PingPongYAML), &dapp)
	dapp.Meta.Name = "controllerpingpong"
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
		diff, err := c.Apps().Update(context.Background(), "", &temp, false)
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
		diff, err := c.Apps().Delete(context.Background(), "controllerpingpong", false)
		diff.Print(os.Stdout)
		if err != nil {
			rest.ERROR(rw, err)
		}
	})

	log.Fatalln(http.ListenAndServe(":8000", mux))
}

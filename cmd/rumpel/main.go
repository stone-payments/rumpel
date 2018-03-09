package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/stone-payments/rumpel/environment"
	"github.com/stone-payments/rumpel/rule"
)

func main() {
	env, err := environment.Read(os.Getenv("RUMPEL_MODE"), os.Args, nil)
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	rls, err := rule.Config(env.RulesConfigPath)
	if err != nil {
		log.Fatalf("Read config rule file error: %v", err)
	}

	t := template.Must(template.New("welcome").Parse(`
   ____                             _
  |  _ \ _   _ _ __ ___  _ __   ___| |
  | |_) | | | | '_ ' _ \| '_ \ / _ \ |
  |  _ /| |_| | | | | | | |_) |  __/ |
  |_| \_\\__,_|_| |_| |_| .__/ \___|_|
                        |_|

  Environment mode:  {{.Name}}
  Rule config path:  {{.RulesConfigPath}}
  Verbose mode:      {{.Verbose}}
  Running on port:   {{.ApplicationPort}}

`))
	if err := t.Execute(os.Stdout, env); err != nil {
		log.Println(err)
	}

	server := http.Server{
		Addr:    env.ApplicationPort,
		Handler: rls.Proxy(env.Verbose),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Panicf("Server error: %v", err)
	}
}

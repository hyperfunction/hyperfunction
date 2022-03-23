package main

import (
	"log"

	"knative.dev/hack/schema/commands"
	"knative.dev/hack/schema/registry"

	"github.com/hyperfunction/hyperfunction/pkg/apis/extensions/v1alpha1"
)

// schema is a tool to dump the schema for Eventing resources.
func main() {
	registry.Register(&v1alpha1.Function{})

	if err := commands.New("github.com/hyperfunction/hyperfunction").Execute(); err != nil {
		log.Fatal("Error during command execution: ", err)
	}
}

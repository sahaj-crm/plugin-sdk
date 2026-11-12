// Example Go extension. Build it into a manifest directory:
//
//	mkdir -p plugins/example-optin
//	go build -o plugins/example-optin/plugin ./examples/go-optin
package main

import (
	"context"
	"log"

	extensionapi "github.com/sahaj-crm/plugin-sdk"
)

type Extension struct{}

func (Extension) Info(context.Context) (*extensionapi.Info, error) {
	return &extensionapi.Info{
		ID:         "com.example.optin",
		Name:       "Example Opt-in",
		Version:    "0.1.0",
		APIVersion: extensionapi.APIVersion,
		Targets:    []string{"course.new", "course.edit"},
		Events:     []string{"course.create", "course.update"},
	}, nil
}

func (Extension) Contribute(
	_ context.Context,
	req *extensionapi.ContributionRequest,
) (*extensionapi.ContributionResponse, error) {
	if req.Action != "new" && req.Action != "edit" {
		return &extensionapi.ContributionResponse{}, nil
	}
	return &extensionapi.ContributionResponse{
		Sections: []extensionapi.Section{{
			Title: "Example Options",
			Fields: []extensionapi.Field{{
				Key:   "opt_in",
				Label: "Enable the example feature for this course",
				Type:  extensionapi.FieldCheckbox,
			}},
		}},
	}, nil
}

func (Extension) Handle(_ context.Context, ev *extensionapi.Event) (*extensionapi.EventResponse, error) {
	log.Printf("example extension: %s %s", ev.Phase, ev.Name)
	return &extensionapi.EventResponse{}, nil
}

func main() {
	extensionapi.Serve(Extension{})
}

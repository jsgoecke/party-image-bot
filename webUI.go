package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// IndexPage is the main page of the app
type IndexPage struct {
	app.Compo
	HumanPrompt string
	HumanImage  string
	AiPrompt    string
	AiImage     string
	PhoneNo     string
	Status      string
}

// Render renders the main page
func (ip *IndexPage) Render() app.UI {
	return app.P().Body(
		app.Style().Text("body { background-color: black; }"),
		app.Script().Src("web/websocket.js"),
		app.Th().Style("display", "flex").Style("justify-content", "center").Style("align-items", "center").Body(
			app.Img().ID("aiImage").Src("web/cat-spin.gif").Style("width", AiImageSize).Style("height", AiImageSize).Alt("AI Prompted Image"),
			app.H2().ID("aiPrompt").Style("text-align", "center"),
		))
}

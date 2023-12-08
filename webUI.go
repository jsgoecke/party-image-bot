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
	return app.Table().Body(
		app.Script().Src("web/websocket.js"),
		app.Img().Src("web/hitchhiker-thumb.jpg").Style("text-align", "center").Style("width", "180px").Style("height", "180px").Alt("Logo"),
		app.Table().Class("centered-table").Body(
			app.Tr().Body(
				app.Td().Body(
					app.H2().Style("text-align", "center").Text("Instructions"),
				),
				app.Td().Body(
					app.H2().Style("text-align", "center").Text("Prompt Sent"),
				),
				app.Td().Body(
					app.H2().Style("text-align", "center").Text("Embellished Prompt"),
				),
			),
			app.Tr().Body(
				app.Th().Body(
					app.H2().Text("Send an SMS to +1.650.284.3515"),
					app.H2().Text("or scan the QR Code"),
					app.Img().Src("web/sms-qr-code.png").Style("height", "240px").Style("width", "240px").Alt("SMS QR Code"),
					app.H2().ID("statusUpdate").Text("").Style("text-align", "center"),
				),
				app.Th().Body(
					app.Img().ID("humanImage").Src("web/cat-spin.gif").Style("width", "640px").Style("height", "640px").Alt("Human Prompted Image"),
				),
				app.Th().Body(
					app.Img().ID("aiImage").Src("web/cat-spin.gif").Style("width", "640px").Style("height", "640px").Alt("AI Prompted Image"),
				),
			),
			app.Tr().Body(
				app.Td().Text(""),
				app.Th().Body(
					app.H3().ID("humanPrompt").Text("").Style("text-align", "center"),
				),
				app.Th().Body(
					app.H3().ID("aiPrompt").Text("").Style("text-align", "center"),
				),
			)),
	)
}

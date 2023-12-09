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
		app.Script().Src("web/websocket.js"),
		app.Table().Class("centered-table").Body(
			app.Tr().Body(
				app.Th().Body(
					app.H2().Text("Send a DALL-E prompt via SMS to +1.650.284.3515"),
					app.H2().Text("or scan the QR Code"),
					app.Img().Src("web/sms-qr-code.png").Style("height", QRCodeImageSize).Style("width", QRCodeImageSize).Alt("SMS QR Code"),
					app.H2().ID("statusUpdate").Text("").Style("text-align", "center"),
					app.H2().ID("aiPrompt").Text("").Style("text-align", "center"),
				),
				app.Th().Body(
					app.Img().ID("aiImage").Src("web/cat-spin.gif").Style("width", AiImageSize).Style("height", AiImageSize).Alt("AI Prompted Image"),
				),
			),
		))
}

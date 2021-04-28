package main

// func svgIconToPng(size int) {
// 	canvas := svg.New()
// 	canvas.Start(size, size)
// 	canvas.Image(0, 0, 512, 512, "src.jpg", "0.50")
// 	canvas.End()
// }

type webmanifest struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Display   string `json:"display"`
	Icons     []icon `json:"icons"`
}

type icon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

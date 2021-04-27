package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

var mainTemplate *template.Template

var rootHtml http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	mainTemplate.Execute(w, mainParams{
		Configs:  configs,
		Statuses: statuses,
	})
}

var setState http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	fmt.Printf("parts: %v\n", parts)
	topic := fmt.Sprintf("cmnd/%s/%s", parts[2], parts[3])
	fmt.Printf("Sending topic %s\n", topic)
	c.Publish(topic, 0, false, parts[4])
	time.Sleep(time.Millisecond * 250)
	w.Header().Add("Location", "/")
	w.WriteHeader(302)
}

type mainParams struct {
	Configs  map[string]tasmotaConfig
	Statuses map[string]map[string]string
}

func initTemplates() {
	const tpl = `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0-beta1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-giJF6kkoqNQ00vy+HMDP7azOuL0xtbfIcaT9wjKHr8RbDVddVHyTfAAsrekwKmP1" crossorigin="anonymous">
			<title>Tasmota switch control</title>
		</head>
		<body>
		<div class="container">
			<div class="container">
				<h1>Tasmota switch control</h1>
				<div class="row" id="lights">
					{{range $topic, $values := .Statuses}}
					<div class="col-sm">
						<div class="card">
							<div class="card-header">
								<a href="http://{{(index $.Configs $topic).IP}}/">{{(index $.Configs $topic).DescriptiveName}}</a>
							</div>
							<div class="card-body">
								<div class="btn-group">
									{{range $switch, $state := $values}}
									<a 
										href="/set/{{$topic}}/{{$switch}}/{{if eq $state "OFF"}}ON{{else}}OFF{{end}}"
										class="btn btn-{{if eq $state "OFF"}}dark{{else}}light{{end}}">
											{{$switch}}
										</a>
									{{end}}
								</div>
							</div>
						</div>
					</div>
					{{end}}
				</div>
			</div>
		</div>

		</body>
	</html>`
	var err error
	mainTemplate, err = template.New("main").Parse(tpl)
	if err != nil {
		panic(err)
	}
}

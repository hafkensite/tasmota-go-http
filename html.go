package main

import (
	"encoding/json"
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

var webmanifestHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(webmanifest{
		Name:      "Tasmota go http",
		ShortName: "Tasmota",
		Display:   "standalone",
		Icons: []icon{icon{
			Src:   "data:image/png;base64,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",
			Sizes: "225x225",
			Type:  "image/png",
		}},
	})
	if err != nil {
		fmt.Print(err)
		w.WriteHeader(503)
	}
	w.Header().Add("Content-Type", "application/manifest+json")
	w.Write(bytes)
}

type mainParams struct {
	Configs  map[string]tasmotaConfig
	Statuses map[string]map[string]string
}

func initTemplates() {
	const tpl = `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link rel="manifest" href="/manifest.json">
			<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
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
									<span 
										id="btn_{{$topic}}_{{$switch}}"
										class="btn">
											{{$switch}}
										</span>
									{{end}}
								</div>
							</div>
						</div>
					</div>
					{{end}}
				</div>
			</div>
		</div>
		<script>
			var first = true;
			var internalstatus = {};
			const socket = new WebSocket('ws://localhost:8080/ws/states');
			socket.addEventListener('message', function (event) {
				console.log('Message from server ', event.data);
				internalstatus = JSON.parse(event.data);

				Object.entries(internalstatus).forEach(([topic,btns]) => {
					Object.entries(btns).forEach(([btn,state]) => {
						console.log(topic,btn,state);
						const el = document.getElementById('btn_' + topic + '_' + btn);
						if (first) {
							el.addEventListener('click', () => {
								const sendmsg = {};
								sendmsg[topic+'/'+btn] = internalstatus[topic][btn] == "OFF"?"ON":"OFF";
								socket.send(JSON.stringify(sendmsg));
							});
						}
						if (state == "OFF") {
							el.classList.add("btn-dark");
							el.classList.remove("btn-light");
						} else {
							el.classList.add("btn-light");
							el.classList.remove("btn-dark");
						}
					});
				});
				first = false;
			});
		</script>
		</body>
	</html>`
	var err error
	mainTemplate, err = template.New("main").Parse(tpl)
	if err != nil {
		panic(err)
	}
}

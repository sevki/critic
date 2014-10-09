package critic

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"text/template"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)

}
func index(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadFile("static/root.tmpl")
	if err != nil {
		w.Write([]byte("template not happy"))
		return
	}

	tmpl := template.New("")
	tmpl.Parse(string(bytes))
	tmpl.Execute(w, nil)

}

func validsize(i image.Image) bool {
	b := i.Bounds()
	if b.Dx()*b.Dy() < 10424*768 {
		return true
	} else {
		return false
	}

}

func upload(w http.ResponseWriter, req *http.Request) {
	f, _, err := req.FormFile("image")
	if err != nil {
		w.Write([]byte("you need to upload an image"))
		return
	}
	defer f.Close()

	i, _, err := image.Decode(f)
	if err != nil {
		w.Write([]byte("thats not a image"))
		return
	}

	if !validsize(i) {
		w.Write([]byte("that image is too big"))
		return
	}
	var plt color.Palette
	plt = palette.Plan9

	colors, img := AnalyzeAndConvert(i, plt)

	bytz, err := ioutil.ReadFile("static/root.tmpl")
	if err != nil {
		w.Write([]byte("template not happy"))
		return
	}

	tmpl := template.New("")
	tmpl.Parse(string(bytz))
	vars := make(map[string]interface{})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}

	vars["colors"] = colors
	vars["image"] = fmt.Sprintf(
		"data:image/png;base64,%s",
		base64.StdEncoding.EncodeToString(buf.Bytes()),
	)
	tmpl.Execute(w, vars)

}

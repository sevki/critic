package critic

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
	if b.Dx() < 1024 && b.Dy() < 768 {
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

	colors := Analyze(i, X11)

	bytes, err := ioutil.ReadFile("static/root.tmpl")
	if err != nil {
		w.Write([]byte("template not happy"))
		return
	}

	tmpl := template.New("")
	tmpl.Parse(string(bytes))
	tmpl.Execute(w, colors)

}

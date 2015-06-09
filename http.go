package main

import (
	"archive/zip"
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"
)

/***************************************************************************/

var templates map[string]*template.Template

/***************************************************************************/

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	s, ok := binData["s/"+path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=86400")

	w.Write([]byte(s))
}

func handleSkins(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	dir := http.Dir(skinDir)

	f, err := dir.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, f)
}

func handleNew(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			debug.PrintStack()
		}
	}()

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ff, _, err := r.FormFile("file")
	if err != nil {
		log.Panicln(err)
	}

	fb, err := ioutil.ReadAll(ff)
	if err != nil {
		log.Panicln(err)
	}

	zr, err := zip.NewReader(bytes.NewReader(fb), int64(len(fb)))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	time.Sleep(1 * time.Second)

	pb, pn, err := PackZip(zr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path, err := ioutil.TempDir(skinDir, "")
	if err != nil {
		log.Panicln(err)
	}

	err = ioutil.WriteFile(filepath.Join(path, "the.zip"), fb, 0644)
	if err != nil {
		log.Panicln(err)
	}

	err = ioutil.WriteFile(filepath.Join(path, pn), pb, 0644)
	if err != nil {
		log.Panicln(err)
	}

	sn := filepath.Base(path)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("/skins/" + sn + "/" + pn))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			debug.PrintStack()
		}
	}()

	var buf bytes.Buffer

	err := templates["index"].Execute(&buf, nil)
	if err != nil {
		log.Panicln(err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write(buf.Bytes())
}

/***********************************/

func initTemplates() {
	it := func(filenames ...string) *template.Template {
		t := template.New("")

		for _, f := range filenames {
			template.Must(t.Parse(binData[f]))
		}
		return t
	}

	templates = make(map[string]*template.Template)

	templates["index"] = it("base.html", "index.html")
}

func initHandlers() {
	ih := func(pattern string, f func(http.ResponseWriter, *http.Request)) {
		h := http.HandlerFunc(f)

		http.Handle(pattern, http.StripPrefix(pattern, h))
	}

	ih("/", handleRoot)
	ih("/s/", handleStatic)
	ih("/skins/", handleSkins)
	ih("/new", handleNew)
}

func listenAndServe(addr string) {
	log.Printf("Listening on http://%v/", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

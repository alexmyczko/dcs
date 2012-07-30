// vim:ts=4:sw=4:noexpandtab
package show

import (
	"fmt"
	"html/template"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var templates = template.Must(template.ParseFiles("templates/show.html"))

func Show(w http.ResponseWriter, r *http.Request) {
	query := r.URL
	filename := query.Query().Get("file")
	line, err := strconv.ParseInt(query.Query().Get("line"), 10, 0)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	log.Printf("Showing file %s, line %d\n", filename, line)

	// TODO: this needs to be a source-backend query instead
	// TODO: path configuration
	file, err := os.Open(`/media/sdg/debian-source-mirror/unpacked/` + filename)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	lines := strings.Split(string(contents), "\n")
	highestLineNr := fmt.Sprintf("%d", len(lines))

	// Since Go templates don’t offer any way to use {{$idx+1}}, we need to
	// pre-calculate line numbers starting from 1 here.
	lineNumbers := make([]int, len(lines))
	for idx, _ := range lines {
		lineNumbers[idx] = idx + 1
	}

	err = templates.ExecuteTemplate(w, "show.html", map[string]interface{} {
		// XXX: Has string(contents) any problems when the file is not valid UTF-8?
		// (while the indexer only cares for UTF-8, an attacker could send us any file path)
		"lines": lines,
		"numbers": lineNumbers,
		"lnrwidth": len(highestLineNr),
		"filename": filename,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
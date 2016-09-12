package main

import (
	"fmt"
	"github.com/micromdm/squirrel/munki/datastore"
	"github.com/micromdm/squirrel/munki/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"
const REPO_PATH string = "/Users/Shared/repo"

type Context struct {
	Static string
	Prods  []ParsedPkginfo
}

// Template expects icon, 'product' name, description, and version, which becomes the download URL
type ParsedPkginfo struct {
	Icon        string
	Name        string
	Descript    string
	Version     string
	DownloadURL string
}

func Parse(in_pkginfo *models.PkgsInfo, ppi *ParsedPkginfo) *ParsedPkginfo {
	ppi.Icon = in_pkginfo.IconName
	ppi.Name = in_pkginfo.Name
	ppi.Descript = in_pkginfo.Description
	ppi.Version = in_pkginfo.Version
	ppi.DownloadURL = in_pkginfo.PackageCompleteURL
	return ppi
}

func Home(w http.ResponseWriter, req *http.Request) {
	repo := datastore.SimpleRepo{Path: REPO_PATH}
	log.Print("found repo: ", repo)
	catalog, err := repo.AllPkgsinfos()
	if err != nil {
		log.Fatal("faux-catalog building error: ", err)
	}
	allcat := *catalog
	prods := make([]ParsedPkginfo, 0)
	for _, pkginfo := range allcat {
		pi := &ParsedPkginfo{}
		ppi := Parse(pkginfo, pi)
		prods = append(prods, *ppi)
	}
	context := Context{Prods: prods}
	render(w, "gomoscar.html", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
	tmpl_list := []string{"templates/gomoscar.html",
		fmt.Sprintf("templates/%s", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc(STATIC_URL, StaticHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// func moscarg(serverURL string) http.HandlerFunc {
//     return func(rw http.ResponseWriter, r *http.Request) {
//         repo := datastore.SimpleRepo{Path: "/Users/Shared/repo"}
//         catalog, err := repo.AllPkgsinfos()
//         if err != nil {
//             rw.WriteHeader(500)
//             return
//         }
//         prods := *catalog
//         title := r.URL.Path[len("/view/"):]
//         p := loadPage(title)
//         t, _ := template.ParseFiles("index.html")
//         t.Execute(rw, p)
//     }
// }
//
// func loadPage(title string) *Page {
//     filename := title + ".txt"
//     body, _ := ioutil.ReadFile(filename)
//     return &Page{Title: title, Body: body}
// }
//
// func main() {
//     repoPath := "/Users/Shared/repo/pkgs"
//     serverURL := "http://localhost:8080/pkgs/"
//     // http.HandleFunc("/view/", moscarg(serverURL))
//     http.HandleFunc("/", moscarg(serverURL))
//     // http.Handle("/pkgs/", http.StripPrefix("/pkgs", http.FileServer(http.Dir(repoPath))))
//     http.ListenAndServe(":8080", nil)
// }

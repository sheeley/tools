package bearexport

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/richardwilkes/toolbox/errs"
)

type Input struct {
	Verbose, IncludeTrashed, SeparateFiles bool
	Outdir                                 string
}

type Output struct {
	Notes []*Note
}

type Note struct {
	Text      string  `json:"text" db:"ZTEXT"`
	ID        string  `json:"id" db:"ZUNIQUEIDENTIFIER"`
	Trashed   bool    `json:"trashed" db:"ZTRASHED"`
	Pinned    bool    `json:"pinned" db:"ZPINNED"`
	HasFiles  bool    `json:"hasFiles" db:"ZHASFILES"`
	HasImages bool    `json:"hasImages" db:"ZHASIMAGES"`
	Created   float64 `json:"created" db:"ZCREATIONDATE"`
	Modified  float64 `json:"modified" db:"ZMODIFICATIONDATE"`
}

// BearRunning indicates whether Bear.app is currently running.
func BearRunning(in *Input) (bool, error) {
	return false, nil
	// cmd := exec.Command("ps", "x")
	// if in.Verbose {
	// 	fmt.Printf("Running %v\n", cmd)
	// }
	// o, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return false, errs.Wrap(err)
	// }
	// if bytes.Contains(o, []byte("Bear.app")) {
	// 	return true, nil
	// }
	// return false, nil
}

// BearExport queries Bear's sqlite database to export each note.
func BearExport(in *Input) (*Output, error) {
	if in.Verbose {
		fmt.Println("Creating Temp File")
	}
	tf, err := ioutil.TempFile("", "bear.sqlite")
	if err != nil {
		return nil, errs.Wrap(err)
	}

	sqliteFile := os.ExpandEnv("$HOME/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite")
	if in.Verbose {
		fmt.Printf("temp file: %s\n", tf.Name())
		fmt.Printf("Opening sqlite db %s\n", sqliteFile)
	}
	inFile, err := os.Open(sqliteFile)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if in.Verbose {
		fmt.Println("Copying db to temp file")
	}
	_, err = io.Copy(tf, inFile)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	inFile.Close()
	tf.Close()

	exportQuery := `
	SELECT ZUNIQUEIDENTIFIER, ZTEXT, ZHASIMAGES, ZHASFILES, ZTRASHED, 
	cast(ZCREATIONDATE+978307200.0 as float) as ZCREATIONDATE, 
	cast(ZMODIFICATIONDATE+978307200.0 as float) as ZMODIFICATIONDATE, 
	ZPINNED 
	FROM ZSFNOTE 
	WHERE ZUNIQUEIDENTIFIER!='SFNoteIntro3' 
	AND ZTRASHED=0;`

	if in.Verbose {
		// countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ZSFNOTE%s;", whereClause)
		// fmt.Println(countQuery)
		fmt.Println(exportQuery)
	}

	db, err := sqlx.Open("sqlite3", tf.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	o := &Output{}

	err = db.Select(&o.Notes, exportQuery)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if in.Verbose {
		fmt.Printf("%d notes\n", len(o.Notes))
	}

	return o, nil
}

func ModifyNotesDates(in *Input, notes []*Note) error {
	for _, n := range notes {
		// fmt.Printf("%s: %v, %v\n", n.ID, n.Created, n.Modified)
		file := path.Join("/Users/sheeley/Library/Mobile Documents/iCloud~md~obsidian/Documents/Notes", n.ID+".md")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// path/to/whatever does not exist
			fmt.Println("skipnoexist", file)
			continue
		}
		created := time.Unix(int64(n.Created), 0)
		modified := time.Unix(int64(n.Modified), 0)
		err := os.Chtimes(file, created, modified)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
	// 	filesToReplace := map[string]string{}
	// for _, n := range notes {
	// 	// fmt.Printf("%s: %v, %v\n", n.ID, n.Created, n.Modified)
	// 	file := path.Join("/Users/sheeley/Library/Mobile Documents/iCloud~md~obsidian/Documents/Notes", n.ID+".md")
	// 	if _, err := os.Stat(file); os.IsNotExist(err) {
	// 		// path/to/whatever does not exist
	// 		fmt.Println("skipnoexist", file)
	// 		continue
	// 	}
	// 	// contents, err := ioutil.ReadFile(file)
	// 	// if err != nil {
	// 	// 	return errs.Wrap(err)
	// 	// }

	// 	// if bytes.HasPrefix(contents, []byte("# ")) {
	// 	// 	lines := bytes.Split(contents, []byte("\n"))
	// 	// 	firstLine := string(bytes.TrimPrefix(lines[0], []byte("# ")))
	// 	// 	newFilename := path.Join("/Users/sheeley/Library/Mobile Documents/iCloud~md~obsidian/Documents/Notes", firstLine+".md")
	// 	// 	if _, err := os.Stat(newFilename); err == nil {
	// 	// 		fmt.Println("skipexist", newFilename)
	// 	// 		continue
	// 	// 	}

	// 	// 	filesToReplace[n.ID] = newFilename
	// 	// }
	// }

	// for _, n := range notes {
	// 	file := path.Join("/Users/sheeley/Library/Mobile Documents/iCloud~md~obsidian/Documents/Notes", n.ID+".md")
	// 	if _, err := os.Stat(file); os.IsNotExist(err) {
	// 		// path/to/whatever does not exist
	// 		fmt.Println("skipnoexist", file)
	// 		continue
	// 	}
	// 	contents, err := ioutil.ReadFile(file)
	// 	if err != nil {
	// 		return errs.Wrap(err)
	// 	}
	// 	contents = append([]byte(fmt.Sprintf(`
	// 		---
	// 		bearID: %s
	// 		---`, n.ID)), contents...)

	// 	for original, replacement := range filesToReplace {
	// 		contents = bytes.ReplaceAll(contents, []byte(original), []byte(replacement))
	// 	}

	// 	newFilename, ok := filesToReplace[n.ID]
	// 	if ok {
	// 		if _, err := os.Stat(newFilename); err == nil {
	// 			fmt.Println("skipexist", newFilename)
	// 			continue
	// 		}
	// 		fmt.Println(n.ID, " => ", newFilename)
	// 	} else {
	// 		newFilename = file
	// 	}

	// 	err = os.WriteFile(newFilename, contents, 0644)
	// 	if err != nil {
	// 		return errs.Wrap(err)
	// 	}
	// 	created := time.Unix(int64(n.Created), 0)
	// 	modified := time.Unix(int64(n.Modified), 0)
	// 	err = os.Chtimes(newFilename, created, modified)
	// 	if err != nil {
	// 		return errs.Wrap(err)
	// 	}
	// }
	// return nil
}

func WriteNotes(in *Input, notes []*Note) error {
	err := os.MkdirAll(in.Outdir, os.ModePerm)
	if err != nil {
		return errs.Wrap(err)
	}

	// if !in.SeparateFiles {
	// 	fp := filepath.Join(in.Outdir, "notes.json")
	// 	err = encode(fp, notes, in.Verbose)
	// 	if err != nil {
	// 		return errs.Wrap(err)
	// 	}
	// 	return nil
	// }
	fmt.Printf("Writing %d notes", len(notes))
	for _, n := range notes {
		if n.ID == "SFNoteIntro3" {
			continue
		}
		fp := filepath.Join(in.Outdir, n.ID+".md")
		err = encode(fp, n, in.Verbose)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

var linkPattern = regexp.MustCompile(`bear:\/\/x-callback-url\/open-note\?id=([A-Z0-9-]+)`)
var nakedLinkPattern = regexp.MustCompile(`\(([A-Z0-9]{8}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{12}-[A-Z0-9]{4}-[A-Z0-9]{16})\)`)

func encode(fp string, data *Note, verbose bool) error {
	if verbose {
		fmt.Printf("writing %s\n", fp)
	}
	f, err := os.Create(fp)
	if err != nil {
		return errs.Wrap(err)
	}
	defer f.Close()
	if data.Pinned {
		_, err := f.WriteString(`
---
pinned: true
---
		`)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	processed := linkPattern.ReplaceAllString(data.Text, "$1.md")
	// processed = strings.ReplaceAll(processed, "bear://", "")
	processed = nakedLinkPattern.ReplaceAllString(processed, "($1.md)")
	_, err = f.WriteString(processed)
	// err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return errs.Wrap(err)
	}
	// if verbose {
	// 	fmt.Printf("closing %s\n", fp)
	// }
	// f.Close()
	return nil
}

package bearexport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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
	Text      string `json:"text" db:"ZTEXT"`
	ID        string `json:"id" db:"ZUNIQUEIDENTIFIER"`
	Trashed   bool   `json:"trashed" db:"ZTRASHED"`
	Pinned    bool   `json:"pinned" db:"ZPINNED"`
	HasFiles  bool   `json:"hasFiles" db:"ZHASFILES"`
	HasImages bool   `json:"hasImages" db:"ZHASIMAGES"`
	Created   string `json:"created" db:"ZCREATIONDATE"`
	Modified  string `json:"modified" db:"ZMODIFICATIONDATE"`
}

// BearRunning indicates whether Bear.app is currently running.
func BearRunning(in *Input) (bool, error) {
	cmd := exec.Command("ps", "x")
	if in.Verbose {
		fmt.Printf("Running %v\n", cmd)
	}
	o, err := cmd.CombinedOutput()
	if err != nil {
		return false, errs.Wrap(err)
	}
	if bytes.Contains(o, []byte("Bear.app")) {
		return true, nil
	}
	return false, nil
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

	whereClause := " WHERE ZTRASHED=0"
	if in.IncludeTrashed {
		whereClause = ""
	}

	exportQuery := fmt.Sprintf("SELECT ZUNIQUEIDENTIFIER, ZTEXT, ZHASIMAGES, ZHASFILES, ZTRASHED, ZCREATIONDATE, ZMODIFICATIONDATE, ZPINNED FROM ZSFNOTE%s;", whereClause)

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

// WriteNotes creates a single json file for each note in dir.
func WriteNotes(in *Input, notes []*Note) error {
	err := os.MkdirAll(in.Outdir, os.ModePerm)
	if err != nil {
		return errs.Wrap(err)
	}

	if !in.SeparateFiles {
		fp := filepath.Join(in.Outdir, "notes.json")
		err = encode(fp, notes, in.Verbose)
		if err != nil {
			return errs.Wrap(err)
		}
		return nil
	}
	fmt.Printf("Writing %d notes", len(notes))
	for _, n := range notes {
		fp := filepath.Join(in.Outdir, n.ID+".json")
		err = encode(fp, n, in.Verbose)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func encode(fp string, data interface{}, verbose bool) error {
	if verbose {
		fmt.Printf("writing %s\n", fp)
	}
	f, err := os.Create(fp)
	if err != nil {
		return errs.Wrap(err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return errs.Wrap(err)
	}
	// if verbose {
	// 	fmt.Printf("closing %s\n", fp)
	// }
	// f.Close()
	return nil
}

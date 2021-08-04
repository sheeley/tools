package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/richardwilkes/toolbox/errs"
)

var (
	timeFormat = time.RFC3339
	dateFormat = "2006-01-02"

	linkPattern          = regexp.MustCompile(`bear:\/\/x-callback-url\/open-note\?id=([A-Z0-9-]+)`)
	nakedLinkPattern     = regexp.MustCompile(`\(([A-Z0-9]{8}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{12}-[A-Z0-9]{4}-[A-Z0-9]{16})\)`)
	tagWithSpacesPattern = regexp.MustCompile(`#\S.+\S#`)
	fileEmbedPattern     = regexp.MustCompile(`\[file:.*\]`)
	imageEmbedPattern    = regexp.MustCompile(`\[image:.*/`)

	namePattern           = regexp.MustCompile("[^a-zA-Z0-9 ]+")
	startsWithWordPattern = regexp.MustCompile(`^\w`)
	multipleSpacePattern  = regexp.MustCompile(`\s{2,}`)
)

type Input struct {
	Verbose bool
	Outdir  string
}

type Note struct {
	Text         string  `db:"ZTEXT"`
	ID           string  `db:"ZUNIQUEIDENTIFIER"`
	Pinned       bool    `db:"ZPINNED"`
	Created      float64 `db:"ZCREATIONDATE"`
	Modified     float64 `db:"ZMODIFICATIONDATE"`
	CombinedTags *string `db:"ZTAGS"`
}

func (n *Note) tags() []string {
	if n.CombinedTags == nil || len(*n.CombinedTags) == 0 {
		return nil
	}
	return strings.Split(*n.CombinedTags, ",")
}

// BearExport queries Bear's sqlite database to export each note.
func BearExport(in *Input) ([]*Note, error) {
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
	SELECT note.ZUNIQUEIDENTIFIER, note.ZTEXT, note.ZPINNED,
	cast(note.ZCREATIONDATE+978307200.0 as float) as ZCREATIONDATE, 
	cast(note.ZMODIFICATIONDATE+978307200.0 as float) as ZMODIFICATIONDATE,
	group_concat(ZSFNOTETAG.ZTITLE) as ZTAGS
	FROM ZSFNOTE as note
	LEFT JOIN Z_7TAGS on Z_7TAGS.Z_7NOTES=note.Z_PK
	LEFT JOIN ZSFNOTETAG on Z_7TAGS.Z_14TAGS=ZSFNOTETAG.Z_PK
	WHERE ZUNIQUEIDENTIFIER!='SFNoteIntro3' 
	AND ZTRASHED=0
	GROUP BY note.ZUNIQUEIDENTIFIER;`

	if in.Verbose {
		fmt.Println(exportQuery)
	}

	db, err := sqlx.Open("sqlite3", tf.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var notes []*Note
	err = db.Select(&notes, exportQuery)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if in.Verbose {
		fmt.Printf("%d notes\n", len(notes))
	}

	return notes, nil
}

func main() {
	in := &Input{}
	flag.BoolVar(&in.Verbose, "v", false, "verbose logging")
	flag.StringVar(&in.Outdir, "o", ".", "output directory")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if !path.IsAbs(in.Outdir) {
		in.Outdir = path.Join(wd, in.Outdir)
	}

	notes, err := BearExport(in)
	if err != nil {
		panic(err)
	}

	// writeCSVSummary(out)

	lookup := createLookup(notes)
	for _, note := range notes {
		noteTags := note.tags()
		var cleanedTags []string
		cleanedTagLookup := map[string]string{}
		for _, tag := range noteTags {
			cleaned := strings.TrimSuffix(strings.ReplaceAll(tag, " ", "_"), "#")
			cleanedTags = append(cleanedTags, cleaned)
			cleanedTagLookup[tag] = cleaned
		}
		tagJSON := ""
		if len(cleanedTags) > 0 {
			tagJSON = "tags: " + toJson(cleanedTags)
		}
		created := time.Unix(int64(note.Created), 0)
		modified := time.Unix(int64(note.Modified), 0)
		contents := fmt.Sprintf(`
---
bearID: %s
created: %s
modified: %s
%s
---
%s
		`, note.ID, created.Format(timeFormat), modified.Format(timeFormat), tagJSON, processText(lookup, note.Text, cleanedTagLookup))

		// if strings.Contains(note.Text, "#articles/product thinking#") {
		// strings.Contains(note.Text, "bear://") {
		// if strings.Contains(note.Text, "image:") {
		// 	fmt.Println(createName(note), contents, "\n\n\n", note.Text)
		// }

		filePath := path.Join(in.Outdir, createName(note))
		if in.Verbose {
			fmt.Printf("writing %s\n", filePath)
		}

		err = ioutil.WriteFile(filePath, []byte(contents), 0644)
		if err != nil {
			panic(err)
		}
		os.Chtimes(filePath, created, modified)
	}
	if in.Verbose {
		fmt.Println("writing summary")
	}
	f, err := os.Create(path.Join(in.Outdir, "export-data.csv"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writeCSVSummary(notes, f)
}

func cleanTag(t string) string {
	return strings.TrimSuffix(strings.ReplaceAll(t, " ", "_"), "#")
}

func toJson(tags []string) string {
	s := &strings.Builder{}
	err := json.NewEncoder(s).Encode(tags)
	if err != nil {
		panic(err)
	}
	return s.String()
}

func processText(lookup map[string]string, text string, tags map[string]string) string {
	text = linkPattern.ReplaceAllString(text, "$1.md")
	text = strings.ReplaceAll(text, "bear://", "")
	text = nakedLinkPattern.ReplaceAllString(text, "($1.md)")
	text = fileEmbedPattern.ReplaceAllString(text, "")
	text = imageEmbedPattern.ReplaceAllString(text, "![[images/")

	for original, replacement := range lookup {
		text = strings.ReplaceAll(text, original, replacement)
	}

	for tag, cleaned := range tags {
		if strings.Contains(tag, " ") {
			text = strings.ReplaceAll(text, "#"+tag+"#", "#"+cleaned)
		}
	}

	text = strings.ReplaceAll(text, ".md.md", ".md")

	// iterate over each line. Remove tags and #. If just whitespace, eliminate
	outText := ""
	for _, line := range strings.Split(text, "\n") {
		workingLine := line
		include := false
		if len(strings.TrimSpace(workingLine)) == 0 {
			include = true
		}

		if !include {
			for tag, cleaned := range tags {
				workingLine = strings.ReplaceAll(workingLine, "#"+tag, "")
				workingLine = strings.ReplaceAll(workingLine, "#"+cleaned, "")
			}
			include = strings.TrimSpace(workingLine) != ""
		}

		if include {
			outText = outText + "\n" + line
		}
	}

	return strings.TrimSpace(outText)
}

func createLookup(notes []*Note) map[string]string {
	lookup := map[string]string{}
	for _, note := range notes {
		lookup[note.ID] = createName(note)
	}
	return lookup
}

func writeCSVSummary(notes []*Note, w io.Writer) {
	c := csv.NewWriter(w)
	for _, n := range notes {
		created := time.Unix(int64(n.Created), 0)
		modified := time.Unix(int64(n.Modified), 0)
		newName := createName(n)
		c.Write([]string{n.ID, newName, created.Format("200601021504.05"), modified.Format("200601021504.05")})
	}
	c.Flush()
}

func createName(note *Note) string {
	lines := strings.Split(note.Text, "\n")
	name := lines[0]
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") || startsWithWordPattern.MatchString(line) {
			name = line
			break
		}
	}
	name = namePattern.ReplaceAllString(name, "")
	name = multipleSpacePattern.ReplaceAllString(name, " ")
	if len(name) > 60 {
		name = name[0:60]
	}
	name = strings.TrimSpace(name)
	created := time.Unix(int64(note.Created), 0)
	// if len(newName) < 10 {
	// 	continue
	// }
	name = fmt.Sprintf("%s %s", created.Format(dateFormat), name)
	return strings.TrimSpace(name) + ".md"
}

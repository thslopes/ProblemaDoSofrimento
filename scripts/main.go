package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"log"

	"github.com/bmaupin/go-epub"
)

const (
	// effectiveGoCoverImg      = "assets/covers/capa.png"
	epubSimpleFilename = "problemaDoSofrimento.epub"
	epubFilename       = "problemaDoSofrimentoComentado.epub"
	epubTitle          = "O Problema do Sofrimento - Uma Perspectiva Bíblica"
	epubTitleFilename  = "title.xhtml"
	mdUrl              = "../docs/README.md"
	epubCSSFile        = "ebub.css"
	epubId             = "e2961d9db6194dc584ad206344c94023tsl"
	epubSimpleId       = "f695849e3d8942adaa44686408869317tsl"
	// preFontFile              = "assets/fonts/SourceCodePro-Regular.ttf"
)

var (
	titleRegex  = regexp.MustCompile(`[^a-zA-Z]`)
	fullVersion = true
	epubCSSPath = ""
)

type epubSection struct {
	title       string
	filename    string
	text        []string
	subSections []*epubSection
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "simple" {
		fullVersion = false
	}
	err := buildEffectiveGo()
	if err != nil {
		log.Printf("Error building Effective Go: %s", err)
	}
}

func buildEffectiveGo() error {
	r, err := os.ReadFile(mdUrl)
	if err != nil {
		return err
	}

	sections := []*epubSection{}
	var section *epubSection
	var subSection *epubSection

	for _, line := range strings.Split(string(r), "\n") {
		switch {
		case line == "## Introdução" && !fullVersion, len(strings.ReplaceAll(line, " ", "")) == 0:
			continue
		case strings.HasPrefix(line, "# "):
			section, line = buildSection(line, "# ", "<h1>%s</h1>")
			sections = append(sections, section)
		case strings.HasPrefix(line, "## "):
			section, line = buildSection(line, "## ", "<h2>%s</h2>")
			sections = append(sections, section)
		case strings.HasPrefix(line, "### "):
			subSection, line = buildSection(line, "### ", "<h3>%s</h3>")
			section.subSections = append(section.subSections, subSection)
		case strings.HasPrefix(line, ">"):
			line = fmt.Sprintf("<h4>%s</h4>", line[1:])
		case fullVersion:
			line = fmt.Sprintf("<p>%s</p>", line)
		default:
			line = ""
		}

		if len(section.subSections) > 0 {
			subSection.text = append(subSection.text, line)
		} else {
			section.text = append(section.text, line)
		}
	}

	e := epub.NewEpub(epubTitle)
	// effectiveGoCoverImgPath, err := filepath.Abs(effectiveGoCoverImg)
	// effectiveGoCoverImgPath, err := e.AddImage(effectiveGoCoverImg, "cover.png")
	// if err != nil {
	// return err
	// }
	// e.SetCover(effectiveGoCoverImgPath, "")

	epubCSSPath, err = e.AddCSS(epubCSSFile, "")
	if err != nil {
		return err
	}

	// _, err = e.AddFont(preFontFile, "")
	// if err != nil {
	// return err
	// }

	// Iterate through each section and add it to the EPUB
	for _, section := range sections {
		sectionContent := ""
		for _, sectionNode := range section.text {
			sectionContent += sectionNode
		}
		sn, err := e.AddSection(sectionContent, section.title, section.filename, epubCSSPath)
		if err != nil {
			return err
		}
		err = addSubSection(section, e, sn)
		if err != nil {
			return err
		}
	}

	id := epubId
	fileName := epubFilename

	if !fullVersion {
		id = epubSimpleId
		fileName = epubSimpleFilename
	}

	e.SetIdentifier(id)
	e.SetAuthor("Thiago Lopes")

	err = e.Write(fileName)
	if err != nil {
		return err
	}

	return nil
}

func addSubSection(section *epubSection, e *epub.Epub, sn string) error {
	for _, subSection := range section.subSections {
		subSectionContent := ""
		for _, subSectionNode := range subSection.text {
			subSectionContent += subSectionNode
		}
		ssn, err := e.AddSubSection(sn, subSectionContent, subSection.title, subSection.filename, epubCSSPath)
		if err != nil {
			return err
		}
		addSubSection(subSection, e, ssn)
	}
	return nil
}

func buildSection(line, prefix, titlePattern string) (*epubSection, string) {
	sectionTitle := strings.Replace(line, prefix, "", 1)
	sectionFilename := titleToFilename(sectionTitle)

	return &epubSection{
		filename:    sectionFilename,
		title:       sectionTitle,
		text:        []string{},
		subSections: []*epubSection{},
	}, fmt.Sprintf(titlePattern, sectionTitle)
}

func titleToFilename(title string) string {
	title = strings.ToLower(title)
	title = titleRegex.ReplaceAllString(title, "")

	return fmt.Sprintf("%s.xhtml", title)
}

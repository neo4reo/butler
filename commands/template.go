package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/netzkern/butler/config"
	survey "gopkg.in/AlecAivazis/survey.v1"
	git "gopkg.in/src-d/go-git.v4"
)

const (
	startDelim = "[["
	endDelim   = "]]"
)

var (
	allowedExtensions = [...]string{".md", ".txt", ".html", ".htm", ".rtf", ".json", ".yml", ".csproj", ".sln"}
	blacklistDirs     = map[string]bool{
		"node_modules":     true,
		"bower_components": true,
		"jspm_packages":    true,
		"dist":             true,
		"logs":             true,
		"bin":              true,
	}
)

type (
	project struct {
		Name     string
		Path     string
		Template string
	}
	Templating struct {
		Templates []config.Template
	}
)

func (t *Templating) cloneRepo(repoURL string, dest string) error {
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	return err
}

func (t *Templating) getTemplateByName(name string) *config.Template {
	for _, tpl := range t.Templates {
		if tpl.Name == name {
			return &tpl
		}
	}

	return nil
}

func (t *Templating) getTemplateOptions() []string {
	tpls := make([]string, 0)

	for _, tpl := range t.Templates {
		tpls = append(tpls, tpl.Name)
	}

	return tpls
}

func (t *Templating) prompts() (*project, error) {
	var simpleQs = []*survey.Question{
		{
			Name:     "Template",
			Validate: survey.Required,
			Prompt: &survey.Select{
				Message: "What system are you using?",
				Options: t.getTemplateOptions(),
			},
		},
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "What is the project name?",
			},
			Validate: survey.Required,
		},
		{
			Name:     "Path",
			Validate: survey.Required,
			Prompt: &survey.Input{
				Message: "What is the destination?",
				Default: "./src",
			},
		},
	}

	var project = &project{}

	// ask the question
	err := survey.Ask(simpleQs, project)

	if err != nil {
		return nil, err
	}

	return project, nil
}

// Run - Runs the command
func (t *Templating) Run() error {
	project, err := t.prompts()

	if err != nil {
		return err
	}

	tpl := t.getTemplateByName(project.Template)

	if tpl != nil {
		err := t.cloneRepo(tpl.Url, project.Path)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("butler: template %s could not be found", project.Template)
	}

	walkErr := filepath.Walk(project.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip blacklisted directories
		if info.IsDir() && blacklistDirs[info.Name()] {
			return filepath.SkipDir
		}

		// ignore hidden dirs
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// skip directories but go ahead with traversing
		if info.IsDir() {
			return nil
		}

		// check for valid file extension
		fileExt := strings.ToLower(info.Name())
		validExt := false
		for _, ext := range allowedExtensions {
			if strings.HasSuffix(fileExt, ext) {
				validExt = true
				break
			}
		}

		if !validExt {
			return nil
		}

		dat, err := ioutil.ReadFile(path)

		tmpl, err := template.New(path).Delims(startDelim, endDelim).Parse(string(dat))

		f, err := os.Create(path)

		defer f.Close()

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("butler: File %s recovered due to invalid template! Error: %s \n", path, r)
				ioutil.WriteFile(path, dat, 0644)
			}
		}()

		if err != nil {
			return err
		}

		var templateData = struct {
			ProjectName string
		}{
			project.Name,
		}

		err = tmpl.Execute(f, templateData)

		if err != nil {
			return err
		}

		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	return nil
}

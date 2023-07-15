/*
Copyright © 2023 Ambor <saltbo@foxmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package decrypter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/spf13/cobra"
)

var qs = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Select the decrypter type",
			Options: []string{"remote", "local"},
			Default: "remote",
		},
		Validate: survey.Required,
	},
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Input the profile name"},
		Validate: survey.ComposeValidators(survey.Required, survey.MaxLength(32)),
	},
}

var remoteQs = []*survey.Question{
	{
		Name:   "baseURL",
		Prompt: &survey.Input{Message: "Input the Base URL"},
		Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
			v := ans.(string)
			if !strings.HasPrefix(v, "http") {
				return fmt.Errorf("url must be started with http(s)://")
			}

			_, err := url.Parse(v)
			return err
		}),
	},
	{
		Name:   "token",
		Prompt: &survey.Input{Message: "Input the Token"},
	},
}

// InstallCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := entity.NewConfig()
		if err != nil {
			return err
		}

		answers := struct {
			Type string
			Name string
		}{}
		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		var props entity.Properties
		if answers.Type == "local" {
			err = survey.AskOne(&survey.Input{Message: "Input your local cmd"}, &props.Cmd)
		} else {
			err = survey.Ask(remoteQs, &props)
		}
		if err != nil {
			return err
		}

		d := entity.NewDecrypter(answers.Type, props)
		if err := cfg.DecrypterInstall(answers.Name, d); err != nil {
			return err
		}

		return cfg.Save()
	},
}

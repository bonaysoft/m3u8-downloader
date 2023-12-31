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
	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := entity.NewConfig()
		if err != nil {
			return err
		}

		if err := cfg.DecrypterUninstall(viper.GetString("name")); err != nil {
			return err
		}

		return cfg.Save()
	},
}

func init() {
	DeleteCmd.Flags().String("name", "", "specify the name")

	_ = viper.BindPFlags(DeleteCmd.Flags())
}

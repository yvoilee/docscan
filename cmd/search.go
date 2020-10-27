/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yvoilee/docscan/pkg/docx"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var dir, out string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the special strings",
	Long: `Search the special strings, default directory is current work directory
default output file is search_result.txt under current work directory.

Example
search key.
search key --dir directory 
search key --o file 
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var dirRoot, output string
		pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		if len(dir) == 0 {
			dirRoot = pwd
		} else {
			dirRoot = dir
		}
		if !IsDir(dirRoot) {
			log.Fatal("dir is not Exists, please check!!")
		}

		output = out
		if len(out) == 0 {
			output = path.Join(pwd, "search_result.txt")
		}
		searchKey := args[0]
		f, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer func() {
			io.WriteString(f, "---------------end\n\n")
			f.Close()
		}()
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("output file %s not exists, create file\n", output)
				f, _ = os.Create(output)
			} else {
				log.Fatal(err)
			}
		}
		io.WriteString(f, fmt.Sprintf("---------------dir: %s\n---------------search: %s \n",
			dirRoot, searchKey))
		SearchDirDoc(dirRoot, searchKey, f)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&dir, "dir", "d", "", "assign search directory")
	searchCmd.Flags().StringVarP(&out, "out", "o", "", "assign search output directory")
}
func SearchDirDoc(dir, searchKey string, outWrite *os.File) {
	filesDoc, _ := GetAllFiles(dir, []string{".doc", ".docx"})
	for _, file := range filesDoc {
		r, err := docx.ReadDocxFile(file)
		if err != nil {
			continue
		}
		doc := r.Editable()
		if doc.Search(searchKey) {
			fmt.Println(file)
			io.WriteString(outWrite, file+"\n")
		}
	}
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
func GetAllFiles(dirPth string, filters []string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth+PthSep+fi.Name(), filters)
		} else {
			// 过滤指定格式
			var ok bool
			for _, filter := range filters {
				ok = strings.HasSuffix(fi.Name(), filter)
				if ok {
					break
				}
			}
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table, filters)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

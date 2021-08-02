package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", show)
	http.HandleFunc("/add", add)
	http.HandleFunc("/remove", remove)
	http.ListenAndServe(":80", nil)
}

func show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<!DOCTYPE html><html lang=\"ja\"><head><meta charset=\"UTF-8\"><title>Todo</title><link rel=\"stylesheet\" href=\"http://18.217.133.253/static/style.css\"></head>")
	fmt.Fprintf(w, "<body>")
	fmt.Fprintf(w, "<div id=\"title\">メモ</div>")
	fmt.Fprintf(w, "<ul id=\"LR\">")

	todoFiles, err := getTodoFile()
	if err != nil {
		fmt.Fprintf(w, "タスクの取得に失敗しました")
	}

	fmt.Fprintf(w, "<li id=\"L\">")
	fmt.Fprintf(w, "<ul id=\"todo-list\">")
	if err := render(w, r, todoFiles); err != nil {
		fmt.Fprintf(w, "タスクの取得に失敗しました")
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</li>")

	memberFiles, err := getMemberFile()
	if err != nil {
		fmt.Fprintf(w, "タスクの取得に失敗しました")
	}

	fmt.Fprintf(w, "<li id=\"R\">")
	fmt.Fprintf(w, "<ul id=\"member-list\">")
	if err := render(w, r, memberFiles); err != nil {
		fmt.Fprintf(w, "タスクの取得に失敗しました")
	}
	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</li>")

	fmt.Fprintf(w, "</ul>")
	fmt.Fprintf(w, "</body>")
	fmt.Fprintf(w, "</html>")
}

func add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	target := r.Form.Get("target")
	task := r.Form.Get("task")

	file, err := os.OpenFile("data/"+target, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Fprintf(w, "タスク追加にエラーが発生しました")
	}
	defer file.Close()

	fmt.Fprintln(file, task)
}

func remove(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	target := r.Form.Get("target")
	displayStr := r.Form.Get("task")
	display, err := strconv.Atoi(displayStr)
	if err != nil {
		fmt.Fprintf(w, "タスク削除にエラーが発生しました")
	}

	file, err := os.OpenFile("data/"+target, os.O_RDWR, 0600)
	if err != nil {
		fmt.Fprintf(w, "タスク削除にエラーが発生しました")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "タスク削除にエラーが発生しました")
	}

	lines := strings.Split(string(data), "\n")

	if display >= len(lines) || display < 1 {
		fmt.Fprintf(w, "指定されたタスクは存在しません")
		return
	}

	task := display - 1
	new_lines := append(lines[:task], lines[task+1:]...)
	new_data := strings.Join(new_lines, "\n")

	ioutil.WriteFile("data/"+target, []byte(new_data), 0600)
}

func getTodoFile() ([]fs.DirEntry, error) {
	files, err := os.ReadDir("data")
	if err != nil {
		return files, err
	}

	result := []fs.DirEntry{}

	for _, v := range files {
		if v.Name() == "決定" {
			result = append(result, v)
		}
	}

	for _, v := range files {
		if v.Name() == "未決定" {
			result = append(result, v)
		}
	}

	return result, err
}

func getMemberFile() ([]fs.DirEntry, error) {
	files, err := os.ReadDir("data")
	if err != nil {
		return files, err
	}

	result := []fs.DirEntry{}

	for i, v := range files {
		if v.Name() == "決定" || v.Name() == "未決定" {
			result = append(files[:i], files[i+1:]...)
			break
		}
	}

	for i, v := range files {
		if v.Name() == "決定" || v.Name() == "未決定" {
			result = append(result[:i], result[i+1:]...)
			break
		}
	}

	return result, err
}

func render(w http.ResponseWriter, r *http.Request, files []fs.DirEntry) error {
	for _, v := range files {
		fmt.Fprintf(w, "<li class=\"todo\">")
		fmt.Fprintf(w, "<div class=\"name\">")
		fmt.Fprintf(w, v.Name())
		fmt.Fprintf(w, "</div>")

		file, err := os.Open("data/" + v.Name())
		if err != nil {
			return err
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		lines := strings.Split(string(data), "\n")

		fmt.Fprintf(w, "<ol class=\"task-list\">")
		for i, str := range lines {
			display := i + 1
			if display == len(lines) {
				break
			}

			fmt.Fprintf(w, "<li class=\"task\">")
			fmt.Fprintf(w, str)
			fmt.Fprintf(w, "</li>")
		}
		fmt.Fprintf(w, "</ol>")
		fmt.Fprintf(w, "</li>")
	}

	return nil
}

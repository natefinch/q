package q

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

const templ = "{{.Id}}:\t{{.Title}}\n"

var (
	out      = template.Must(template.New("out").Parse(templ))
	filename = filepath.Join(os.Getenv("HOME"), ".config", "q")
)

type Todo struct {
	Title string `toml:"title"`
	Id    int    `toml:"id"`
}

type todos struct {
	Todos []Todo `toml:"todo"`
}

func Add(t Todo) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	ts := todos{}
	m, err := toml.DecodeReader(f, &ts)
	if err != nil {
		return fmt.Errorf("Error reading todos from %q: %s", filename, err)
	}
	if len(m.Undecoded()) > 0 {
		return fmt.Errorf("Unknown data in %s: %v", filename, m.Undecoded())
	}

	id := 0
	for _, t := range ts.Todos {
		if t.Id > id {
			id = t.Id
		}
	}
	t.Id = id + 1
	ts.Todos = append(ts.Todos, t)
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("Error rewriting todo file: %s", err)
	}

	if err := toml.NewEncoder(f).Encode(ts); err != nil {
		return fmt.Errorf("Error writing todos: %s", err)
	}
	fmt.Printf("Task %d created.", t.Id)
	return nil
}

func List() error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	ts := todos{}
	m, err := toml.DecodeReader(f, &ts)
	if err != nil {
		return fmt.Errorf("Error reading todos from %q: %s", filename, err)
	}
	if len(m.Undecoded()) > 0 {
		return fmt.Errorf("Unknown data in %s: %v", filename, m.Undecoded())
	}

	for _, t := range ts.Todos {
		if err := out.Execute(os.Stdout, t); err != nil {
			return fmt.Errorf("Error rendering todo: %s", err)
		}
	}
	return nil
}

func Delete(id int) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	ts := todos{}
	m, err := toml.DecodeReader(f, &ts)
	if err != nil {
		return fmt.Errorf("Error reading todos from %q: %s", filename, err)
	}
	if len(m.Undecoded()) > 0 {
		return fmt.Errorf("Unknown data in %s: %v", filename, m.Undecoded())
	}
	l := len(ts.Todos)
	for n, t := range ts.Todos {
		if t.Id == id {
			if n == l-1 {
				ts.Todos = ts.Todos[:n]
			} else {
				ts.Todos = append(ts.Todos[:n], ts.Todos[n+1:]...)
			}
		}
	}

	if l == len(ts.Todos) {
		return fmt.Errorf("Task %d not found.", id)
	}

	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("Error rewriting todo file: %s", err)
	}

	if err := toml.NewEncoder(f).Encode(ts); err != nil {
		return fmt.Errorf("Error writing todos: %s", err)
	}
	return nil
}

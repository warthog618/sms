package main

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	switch os.Getenv("TEST_ENTRY") {
	case "main":
		main()
	default:
		os.Exit(m.Run())
	}
}

type count struct {
	Encoding  string
	Messages  int
	Tlen      int
	Llen      int
	Pdulen    int
	Remaining int
}

type testPattern struct {
	name    string
	msg     string
	success bool
	out     count
}

var outTemplate = `encoding: {{.Encoding}}
messages: {{.Messages}}
total length: {{.Tlen}}
last PDU length: {{.Llen}}
per_message: {{.Pdulen}}
remaining: {{.Remaining}}
`

func TestExec(t *testing.T) {
	env := []string{}
	tmpl, err := template.New("test").Parse(outTemplate)
	if err != nil {
		panic(err)
	}
	for _, p := range tests {
		f := func(t *testing.T) {
			cmd := exec.Command(os.Args[0])
			cmd.Args = append([]string{"smscounter", "-message"}, p.msg)
			cmd.Env = append(env, "TEST_ENTRY=main")
			stdout, err := cmd.Output()
			if e, ok := err.(*exec.ExitError); ok {
				assert.Equal(t, p.success, e.Success())
			} else {
				var out bytes.Buffer
				err := tmpl.Execute(&out, p.out)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, out.String(), string(stdout))
			}
		}
		t.Run(p.name, f)
	}
}

var tests = []testPattern{
	{"std", "content of the SMS", true, count{"7BIT", 1, 18, 18, 160, 142}},
}

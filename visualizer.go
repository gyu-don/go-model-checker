package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func (wld world) label() string {
	strs := []string{}

	vnames := []string{}
	for name := range wld.environment.variables {
		vnames = append(vnames, string(name))
	}
	sort.Strings(vnames)
	for _, name := range vnames {
		val := wld.environment.variables[varName(name)]
		strs = append(strs, fmt.Sprintf("%s=%d", name, val))
	}
	lnames := []string{}
	for name := range wld.environment.locks {
		lnames = append(lnames, string(name))
	}
	sort.Strings(lnames)
	for _, name := range lnames {
		pname := wld.environment.locks[lockName(name)]
		strs = append(strs, fmt.Sprintf("%s[%s]", name, pname))
	}
	return strings.Join(strs, "\n")
}

func (model kripkeModel) WriteAsDot(w io.Writer, result *verificationResult) {
	fmt.Fprintln(w, "digraph {")
	for id, wld := range model.worlds {
		fmt.Fprintf(w, "    %d [ label = \"%s\" ];\n", id, wld.label())
		if id == model.initial {
			fmt.Fprintf(w, "    %d [ penwidth = 5 ];\n", id)
		}
		if result != nil && result.targets.member(id) {
			fmt.Fprintf(w, "    %d [ style = filled, fillcolor = gray ]; \n", id)
		}
	}
	for from, tos := range model.accessible {
		for _, to := range tos {
			fmt.Fprintf(w, "    %d -> %d;\n", from, to)
		}
	}
	fmt.Fprintln(w, "}")
}

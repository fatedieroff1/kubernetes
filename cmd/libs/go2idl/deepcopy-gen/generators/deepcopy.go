/*
Copyright 2015 The Kubernetes Authors.

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

package generators

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"k8s.io/kubernetes/cmd/libs/go2idl/args"
	"k8s.io/kubernetes/cmd/libs/go2idl/generator"
	"k8s.io/kubernetes/cmd/libs/go2idl/namer"
	"k8s.io/kubernetes/cmd/libs/go2idl/types"
	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/golang/glog"
)

// Constraints is a set of optional limitations on what deep copy will generate.
type Constraints struct {
	// PackageConstraints is an optional set of package prefixes that constrain which types
	// will have inline deep copy methods generated for. Any type outside of these packages
	// (if specified) will not have a function generated and will result in a call to the
	// cloner.DeepCopy method.
	PackageConstraints []string
}

// TODO: This is created only to reduce number of changes in a single PR.
// Remove it and use PublicNamer instead.
func deepCopyNamer() *namer.NameStrategy {
	return &namer.NameStrategy{
		Join: func(pre string, in []string, post string) string {
			return strings.Join(in, "_")
		},
		PrependPackageNames: 1,
	}
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": deepCopyNamer(),
		"raw":    namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		glog.Fatalf("Failed loading boilerplate: %v", err)
	}

	initInputs := sets.NewString()
	explicitInputs := sets.NewString()
	inputs := sets.NewString()
	for _, s := range arguments.InputDirs {
		switch {
		case strings.HasPrefix(s, "+"):
			// packages with '+' prefix get functions generated for everything except gencopy=false, but
			// no init function
			s = strings.TrimPrefix(s, "+")
			inputs.Insert(s)
		case strings.HasPrefix(s, "-"):
			// packages with '-' prefix only get functions generated for those with gencopy=true
			s = strings.TrimPrefix(s, "-")
			inputs.Insert(s)
			explicitInputs.Insert(s)
		default:
			inputs.Insert(s)
			initInputs.Insert(s)
		}
	}

	var restrictRange []string
	if c, ok := arguments.CustomArgs.(Constraints); ok {
		restrictRange = c.PackageConstraints
	}

	packages := generator.Packages{}
	header := append([]byte(
		`// +build !ignore_autogenerated

`), boilerplate...)
	header = append(header, []byte(
		`
// This file was autogenerated by deepcopy-gen. Do not edit it manually!

`)...)
	for _, p := range context.Universe {
		copyableType := false
		for _, t := range p.Types {
			if copyableWithinPackage(t, explicitInputs.Has(t.Name.Package)) && inputs.Has(t.Name.Package) {
				copyableType = true
			}
		}
		if copyableType {
			// TODO: replace this with a more sophisticated algorithm that generates private copy methods
			// (like auto_DeepCopy_...) for any type that is outside of the PackageConstraints. That would
			// avoid having to make a reflection call.
			canInlineTypeFn := func(c *generator.Context, t *types.Type) bool {
				// types must be public structs or have a custom DeepCopy_<method> already defined
				if !copyableWithinPackage(t, explicitInputs.Has(t.Name.Package)) && !publicCopyFunctionDefined(c, t) {
					return false
				}

				// only packages within the restricted range can be inlined
				for _, s := range restrictRange {
					if strings.HasPrefix(t.Name.Package, s) {
						return true
					}
				}
				return false
			}

			path := p.Path
			packages = append(packages,
				&generator.DefaultPackage{
					PackageName: strings.Split(filepath.Base(path), ".")[0],
					PackagePath: path,
					HeaderText:  header,
					GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
						generators = []generator.Generator{}
						generators = append(
							generators, NewGenDeepCopy("deep_copy_generated", path, initInputs.Has(path), explicitInputs.Has(path), canInlineTypeFn))
						return generators
					},
					FilterFunc: func(c *generator.Context, t *types.Type) bool {
						return t.Name.Package == path
					},
				})
		}
	}
	return packages
}

// CanInlineTypeFunc should return true if the provided type can be converted to a function call
type CanInlineTypeFunc func(*generator.Context, *types.Type) bool

const (
	apiPackagePath        = "k8s.io/kubernetes/pkg/api"
	conversionPackagePath = "k8s.io/kubernetes/pkg/conversion"
)

// genDeepCopy produces a file with autogenerated deep-copy functions.
type genDeepCopy struct {
	generator.DefaultGen

	targetPackage      string
	imports            namer.ImportTracker
	typesForInit       []*types.Type
	generateInitFunc   bool
	requireExplicitTag bool
	canInlineTypeFn    CanInlineTypeFunc

	context *generator.Context

	globalVariables map[string]interface{}
}

func NewGenDeepCopy(sanitizedName, targetPackage string, generateInitFunc, requireExplicitTag bool, canInlineTypeFn CanInlineTypeFunc) generator.Generator {
	return &genDeepCopy{
		DefaultGen: generator.DefaultGen{
			OptionalName: sanitizedName,
		},
		targetPackage:      targetPackage,
		imports:            generator.NewImportTracker(),
		typesForInit:       make([]*types.Type, 0),
		generateInitFunc:   generateInitFunc,
		requireExplicitTag: requireExplicitTag,
		canInlineTypeFn:    canInlineTypeFn,
	}
}

func (g *genDeepCopy) Namers(c *generator.Context) namer.NameSystems {
	// Have the raw namer for this file track what it imports.
	return namer.NameSystems{"raw": namer.NewRawNamer(g.targetPackage, g.imports)}
}

func (g *genDeepCopy) Filter(c *generator.Context, t *types.Type) bool {
	// Filter out all types not copyable within the package.
	copyable := copyableWithinPackage(t, g.requireExplicitTag)
	if copyable {
		g.typesForInit = append(g.typesForInit, t)
	}
	return copyable
}

// publicCopyFunctionDefined returns true if a DeepCopy function has already been defined in a given
// package, which allows more efficient deep copy implementations to be defined by the caller.
func publicCopyFunctionDefined(c *generator.Context, t *types.Type) bool {
	p, ok := c.Universe[t.Name.Package]
	if !ok {
		return false
	}
	return p.Functions["DeepCopy_"+path.Base(t.Name.Package)+"_"+t.Name.Name] != nil
}

func copyableWithinPackage(t *types.Type, explicitCopyRequired bool) bool {
	tag := types.ExtractCommentTags("+", t.CommentLines)["gencopy"]
	if tag == "false" {
		return false
	}
	if explicitCopyRequired && tag != "true" {
		return false
	}
	// TODO: Consider generating functions for other kinds too.
	if t.Kind != types.Struct {
		return false
	}
	// Also, filter out private types.
	if namer.IsPrivateGoName(t.Name.Name) {
		return false
	}
	return true
}

func (g *genDeepCopy) isOtherPackage(pkg string) bool {
	if pkg == g.targetPackage {
		return false
	}
	if strings.HasSuffix(pkg, "\""+g.targetPackage+"\"") {
		return false
	}
	return true
}

func (g *genDeepCopy) Imports(c *generator.Context) (imports []string) {
	importLines := []string{}
	for _, singleImport := range g.imports.ImportLines() {
		if g.isOtherPackage(singleImport) {
			importLines = append(importLines, singleImport)
		}
	}
	return importLines
}

func (g *genDeepCopy) withGlobals(args map[string]interface{}) map[string]interface{} {
	for k, v := range g.globalVariables {
		if _, ok := args[k]; !ok {
			args[k] = v
		}
	}
	return args
}

func argsFromType(t *types.Type) map[string]interface{} {
	return map[string]interface{}{
		"type": t,
	}
}

func (g *genDeepCopy) funcNameTmpl(t *types.Type) string {
	tmpl := "DeepCopy_$.type|public$"
	g.imports.AddType(t)
	if t.Name.Package != g.targetPackage {
		tmpl = g.imports.LocalNameOf(t.Name.Package) + "." + tmpl
	}
	return tmpl
}

func (g *genDeepCopy) Init(c *generator.Context, w io.Writer) error {
	g.context = c
	cloner := c.Universe.Type(types.Name{Package: conversionPackagePath, Name: "Cloner"})
	g.imports.AddType(cloner)
	g.globalVariables = map[string]interface{}{
		"Cloner": cloner,
	}
	if !g.generateInitFunc {
		// TODO: We should come up with a solution to register all generated
		// deep-copy functions. However, for now, to avoid import cycles
		// we register only those explicitly requested.
		return nil
	}
	scheme := c.Universe.Variable(types.Name{Package: apiPackagePath, Name: "Scheme"})
	g.imports.AddType(scheme)
	g.globalVariables["scheme"] = scheme

	sw := generator.NewSnippetWriter(w, c, "$", "$")
	sw.Do("func init() {\n", nil)
	sw.Do("if err := $.scheme|raw$.AddGeneratedDeepCopyFuncs(\n", map[string]interface{}{
		"scheme": scheme,
	})
	for _, t := range g.typesForInit {
		sw.Do(fmt.Sprintf("%s,\n", g.funcNameTmpl(t)), argsFromType(t))
	}
	sw.Do("); err != nil {\n", nil)
	sw.Do("// if one of the deep copy functions is malformed, detect it immediately.\n", nil)
	sw.Do("panic(err)\n", nil)
	sw.Do("}\n", nil)
	sw.Do("}\n\n", nil)
	return sw.Error()
}

func (g *genDeepCopy) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	funcName := g.funcNameTmpl(t)
	sw.Do(fmt.Sprintf("func %s(in $.type|raw$, out *$.type|raw$, c *$.Cloner|raw$) error {\n", funcName), g.withGlobals(argsFromType(t)))
	g.generateFor(t, sw)
	sw.Do("return nil\n", nil)
	sw.Do("}\n\n", nil)
	return sw.Error()
}

// we use the system of shadowing 'in' and 'out' so that the same code is valid
// at any nesting level. This makes the autogenerator easy to understand, and
// the compiler shouldn't care.
func (g *genDeepCopy) generateFor(t *types.Type, sw *generator.SnippetWriter) {
	var f func(*types.Type, *generator.SnippetWriter)
	switch t.Kind {
	case types.Builtin:
		f = g.doBuiltin
	case types.Map:
		f = g.doMap
	case types.Slice:
		f = g.doSlice
	case types.Struct:
		f = g.doStruct
	case types.Interface:
		f = g.doInterface
	case types.Pointer:
		f = g.doPointer
	case types.Alias:
		f = g.doAlias
	default:
		f = g.doUnknown
	}
	f(t, sw)
}

func (g *genDeepCopy) doBuiltin(t *types.Type, sw *generator.SnippetWriter) {
	sw.Do("*out = in\n", nil)
}

func (g *genDeepCopy) doMap(t *types.Type, sw *generator.SnippetWriter) {
	sw.Do("*out = make($.|raw$)\n", t)
	if t.Key.IsAssignable() {
		sw.Do("for key, val := range in {\n", nil)
		if t.Elem.IsAssignable() {
			sw.Do("(*out)[key] = val\n", nil)
		} else {
			if g.canInlineTypeFn(g.context, t.Elem) {
				sw.Do("newVal := new($.|raw$)\n", t.Elem)
				funcName := g.funcNameTmpl(t.Elem)
				sw.Do(fmt.Sprintf("if err := %s(val, newVal, c); err != nil {\n", funcName), argsFromType(t.Elem))
				sw.Do("return err\n", nil)
				sw.Do("}\n", nil)
				sw.Do("(*out)[key] = *newVal\n", nil)
			} else {
				sw.Do("if newVal, err := c.DeepCopy(val); err != nil {\n", nil)
				sw.Do("return err\n", nil)
				sw.Do("} else {\n", nil)
				sw.Do("(*out)[key] = newVal.($.|raw$)\n", t.Elem)
				sw.Do("}\n", nil)
			}
		}
	} else {
		// TODO: Implement it when necessary.
		sw.Do("for range in {\n", nil)
		sw.Do("// FIXME: Copying unassignable keys unsupported $.|raw$\n", t.Key)
	}
	sw.Do("}\n", nil)
}

func (g *genDeepCopy) doSlice(t *types.Type, sw *generator.SnippetWriter) {
	sw.Do("*out = make($.|raw$, len(in))\n", t)
	if t.Elem.Kind == types.Builtin {
		sw.Do("copy(*out, in)\n", nil)
	} else {
		sw.Do("for i := range in {\n", nil)
		if t.Elem.IsAssignable() {
			sw.Do("(*out)[i] = in[i]\n", nil)
		} else if g.canInlineTypeFn(g.context, t.Elem) {
			funcName := g.funcNameTmpl(t.Elem)
			sw.Do(fmt.Sprintf("if err := %s(in[i], &(*out)[i], c); err != nil {\n", funcName), argsFromType(t.Elem))
			sw.Do("return err\n", nil)
			sw.Do("}\n", nil)
		} else {
			sw.Do("if newVal, err := c.DeepCopy(in[i]); err != nil {\n", nil)
			sw.Do("return err\n", nil)
			sw.Do("} else {\n", nil)
			sw.Do("(*out)[i] = newVal.($.|raw$)\n", t.Elem)
			sw.Do("}\n", nil)
		}
		sw.Do("}\n", nil)
	}
}

func (g *genDeepCopy) doStruct(t *types.Type, sw *generator.SnippetWriter) {
	for _, m := range t.Members {
		t := m.Type
		if t.Kind == types.Alias {
			copied := *t.Underlying
			copied.Name = t.Name
			t = &copied
		}
		args := map[string]interface{}{
			"type": t,
			"name": m.Name,
		}
		switch t.Kind {
		case types.Builtin:
			sw.Do("out.$.name$ = in.$.name$\n", args)
		case types.Map, types.Slice, types.Pointer:
			sw.Do("if in.$.name$ != nil {\n", args)
			sw.Do("in, out := in.$.name$, &out.$.name$\n", args)
			g.generateFor(t, sw)
			sw.Do("} else {\n", nil)
			sw.Do("out.$.name$ = nil\n", args)
			sw.Do("}\n", nil)
		case types.Struct:
			if g.canInlineTypeFn(g.context, t) {
				funcName := g.funcNameTmpl(t)
				sw.Do(fmt.Sprintf("if err := %s(in.$.name$, &out.$.name$, c); err != nil {\n", funcName), args)
				sw.Do("return err\n", nil)
				sw.Do("}\n", nil)
			} else {
				sw.Do("if newVal, err := c.DeepCopy(in.$.name$); err != nil {\n", args)
				sw.Do("return err\n", nil)
				sw.Do("} else {\n", nil)
				sw.Do("out.$.name$ = newVal.($.type|raw$)\n", args)
				sw.Do("}\n", nil)
			}
		default:
			sw.Do("if in.$.name$ == nil {\n", args)
			sw.Do("out.$.name$ = nil\n", args)
			sw.Do("} else if newVal, err := c.DeepCopy(in.$.name$); err != nil {\n", args)
			sw.Do("return err\n", nil)
			sw.Do("} else {\n", nil)
			sw.Do("out.$.name$ = newVal.($.type|raw$)\n", args)
			sw.Do("}\n", nil)
		}
	}
}

func (g *genDeepCopy) doInterface(t *types.Type, sw *generator.SnippetWriter) {
	// TODO: Add support for interfaces.
	g.doUnknown(t, sw)
}

func (g *genDeepCopy) doPointer(t *types.Type, sw *generator.SnippetWriter) {
	sw.Do("*out = new($.Elem|raw$)\n", t)
	if t.Elem.IsAssignable() {
		sw.Do("**out = *in", nil)
	} else if g.canInlineTypeFn(g.context, t.Elem) {
		funcName := g.funcNameTmpl(t.Elem)
		sw.Do(fmt.Sprintf("if err := %s(*in, *out, c); err != nil {\n", funcName), argsFromType(t.Elem))
		sw.Do("return err\n", nil)
		sw.Do("}\n", nil)
	} else {
		sw.Do("if newVal, err := c.DeepCopy(*in); err != nil {\n", nil)
		sw.Do("return err\n", nil)
		sw.Do("} else {\n", nil)
		sw.Do("**out = newVal.($.|raw$)\n", t.Elem)
		sw.Do("}\n", nil)
	}
}

func (g *genDeepCopy) doAlias(t *types.Type, sw *generator.SnippetWriter) {
	// TODO: Add support for aliases.
	g.doUnknown(t, sw)
}

func (g *genDeepCopy) doUnknown(t *types.Type, sw *generator.SnippetWriter) {
	sw.Do("// FIXME: Type $.|raw$ is unsupported.\n", t)
}

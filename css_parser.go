package css_parser

import (
	"sync"
	"strings"
)

type CSSKeyVals map[string]string

type CSS struct {
	Mu sync.RWMutex
	Element map[string] CSSKeyVals
	ID map[string] CSSKeyVals
	Class map[string] CSSKeyVals
}

func NewCSS (s string) *CSS {
	css := new(CSS)
	css.Element = make(map[string]CSSKeyVals)
	css.ID = make(map[string]CSSKeyVals)
	css.Class = make(map[string]CSSKeyVals)

	css.Parse(s)

	return css
}

func (css CSS) Parse (s string) {
	s1 := strings.Split(s, "}")
	for _,e := range s1 {
		keyvals := strings.Split(e, "{")
		if len(keyvals) != 2 {
			continue
		}

		key := strings.TrimSpace(keyvals[0])
		classification := classify(key)
		switch classification {
		case 1:
			css.Mu.Lock()
			css.Class[key] = makeKeyValMap(keyvals[1])
			css.Mu.Unlock()
		case 2:
			css.Mu.Lock()
			css.ID[key] = makeKeyValMap(keyvals[1])
			css.Mu.Unlock()
		case 3:
			css.Mu.Lock()
			css.Element[key] = makeKeyValMap(keyvals[1])
			css.Mu.Unlock()
		default:
			continue
		}
	}
}

func (css CSS) GetElement (element, property string) string {
	css.Mu.RLock()
	value, ok := css.Element[element][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func (css CSS) GetClass (class, property string) string {
	css.Mu.RLock()
	value, ok := css.Class[class][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func (css CSS) GetID (id, property string) string {
	css.Mu.RLock()
	value, ok := css.ID[id][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func classify (s string) int {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 0
	}

	switch s[0] {
	case 46:
		return 1
	case 35:
		return 2
	default:
		return 3
	}
}

func makeKeyValMap (s string) map[string]string {
	kvMap := make(map[string]string)
	a1 := strings.Split(s, ";")
	for _,e := range a1 {
		val := removeComments(e)
		a2 := strings.Split(val, ":")
		if len(a2) != 2 {
			continue
		}

		kvMap[strings.TrimSpace(a2[0])] = strings.TrimSpace(a2[1])
	}

	return kvMap 
}

func removeComments (s string) string {
	start := strings.Index(s, "/*")
	if start == -1 || len(s) < 4 {
		return s
	}

	end := strings.Index(s, "*/")
	s = s[:start] + s[end+2:]

	return removeComments(s)
}
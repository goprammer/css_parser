package css_parser

import (
	"sync"
	"strconv"
	"strings"
)

type CSSKeyVals map[string]string

type Condition struct {
	MinWidth int
	MaxWidth int
	ConditionalCSS *CSS
}

type CSS struct {
	Mu sync.RWMutex
	Element map[string] CSSKeyVals
	ID map[string] CSSKeyVals
	Class map[string] CSSKeyVals
	MediaQueries []*Condition
}

func NewCSS (s string) *CSS {
	css := new(CSS)
	css.Element = make(map[string]CSSKeyVals)
	css.ID = make(map[string]CSSKeyVals)
	css.Class = make(map[string]CSSKeyVals)
	css.MediaQueries = make([]*Condition, 0)
	css.Parse(s)
	return css
}

func NewMediaQuery (head, body string) *Condition {
	condition := new(Condition)
	i1 := strings.Index(head, "(")
	i2 := strings.Index(head, ")")
	for i1 != -1 && i2 != -1 {
		rule := head[i1+1:i2]
		arr := strings.Split(rule, ":")
		if len(arr) == 2 {
			key := strings.ToLower(strings.TrimSpace(arr[0]))
			val := strings.TrimSpace(arr[1])
			if key == "min-width" {
				condition.MinWidth = extractNumbers(val)
			} else if key == "max-width" {
				condition.MaxWidth = extractNumbers(val)
			}
		}
		
		head = head[i2+1:]
		i1 = strings.Index(head, "(")
		i2 = strings.Index(head, ")")
	}

	condition.ConditionalCSS = NewCSS(body)

	return condition
}

func (css *CSS) Parse (s string) {
	s1 := extractNamespaces(s)
	for _,e := range s1 {
		keyvals := strings.Split(e, "{")
		if len(keyvals) < 2 {
			continue
		}

		if len(keyvals) == 3 && strings.Contains(keyvals[0], "@media") {
			css.MediaQueries = append(css.MediaQueries, NewMediaQuery(keyvals[0], keyvals[1] + "{" + keyvals[2]))
			continue
		}
		
		key := strings.TrimSpace(removeComments(keyvals[0]))
		
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

func (css *CSS) getID (id, property string) string {
	css.Mu.RLock()
	value, ok := css.ID[id][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func (css *CSS) getClass (class, property string) string {
	css.Mu.RLock()
	value, ok := css.Class[class][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func (css *CSS) getElement (element, property string) string {
	css.Mu.RLock()
	value, ok := css.Element[element][property]
	css.Mu.RUnlock()
	if !ok {
		return ""
	}

	return value
}

func (css *CSS) getNormal (id, class, element, property string) string {
	if id != "" {
		return css.getID(id, property)
	} else if class != "" {
		return css.getClass(class, property)
	} else if element != "" {
		return css.getElement(element, property)
	}

	return ""
}

func (css *CSS) Get (id, class, element, property, width string) string {
	if width != "" {
		answer := ""
		nWidth := extractNumbers(width)

		for _,e := range css.MediaQueries {
			if e.MaxWidth > nWidth && e.MinWidth < nWidth {
				answer = e.ConditionalCSS.Get(id, class, element, property, "")
			}
		}

		if answer != "" {
			return answer
		}
	}

	return css.getNormal(id, class, element, property)
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

		property := strings.TrimSpace(a2[0])
		pVal := strings.TrimSpace(a2[1])
		if property == "padding" {
			undoShorthand(kvMap, property, pVal)
		} else if property == "margin" {
			undoShorthand(kvMap, property, pVal)
		} else {
			kvMap[property] = pVal
		}
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

func undoShorthand (kvMap map[string]string, shorthandProperty, val string) {
	if len(val) < 2 {
		return
	}

	if val[len(val)-1] == 59 {
		val = val[:len(val)-1]
	}
	
	important := ""
	i := strings.Index(val, "!important")
	if i != -1 {
		val = val[:i] + val[i+10:]
		important = " !important"
	}
	
	arr1 := strings.Split(val, " ")
	arr2 := make([]string, 0)
	for _,e := range arr1 {
		if e != "" {
			arr2 = append(arr2, strings.TrimSpace(e))
		}
	}

	l := len(arr2)
	if l == 0 || l > 4 {
		return
	}

	switch l {
	case 1:
		arr2 = append(arr2, arr1[0])
		arr2 = append(arr2, arr1[0])
		arr2 = append(arr2, arr1[0])
	case 2:
		arr2 = append(arr2, arr1[0])
		arr2 = append(arr2, arr1[1])
	case 3:
		arr2 = append(arr2, arr1[1])
	}

	kvMap[shorthandProperty + "-top"] = arr2[0] + important
	kvMap[shorthandProperty + "-right"] = arr2[1] + important
	kvMap[shorthandProperty + "-bottom"] = arr2[2] + important
	kvMap[shorthandProperty + "-left"] = arr2[3] + important
}

func extractNumbers (s string) int {
	res := new(strings.Builder)
	units := new(strings.Builder)
	for i := 0; i < len(s); i++ {
		if s[i] > 47 && s[i] < 58 {
			res.WriteByte(s[i])	
		} else {
			units.WriteByte(s[i])
		}
	}
	
	n, err := strconv.Atoi(res.String())
	if err != nil {
		return 0
	}

	return n
}

func extractNamespace (s string) (int, int) {
	depth := 0
	i1 := -1
	i2 := -1

	for i,e := range s {
		if e == 123 {
			if depth == 0 {
				i1 = i
			}

			depth++
		} else if e == 125 {
			if depth == 1 {
				i2 = i
				break
			}

			depth--
		}
	}

	return i1, i2
}

func extractNamespaces (s string) []string {
	answer := make([]string, 0)
	i2 := len(s)
	for {
		_, i2 = extractNamespace(s)
		if i2 == -1 {
			break
		}

		answer = append(answer, s[:i2])
		s = s[i2+1:]
	}

	return answer
}
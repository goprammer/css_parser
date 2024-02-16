package css_parser

import (
	"fmt"
	"sync"
	"strconv"
	"strings"
)

// Build a square bracket attribute matching system as well.
// Test on what media queries with % are doing.

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

func NewCSS () *CSS {
	css := new(CSS)
	css.Element = make(map[string]CSSKeyVals)
	css.ID = make(map[string]CSSKeyVals)
	css.Class = make(map[string]CSSKeyVals)
	css.MediaQueries = make([]*Condition, 0)
	
	return css
}

func NewMediaQuery (head, namespace string) *Condition {
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

	condition.ConditionalCSS = NewCSS()
	condition.ConditionalCSS.Parse(removeCurlyBrackets(namespace))

	return condition
}

func (css *CSS) Parse (s string) {
	namespaces := make([]string, 0)
	for {
		_, i2 := extractNamespace(s)
		if i2 == -1 {
			break
		}

		namespaces = append(namespaces, s[:i2+1])
		s = s[i2+1:]
	}

	for _,e := range namespaces {
		i := strings.Index(e, "{")
		if i == -1 {
			continue
		}

		selector := strings.TrimSpace(removeComments(e[:i]))
		namespace := e[i:]

		if strings.Contains(selector, "@media") {
			// selector is actually a media query
			css.MediaQueries = append(css.MediaQueries, NewMediaQuery(selector, namespace))
			continue
		}
		
		selectors := splitSelectorString(selector)

		for _,e := range selectors {
			css.AppendKeyVals(e, namespace)
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
	answer := ""
	if id != "" {
		if answer = css.getID(removeClassifier(id), property); answer != "" {
			return answer
		}
	}

	if class != "" {
		if answer = css.getClass(removeClassifier(class), property); answer != "" {
			return answer
		}
	}

	if element != "" {
		if answer = css.getElement(element, property); answer != "" {
			return answer
		}
	}

	return answer
}

func (css *CSS) Get (id, class, element, property, width string) string {
	if width != "" {
		answer := ""
		nWidth := extractNumbers(width)

		for _,e := range css.MediaQueries {
			if e.MaxWidth >= nWidth && e.MinWidth <= nWidth {
				if possibleAnswer := e.ConditionalCSS.Get(id, class, element, property, ""); possibleAnswer != "" {
					answer = possibleAnswer
				}
			}
		}

		if answer != "" {
			return answer
		}
	}

	return css.getNormal(id, class, element, property)
}

func (css *CSS) PrintAll ()  {
	css.Mu.RLock()
	defer css.Mu.RUnlock()
	printType(css.ID)
	printType(css.Class)
	printType(css.Element)
	for _,e := range css.MediaQueries {
		fmt.Println("\n@media min-width:", e.MinWidth, " max-width:",e.MaxWidth)
		e.ConditionalCSS.PrintAll()
	}
}

func printType (cssType map[string]CSSKeyVals) {
	for i,e := range cssType {
		fmt.Println("-----", i, "-----")
		for j,f := range e {
			fmt.Println(j,":",f)
		}
		fmt.Println()
	}
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

// splitSelectorString() currently only reads simple selectors. 
// Ex simple: #idName

// splitSelectorString() ignores attribute, compound and combinator selectors.
// Ex attribute: [href^="https"] (only tags with href attribute whose value starts with "https")
// Ex compound: div.className (only div elements with class="className")
// Ex combinator: .footer td (only td element within footer class)

// See full reference:
// https://www.w3schools.com/cssref/css_selectors.php

func splitSelectorString (s string) []string {
	result := make([]string, 0)
	selectors := strings.Split(s, ",")
	if len(selectors) == 1 {
		return selectors
	}

	for i := 0; i < len(selectors); i++ {
		selectors[i] = removeCommas(strings.TrimSpace(selectors[i]))
		if selectors[i] != "" {
			result = append(result, selectors[i])
		}
	}

	return result
}

func removeClassifier (s string) string {
	if len(s) == 0 {
		return ""
	}

	if s[0] == 46 || s[0] == 35 {
		return s[1:]
	}

	return s
}

func removeCommas (s string) string {
	res := new(strings.Builder)
	for i := 0; i < len(s); i++ {
		if s[i] != 44 {
			res.WriteByte(s[i])
		}
	}

	return res.String()
}

func (css *CSS) AppendKeyVals (selector, keyvalStr string) {
	kvMap := make(map[string]string)
	
	switch classify(selector) {
	case 1:
		css.Mu.Lock()
		tmp_kv_Map, ok := css.Class[removeClassifier(selector)]
		if ok {
			kvMap = tmp_kv_Map
		}
		css.Class[removeClassifier(selector)] = makeKeyValMap(kvMap, keyvalStr)
		css.Mu.Unlock()
	case 2:
		css.Mu.Lock()
		tmp_kv_Map, ok := css.ID[removeClassifier(selector)]
		if ok {
			kvMap = tmp_kv_Map
		}
		css.ID[removeClassifier(selector)] = makeKeyValMap(kvMap, keyvalStr)
		css.Mu.Unlock()
	case 3:
		css.Mu.Lock()
		tmp_kv_Map, ok := css.Element[removeClassifier(selector)]
		if ok {
			kvMap = tmp_kv_Map
		}
		css.Element[removeClassifier(selector)] = makeKeyValMap(kvMap, keyvalStr)
		css.Mu.Unlock()
	}
}

func makeKeyValMap (kvMap map[string]string, s string) map[string]string {
	s = removeCurlyBrackets(s)
	
	a1 := strings.Split(s, ";")
	for _,e := range a1 {
		val := strings.TrimSpace(removeComments(e))
		
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
	if len(val) == 0 {
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

func removeCurlyBrackets (s string) string {
	i1 := 0
	i2 := 0

	for i := 0; i < len(s)-1; i++ {
		if s[i] == 123 {
			i1 = i
			break
		}
	}

	for i := len(s)-1; i >=0; i-- {
		if s[i] == 125 {
			i2 = i
			break
		}
	}

	if i1 < i2 {
		return s[i1+1:i2]
	}
	
	return s	
}
package css_parser

import (
	"testing"
)

func TestParser (t *testing.T) {
	css := NewCSS()
	css.Parse(TestStr)

	if css.Get("comment", "", "img", "margin-left", "") != "21px" {
		t.Fatal("Comments were not handled correctly.")
	}

	if css.Get("", "", "img", "margin-left", "") != "13px !important" {
		t.Fatal("Important keyword was not handled correctly.")
	}

	if css.Get("img", "img", "img", "margin-left", "") != "13px" {
		t.Fatal("ID, Class, Element hierarchy was not handled correctly.")
	}

	if css.Get("", "img", "img", "margin-left", "") != "3px" {
		t.Fatal("ID, Class, Element hierarchy was not handled correctly.")
	}
	
	if css.Get("", "img", "img", "margin-left", "700px") != "11px" {
		t.Fatal("Matching media query was not handled correctly.")
	}

	if css.Get("", "img", "img", "margin-left", "900px") != "3px" {
		t.Fatal("Non-matching media query was not handled correctly. It should have reverted to standard css.")
	}

	if css.Get("", "img", "img", "margin-left", "300px") != "18px" {
		t.Fatal("2nd matching media query was not handled correctly.")
	}

	if css.Get("", "", "body", "padding-left", "") != "0" {
		t.Fatal("Original value was overwritten by a non-matching media query.")
	}

	if css.Get("", "", "body", "padding-left", "360") != "10px !important" {
		t.Fatal("Media query was not handled correctly.")
	}

	if css.Get("", "double", "div", "position", "511px") != "fixed" {
		t.Fatal("Media query with 'and' condition was not handled correctly.")
	}

	if css.Get("", "desktop_hide", "div", "display", "1234px") != "none" {
		t.Fatal("Multiple selectors were not handled correctly.")
	}

	if css.Get("laptop_hide", "", "div", "display", "1234px") != "none" {
		t.Fatal("Multiple selectors were not handled correctly ahh.")
	}

	if css.Get("", "body", "body", "padding-left", "600") != "10px !important" {
		t.Fatal("A non-matching selector failed to pass the search down to the next lower selector in hierarchy.")
	}
}
package css_parser

import (
	"testing"
)

func TestParser (t *testing.T) {
	testStr := `#comment /*Comment*/
				{
					margin: 21px
				} 

				img {
					margin: 10px 13px 15px !important;
				}

				#img{margin:13px}.img{margin-left:3px/*Another Comment*/;margin-right:1px;}

				@media (max-width:800px){
						.img{margin-left:11px
					}
				}

				@media(max-width:400px){.img{margin-left:18px}}`

	css := NewCSS(testStr)

	if css.Get("#comment", "", "img", "margin-left", "") != "21px" {
		t.Fatal("Comments were not handled correctly.")
	}

	if css.Get("", "", "img", "margin-left", "") != "13px !important" {
		t.Fatal("Important keyword was not handled correctly.")
	}

	if css.Get("#img", ".img", "img", "margin-left", "") != "13px" {
		t.Fatal("ID, Class, Element hierarchy was not handled correctly.")
	}

	if css.Get("", ".img", "img", "margin-left", "") != "3px" {
		t.Fatal("ID, Class, Element hierarchy was not handled correctly.")
	}
	
	if css.Get("", ".img", "img", "margin-left", "700px") != "11px" {
		t.Fatal("Matching media query was not handled correctly.")
	}

	if css.Get("", ".img", "img", "margin-left", "900px") != "3px" {
		t.Fatal("Non-matching media query was not handled correctly. It should have reverted to standard css.")
	}

	if css.Get("", ".img", "img", "margin-left", "300px") != "18px" {
		t.Fatal("2nd matching media query was not handled correctly.")
	}
}
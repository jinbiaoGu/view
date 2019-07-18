package view

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	Render := New(nil, "test","/var/work/javascript/website","view")

	fmt.Printf("%v \n",Render.ViewPaths)

	//request := httptest.NewRequest("GET", "/test", nil)
	//responseWriter := httptest.NewRecorder()
	//var context interface{}
	//
	//tmpl := Render.Layout("layout_for_test")
	//tmpl.Execute("test", context, request, responseWriter)
	//
	//if textproto.TrimString(responseWriter.Body.String()) != "Template for test" {
	//	t.Errorf("The template isn't rendered")
	//}
}

// anyxml - marshal an XML document from almost any Go variable
// Marshal XML from map[string]interface{}, arrays, slices, alpha/numeric, etc.
//
// Wraps xml.Marshal with functionality in github.com/clbanning/mxj to create
// a more genericized XML marshaling capability. Note: unmarshaling the resultant
// XML may not return the original value, since tag labels may have been injected
// to create the XML representation of the value.
//
// See mxj package documentation for more information.  See anyxml_test.go for
// examples or just try Xml() or XmlIndent().
/*
 Encode an arbitrary JSON object.
	package main

	import (
		"encoding/json"
		"fmt"
		"github.com/clbanning/anyxml"
	)

	func main() {
		jsondata := []byte(`[
			{ "somekey":"somevalue" },
			"string",
			3.14159265,
			true
		]`)
		var i interface{}
		err := json.Unmarshal(jsondata, &i)
		if err != nil {
			// do something
		}
		x, err := anyxml.XmlIndent(i, "", "  ", "mydoc")
		if err != nil {
			// do something else
		}
		fmt.Println(string(x))
	}

	output:
		<mydoc>
		  <somekey>somevalue</somekey>
		  <element>string</element>
		  <element>3.14159265</element>
		  <element>true</element>
		</mydoc>
*/
package anyxml

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"time"
)

// Encode arbitrary value as XML.  Note: there are no guarantees.
func Xml(v interface{}, rootTag ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.Marshal(v)
	}

	var err error
	s := new(string)
	p := new(pretty)

	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">"
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndent(false, s, tag, val, p)
					}
				} else {
					err = mapToXmlIndent(false, s, "element", vv, p)
				}
			default:
				err = mapToXmlIndent(false, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		b, err = anyxml(v.(map[string]interface{}), rootTag...)
	case []map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的func Xml() case5")
		for _, vv := range v.([]map[string]interface{}) {
			//			fmt.Println("===============================>进入anyxml.go的func Xml() case5 vv %#v", reflect.TypeOf(vv))
			b, err = anyxml(vv, rootTag...)
			ss += string(b)
			if err != nil {
				break
			}
		}
		b = []byte(ss)
	default:
		err = mapToXmlIndent(false, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

// Encode arbitrary value as XML.  Note: there are no guarantees.
func XmlWithDateFormat(dateFormat string, v interface{}, rootTag ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.Marshal(v)
	}

	var err error
	s := new(string)
	p := new(pretty)

	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">"
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						//						fmt.Println("===============================>进入anyxml.go的func Xml() case1")
						err = mapToXmlIndentWithDateFormat(dateFormat, false, s, tag, val, p)
					}
				} else {
					//					fmt.Println("===============================>进入anyxml.go的func Xml() case2")
					err = mapToXmlIndentWithDateFormat(dateFormat, false, s, "element", vv, p)
				}
			default:
				//				fmt.Println("===============================>进入anyxml.go的func Xml() case3")
				err = mapToXmlIndentWithDateFormat(dateFormat, false, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		b, err = anyxmlWithDateFormat(dateFormat, v.(map[string]interface{}), rootTag...)
	case []map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的func Xml() case5")
		for _, vv := range v.([]map[string]interface{}) {
			//			fmt.Println("===============================>进入anyxml.go的func Xml() case5 vv %#v", reflect.TypeOf(vv))
			b, err = anyxmlWithDateFormat(dateFormat, vv, "element")
			ss += (string(b) + "\n")
			if err != nil {
				break
			}
		}
		b = []byte(ss)
	default:
		err = mapToXmlIndentWithDateFormat(dateFormat, false, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

// Encode an arbitrary value as a pretty XML string. Note: there are no guarantees.
func XmlIndent(v interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
	}

	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}

	var ss string
	var b []byte

	switch v.(type) {

	case []interface{}:
		ss = "<" + rt + ">\n"
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				//				fmt.Println("===============================>进入anyxml.go的func XmlIndent() case1")
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {

						err = mapToXmlIndent(true, s, tag, val, p)
					}
				} else {
					p.start = 1 // we're 1 tag in to the doc
					err = mapToXmlIndent(true, s, "element", vv, p)
					*s += "\n"
				}
			case []map[string]interface{}:
				*s += p.padding+"<element>\n"+p.padding
				fmt.Println("===============================>进入anyxml.go的func XmlIndent() case5 vv %#v", reflect.TypeOf(vv))
				for _, vvv := range vv.([]map[string]interface{}) {
					fmt.Println("===============================>进入anyxml.go的func XmlIndent() case5 for vv %#v", reflect.TypeOf(vv))
					fmt.Println("===============================>进入anyxml.go的func XmlIndent() case5 for rootTag %#v", rootTag)

//					bs, err := anyxmlIndent(vvv, prefix, indent, "element")
//					*s += (string(bs) + "\n")
					err = mapToXmlIndent(true, s, "element", vvv, p)
					*s += "\n"
					if err != nil {
						break
					}
				}
				*s += "</element>\n"
			default:
				fmt.Println("===============================>进入anyxml.go的func XmlIndent() case1 %#v", reflect.TypeOf(vv))
				p.start = 0
				err = mapToXmlIndent(true, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的func XmlIndent() case3")
		b, err = anyxmlIndent(v.(map[string]interface{}), prefix, indent, rootTag...)
	case []map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的func XmlIndent() case5")
		for _, vv := range v.([]map[string]interface{}) {
			//			fmt.Println("===============================>进入anyxml.go的func XmlIndent() case5 vv %#v", reflect.TypeOf(vv))
			b, err = anyxmlIndent(vv, prefix, indent, rootTag...)
			ss += (string(b) + "\n")
			if err != nil {
				break
			}
		}
		b = []byte(ss)
	default:
		//		fmt.Println("===============================>进入anyxml.go的func XmlIndent() case4")
		//		fmt.Println("===============================>进入anyxml.go的func XmlIndent() v.(type) %#v", reflect.TypeOf(v))
		err = mapToXmlIndent(true, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

// Encode an arbitrary value as a pretty XML string. Note: there are no guarantees.
func XmlIndentWithDateFormat(dateFormat string, v interface{}, prefix, indent string, rootTag ...string) ([]byte, error) {
	//	fmt.Println("===============================>进入anyxml.go的func XmlIndentWithDateFormat()")
	//	fmt.Println("===============================>reflect.TypeOf(v).Kind()：%#v", reflect.TypeOf(v).Kind())

	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return xml.MarshalIndent(v, prefix, indent)
		//		fmt.Println("===============================>进入anyxml.go的if判断")

	}

	//	if reflect.TypeOf(v).Kind() == reflect.Ptr {
	//		v = Struct2MapWithDateFormat(dateFormat, v)
	//	}
	//	fmt.Println("===============================>出anyxml.go的if判断")
	var err error
	s := new(string)
	p := new(pretty)
	p.indent = indent
	p.padding = prefix

	var rt string
	if len(rootTag) == 1 {
		rt = rootTag[0]
	} else {
		rt = DefaultRootTag
	}

	var ss string
	var b []byte
	switch v.(type) {
	case []interface{}:
		ss = "<" + rt + ">\n"
		p.Indent()
		for _, vv := range v.([]interface{}) {
			switch vv.(type) {
			case map[string]interface{}:
				//				fmt.Println("===============================>进入anyxml.go的case1")
				m := vv.(map[string]interface{})
				if len(m) == 1 {
					for tag, val := range m {
						err = mapToXmlIndentWithDateFormat(dateFormat, true, s, tag, val, p)
					}
				} else {
					p.start = 1 // we're 1 tag in to the doc
					err = mapToXmlIndentWithDateFormat(dateFormat, true, s, "element", vv, p)
					*s += "\n"
				}
			default:
				//				fmt.Println("===============================>进入anyxml.go的case2")
				p.start = 0
				err = mapToXmlIndentWithDateFormat(dateFormat, true, s, "element", vv, p)
			}
			if err != nil {
				break
			}
		}
		ss += *s + "</" + rt + ">"
		b = []byte(ss)
	case map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的case3")
		b, err = anyxmlIndentWithDateFormat(dateFormat, v.(map[string]interface{}), prefix, indent, rootTag...)
	case []map[string]interface{}:
		//		fmt.Println("===============================>进入anyxml.go的func Xml() case5")
		for _, vv := range v.([]map[string]interface{}) {
			//			fmt.Println("===============================>进入anyxml.go的func Xml() case5 vv %#v", reflect.TypeOf(vv))
			b, err = anyxmlIndentWithDateFormat(dateFormat, vv, prefix, indent, rootTag...)
			ss += (string(b) + "\n")
			if err != nil {
				break
			}
		}
		b = []byte(ss)
	default:
		//		fmt.Println("===============================>进入anyxml.go的case4")
		err = mapToXmlIndentWithDateFormat(dateFormat, true, s, rt, v, p)
		b = []byte(*s)
	}

	return b, err
}

func Struct2MapWithDateFormat(dateFormat string, obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	fmt.Println("===============================>查看t：%#v", t)
	var data = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		fmt.Println("===============================>t.Field(%d).Type %#v:", i, t.Field(i).Type)
		if t.Field(i).Type == reflect.TypeOf(time.Now()) {
			fmt.Println("===============================> %#v:", t.Field(i).Type)
			data[t.Field(i).Name] = (v.Field(i).Interface().(time.Time)).Format(dateFormat)
		} else {
			data[t.Field(i).Name] = v.Field(i).Interface()
		}

	}
	return data
}

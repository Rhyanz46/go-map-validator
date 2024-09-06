//package main
//
//import "fmt"
//
//type ChainerType interface {
//	GetParentName() string
//	Next() ChainerType
//	Back() ChainerType
//	SetName(name string) ChainerType
//	GetKey() string
//	GetParentNames(backCount int64) []string
//}
//
//type chainState struct {
//	totalChild int
//	name       string
//	parent     *chainState
//}
//
//func (cs *chainState) Next() ChainerType {
//	return &chainState{
//		parent: cs,
//	}
//}
//
//func (cs *chainState) Back() ChainerType {
//	//if cs.parent == nil {
//	//	return cs.parent
//	//}
//	if cs.parent == nil {
//		panic("cant back")
//	}
//	//*cs = *cs.parent
//	return cs.parent
//}
//
//func (cs *chainState) GetKey() string {
//	return cs.name
//}
//
//func (cs *chainState) GetParentName() string {
//	if cs.parent == nil {
//		return ""
//	}
//	return cs.parent.name
//}
//
//func (cs *chainState) SetName(name string) ChainerType {
//	cs.name = name
//	return cs
//}
//
//func (cs *chainState) GetParentNames(backCount int64) []string {
//	var names []string
//	for {
//		if cs.parent == nil {
//			break
//		}
//		cs.Back()
//		names = append(names, cs.name)
//	}
//	return names
//}
//
//func NewChainer() ChainerType {
//	return &chainState{}
//}
//
//func main() {
//	data := NewChainer().SetName("first").
//		Next().SetName("second").
//		Next().SetName("third").
//		Next().SetName("fourth").
//		Next().SetName("fifth").
//		Next().SetName("sixth").
//		Next().SetName("seventh").
//		Next().SetName("eighth").
//		Next().SetName("ninth").
//		Next().SetName("tenth").
//		Next().SetName("eleventh")
//
//	for k, v := range data.GetParentNames(0) {
//		fmt.Println(k, v)
//	}
//}

//======

//package main
//
//import "fmt"
//
//type ChainerType interface {
//	GetParentName() string
//	Next() ChainerType
//	Back() ChainerType
//	SetName(name string) ChainerType
//	GetKey() string
//	GetParentNames() []string
//}
//
//type chainState struct {
//	name   string
//	parent *chainState
//}
//
//func (cs *chainState) Next() ChainerType {
//	return &chainState{
//		parent: cs,
//	}
//}
//
//func (cs *chainState) Back() ChainerType {
//	if cs.parent == nil {
//		panic("can't back")
//	}
//	return cs.parent
//}
//
//func (cs *chainState) GetKey() string {
//	return cs.name
//}
//
//func (cs *chainState) GetParentName() string {
//	if cs.parent == nil {
//		return ""
//	}
//	return cs.parent.name
//}
//
//func (cs *chainState) SetName(name string) ChainerType {
//	cs.name = name
//	return cs
//}
//
//func (cs *chainState) GetParentNames() []string {
//	var names []string
//	current := cs
//	for current.parent != nil {
//		current = current.Back().(*chainState)
//		names = append(names, current.GetKey())
//	}
//	return names
//}
//
//func NewChainer() ChainerType {
//	return &chainState{}
//}
//
//func main() {
//	data := NewChainer().
//		SetName("first").
//		Next().SetName("second").
//		Next().SetName("third").
//		Next().SetName("fourth").
//		Next().SetName("fifth").
//		Next().SetName("sixth").
//		Next().SetName("seventh").
//		Next().SetName("eighth").
//		Next().SetName("ninth").
//		Next().SetName("tenth").
//		Next().SetName("eleventh")
//
//	for k, v := range data.GetParentNames() {
//		fmt.Println(k, v)
//	}
//}

//package main
//
//import "fmt"
//
//type ChainerType interface {
//	GetParentName() string
//	Next() ChainerType
//	Back() ChainerType
//	Forward() ChainerType
//	SetName(name string) ChainerType
//	GetKey() string
//	GetParentNames() []string
//}
//
//type chainState struct {
//	name   string
//	parent *chainState
//	child  *chainState
//}
//
//func (cs *chainState) Next() ChainerType {
//	newChild := &chainState{
//		parent: cs,
//	}
//	cs.child = newChild
//	return newChild
//}
//
//func (cs *chainState) Back() ChainerType {
//	if cs.parent == nil {
//		panic("can't back")
//	}
//	return cs.parent
//}
//
//func (cs *chainState) Forward() ChainerType {
//	if cs.child == nil {
//		panic("can't forward")
//	}
//	return cs.child
//}
//
//func (cs *chainState) GetKey() string {
//	return cs.name
//}
//
//func (cs *chainState) GetParentName() string {
//	if cs.parent == nil {
//		return ""
//	}
//	return cs.parent.name
//}
//
//func (cs *chainState) SetName(name string) ChainerType {
//	cs.name = name
//	return cs
//}
//
//func (cs *chainState) GetParentNames() []string {
//	var names []string
//	current := cs
//	for current.parent != nil {
//		current = current.Back().(*chainState)
//		names = append(names, current.GetKey())
//	}
//	return names
//}
//
//func NewChainer() ChainerType {
//	return &chainState{}
//}
//
//func main() {
//	chain := NewChainer().
//		SetName("first").
//		Next().SetName("second").
//		Next().SetName("third").
//		Next().SetName("fourth").
//		Next().SetName("fifth").
//		Next().SetName("sixth").
//		Next().SetName("seventh").
//		Next().SetName("eighth").
//		Next().SetName("ninth").
//		Next().SetName("tenth").
//		Next().SetName("eleventh")
//
//	// Cetak nama-nama parent
//	for k, v := range chain.GetParentNames() {
//		fmt.Println("Parent", k, ":", v)
//	}
//
//	// Bergerak mundur dua kali
//	chain = chain.Back().Back()
//	fmt.Println("Setelah mundur dua kali:", chain.GetKey())
//
//	// Bergerak maju satu kali
//	chain = chain.Forward()
//	fmt.Println("Setelah maju satu kali:", chain.GetKey())
//}

package main

func main() {
	//// Membuat root dari tree
	//chain := map_validator.NewChainer().SetName("root")
	//
	//// Menambahkan child ke root
	//firstChild := chain.AddChild().SetName("firstChild")
	//secondChild := chain.AddChild().SetName("secondChild")
	//
	//// Menambahkan grandchildren ke firstChild
	//firstChild.AddChild().SetName("firstChild_firstGrandchild")
	//firstChild_secondGrandchild := firstChild.AddChild().SetName("firstChild_secondGrandchild")
	//mantap := firstChild_secondGrandchild.AddChild().SetName("mantap")
	//firstChild_secondGrandchild.AddChild().SetName("keren")
	//mantap.AddChild().SetName("mantap1")
	//
	//// Menambahkan grandchildren ke secondChild
	//secondChild.AddChild().SetName("secondChild_firstGrandchild")
	//secondChild.AddChild().SetName("secondChild_secondGrandchild")
	//
	//// Cetak nama-nama parent
	//fmt.Println("Parents of firstChild_firstGrandchild:")
	//for k, v := range firstChild.Next(0).GetParentNames() {
	//	fmt.Println(k, v)
	//}
	//
	//// Bergerak ke parent dari firstGrandchild dari firstChild
	//current := firstChild.Next(0).Back()
	//fmt.Println("Back to parent:", current.GetKey())
	//
	//// Bergerak ke firstGrandchild dari firstChild lagi
	//current = firstChild.Forward(0)
	//fmt.Println("Forward to first grandchild:", current.GetKey())
	//
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//chain.PrintHierarchyWithSeparator("/", "")
	//
	//aa := map_validator.NewChainer().SetName("module")
	//aa.LoadTreeFromMap(map[string]interface{}{
	//	"math": map[string]interface{}{
	//		"add":      map[string]interface{}{"validateInput": nil, "computeResult": nil},
	//		"subtract": map[string]interface{}{"verifyInput": nil, "calculateDifference": nil},
	//		"multiply": map[string]interface{}{
	//			"nama":        "arian",
	//			"kelas":       12,
	//			"verifyInput": nil,
	//			"orang": map[string]interface{}{
	//				"nama":    "arian",
	//				"kelas":   12,
	//				"age":     12,
	//				"married": false,
	//			},
	//		},
	//		"divide": nil,
	//	},
	//	"string": map[string]interface{}{
	//		"concat": nil,
	//		"split":  nil,
	//	},
	//})
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//aa.PrintHierarchyWithSeparator(".", "")
}

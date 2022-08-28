// Reflect 사용법 https://johngrib.github.io/wiki/golang-reflect/
package test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"testing"
)

// switch Type() check
func typeSwitch(value interface{}) {
	switch value.(type) {
	case int:
		fmt.Println("type:", reflect.ValueOf(value).Type())
		fmt.Println(value, "is int.")
	case string:
		fmt.Println("type:", reflect.ValueOf(value).Type())
		fmt.Println(value, "is string.")
	case map[string]int:
		fmt.Println("type:", reflect.ValueOf(value).Type())
		fmt.Println(value, "is map.")
	default:
		fmt.Println("type:", reflect.ValueOf(value).Type())
	}
}

func TestTypeCheck(t *testing.T) {
	typeSwitch(42)
	typeSwitch("42")
	typeSwitch([]float32{3.14})
	typeSwitch(map[string]int{"Hi": 26})
}

// struct 메타 데이터
type Person struct {
	Name   string `json:"name"`
	Nation string `json:"country"`
	Zip    int    `json:"zipCode"`
}

type Dog struct {
	Name string  `json:"name"`
	Like *Person `json:"person"`
}

func metaData(anything interface{}) {
	target := reflect.ValueOf(anything)
	elements := target.Elem()

	fmt.Printf("Type: %s\n", target.Type()) // 구조체 타입명

	for i := 0; i < elements.NumField(); i++ {
		mValue := elements.Field(i)
		mType := elements.Type().Field(i)
		tag := mType.Tag

		fmt.Printf("%10s %10s ==> %10v, json: %10s\n",
			mType.Name,         // 이름
			mType.Type,         // 타입
			mValue.Interface(), // 값
			tag.Get("json")) // json 태그
	}
}

func TestStructMetaData(t *testing.T) {
	john := &Person{
		Name:   "JohnGrib",
		Nation: "Korea",
		Zip:    12345,
	}
	metaData(john)
	metaData(&Dog{
		Name: "Wolfy",
		Like: john,
	})
}

// Method 메타 데이터
type Temp struct{}

func (t *Temp) Prints(msg string, num int) (string, error) {
	return fmt.Sprintf("%s %d", msg, num), nil
}

func (t *Temp) Add(a, b int) int {
	return a + b
}

func (t *Temp) sub(a, b int) int {
	return a - b
}

// input parameter 정보 출력
func methodInMeta(method reflect.Value) {
	methodType := method.Type()
	fmt.Println(" input types:", methodType.NumIn())

	for j := 0; j < methodType.NumIn(); j++ {
		param := methodType.In(j)
		fmt.Printf("  in #%d : %s\n", j, param.Name())
	}
}

// output return 값 정보 출력
func methodOutMeta(method reflect.Value) {
	methodType := method.Type()
	fmt.Println(" output types:", methodType.NumOut())

	for j := 0; j < methodType.NumOut(); j++ {
		param := methodType.Out(j)
		fmt.Printf("  out #%d : %s\n", j, param.Name())
	}
}

func methodMetaData(method reflect.Value) {
	fmt.Println(method.Type())
	methodInMeta(method)
	methodOutMeta(method)
	fmt.Println("")
}

func TestGetMethodMetaInfo(t *testing.T) {
	var temp Temp

	value := reflect.ValueOf(&temp)
	fmt.Println(value.Type())
	fmt.Println("")

	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		methodMetaData(method)
	}

	// 메소드 이름으로도 찾을 수 있다
	fmt.Println("? add Method ?")
	addMethod := value.MethodByName("Add")
	if addMethod.IsValid() {
		methodMetaData(addMethod)
	}

	fmt.Println("? sub Method ?")
	subMethod := value.MethodByName("sub")
	if subMethod.IsValid() {
		fmt.Println(subMethod.IsValid())
	} else {
		fmt.Println("sub is private method")
	}
}

// Method Call
func printRes(ret []reflect.Value) {
	//t.Log("returns %d values\n", len(ret))
	fmt.Println(fmt.Sprintf("returns %d values\n", len(ret)))
	for i, retValue := range ret {
		//t.Log("  %d => %v (type: %s)\n", i, retValue, retValue.Type())
		fmt.Println(fmt.Sprintf("%d => %v (type: %s\n", i, retValue, retValue.Type()))
	}
}

func TestCallMethos(t *testing.T) {
	var temp Temp

	ret1 := reflect.ValueOf(&temp).
		MethodByName("Prints").
		Call(
			[]reflect.Value{
				reflect.ValueOf("hello"),
				reflect.ValueOf(3),
			},
		)
	printRes(ret1)

	ret2 := reflect.ValueOf(&temp).
		MethodByName("Add").
		Call(
			[]reflect.Value{
				reflect.ValueOf(37),
				reflect.ValueOf(5),
			},
		)
	printRes(ret2)
}

// Package 파일 찾기
func TestFindPackageFiles(t *testing.T) {
	//const pkgName = "calc"
	const pkgName = "../test"

	//var pkgs map[string]*ast.Package
	pkgs, err := parser.ParseDir(token.NewFileSet(), pkgName, nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		fmt.Println("package:", pkg.Name)

		for _, file := range pkg.Files {
			fmt.Println("\t file:", file.Name) // package 아래의 파일 이름

			fmt.Println("\t funtions names:") // package 아래의 함수 이름들
			for _, decl := range file.Decls {
				if function, ok := decl.(*ast.FuncDecl); ok {
					fmt.Println(function.Name)
				}
			}
		}
	}
}

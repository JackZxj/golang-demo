package main
 
import (
    "io/ioutil" //io 工具包
    "fmt"
    "os"
    "strings"
)
 
func main() {
    str := `a\naa
aaa`
    newstr := strings.ReplaceAll(str, "\\n", string(byte(10)))

    var d1 = []byte(str)
    var d2 = []byte(newstr)

    fmt.Printf("before: %v %s\n", d1, str)
    fmt.Printf("after: %v %s\n", d2, newstr)

    wr := func(source []byte, filename string) {
        ioutil.WriteFile(filename, source, 0666)

        file, _ := os.Open(filename)
        b, _ := ioutil.ReadAll(file)
        fmt.Printf("read from %s: %s\n", filename, string(b))
    }
    wr(d1, "str.txt")
    wr(d2, "newstr.txt")
}

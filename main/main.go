package main
import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"sort"
)

func readLines(name string, startLine int, endLine int) (line []string) {
	a := make([]string, 1)
	inFile, _ := os.Open(name)
  	defer inFile.Close()
  	sc := bufio.NewScanner(inFile)
    lastLine := 0
    for sc.Scan() {
        lastLine++
        if lastLine >= startLine &&  lastLine<endLine{
       		a = append(a,sc.Text())     
        }else if lastLine >=endLine{
        	break
        }
    }
    return a
}

func writeToFile(name string, list []string){
	file, _ := os.Create(name)

	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range list {
		fmt.Fprintln(w, line)
	}
	
}



func main(){
	start := 0
	const jump = 10000000
	const end = 200000000
	for start < end{
		a := readLines("/home/prajakt/Downloads/64/val",start,start + jump)
		sort.Strings(a)
		writeToFile(strconv.Itoa(start),a)
		start += jump
	}
	fmt.Println("done")
}


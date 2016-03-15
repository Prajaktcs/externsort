package main
import (
	"strconv"
	"fmt"
	"bufio"
	"os"
	"sort"
	"strings"
)

func readLines(name string, startLine int, endLine int) (line []string) {
	a := make([]string, 0)
	inFile, _ := os.Open(name)
  	defer inFile.Close()
  	sc := bufio.NewScanner(inFile)
    lastLine := 0
    for sc.Scan() {
        if lastLine >= startLine &&  lastLine<endLine{
       		a = append(a,sc.Text()) 
       		
        }else if lastLine >=endLine{
        	break
        }
        lastLine++
    }
    return a
}

func writeToFile(name string, list []string){
	file, _ := os.Create(name)
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range list {
		fmt.Fprintln(w, line)
		w.Flush()
	}
	
}

func find_min(values []string) (min_val string,index int, flag bool){
	min := values[0]
	min_index := 0
	flag = true
	for i:=1; i< len(values); i++{
		if(values[i]== ""){

		}else if(strings.Compare(min,values[i]) == 1){
			min = values[i]
			min_index = i	
			flag = false
		}else{
			flag = false
		}
	}
	return min, min_index, flag
}

func writeToMergedFile(min string){
	file, _ := os.OpenFile("final",  os.O_RDWR|os.O_APPEND, 0666)
	defer file.Close()
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, min)
	w.Flush()
	
}	


func mergeList(names []string){
	//fileHandlerList := make([]*os.File,0)
	readerList := make([]*bufio.Reader,0)
	
	for _,name := range names{
		//append(fileHandlerList,)
		file,_ := os.Open(name)
		readerList = append(readerList,bufio.NewReader(file))
	}
	_, _ = os.Create("final")
	//creating an array whose length is the number of files
	min_val_slice := make([]string,0)
	for _,read := range readerList{
		line,_,_ := read.ReadLine()
		min_val_slice = append(min_val_slice,string(line))
	}
	for true{
		min, min_index, flag := find_min(min_val_slice)
		
		//fmt.Println(min_index)
		writeToMergedFile(min)
		line,_,err := readerList[min_index].ReadLine()
		if err == nil{
			min_val_slice[min_index] = string(line)
		}else{
			min_val_slice[min_index] = ""
		}
		if(flag){
			break
		}
	}
	
}


func main(){
	names := make([]string, 0)
	start := 0
	const jump = 10000000
	const end = 100000000
	for start < end{
		mid := start + jump
		a := readLines("/home/prajakt/Downloads/64/val",start,mid)
		sort.Strings(a)
		writeToFile(strconv.Itoa(start),a)
		names = append(names,strconv.Itoa(start))
		start += jump
	}
	mergeList(names)
	fmt.Println("done")
}

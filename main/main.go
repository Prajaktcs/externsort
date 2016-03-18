package main
import (
	"fmt"
	"flag"
	"bufio"
	"sort"
	"os"
	"strings"
	"strconv"
	"bytes"
	"time"
)
/**

*/
func writeToFile(list []string,name string){
	file, _ := os.Create(name)
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, line := range list {
		fmt.Fprintln(w, line)
		w.Flush()
	}
	
}
func display_list(values []string) {
	for _,line := range values{
		fmt.Println(line)
	}
}

func fill_min() (string){
	var buffer bytes.Buffer
	for i := 0; i < 100; i++ {
        buffer.WriteString("~")
    }
    return buffer.String()
}
/*

*/
func sortAndWriteToFile(list []string,name string) {
	sort.Strings(list) //Sort the strings
	writeToFile(list,name)
}
//find the minumum in the array
func find_min(values []string) (min_val string,index int, flag bool){
	min := fill_min()
	min_index := 0
	flag = true
	for i:=0; i< len(values); i++{
		if(values[i]== ""){
			//display_list(values)
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

func fillBuffer(names []*bufio.Reader, bufferSize int) ([][]string){
	bufferList := make([][]string,0)
	for _,name := range names{
		currentLine := 0
		//Asuming the buffersize is less than the file size
		temp := make([]string,0)
		for currentLine < bufferSize{
			text,_,_ := name.ReadLine()
			temp = append(temp, string(text))
			currentLine++
		}
		bufferList = append(bufferList,temp)
	}
	return bufferList
}

func loadNextSlot(reader *bufio.Reader, bufferSize int) ([]string){
	bufferList := make([]string, 0)
	currentLine := 0
	text,_,err:=reader.ReadLine()
	bufferList = append(bufferList,string(text))
	if err == nil{
		for currentLine < bufferSize{
			currentLine++
			text,_,err = reader.ReadLine()
			if err !=nil{
				break
			}
			bufferList = append(bufferList,string(text))
		}	
	}
	return bufferList
}

func processBufferList(bufferList [][]string, readerList []*bufio.Reader, bufferSize int) {
	//Creating an array to hold the first value for all the lists
	currentValuesFromLists := make([]string,0)
	currentIndicesOfLists := make([]int,0)
	trackEmptyFiles := make([]int,0)
	for _, element := range bufferList{
		currentValuesFromLists = append(currentValuesFromLists,string(element[0]))
		currentIndicesOfLists = append(currentIndicesOfLists,0) //Appending 0 as this is the first index
		trackEmptyFiles = append(trackEmptyFiles,0)
	}
	for true{
		min, min_index, flag := find_min(currentValuesFromLists)
		if(flag){
			break
		}
		writeToMergedFile(min)
		
		//Get the next element from the array
		currentIndicesOfLists[min_index] = currentIndicesOfLists[min_index] +1 //increasing index to get the next element
		if len(bufferList[min_index]) > currentIndicesOfLists[min_index]{
			currentValuesFromLists[min_index] = string(bufferList[min_index][currentIndicesOfLists[min_index]])
		}else{

			// Checking if the file is already empty
			if trackEmptyFiles[min_index] != 1{
				currentIndicesOfLists[min_index] = 0
				a:=  loadNextSlot(readerList[min_index],bufferSize)
				//If the end of the file has reached, then place an empty string in array
				if(len(a)==0){
					currentValuesFromLists[min_index] = ""
					trackEmptyFiles[min_index] = 1
				}else{
					bufferList[min_index] = a
					currentValuesFromLists[min_index] = string(bufferList[min_index][0])
				}
			}else{
				//fmt.Println(min_index)
			}
		}

	}
}

/*
	Function to manage the merging of all the lists
*/
func mergeList(names []string, buffer int){
	//fileHandlerList := make([]*os.File,0)
	readerList := make([]*bufio.Reader,0)
	
	for _,name := range names{
		//append(fileHandlerList,)
		file,_ := os.Open(name)
		readerList = append(readerList,bufio.NewReader(file))
	}
	_, _ = os.Create("final")
	//creating an array whose length is the number of files
	bufferList := fillBuffer(readerList, buffer)
	processBufferList(bufferList,readerList,buffer)
	
}


/*
 function that defines the process for the sort
*/
func process(filePath string, jump int, end int, bufferSize int) (int,error){
	names := make([]string,0)
	currentLine := 0
	tempList := make([]string, 0)
	//Creating a handler for the file to be read
	//If error is found return 1 and the error
	inFile, err := os.Open(filePath)
	if err != nil {
		return 1 ,err
	}
	defer inFile.Close()
  	sc := bufio.NewScanner(inFile)
  	//Running the loop till the file is not over, or the end is not met
	for currentLine<end && sc.Scan(){
		if(currentLine%jump==0 && currentLine !=0 ){
			//fmt.Println(currentLine)
			sortAndWriteToFile(tempList, strconv.Itoa(currentLine))
			names = append(names,strconv.Itoa(currentLine))
			tempList = make([]string, 0)
		}
		tempList = append(tempList,sc.Text())
		currentLine ++
	}
	sortAndWriteToFile(tempList, strconv.Itoa(currentLine))
	names = append(names,strconv.Itoa(currentLine))
	tempList = make([]string, 0)
	//Jump will be equal to end if the experiment is for 1 GB
	if(jump!=end){
		mergeList(names, bufferSize)
	}
	return 0,err
}

/*
	Takes 3 parameters:-
	1. File path
	2. Numbers to sort in one go
	3. Number of lines in the file
*/
func main(){
	filePath := flag.String("filePath", "","Specify the file path that contains the data")
	jump := flag.Int("jump",0,"Specify the amount of lines each intermidiate file will have")
	end := flag.Int("end",0,"the size of the input file")
	bufferSize := flag.Int("buffer",0,"size of the buffer while merging")
	flag.Parse()
	if *filePath == "" || *jump == 0 || *end == 0  || *bufferSize ==0{
		fmt.Println("Invalid parameters. Please run again with right parameters")
	}else{
		start := time.Now()
		process(*filePath,*jump,*end,*bufferSize)
		end := time.Since(start).String()
		fmt.Println(end)
	}
}

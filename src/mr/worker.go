package mr

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}


//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// Your worker implementation here.

	// uncomment to send the Example RPC to the coordinator.
	files := Call()

	intermediate := []KeyValue{}
	for _,filename := range files {
		file ,err := os.Open(filename.Key)
		if err != nil {
			log.Fatal("cannot open file %v", filename.Key)
		}
		content ,err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("cannot read %v", filename.Key)
		}
		file.Close()

		kva := mapf(filename.Key,string(content))
		intermediate = append(intermediate, kva...)
 	}

	sort.Sort(ByKey(intermediate))

	oname := "mr-out-" + string(len(intermediate))
	ofile, _ := os.Create(oname)

	for i := 0 ; i < len(intermediate) ;{
		j := i+1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i ; k<j ; k++ {
			values = append(values,intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key,values)

		fmt.Fprintf(ofile, "%v %v\n", intermediate[i].Key, output)

	}

	ofile.Close()

}

//
// example function to show how to make an RPC call to the coordinator.
//
// the RPC argument and reply types are defined in rpc.go.
//
func Call() []KeyValue{

	// declare an argument structure.
	args := Args{}

	// fill in the argument(s).
	args.argsname = "???"

	// declare a reply structure.
	reply := Reply{}

	// send the RPC request, wait for the reply.
	// the "Coordinator.Example" tells the
	// receiving server that we'd like to call
	// the Example() method of struct Coordinator.
	ok := call("Coordinator.GetFilenames", &args, &reply)
	if ok {
		// reply.Y should be 100.
		fmt.Printf("reply.files %v\n", reply.files)
	} else {
		fmt.Printf("call failed!\n")
	}

	return reply.files
}

//
// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	// sockname := coordinatorSock()
	// c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

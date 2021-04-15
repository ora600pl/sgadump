package main

import (
  "os"
  "fmt"
  "strconv"
  "encoding/binary"
)

func main() {
  var shmfile string
  var blocksize int64
  var data_object_id uint64
  var fname string
  if len(os.Args) < 9 {
      fmt.Printf("file_sgadump by Kamil Stawiarski (@ora600pl) - dumps database blocks from SGA VM memory dump\n")
      fmt.Printf("Usage: file_sgadump -b block_size -d data_object_id -s shmfile -o output_file_name\n")
      return
  }
  for i := 0; i < len(os.Args) ; i++ {
      if os.Args[i] == "-s" {
          shmfile = os.Args[i+1]
      } else if os.Args[i] == "-d" {
          data_object_id, _ = strconv.ParseUint(os.Args[i+1], 10, 32)
     } else if os.Args[i] == "-b" {
          blocksize, _ = strconv.ParseInt(os.Args[i+1], 10, 32)
     } else if os.Args[i] == "-o" {
          fname = os.Args[i+1]
     }
  }

  file, err := os.Create(fname)
  if err != nil {
    panic(err)
  }
  defer file.Close()


  if segment, err := os.Open(shmfile) ; err == nil {
      defer segment.Close()
      foundBlocks := 0
      fs, _ := segment.Stat()
      blocks := fs.Size() / blocksize
      var i int64
      for i = 0; i < blocks ; i++ {
	  block_data := make([]byte, blocksize)
          segment.Read(block_data)
          if block_data[0] == 6 {
              objd := uint64(binary.LittleEndian.Uint32(block_data[24:28]))
              if data_object_id == objd {
                file.Write(block_data)
                foundBlocks++
              }
          }
       }
       fmt.Printf("shmid = %s\t size is %d, blocks = %d\n", shmfile, fs.Size(), blocks)
       fmt.Printf("Dumped %d blocks to %s\n", foundBlocks, fname)
   }
}

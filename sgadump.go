package main

import (
  "github.com/ghetzel/shmtool/shm"
  "os"
  "fmt"
  "strconv"
  "encoding/binary"
)

func main() {
  var shmid int
  var blocksize int64
  var data_object_id uint64
  var fname string
  if len(os.Args) < 9 {
      fmt.Printf("sgadump by Kamil Stawiarski (@ora600pl) - dumps database blocks from SGA.\n")
      fmt.Printf("Usage: sgadump -b block_size -d data_object_id -s shmid -o output_file_name\n")
      return
  }
  for i := 0; i < len(os.Args) ; i++ {
      if os.Args[i] == "-s" {
          shmid, _ = strconv.Atoi(os.Args[i+1])
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


  if segment, err := shm.Open(shmid) ; err == nil {
    defer segment.Destroy()
    if segmentAddress, err := segment.Attach(); err == nil {
      foundBlocks := 0
      defer segment.Detach(segmentAddress)
      blocks := segment.Size / blocksize
      var i int64
      for i = 0; i < blocks ; i++ {
          block_data, err := segment.ReadChunk(blocksize, i*blocksize)
          if err == nil && block_data[0] == 6 {
              objd := uint64(binary.LittleEndian.Uint32(block_data[24:28]))
              if data_object_id == objd {
                file.Write(block_data)
                foundBlocks++
              }
          }
      }
      fmt.Printf("shmid = %d\t size is %d, blocks = %d\n", shmid, segment.Size, blocks)
      fmt.Printf("Dumped %d blocks to %s\n", foundBlocks, fname)
   }
  }
}


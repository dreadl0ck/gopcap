# GOPCAP

This package provides support for reading Packet Capture (PCAP) files efficiently in pure Go.

It does not use reflection to parse the file and packet headers,
and therefore is much faster than other implementations.

It will be extended with support for writing PCAP in the near future.

The following API is exported by the package:

```go
    func Count(path string) int64
    type FileHeader
    type PacketHeader
    type Reader
        func Open(filename string) (*Reader, error)
        func (r *Reader) Close() error
        func (r *Reader) ReadNextPacket() (PacketHeader, []byte, error)
        func (r *Reader) ReadNextPacketHeader() (PacketHeader, []byte, error)
    type Writer
        func (r *Writer) Close() error
        func (r *Writer) Open(filename string) error
```

Data structures for the PCAP file header and packet headers are provided:

```go
    type FileHeader struct {
        // magic number
        MagicNumber uint32
        // major version number
        VersionMajor uint16
        // minor version number
        VersionMinor uint16
        // GMT to local correction
        Thiszone int32
        // accuracy of timestamps
        Sigfigs uint32
        // max length of captured packets, in octets
        Snaplen uint32
        // data link type
        Network uint32
    }

    type PacketHeader struct {
        // timestamp seconds
        TsSec int32
        // timestamp microseconds
        TsUsec int32
        // number of octets of packet saved in file
        CaptureLen int32
        // actual length of packet
        OriginalLen int32
    }
```

## Usage

Reading a PCAP file from disk can be done by calling gopcap.Open and looping on r.ReadNextPacket().
Additionally there is a function called gopcap.Count that will loop over all packets and return the total count.

```go
    // get total packet count
    fmt.Println("total:", gopcap.Count(path))

    // create reader
    r, err := gopcap.Open(path)
    if err != nil {
        panic(err)
    }
    defer r.Close()

    // loop over packets
    for {
        h, data, err := r.ReadNextPacket()
        if err != nil {
            if err == io.EOF {
                println("EOF")
                break
            }
            panic(err)
        }
        fmt.Println(h, len(data))
    }
```

## Benchmarks

# GOPCAP

This package provides support for reading Packet Capture (PCAP) files efficiently in pure Go.

**The problem:**

Existing implementations use **binary.Read** for parsing the binary data into the expected structure.
However, the function makes use of reflection, which hurts performance.

This implementation does not use **binary.Read** to parse the file and packet headers,
and therefore is much faster.

It will be extended with support for writing PCAP in the near future.

The following API is exported by the package:

```golang
func Count(path string) int64
type FileHeader
type PacketHeader
type Reader
    func Open(filename string) (*Reader, error)
    func (r *Reader) Close() error
    func (r *Reader) ReadNextPacket() (PacketHeader, []byte, error)
    func (r *Reader) ReadNextPacketHeader() (PacketHeader, []byte, error)
```

Data structures for the PCAP file header and packet headers are provided:

```golang
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

Reading a PCAP file from disk can be done by calling **gopcap.Open(path)** and looping on **r.ReadNextPacket()** of the returned reader.
Additionally there is a function called **gopcap.Count(path)** that returns the total count of packets in the file (useful for displaying progress).

```golang
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

There are a few pure go implementations for parsing PCAP files available.
The following have not been evaluated in the benchmark, for the given reason(s):

- https://github.com/davecheney/pcap, Latest commit 10760a1  on Aug 19, 2012, fails to compile with various errors
- https://github.com/Lukasa/gopcap, limited in functionality, API only allows to parse a file completely, which is unpractical for big files.

These implementation are included in the benchmarks, in the order they are listed:

- https://github.com/dreadl0ck/gopcap (this package)
- https://github.com/github.com/0intro/pcap
- https://github.com/google/gopacket/pcap
- https://godoc.org/go.universe.tf/netboot
- https://github.com/github.com/miekg/pcap

The benchmark code fetches a single packet in a loop, and discards all data that is not needed.
Make sure the PCAP file for the test has enough packets to be read in one call, otherwise the tests wont produce meaningful results.

I didn't include the dump in the repo because it is around 1.0G in size.
The used PCAP file (**maccdc2012_00000.pcap**) is from the National CyberWatch Mid-Atlantic Collegiate Cyber Defense Competition (MACCDC),
which can be downloaded here: https://www.netresec.com/?page=MACCDC

    $ go test -run=XXX -bench=.
    goos: darwin
    goarch: amd64
    pkg: github.com/dreadl0ck/gopcap
    BenchmarkReadPcap-12            	 5000000	       297 ns/op
    BenchmarkReadPcap0Intro-12      	  200000	      6086 ns/op
    BenchmarkReadPcapGoPacket-12    	 3000000	       437 ns/op
    BenchmarkReadPcapNetboot-12     	 2000000	       809 ns/op
    BenchmarkReadPcapMiekg-12       	 2000000	       909 ns/op
    PASS
    ok  	github.com/dreadl0ck/gopcap	9.854s

This implementation achieves 297 ns/op which is the fastest of all compared.
The gopacket pcap library and fork from miekg both use C bindings.
For the benchmark of the gopacket pcap implementation, the **ZeroCopyReadPacketData()** function was used.

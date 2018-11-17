# GOPCAP

[![Go Report Card](https://goreportcard.com/badge/github.com/dreadl0ck/gopcap)](https://goreportcard.com/report/github.com/dreadl0ck/gopcap)
[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://raw.githubusercontent.com/dreadl0ck/gopcap/master/docs/LICENSE)
[![Golang](https://img.shields.io/badge/Go-1.10-blue.svg)](https://golang.org)
![Linux](https://img.shields.io/badge/Supports-Linux-green.svg)
![macOS](https://img.shields.io/badge/Supports-macOS-green.svg)
![Windows](https://img.shields.io/badge/Supports-Windosa-green.svg)

This package provides support for reading Packet Capture (PCAP) files efficiently in pure Go
and provides benchmarks against other packages offering the same functionality.

**Update**: I just discovered the go implementation from gopacket https://github.com/google/gopacket/pcapgo,
which is comparatively fast and has support for both reading and writing PCAP and PCAP-NG! Better use this library instead.
I've added the gopacket/pcapgo implementation to the benchmarks as well.

However, in my benchmarks **this implementation still is slightly faster**.
So, if performance is crucial and you only need to read PCAP, this lib might still be interesting.

**The problem:**

Existing implementations often use **binary.Read** for parsing the binary data into the expected structure.
However, the function makes use of reflection, which hurts performance.

This implementation does not use **binary.Read** to parse the file and packet headers,
and therefore is much faster.

The following API is exported by the package:

```golang
func Count(path string) int64
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
- https://github.com/google/gopacket/pcapgo

The benchmark code fetches a single packet in a loop, and discards all data that is not needed.
Make sure the PCAP file for the test is big enough, otherwise the tests wont produce meaningful results.

I didn't include the dump in the repo because it is around 1.0G in size.
The used PCAP file (**maccdc2012_00000.pcap**) is from the National CyberWatch Mid-Atlantic Collegiate Cyber Defense Competition (MACCDC),
which can be downloaded here: https://www.netresec.com/?page=MACCDC

Here are the results of a few runs on my development machine (MacBook Pro 2018, 32GB RAM, 2,9 GHz Intel Core i9)

    $ go test -bench=.
    goos: darwin
    goarch: amd64
    pkg: github.com/dreadl0ck/gopcap
    BenchmarkReadPcap-12                  	10000000	       184 ns/op
    BenchmarkReadPcap0Intro-12            	  300000	      5284 ns/op
    BenchmarkReadPcapGoPacket-12          	 5000000	       513 ns/op
    BenchmarkReadPcapNetboot-12           	 2000000	       819 ns/op
    BenchmarkReadPcapMiekg-12             	 2000000	       857 ns/op
    BenchmarkReadPcapGopacketPcapGo-12    	 5000000	       282 ns/op
    PASS
    ok  	github.com/dreadl0ck/gopcap	13.138s

    $ go test -bench=.
    goos: darwin
    goarch: amd64
    pkg: github.com/dreadl0ck/gopcap
    BenchmarkReadPcap-12                  	10000000	       145 ns/op
    BenchmarkReadPcap0Intro-12            	  300000	      4211 ns/op
    BenchmarkReadPcapGoPacket-12          	 5000000	       292 ns/op
    BenchmarkReadPcapNetboot-12           	 3000000	       868 ns/op
    BenchmarkReadPcapMiekg-12             	 2000000	       900 ns/op
    BenchmarkReadPcapGopacketPcapGo-12      10000000	       252 ns/op
    PASS
    ok  	github.com/dreadl0ck/gopcap	13.255s

    $ go test -bench=.
    goos: darwin
    goarch: amd64
    pkg: github.com/dreadl0ck/gopcap
    BenchmarkReadPcap-12                    10000000	       186 ns/op
    BenchmarkReadPcap0Intro-12            	  300000	      5281 ns/op
    BenchmarkReadPcapGoPacket-12          	 3000000	       440 ns/op
    BenchmarkReadPcapNetboot-12           	 2000000	       801 ns/op
    BenchmarkReadPcapMiekg-12             	 2000000	       801 ns/op
    BenchmarkReadPcapGopacketPcapGo-12    	 5000000	       299 ns/op
    PASS
    ok  	github.com/dreadl0ck/gopcap	12.524s

It seems this implementation is the fastest of all compared.
The gopacket pcap library and fork from miekg both use C bindings.
For the benchmark of the gopacket pcap implementation, the **ZeroCopyReadPacketData()** function was used.

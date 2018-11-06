/*
 * GOPACP - PCAP file parsing in Golang
 * Copyright (c) 2017 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package gopcap

import "io"

/////////////////////////////
// Utils
/////////////////////////////

// Count packets in file
// TODO: dont actually read anything... just seek to next offset
func Count(path string) int64 {

	var (
		numPackets int64
		r, err     = Open(path)
	)
	if err != nil {
		println("failed to open file")
		panic(err)
	}

	for {
		_, _, err := r.ReadNextPacket()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		numPackets++
	}

	err = r.Close()
	if err != nil {
		println("failed to close file")
		panic(err)
	}

	return numPackets
}

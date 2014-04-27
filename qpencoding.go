// +build ignore

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bom-d-van/qpencoding"
)

var (
	decode  bool
	encode  bool
	version bool
	help    bool
)

func init() {
	flag.BoolVar(&decode, "decode", false, "Decode Quoted-Printable encoded file")
	flag.BoolVar(&decode, "d", false, "Decode Quoted-Printable encoded file")

	flag.BoolVar(&encode, "encode", false, "Encode file into Quoted-Printable")
	flag.BoolVar(&encode, "e", false, "Encode file into Quoted-Printable")

	flag.BoolVar(&version, "version", false, "Print version number")
	flag.BoolVar(&version, "v", false, "Print version number")

	flag.BoolVar(&help, "help", false, "Print this message")
	flag.BoolVar(&help, "h", false, "Print this message")
}

func main() {
	flag.Parse()
	if version {
		println("1.0")
		return
	}

	if len(flag.Args()) == 0 || help {
		fmt.Println(`qpencoding  --  Encode/decode file as Quoted-Printable (RFC 2045).  Call:
                qpencoding [-e / -d] [file]

Options:
           -d, --decode      Decode Quoted-Printable encoded file
           -e, --encode      Encode file into Quoted-Printable
           -h, --help        Print this message
           -v, --version     Print version number`)

		return
	}

	file, err := os.OpenFile(flag.Args()[0], os.O_RDONLY, 777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if decode {
		r := qpencoding.NewReader(file)
		cnt, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Printf("%s", string(cnt))

		return
	}

	if encode {
		w := qpencoding.NewWriter(os.Stdout)
		_, err := io.Copy(w, file)
		if err != nil {
			fmt.Println(err)
			return
		}

		return
	}
}

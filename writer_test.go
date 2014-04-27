package qpencoding

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

var wcases = []struct {
	in, processed, want string
}{
	{
		in:        "foo bar foo bar",
		processed: "foo bar foo bar",
		want:      "foo bar foo bar=\r\n",
	},
	{
		in:        "\tfoo bar foo bar",
		processed: "\tfoo bar foo bar",
		want:      "\tfoo bar foo bar=\r\n",
	},
	{
		in:        "foo\rbar",
		processed: "foo\r\nbar",
		want:      "foo\r\nbar=\r\n",
	},
	{
		in:        "foo\nbar",
		processed: "foo\r\nbar",
		want:      "foo\r\nbar=\r\n",
	},

	// CRLF Cases
	{
		in:        "foo\n\rbar",
		processed: "foo\r\nbar",
		want:      "foo\r\nbar=\r\n",
	},
	{
		in:        "foo\r\nbar",
		processed: "foo\r\nbar",
		want:      "foo\r\nbar=\r\n",
	},

	{
		in:        "foo bar=foo bar",
		processed: "foo bar=foo bar",
		want:      "foo bar=3Dfoo bar=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010101",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010101",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010=\r\n101=\r\n",
	},
	{
		in:        "0101010101010101010101010101010101010101010101010101010101010101010101\t 101",
		processed: "0101010101010101010101010101010101010101010101010101010101010101010101\t 101",
		want:      "0101010101010101010101010101010101010101010101010101010101010101010101\t 101=\r\n",
	},

	// Tabs And Spaces At The End Of Line
	{
		in:        "foo bar   \nfoo bar   ",
		processed: "foo bar   \r\nfoo bar   ",
		want:      "foo bar  =20\r\nfoo bar   =\r\n",
	},
	{
		in:        "01010101010101010101010101010101010101010101010101010101010101010101010101 ",
		processed: "01010101010101010101010101010101010101010101010101010101010101010101010101 ",
		want:      "01010101010101010101010101010101010101010101010101010101010101010101010101 =\r\n",
	},

	{
		in:        "0101010101010101010101010101010101010101010101010101010101010101010101010101",
		processed: "0101010101010101010101010101010101010101010101010101010101010101010101010101",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010=\r\n1=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010 something",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010 something",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010=\r\n something=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010\nsomething",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010\rsomething",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010\r\nsomething=\r\n",
	},
	{
		in:        "010101010101010101010101010101010101010101010101010101010101010101010101010 \nsomething",
		processed: "010101010101010101010101010101010101010101010101010101010101010101010101010 \r\nsomething",
		want:      "010101010101010101010101010101010101010101010101010101010101010101010101010=\r\n=20\r\nsomething=\r\n",
	},
	{
		in:        "01010101010101010101010101010101010101010101010101010101010101010101010101 \nsomething",
		processed: "01010101010101010101010101010101010101010101010101010101010101010101010101 \r\nsomething",
		want:      "01010101010101010101010101010101010101010101010101010101010101010101010101=\r\n=20\r\nsomething=\r\n",
	},
	{
		in:        "01010101010101010101010101010101010101010101010101010101010101010101010101\t\nsomething",
		processed: "01010101010101010101010101010101010101010101010101010101010101010101010101\t\r\nsomething",
		want:      "01010101010101010101010101010101010101010101010101010101010101010101010101=\r\n=09\r\nsomething=\r\n",
	},

	{
		in:        "我只愿面朝大海，春暖花开",
		processed: "我只愿面朝大海，春暖花开",
		want:      "=E6=88=91=E5=8F=AA=E6=84=BF=E9=9D=A2=E6=9C=9D=E5=A4=A7=E6=B5=B7=EF=BC=8C=E6=\r\n=98=A5=E6=9A=96=E8=8A=B1=E5=BC=80=\r\n",
	},
}

func TestWriter(t *testing.T) {
	for _, c := range wcases {
		buf := bytes.NewBuffer([]byte{})
		NewWriter(buf).Write([]byte(c.in))
		if buf.String() != c.want {
			t.Errorf("In: %q\nGot: %q\nWant: %q\n", c.in, buf.String(), c.want)
		}
	}
}

func TestReadWrite(t *testing.T) {
	for _, c := range wcases {
		buf := bytes.NewBuffer([]byte{})
		NewWriter(buf).Write([]byte(c.in))
		processed, err := ioutil.ReadAll(NewReader(buf))
		if err != nil {
			t.Errorf("Processing Error: %s (In: %q)", err, c.in)
		}
		if string(processed) != c.processed {
			t.Errorf("In: %q\nGot: %q\nWant: %q\n", c.in, string(processed), c.processed)
		}
	}
}

func TestFullArticleReadWrite(t *testing.T) {
	article, err := ioutil.ReadFile("the_old_man_and_the_sea_except.txt")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	NewWriter(buf).Write(article)
	in := string(article)
	processed, err := ioutil.ReadAll(NewReader(buf))
	if err != nil {
		t.Errorf("Processing Error: %s (In: %q)", err, in)
	}
	want := strings.Replace(in, "\n", "\r\n", -1)
	if string(processed) != want {
		t.Errorf("In: %q\nGot: %q\nWant: %q\n", in, string(processed), want)
	}
}

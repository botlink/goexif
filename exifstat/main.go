package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/botlink/goexif/exif"
	"github.com/botlink/goexif/mknote"
	"github.com/botlink/goexif/tiff"
)

var mnote = flag.Bool("mknote", false, "try to parse makernote data")
var thumb = flag.Bool("thumb", false, "dump thumbail data to stdout (for first listed image file)")

func main() {
	flag.Parse()
	fnames := flag.Args()

	if *mnote {
		exif.RegisterParsers(mknote.All...)
	}

	for _, name := range fnames {
		f, err := os.Open(name)
		if err != nil {
			log.Printf("err on %v: %v", name, err)
			continue
		}

		x, err := exif.Decode(f)
		if err != nil {
			log.Printf("err on %v: %v", name, err)
			continue
		}

		if *thumb {
			data, err := x.JpegThumbnail()
			if err != nil {
				log.Fatal("no thumbnail present")
			}
			if _, err := os.Stdout.Write(data); err != nil {
				log.Fatal(err)
			}
			return
		}

		tg, _ := x.Get("Orientation")
		fmt.Printf("Orientation %v Offset %v\n", tg, tg.ValOffset)
		//x.Walk(Walker{})
	}
}

type Walker struct{}

func (_ Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %v: %v %v\n", name, string(data), tag.ValOffset)
	return nil
}

// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var data = []struct {
	URL   string
	Sub   string
	Store string
	Type  func(io.ReadCloser, string, string) error
}{
	{
		URL:   "http://dffd.bay12games.com/download.php?id=4&f=emeecamo_ametha.zip",
		Sub:   "Emeecamo_Ametha",
		Store: "dffd_0000004",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=573&f=worldofusmzaheatmenace.zip",
		Sub:   "worldofusmzaheatmenace/region9001",
		Store: "dffd_0000573",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=3810&f=region1_elfy-dwarves.zip",
		Sub:   "region1_elvy-dwarves",
		Store: "dffd_0003810",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=5154&f=Akrulatol.zip",
		Sub:   "Akrulatol",
		Store: "dffd_0005154",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=5574&f=region1.zip",
		Sub:   "region1",
		Store: "dffd_0005574",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=5930&f=rampage.zip",
		Sub:   "rampage",
		Store: "dffd_0005930",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=5994&f=region2.zip",
		Sub:   "region2",
		Store: "dffd_0005994",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=6331&f=save.zip",
		Sub:   "region2",
		Store: "dffd_0006331",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=6808&f=region3.zip",
		Sub:   "region3",
		Store: "dffd_0006808",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=7554&f=region4.zip",
		Sub:   "region4",
		Store: "dffd_0007554",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=8345&f=batman.zip",
		Sub:   "batman",
		Store: "dffd_0008345",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=10759&f=pocket_maxmonsters.zip",
		Sub:   "region5",
		Store: "dffd_0010759",
		Type:  Zip,
	},
	{
		URL:   "http://dffd.bay12games.com/download.php?id=10619&f=save.zip",
		Sub:   "save/region1",
		Store: "dffd_0010619",
		Type:  Zip,
	},
}

func Zip(r io.ReadCloser, sub, store string) error {
	b, err := ioutil.ReadAll(r)
	r.Close()
	if err != nil {
		return err
	}

	z, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}

	for _, f := range z.File {
		if f.Mode().IsRegular() && strings.HasPrefix(filepath.ToSlash(f.Name), sub) {
			path := filepath.Join("testdata", store, filepath.FromSlash(strings.TrimPrefix(filepath.ToSlash(f.Name), sub)))
			err = os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return err
			}
			r, err := f.Open()
			w, err := os.Create(path)
			if err != nil {
				r.Close()
				return err
			}
			_, err = io.Copy(w, r)
			r.Close()
			w.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	for _, d := range data {
		if _, err := os.Stat(filepath.Join("testdata", d.Store)); err == nil {
			continue
		}

		fmt.Println("downloading", d.URL)
		resp, err := http.Get(d.URL)
		if err != nil {
			panic(err)
		}
		err = d.Type(resp.Body, d.Sub, d.Store)
		if err != nil {
			panic(err)
		}
	}
}

// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xferspdy

import (
	"io"
	"os"

	"github.com/golang/glog"
)

//Patch is a wrapper on PatchFile (current version only supports patching of local files)
func Patch(delta []Block, sign Fingerprint, t io.Writer) error {
	return PatchFile(delta, sign.Source, t)
}

// PatchFile takes a source file and Diff as input, and writes out to the Writer.
// The source file would normally be the base version of the file  and
// the Diff is the delta computed by using the Fingerprint generated for the base file and the new version of the file
func PatchFile(delta []Block, source string, t io.Writer) error {
	s, err := os.Open(source)
	defer s.Close()
	if err != nil {
		return err
	}
	return PatchOpenedFile(delta, s, t)
}

func PatchOpenedFile(delta []Block, source *os.File, output io.Writer) error {
	wptr := int64(0)
	for _, block := range delta {
		if block.HasData {
			glog.V(3).Infof("Writing RawBytes block , wptr=%v , num bytes = %v \n", wptr, len(block.RawBytes))
			_, err := output.Write(block.RawBytes)
			glog.V(4).Infof("Writing bytes = %v \n", block.RawBytes)
			if err != nil {
				return err
			}
			wptr += int64(len(block.RawBytes))
		} else {
			source.Seek(block.Start, 0)
			ds := block.End - block.Start
			glog.V(3).Infof("Writing RawBytes block, Block=%v\n , wptr=%v , num bytes = %v \n", block, wptr, ds)
			if _, err := io.CopyN(output, source, block.End-block.Start); err != nil {
				return err
			}
			wptr += ds
		}
	}
	return nil
}

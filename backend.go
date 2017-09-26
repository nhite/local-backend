package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "github.com/nhite/pb-backend"
)

type backend struct {
	workingDir    string
	permission    os.FileMode
	fileExtension string
}

func (b *backend) Store(stream pb.Backend_StoreServer) error {
	if debug {
		log.Println("[store] =>")
	}
	body, err := stream.Recv()
	if err == io.EOF || err == nil {
		if debug {
			log.Printf("[store] ==> received flow with ID %v", body.GetID().GetID())
		}

		// By now let's use a gobencoding for flexibility
		var content bytes.Buffer

		// Create an encoder and send a value.
		enc := gob.NewEncoder(&content)
		err = enc.Encode(body)
		if err != nil {
			return err
		}
		err := ioutil.WriteFile(filepath.Join(b.workingDir, body.GetID().GetID()+b.fileExtension), content.Bytes(), b.permission)
		if err != nil {
			return err
		}
		if debug {
			log.Printf("[store] ==> content written in ", filepath.Join(b.workingDir, body.GetID().GetID()+b.fileExtension))
		}
	}
	if err != nil {
		return err
	}
	return stream.SendAndClose(&pb.Error{})
}
func (b *backend) Fetch(id *pb.ElementID, stream pb.Backend_FetchServer) error {
	var element *pb.Element
	body, err := ioutil.ReadFile(filepath.Join(b.workingDir, id.GetID()+b.fileExtension))
	if err != nil {
		return err
	}
	content := bytes.NewBuffer(body)
	// Create a decoder and receive a value.
	dec := gob.NewDecoder(content)
	err = dec.Decode(&element)
	if err != nil {
		return err
	}

	return stream.Send(element)
}

// TODO implement the list function
func (b *backend) List(context.Context, *pb.Pagination) (*pb.Elements, error) {
	return nil, nil
}

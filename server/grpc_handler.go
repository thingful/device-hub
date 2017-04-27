// Copyright © 2017 thingful

package server

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"

	hashids "github.com/speps/go-hashids"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/proto"
	context "golang.org/x/net/context"
)

type handler struct {
	manager *manager
	store   *store
}

func (s *handler) Create(ctx context.Context, request *proto.CreateRequest) (*proto.CreateReply, error) {

	hash, err := hash(request)

	if err != nil {
		return &proto.CreateReply{
			Error: err.Error(),
		}, nil
	}

	var bucket bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = listenersBucket

		exists := hub.IsListenerRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil

		}

	case "endpoint":
		bucket = endpointsBucket

		exists := hub.IsEndpointRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil
		}

	default:
		return &proto.CreateReply{
			Error: fmt.Sprintf("type : %s not registered", request.Type),
		}, nil
	}

	err = s.store.Insert(bucket, hash, proto.Endpoint{
		Uid:           string(hash),
		Type:          request.Type,
		Kind:          request.Kind,
		Configuration: request.Configuration,
	})

	if err != nil {
		return &proto.CreateReply{
			Error: err.Error(),
		}, nil
	}

	return &proto.CreateReply{
		Uid: string(hash),
		Ok:  true,
	}, nil
}

func hash(data interface{}) ([]byte, error) {

	j, err := json.Marshal(data)

	if err != nil {
		return []byte{}, err
	}

	checksum := crc32.ChecksumIEEE(j)
	h := hashids.New()

	uid, err := h.Encode([]int{int(checksum)})

	if err != nil {
		return []byte{}, err
	}
	return []byte(uid), nil

}

func (s *handler) Delete(ctx context.Context, request *proto.DeleteRequest) (*proto.DeleteReply, error) {
	return nil, nil
}

func (s *handler) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetReply, error) {

	keys := strings.Split(request.Filter, ",")

	all := []*proto.Endpoint{}

	for _, key := range keys {

		switch strings.ToLower(key) {
		case "listener", "l":

			listeners := []*proto.Endpoint{}

			err := s.store.List(listenersBucket, &listeners)

			if err != nil {
				return nil, err
			}
			all = append(all, listeners...)

		case "endpoint", "e":

			endpoints := []*proto.Endpoint{}

			err := s.store.List(endpointsBucket, &endpoints)

			if err != nil {
				return nil, err
			}
			all = append(all, endpoints...)

		default:
			return &proto.GetReply{
				Error: fmt.Sprintf("filter of type : %s not registered", key),
			}, nil
		}
	}

	return &proto.GetReply{Ok: true, Endpoints: all}, nil
}

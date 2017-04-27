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
			Ok:    false,
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
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil

		}

	case "endpoint":
		bucket = endpointsBucket

		exists := hub.IsEndpointRegistered(request.Kind)

		if !exists {
			return &proto.CreateReply{
				Ok:    false,
				Error: fmt.Sprintf("kind : %s not registered", request.Kind),
			}, nil
		}
	case "profile":
		bucket = endpointsBucket

	default:
		return &proto.CreateReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not registered", request.Type),
		}, nil
	}

	err = s.store.Insert(bucket, hash, proto.Entity{
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

	var bucket bucket

	switch strings.ToLower(request.Type) {
	case "listener":
		bucket = listenersBucket
	case "endpoint":
		bucket = endpointsBucket
	case "profile":
		bucket = profilesBucket
	default:
		return &proto.DeleteReply{
			Ok:    false,
			Error: fmt.Sprintf("type : %s not found", request.Type),
		}, nil
	}

	//TODO : error if any running pipes
	err := s.store.Delete(bucket, request.Uid)

	if err != nil {
		return &proto.DeleteReply{
			Ok:    false,
			Error: err.Error(),
		}, nil
	}

	return &proto.DeleteReply{Ok: true}, nil
}

func (s *handler) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetReply, error) {

	keys := strings.Split(request.Filter, ",")

	all := []*proto.Entity{}

	for _, key := range keys {

		switch strings.ToLower(key) {
		case "listener", "l":

			listeners := []*proto.Entity{}

			err := s.store.List(listenersBucket, &listeners)

			if err != nil {
				return nil, err
			}
			all = append(all, listeners...)

		case "endpoint", "e":

			endpoints := []*proto.Entity{}

			err := s.store.List(endpointsBucket, &endpoints)

			if err != nil {
				return nil, err
			}
			all = append(all, endpoints...)
		case "profile", "p":

			profiles := []*proto.Entity{}

			err := s.store.List(profilesBucket, &profiles)

			if err != nil {
				return nil, err
			}
			all = append(all, profiles...)

		default:
			return &proto.GetReply{
				Ok:    false,
				Error: fmt.Sprintf("filter of type : %s not registered", key),
			}, nil
		}
	}

	return &proto.GetReply{Ok: true, Entities: all}, nil
}

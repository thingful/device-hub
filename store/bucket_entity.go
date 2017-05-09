// Copyright © 2017 thingful

package store

import "github.com/thingful/device-hub/proto"

type entityBucket struct {
	bucket
}

func (b entityBucket) One(uid string) (*proto.Entity, error) {

	out := proto.Entity{}
	err := b.store.One(b.bucket, []byte(uid), &out)

	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (b entityBucket) Many(uids []string) ([]*proto.Entity, error) {

	out := make([]*proto.Entity, len(uids), len(uids))

	for i, _ := range uids {
		err := b.store.One(b.bucket, []byte(uids[i]), &out[i])

		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (b entityBucket) Delete(uid string) error {
	return b.store.Delete(b.bucket, []byte(uid))
}

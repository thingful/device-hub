// Copyright © 2017 thingful

package registry

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/proto"
)

func TestConfigurationFile(t *testing.T) {

	//for each file in ./test-configurations/
	folder := "./test-configurations/"
	listing, err := ioutil.ReadDir(folder)
	assert.Nil(t, err)

	for _, fi := range listing {

		folderPath := path.Join(folder, fi.Name())

		dm := iocodec.DefaultDecoders["yaml"]

		f, err := os.Open(folderPath)

		assert.Nil(t, err)

		in := dm.NewDecoder(f)
		entity := proto.Entity{}

		err = in.Decode(&entity)

		assert.Nil(t, err)
		assert.NotEmpty(t, entity.Kind)
		assert.NotEmpty(t, entity.Type)

		register := New()
		listener.Register(register)
	}

	//load file

	//parse and describe

}

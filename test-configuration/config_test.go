// Copyright © 2017 thingful

package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/fiorix/protoc-gen-cobra/iocodec"
	"github.com/stretchr/testify/assert"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/endpoint"
	"github.com/thingful/device-hub/listener"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/registry"
)

// TestConfigurationFilesAreValid loads all of the testing configuration files and checks them for errors
// The idea is that all of the files should be valid by default
func TestConfigurationFilesAreValid(t *testing.T) {

	register := registry.Default
	listener.Register(register)
	endpoint.Register(register)

	//for each file in ./test-configurations/
	folder := "./samples/"
	listing, err := ioutil.ReadDir(folder)

	for _, fi := range listing {

		assert.Nil(t, err)
		folderPath := path.Join(folder, fi.Name())

		if path.Ext(folderPath) == ".yaml" {

			dm := iocodec.DefaultDecoders["yaml"]

			f, err := os.Open(folderPath)

			assert.Nil(t, err)

			in := dm.NewDecoder(f)
			entity := proto.Entity{}

			err = in.Decode(&entity)

			assert.Nil(t, err)
			assert.NotEmpty(t, entity.Kind)
			assert.NotEmpty(t, entity.Type)

			var params describe.Parameters

			switch strings.ToLower(entity.Type) {

			case "listener":
				params, err = register.DescribeListener(entity.Kind)
			case "endpoint":
				params, err = register.DescribeEndpoint(entity.Kind)
			}

			assert.Nil(t, err)

			_, err = describe.NewValues(entity.Configuration, params)
			assert.Nil(t, err)

		}
	}

}

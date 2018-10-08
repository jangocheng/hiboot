// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/factory/autoconfigure"
	"reflect"
)

// appendParam is the common func to append meta data to meta data slice
func appendParam(eliminator string, container []*factory.MetaData, params ...interface{}) (retVal []*factory.MetaData, err error) {

	retVal = container

	// parse meta data
	metaData := factory.NewMetaData(params...)

	// append meta data
	if metaData.Object != nil {
		retVal = append(retVal, metaData)
		//return
	}
	//err = ErrInvalidObjectType
	return
}

// appendParams is the common func to append params to component or configuration containers
func appendParams(eliminator string, container []*factory.MetaData, params ...interface{}) (retVal []*factory.MetaData, err error) {
	retVal = container
	if len(params) == 0 || params[0] == nil {
		err = ErrInvalidObjectType
		return
	}

	if len(params) > 1 && reflect.TypeOf(params[0]).Kind() != reflect.String {
		for _, param := range params {
			retVal, err = appendParam(eliminator, retVal, param)
		}
	} else {
		retVal, err = appendParam(eliminator, retVal, params...)
	}
	return
}

// Register register a struct instance or constructor (func), so that it will be injectable.
func Register(params ...interface{}) (err error) {
	componentContainer, err = appendParams("", componentContainer, params...)
	return
}

// AutoConfiguration register auto configuration struct
func AutoConfiguration(params ...interface{}) (err error) {
	configContainer, err = appendParams(autoconfigure.PostfixConfiguration, configContainer, params...)
	return
}

// Component register all component into container
var Component = Register
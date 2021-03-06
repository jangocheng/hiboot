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

// Package at provides annotations for web RestController
package at

// RestController is the annotation that declare current controller is the RESTful Controller
type RestController interface{}

// JwtRestController is the annotation that declare current controller is the RESTful Controller with JWT support
type JwtRestController interface{}

// ContextPath is the annotation that set the context path of a controller
type ContextPath interface{}

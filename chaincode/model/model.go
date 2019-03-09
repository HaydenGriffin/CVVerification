// Copyright 2018 Antoine CHABERT, toHero.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

// Actor metadata used for an admin and a consumer
type Actor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Available actor type
const (
	ActorAttribute = "actor"
	ActorApplicant = "applicant"
	ActorVerifier = "verifier"
	ActorEmployer = "employer"
	ActorAdmin = "admin"
)

// Admin that manage resources available
type Admin struct {
	Actor
}

type Applicant struct {
	Actor
}

type UserProfile struct {
	Username	string	`json:"Name"`
	CVHistory []string `json:"CVHistory"`
	Ratings map[string] []CVRating
}

type CVObject struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`
	Speciality	string	`json:"Speciality"`
	CV	string	`json:"CV"`
	CVDate	string	`json:"CVDate"`
}

type CVRating struct {
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}

// List of object type stored in the ledger
const (
	ObjectTypeAdmin            = "admin"
	ObjectTypeConsumer         = "consumer"
	ObjectTypeResource         = "resource"
	ObjectTypeResourcesDeleted = "resources-deleted"
)

// List of available filter for query resources
const (
	ResourcesFilterAll             = "all"
	ResourcesFilterOnlyAvailable   = "only-available"
	ResourcesFilterOnlyUnavailable = "only-unavailable"
)

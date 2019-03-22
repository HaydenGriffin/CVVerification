/**
  @Author : Hayden Griffin
*/

package main

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



/**
  @Author : Hayden Griffin
*/

package main

import "time"

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

type CVHistory struct {
	Transaction string    `json:"transaction"`
	CV    CVObject  `json:"value"`
	Time        time.Time `json:"time"`
	Deleted     bool      `json:"deleted"`
}

type CVHistoryDetails []CVHistory

type CVRating struct {
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}



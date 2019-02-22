/**
  @Author : Hayden Griffin
*/

package main

/**

 */
type CVObject struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`
	Speciality	string	`json:"Speciality"`
	CV	string	`json:"CV"`
	CVHash	string	`json:"CVHash"`
	CVDate	string	`json:"CVDate"`
	History	[]HistoryItem
	Rating []CVRating
}

type CVRating struct {
	Name string `json:"Name"`
	Comment string `json:"Comment"`
	Rating int `json:"Rating"`
}

type HistoryItem struct {
	TxId	string
	CV	CVObject
}

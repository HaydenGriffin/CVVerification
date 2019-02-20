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
	CVHash	string	`json:"CVHash"`
	CVDate	string	`json:"CVDate"`
	History	[]HistoryItem
}

type HistoryItem struct {
	TxId	string
	Resume	CVObject
}

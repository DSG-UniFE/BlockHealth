

package main

import (
	"fmt"
	"encoding/json"
	"bytes"
	"strconv"
//	"time"
//	"container/list" //DA IMPORTARE PROBAB.
	"crypto/md5"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)


type TestsContract struct {
}


// Define Status codes for the response
const (
	OK    = 200
	ERROR = 500
)

// KEY = [Nome azienda; Nome e cognome persona; timestamp dell'esame] ; valore = lista dinamica di hash
// Altri parametri: tipo dell'esame; risultato (positivo/negativo); More (altro): qualsiasi altra cosa. 
type Test struct {
	Company         string		`json:"company"`
	Name		string		`json:"name"` 
   	Surname		string		`json:"surname"`
	Timestamp	string 		`json:"timestamp"`
	HashSlice       [][16]byte 	`json:"hash_slice"`
} //non posso usare così hashSlice, rivedere questa parte e testare con go test


// Init is called when the smart contract is instantiated
func (s *TestsContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Chaincode istanziato")
	return shim.Success(nil)
}



func (s *TestsContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "addtest" {			
		return s.addTest(APIstub, args)
	} else if function == "adddynamictest" {			
		return s.addDynamicTest(APIstub, args)
	} else if function == "gettests" {
		return s.getTests(APIstub, args)
	} else if function == "initledger" {
		return s.initLedger(APIstub, args)
	} 
	return shim.Error("Invalid Smart Contract function name on Test smart contract.")
}





func (s *TestsContract) initLedger (APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 0{ 
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}

	// creo due strutture dati di due esami
	var tests [2]Test

	//Esame 1
	tests[0].Company = "companyA" 
	tests[0].Name = "Person1"
	tests[0].Surname = "Surname1"
	tests[0].Timestamp = "0000000001"
	// dati non salvati in bc per privacy
	test_type0 := "covid19"
	result0 := "positive"
	more0 := "asymptomatic"


	//Esame 2
	tests[1].Company = "companyB" 
	tests[1].Name = "Person2"
	tests[1].Surname = "Surname2"
	tests[1].Timestamp = "0000000002"
	test_type1 := "pneumonia"
	result1 := "negative"
	more1 := "cough"


	
	//suddivido in 2 sotto gruppi di info per il primo esame (uso le slice)
	// solo esame senza esito e altro
	t1_sub_info := []string{tests[0].Company, tests[0].Name, tests[0].Surname, tests[0].Timestamp, test_type0}
	// tutti
	t1_sub_info2 := []string{tests[0].Company, tests[0].Name, tests[0].Surname, tests[0].Timestamp, test_type0, result0, more0}

	// preparo i dati per la funzione hash md5
	t1_sub_info_data,_ := json.Marshal(t1_sub_info)
	t1_sub_info2_data,_ := json.Marshal(t1_sub_info2)

	// inserisco gli hash md5 del primo esame
	tests[0].HashSlice = append(tests[0].HashSlice, md5.Sum(t1_sub_info_data), md5.Sum(t1_sub_info2_data))



	// suddivido in 3 sotto gruppi di info del secondo esame
	t2_sub_info := []string{tests[1].Company, tests[1].Name, tests[1].Surname, tests[1].Timestamp, test_type1}
	t2_sub_info2 := []string{tests[1].Company, tests[1].Name, tests[1].Surname, tests[1].Timestamp, test_type1, result1}
	t2_sub_info3 := []string{tests[1].Company, tests[1].Name, tests[1].Surname, tests[1].Timestamp, test_type1, result1, more1}

	// preparo i dati per la funzione hash md5
	t2_sub_info_data,_ := json.Marshal(t2_sub_info)
	t2_sub_info2_data,_ := json.Marshal(t2_sub_info2)
	var t2_sub_info3_data []byte 
	t2_sub_info3_data,_ = json.Marshal(t2_sub_info3)

	// inserisco gli hash md5 del secondo esame
	tests[1].HashSlice = append(tests[1].HashSlice, md5.Sum(t2_sub_info_data), md5.Sum(t2_sub_info2_data), md5.Sum(t2_sub_info3_data))



	// Salvo i due oggetti Test nel ledger
	var results []Test
	i := 0
	for i < len(tests) {
	//	fmt.Println("i is ", i)
		testAsBytes, _ := json.Marshal(tests[i])
		key := tests[i].Company + "-"+tests[i].Name+"-"+tests[i].Surname+"-"+tests[i].Timestamp
		APIstub.PutState(key, testAsBytes)
		fmt.Println("Added", tests[i])
		results = append(results, tests[i])
		i = i + 1
	}

	resultsAsBytes, _ := json.Marshal(results)
	// ritorno il risultato
	return shim.Success(resultsAsBytes)
}





func (s *TestsContract) addTest (APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 7{ 
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	// key,
	company := args[0]
	name := args[1]
	surname := args[2]
	timestamp := args[3]
	// salvare nel Ledger
	var hash_slice [][16]byte
	// altre info
	test_type := args[4]
	result := args[5]
	more := args[6]

	sub_info := []string{company, name, surname, timestamp, test_type}
	sub_info2 := []string{company, name, surname, timestamp, test_type, result, more}

	// preparo i dati per la funzione hash md5
	sub_info_data,_ := json.Marshal(sub_info)
	sub_info2_data,_ := json.Marshal(sub_info2)

	// inserisco gli hash md5
	hash_slice = append(hash_slice, md5.Sum(sub_info_data), md5.Sum(sub_info2_data))



	//ADD TEST
	key := company+"-"+name+"-"+surname+"-"+timestamp
	getState, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error from getState into addTest: %s", err.Error()))
	}
	//Se non è già presente questo oggetto con questa chiave, then create new test
	if bytes.Equal(getState,[]byte("")) {
		test := Test{company, name, surname, timestamp, hash_slice}
		testAsBytes, marshalErr := json.Marshal(test)
		if marshalErr != nil {
			return shim.Error(fmt.Sprintf("Could not marshal new %s test: %s", key, marshalErr.Error()))
		}
		putErr := APIstub.PutState(name, testAsBytes)
		if putErr != nil {
			return shim.Error(fmt.Sprintf("Could not put new %s test in the ledger: %s", key, putErr.Error()))
		}

		fmt.Println("Added new test: ", test)
		return shim.Success([]byte(fmt.Sprintf("Successfully added %s test",  key )))
	}
	return shim.Error("Error in addTest: Test already exist with this key.")
}














func (s *TestsContract) addDynamicTest (APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if (len(args) < 8  || len(args) >15 ) {
		return shim.Error("Incorrect number of arguments. Expecting at least 8; at most 15")
	}

	// key, salvare nel ledger
	company := args[0]
	name := args[1]
	surname := args[2]
	timestamp := args[3]
	// salvare nel Ledger
	var hash_slice [][16]byte
	// altre info
	test_type := args[4]
	result := args[5]
	more := args[6:len(args)-1] //dal sesto al penultimo (ultimo non compreso)
	n_hash,_ := strconv.Atoi(args[len(args)-1]) //8 args -> n_hash in più = 0
	
	if (n_hash > len(more) + 2) {
		return shim.Error("Incorrect number of arguments. Expecting hash number is less or equal then number of params")
	}

	// key + test type
	sub_info := []string{company, name, surname, timestamp, test_type}
	sub_info_data,_ := json.Marshal(sub_info)
	hash_slice = append(hash_slice, md5.Sum(sub_info_data))

	// key + test type + test result 
	sub_info = append(sub_info, result)
	sub_info_data,_ = json.Marshal(sub_info)
	hash_slice = append(hash_slice, md5.Sum(sub_info_data))

	// other hash
	for i := 0; i < n_hash; i++ {
		sub_info = append(sub_info, more[i])
		sub_info_data,_ = json.Marshal(sub_info)
		hash_slice = append(hash_slice, md5.Sum(sub_info_data))
	}
	
	//ADD TEST
	key := company+"-"+name+"-"+surname+"-"+timestamp
	getState, err := APIstub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error from getState into addTest: %s", err.Error()))
	}
	//Se non è già presente questo oggetto con questa chiave, then create new test
	if bytes.Equal(getState,[]byte("")) {
		test := Test{company, name, surname, timestamp, hash_slice}
		testAsBytes, marshalErr := json.Marshal(test)
		if marshalErr != nil {
			return shim.Error(fmt.Sprintf("Could not marshal new %s test: %s", key, marshalErr.Error()))
		}
		putErr := APIstub.PutState(name, testAsBytes)
		if putErr != nil {
			return shim.Error(fmt.Sprintf("Could not put new %s test in the ledger: %s", key, putErr.Error()))
		}

		fmt.Println("Added new test: ", test)
		return shim.Success([]byte(fmt.Sprintf("Successfully added %s test",  key )))
	}
	return shim.Error("Error in addTest: Test already exist with this key.")
}










func (s *TestsContract) getTests(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// Check we have a valid number of args
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments, expecting 0")
	}


	resultsIterator, err := APIstub.GetStateByRange("","")
	if err != nil {
		fmt.Println("Errore getStateByRange")
		return shim.Error(fmt.Sprintf("Errore getStateByRange -> %s",err.Error()))
	}
	defer resultsIterator.Close()

	var tests []Test //[]byte //il risultato di tutte le macchine
	var test Test //byte //variabile machine temporanea per poi assegnarla all'array con append
 
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		json.Unmarshal(queryResponse.Value, &test)
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		return shim.Error("Errore, tests array is empty.")
	} else {
		fmt.Println("Len Tests array: ",len(tests))
	}

	testsAsBytes, _ := json.Marshal(tests)
	return shim.Success(testsAsBytes)
}












// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(TestsContract))
	if err != nil {
		fmt.Printf("Error creating new Tests Smart Contract: %s", err)
	}
}




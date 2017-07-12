package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("FireArms")

//ALL_ELEMENENTS Key to refer the list of application
const ALL_ELEMENENTS = "ALL_APP"

//FireArms Chaincode default interface
type FireArms struct {
}

//Append a new fireArms appid to the master list
func updateMasterRecords(stub shim.ChaincodeStubInterface, appId string) error {
	var recordList []string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)

	err := json.Unmarshal(recBytes, &recordList)
	if err != nil {
		return errors.New("Failed to unmarshal updateMasterReords ")
	}
	recordList = append(recordList, appId)
	bytesToStore, _ := json.Marshal(recordList)
	logger.Info("After addition" + string(bytesToStore))
	stub.PutState(ALL_ELEMENENTS, bytesToStore)
	return nil
}

// Creating a new fireArm Application
func createApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("createfireArm application called")
	var id string
	var data map[string]string
	valAsbytes, err := stub.GetState("id")
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for id\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println((string)(valAsbytes))
	id = (string)(valAsbytes)

	fmt.Println("new id is" + id)
	uniqueId, _ := strconv.Atoi(id)
	appId := "AppId" + id
	newid := uniqueId + 1
	stub.PutState("id", []byte(strconv.Itoa(newid)))
	payload := args[0]
	json.Unmarshal([]byte(payload), &data)
	email := data["appemail"]
	fmt.Println("email is: " + email)
	stub.PutState(email, []byte(appId))
	fmt.Println("new Payload is " + payload)

	stub.PutState(appId, []byte(payload))

	updateMasterRecords(stub, appId)
	logger.Info("Created the FireArms")

	return nil, nil
}

//Validate a input string as number or not
func validateNumber(str string) float64 {
	if netCharge, err := strconv.ParseFloat(str, 64); err == nil {
		return netCharge
	}
	return float64(-1.0)
}

//Update the existing record with the mofied key value pair
func updateRecord(existingRecord map[string]string, fieldsToUpdate map[string]string) (string, error) {
	for key, value := range fieldsToUpdate {

		existingRecord[key] = value
	}
	outputMapBytes, _ := json.Marshal(existingRecord)
	logger.Info("updateRecord: Final json after update " + string(outputMapBytes))
	return string(outputMapBytes), nil
}

func probe() []byte {
	ts := time.Now().Format(time.UnixDate)
	output := "{\"status\":\"Success\",\"ts\" : \"" + ts + "\" }"
	return []byte(output)
}

// Init initializes the smart contracts
func (t *FireArms) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Init called")
	//Place an empty arry
	stub.PutState(ALL_ELEMENENTS, []byte("[]"))
	stub.PutState("id", []byte("1"))
	stub.PutState("license", []byte("1000"))
	return nil, nil
}

// Invoke entry point
func (t *FireArms) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Invoke called")

	if function == "createApplication" {
		createApplication(stub, args)
	} else if function == "updateApplication" {
		updateApplication(stub, args)
	}

	return nil, nil
}

// Query the rcords form the  smart contracts
func (t *FireArms) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Query called")
	if function == "getAppById" {
		return getAppById(stub, args[0])
	} else if function == "getAppByEmailId" {
		return getAppByEmailId(stub, args[0])
	} else if function == "getAllApp" {
		return getAllApp(stub)
	} else if function == "getAllAppByStatus" {
		return getAllAppByStatus(stub, args)
	} else if function == "getAppForRefree" {
		return getAppForRefree(stub, args[0])
	} else if function == "getLicenseByLicenseId" {
		return getLicenseByLicenseId(stub, args[0])
	}

	return nil, nil
}

//Get a single Application
func getAppById(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	logger.Info("getAppById called with AppId: " + args)

	var outputRecord map[string]string
	appid := args //AppId
	recBytes, _ := stub.GetState(appid)
	json.Unmarshal(recBytes, &outputRecord)
	outputBytes, _ := json.Marshal(outputRecord)
	logger.Info("Returning records from getAppId " + string(outputBytes))
	return outputBytes, nil
}

func getLicenseByLicenseId(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	logger.Info("getLicenseByLicenseId called with getLicenseByLicenseId: " + args)

	var outputRecord map[string]string
	lId := args //AppId
	recBytes, _ := stub.GetState(lId)
	json.Unmarshal(recBytes, &outputRecord)
	outputBytes, _ := json.Marshal(outputRecord)
	logger.Info("Returning records from getLicenseByLicenseId " + string(outputBytes))
	return outputBytes, nil
}

//Get a single Application on the basis of email id
func getAppByEmailId(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	logger.Info("getAppById called with AppId: " + args)

	id, err := stub.GetState(args)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for id\"}"
		return nil, errors.New(jsonResp)
	}
	if (string)(id) == "" {
		jsonResp := "{\"Error\":\"Failed to get state for this email " + args + "\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("The email id is" + (string)(id))

	var outputRecord map[string]string
	appid := (string)(id) //AppId
	recBytes, _ := stub.GetState(appid)
	json.Unmarshal(recBytes, &outputRecord)
	outputRecord["applicationNumber"] = appid
	outputBytes, _ := json.Marshal(outputRecord)
	logger.Info("Returning records from getAppId " + string(outputBytes))
	return outputBytes, nil
}

//Get all the Application based on the status
func getAllAppByStatus(stub shim.ChaincodeStubInterface, status []string) ([]byte, error) {
	logger.Info("; called")
	var recordList []string
	var allApp []map[string]string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)
	json.Unmarshal(recBytes, &recordList)

	for _, value := range recordList {
		logger.Info("inside getAllAppByStatus range func")
		recBytes, _ := getAppById(stub, value)
		var record map[string]string
		json.Unmarshal(recBytes, &record)
		for _, data := range status {
			if record["status"] == data {
				record["applicationNumber"] = value
				allApp = append(allApp, record)
			}
		}
	}
	outputBytes, _ := json.Marshal(allApp)
	return outputBytes, nil
}

//get the application for a refree
func getAppForRefree(stub shim.ChaincodeStubInterface, refreeEmail string) ([]byte, error) {
	logger.Info("getAppForRefree called")
	var recordList []string
	var allApp []map[string]string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)
	json.Unmarshal(recBytes, &recordList)

	for _, value := range recordList {
		logger.Info("inside getAppForRefree range func")
		recBytes, _ := getAppById(stub, value)
		var record map[string]string
		json.Unmarshal(recBytes, &record)

		if record["referee1email"] == refreeEmail {
			record["applicationNumber"] = value
			allApp = append(allApp, record)
		} else if record["referee2email"] == refreeEmail {
			record["applicationNumber"] = value
			allApp = append(allApp, record)
		}
	}
	outputBytes, _ := json.Marshal(allApp)
	return outputBytes, nil
}

//get all applications
func getAllApp(stub shim.ChaincodeStubInterface) ([]byte, error) {
	logger.Info("getAllApp called")
	var recordList []string
	var allApp []map[string]string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)
	json.Unmarshal(recBytes, &recordList)

	for _, value := range recordList {
		logger.Info("inside getallApp range func")
		recBytes, _ := getAppById(stub, value)

		var record map[string]string
		json.Unmarshal(recBytes, &record)
		record["applicationNumber"] = value
		allApp = append(allApp, record)
	}
	outputBytes, _ := json.Marshal(allApp)

	return outputBytes, nil
}

//update application
func updateApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var existingRecMap map[string]string
	var updatedFields map[string]string

	var appId string
	logger.Info("update Application called ")

	payload := args[0]
	logger.Info("update Application payload  " + payload)
	json.Unmarshal([]byte(payload), &updatedFields)

	for key, value := range updatedFields {
		if key == "applicationNumber" {
			appId = value
			recBytes, _ := stub.GetState(value)
			if recBytes == nil {
				jsonResp := "{\"Error\":\"No records available for this id " + key + "\"}"
				return nil, errors.New(jsonResp)
			}
			json.Unmarshal(recBytes, &existingRecMap)

		}
	}
	//now generate the license no., create a seperate json for license and deploy
	if updatedFields["status"] == "ISSUED" {

		valAsbytes, err := stub.GetState("license")
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for license\"}"
			return nil, errors.New(jsonResp)
		}
		fmt.Println((string)(valAsbytes))
		id := (string)(valAsbytes)

		fmt.Println("new license no. is" + id)
		uniqueId, _ := strconv.Atoi(id)
		license := "L" + id
		newLicenseid := uniqueId + 1
		stub.PutState("license", []byte(strconv.Itoa(newLicenseid)))

		licenseJson, _ := createLicenseJson(existingRecMap)
		stub.PutState(license, []byte(licenseJson))
	}

	updatedReord, _ := updateRecord(existingRecMap, updatedFields)
	stub.PutState(appId, []byte(updatedReord))
	return nil, nil
}

//Update the existing record with the mofied key value pair
func createLicenseJson(existingRecMap map[string]string) (string, error) {

	var licenseRecord map[string]string

	for key, value := range existingRecMap {

		if key == "fname" || key == "lname" || key == "email" || key == "phone" || key == "gender" || key == "firearms" || key == "applicationNumber" {

			licenseRecord[key] = value

		}

	}

	licenseRecord["weaponname"] = ""
	licenseRecord["dateofpurchase"] = ""

	outputMapBytes, _ := json.Marshal(licenseRecord)
	logger.Info("createLicenseJson : Final json after create is  " + string(outputMapBytes))
	return string(outputMapBytes), nil
}

//Main method
func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(FireArms))
	if err != nil {
		fmt.Printf("Error starting FireArms: %s", err)
	}
}

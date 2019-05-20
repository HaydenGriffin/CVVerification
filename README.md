# Hayden Griffin FYP
#### Verifiable Curriculum Vitae via Blockchain
The files contained within the fabric-network/bin subdirectory are obtained from Hyperledger Fabric and are used to generate the cryptographic material for the Fabric Network.

##### Setup
Before launching the application, please make sure that the following software has been downloaded and installed:
* Docker (https://www.docker.com/get-started)
* Golang (https://golang.org/dl/)
* MySQL Workbench (https://dev.mysql.com/downloads/workbench/)
 
Please extract the whole CVVerification project folder to a location within your GOPATH, or reset your GOPATH to the location of the CVVerification project folder.

A database creation script has been included, located at app/database/database_creation_script.sql. After MySQL Workbench has been installed, run this script within the MySQL Workbench application to create the necessary database tables.


##### Launching the application
To launch the application for the first time, execute the command 'make rebuild' from the terminal. This command will launch the Fabric network containers, retrieve all Go package dependencies, install the Chaincode on the peers, register the users, cleardown the database tables and then launch the Web Server.

To continue the application from the previous state, execute the command 'make'. This command will relaunch the Fabric network containers, retrieve all Go package dependencies and then launch the Web Server. The user data from previous executions will still be available and all ledger data is saved.

To access the application, open a browser window and navigate to http://localhost:3000/. A dialog box will be displayed requesting a username and password. The usernames for the  accounts for the application are as follows:
* admin1
* applicant1
* applicant2
* verifier1
* verifier2
* employer1

The password for all of the accounts is 'password'.

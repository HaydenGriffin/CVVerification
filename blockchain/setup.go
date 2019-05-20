/**
  author: Hayden Griffin
 */

package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	caMsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// FabricSetup struct; contains all the variables needed to launch the SDK
type FabricSetup struct {
	ConfigFile       string
	OrgID            string
	OrdererID        string
	ChannelID        string
	ChaincodeID      string
	ChannelConfig    string
	ChaincodeGoPath  string
	ChaincodeVersion string
	ChaincodePath    string
	OrgAdmin         string
	OrgMspID         string
	OrdererOrgID     string
	OrgName          string
	UserName         string
	CaID			 string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
	caClient        *caMsp.Client
}

// User stuct that allow a registered user to query and invoke the blockchain
type User struct {
	Username        string
	Fabric          *FabricSetup
	ChannelClient   *channel.Client
	SigningIdentity msp.SigningIdentity
}

// Initialize and configure the SDK
func (setup *FabricSetup) Initialize() error {
	fmt.Println("Initialising SDK")

	// Use the created config file to configure the SDK
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "SDK could not be configured")
	}
	setup.sdk = sdk

	caClient, err := caMsp.New(sdk.Context())
	if err != nil {
		return fmt.Errorf("failed to create new CA client: %v", err)
	}

	setup.caClient = caClient

	fmt.Println("SDK initialised")
	return nil
}

func (setup *FabricSetup) InstallChaincode() (*channel.Client, error) {

	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))

	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	setup.admin = resMgmtClient
	fmt.Println("Resource management client created")

	// Create channel
	err = setup.createChannel(resMgmtClient)
	if err != nil {
		return nil, fmt.Errorf("unable to create the channel: %v", err)
	}

	// Make admin user join the previously created channel
	if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		return nil, errors.WithMessage(err, "failed to make admin join channel")
	}

	fmt.Println("Channel joined")

	fmt.Printf("Install chaincode...\n")

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return nil, err
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{
		Name: setup.ChaincodeID,
		Path: setup.ChaincodePath,
		Version: setup.ChaincodeVersion,
		Package: ccPkg,
	}
	_, err = setup.admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return nil, err
	}

	fmt.Println("Chaincode installed")

	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.cvverification.com"})

	resp, err := setup.admin.InstantiateCC(
		setup.ChannelID,
		resmgmt.InstantiateCCRequest{
			Name: setup.ChaincodeID,
			Path: setup.ChaincodeGoPath,
			Version: setup.ChaincodeVersion,
			Args: [][]byte{[]byte("init")},
			Policy: ccPolicy,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Chaincode '%s' (version '%s') instantiated with transaction ID '%s'\n", setup.ChaincodeID, setup.ChaincodeVersion, resp.TransactionID)

	return setup.client, nil
}

// LogUser allow to login a user using credentials provided and retrieve the blockchain user related
func (setup *FabricSetup) LogUser(username, password string) (*User, error) {

	err := setup.caClient.Enroll(username, caMsp.WithSecret(password))
	if err != nil {
		return nil, fmt.Errorf("failed to enroll identity '%s': %v", username, err)
	}

	var user User
	user.Username = username
	user.Fabric = setup

	user.SigningIdentity, err = setup.caClient.GetSigningIdentity(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get signing identity for '%s': %v", username, err)
	}

	clientChannelContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(username), fabsdk.WithOrg(setup.OrgID), fabsdk.WithIdentity(user.SigningIdentity))

	// Channel client is used to query and execute transactions
	user.ChannelClient, err = channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("failed to create new channel client for '%s': %v", username, err)
	}

	return &user, nil
}

// RegisterUser register a user to the Fabric CA client and into the blockchain using invoke on the chaincode
func (setup *FabricSetup) RegisterUser(username, password, userType string) error {
	fmt.Printf("Register user '%s'... \n", username)
	_, err := setup.caClient.Register(&caMsp.RegistrationRequest{
		Name:           username,
		Secret:         password,
		Type:           "user",
		MaxEnrollments: -1,
		Affiliation:    "org1",
		Attributes: []caMsp.Attribute{
			{
				Name:  "actor",
				Value: userType,
				ECert: true,
			},
		},
		CAName: setup.CaID,
	})
	if err != nil {
		return fmt.Errorf("unable to register user '%s': %v", username, err)
	}

	u, err := setup.LogUser(username, password)
	if err != nil {
		return fmt.Errorf("unable to log user '%s' after registration: %v", username, err)
	}

	err = u.UpdateRegister()
	if err != nil {
		return fmt.Errorf("unable to add the user '%s' in the ledger: %v", username, err)
	}

	fmt.Printf("User '%s' registered.\n", username)

	return nil
}


// createChannel internal method that allow to create a channel in the blockchain network
func (setup *FabricSetup) createChannel(resMgmtClient *resmgmt.Client) error {
	fmt.Printf("Creating channel...\n")

	mspClient, err := mspclient.New(setup.sdk.Context(), mspclient.WithOrg(setup.OrgID))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return err
	}

	req := resmgmt.SaveChannelRequest{
		ChannelID: setup.ChannelID,
		ChannelConfigPath: setup.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}

	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Printf("Channel '%s' created with transaction ID '%s'\n", setup.ChannelID, txID.TransactionID)
	return nil
}


func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}
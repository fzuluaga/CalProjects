package client

// CS 161 Project 2

// Only the following imports are allowed! ANY additional imports
// may break the autograder!
// - bytes
// - encoding/hex
// - encoding/json
// - errors
// - fmt
// - github.com/cs161-staff/project2-userlib
// - github.com/google/uuid
// - strconv
// - strings

import (
	// "bytes"
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	// "strings"
	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// HELPER FUNCTIONS :

// Computes the HMAC for User struct --> doesnt cover encryption
func UserMAC(user User, salt []byte, key []byte) (HMAC []byte, err error) {
	userTemp := user.MAC
	user.MAC = nil
	marshalledStruct, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	encStruct := userlib.SymEnc(key, salt, marshalledStruct)
	userMAC, err := userlib.HMACEval(key, encStruct)
	if err != nil {
		return nil, err
	}
	user.MAC = userTemp
	return userMAC, nil
}

// Computes the HMAC for AppendBlock struct --> doesnt cover encryption
func AppendMAC(append AppendBlock, salt []byte, blockKey []byte) (HMAC []byte, err error) {
	appendTemp := append.MAC
	append.MAC = nil
	marshalledStruct, err := json.Marshal(append)
	if err != nil {
		return nil, err
	}
	encStruct := userlib.SymEnc(blockKey, salt, marshalledStruct)
	appendMAC, err := userlib.HMACEval(blockKey, encStruct)
	if err != nil {
		return nil, err
	}
	append.MAC = appendTemp
	return appendMAC, nil
}

// Computes the HMAC for FileInfo struct --> doesnt cover encryption
func FileMAC(file FileInfo, salt []byte, key []byte) (HMAC []byte, err error) {
	structTemp := file.MAC
	file.MAC = nil
	marshalledStruct, err := json.Marshal(file)
	if err != nil {
		return nil, err
	}
	encStruct := userlib.SymEnc(key, salt, marshalledStruct)
	fileMAC, err := userlib.HMACEval(key, encStruct)
	if err != nil {
		return nil, err
	}
	file.MAC = structTemp
	return fileMAC, nil
}

// Computes the HMAC for the Certificate struct --> doesnt cover encryption NEED FIXING !!!!
func CertMAC(cert Certificates, key []byte) (HMAC []byte, err error) {
	structTemp := cert.MAC
	cert.MAC = nil
	marshalledCert, err := json.Marshal(cert)
	if err != nil {
		return nil, err
	}
	encCert := userlib.SymEnc(key, cert.Salt, marshalledCert)
	CertMAC, err := userlib.HMACEval(key, encCert)
	if err != nil {
		return nil, err
	}
	cert.MAC = structTemp
	return CertMAC, nil
}

func AppendDataMAC(appendData AppendData, key []byte) (HMAC []byte, err error) {
	structTemp := appendData.MAC
	appendData.MAC = nil
	marshalledAppendData, err := json.Marshal(appendData)
	if err != nil {
		return nil, err
	}
	encData := userlib.SymEnc(key, appendData.Salt, marshalledAppendData)
	CertData, err := userlib.HMACEval(key, encData)
	if err != nil {
		return nil, err
	}
	appendData.MAC = structTemp
	return CertData, err
}

func hybridGetEncKey(publicKey userlib.PKEEncKey, symKey []byte) (encSymKey []byte, err error) {
	// encrypt symKey with public key
	encSymKey, err = userlib.PKEEnc(publicKey, symKey)
	if err != nil {
		return nil, err
	}
	return encSymKey, nil
}

func hybridGetSymKey(privateKey userlib.PKEDecKey, encSymKey []byte) (symKey []byte, err error) {
	// decrypt symKey with private key
	symKey, err = userlib.PKEDec(privateKey, encSymKey)
	if err != nil {
		return nil, err
	}
	return symKey, nil
}

func getCertStructKeyUUID(sender string, recipient string, certPtr userlib.UUID) (certStructKeyUUID uuid.UUID, err error) {
	certStructKeyUUIDBytes := userlib.Hash([]byte(sender + " to " + recipient + certPtr.String()))
	certStructKeyUUID, err = uuid.FromBytes(certStructKeyUUIDBytes[:16])
	if err != nil {
		return uuid.Nil, err
	}
	return certStructKeyUUID, nil
}

// called by sender to create a encrypted cert struct
func (userdata *User) certificateEncryption(sender string, recipient string, fileName string, cert Certificates) (encCertStruct []byte, encCertStructUUID uuid.UUID, err error) {
	symKey := userlib.RandomBytes(16)
	encKey, exists := userlib.KeystoreGet(recipient + " encKey")
	if !exists {
		return nil, uuid.Nil, errors.New("Could not find key")
	}
	// encrypt the symKey with recipient's public key
	encSymKey, err := hybridGetEncKey(encKey, symKey)
	if err != nil {
		return nil, uuid.Nil, err
	}
	// sign the enc key and save the signature in keySig
	keySig, err := userlib.DSSign(userdata.SignKey, encSymKey)
	if err != nil {
		return nil, uuid.Nil, err
	}
	// store signature in Datastore w/ arbitrary UUID
	signatureUUID := uuid.New()
	userlib.DatastoreSet(signatureUUID, keySig)

	// update this info in the Cert struct before marshal
	cert.SignatureUUID = signatureUUID
	cert.MAC, err = CertMAC(cert, symKey)
	if err != nil {
		return nil, uuid.Nil, err
	}

	// marshal the Cert Struct
	certBytes, err := json.Marshal(cert)
	if err != nil {
		return nil, uuid.Nil, err
	}
	// encrypt the cert struct with symKey
	encCert := userlib.SymEnc(symKey, cert.Salt, certBytes)

	encCertStructUUID = uuid.New()
	// Store encrypted struct at encCertStructUUID
	userlib.DatastoreSet(encCertStructUUID, encCert)

	// get CertStruct Key UUID
	structKeyUUID, err := getCertStructKeyUUID(sender, recipient, encCertStructUUID)
	if err != nil {
		return nil, uuid.Nil, err
	}
	// Store encrypted key at structKeyUUID
	userlib.DatastoreSet(structKeyUUID, encSymKey)

	userdata.Invites[fileName] = sender

	return encCert, encCertStructUUID, nil
}

func (userdata *User) certificateReencryption(sender string, recipient string, fileName string, certUUID uuid.UUID, cert Certificates) (encCertStruct []byte, err error) {
	structKeyUUID, err := getCertStructKeyUUID(sender, recipient, certUUID)
	if err != nil {
		return nil, err
	}
	// grab the decryption key
	decKey := userdata.DecryptKey
	encSymKey, exists := userlib.DatastoreGet(structKeyUUID)
	if !exists {
		return nil, errors.New("Error finding Encrypted SymKey in Datastore")
	}
	// decrypt the encSymKey with private key
	symKey, err := hybridGetSymKey(decKey, encSymKey)
	if err != nil {
		return nil, err
	}
	// REMAC certificate struct and update it in Datastore
	cert.MAC, err = CertMAC(cert, symKey)
	if err != nil {
		return nil, err
	}
	certBytes, err := json.Marshal(cert)
	if err != nil {
		return nil, err
	}
	// encrypt the Certificate Struct
	encCert := userlib.SymEnc(symKey, cert.Salt, certBytes)

	// Store encrypted struct at encCertStructUUID
	userlib.DatastoreSet(certUUID, encCert)

	userdata.Invites[fileName] = sender

	return encCert, nil
}

// called by the recipient to decrypt the certificate struct
func (userdata *User) certificateDecryption(sender string, recipient string, fileName string, certPtr uuid.UUID) (decCertStruct []byte, err error) {
	userdata, err = usernameToUserStruct(userdata.Username)
	if err != nil {
		return nil, err
	}

	// if user owns the file
	if sender == "" && fileName != "" {
		_, exists := userdata.Invites[fileName]
		if !exists {
			return nil, errors.New("File Name not in Invites")
		}
		sender = userdata.Invites[fileName]
	}

	structKeyUUID, err := getCertStructKeyUUID(sender, recipient, certPtr)
	if err != nil {
		return nil, err
	}

	// grab the decryption key
	decKey := userdata.DecryptKey
	encSymKey, exists := userlib.DatastoreGet(structKeyUUID)
	if !exists {
		return nil, errors.New("Error finding Encrypted SymKey in Datastore")
	}
	// decrypt the encSymKey with private key
	symKey, err := hybridGetSymKey(decKey, encSymKey)
	if err != nil {
		return nil, err
	}
	// grab the encrypted cert struct from Datastore
	encCertStruct, exists := userlib.DatastoreGet(certPtr)
	if !exists {
		return nil, errors.New("Error finding Encrypted Certificate Struct in Datastore")
	}
	// decrypt the cert struct with symKey
	decCertStructBytes := userlib.SymDec(symKey, encCertStruct)
	var certStruct Certificates
	err = json.Unmarshal(decCertStructBytes, &certStruct)
	if err != nil {
		return nil, err
	}
	// verify no tampering with MAC
	certStructMAC, err := CertMAC(certStruct, symKey)
	if err != nil {
		return nil, err
	}
	// verify the MAC
	hmacCheck := userlib.HMACEqual(certStructMAC, certStruct.MAC)
	if !(hmacCheck) {
		return nil, errors.New("Error verifying MAC of Certificate Struct")
	}
	// verify no tampering with File
	fileInfoUUID := certStruct.FileInfo
	encFileInfo, exists := userlib.DatastoreGet(fileInfoUUID)
	if !exists {
		return nil, errors.New("Error finding FileInfo Struct in Datastore")
	}
	// use Access Token to decrypt fileinfo struct
	accessToken := certStruct.AccessToken
	decFileInfo := userlib.SymDec(accessToken, encFileInfo)
	// need to error if the accessToken is wrong
	var fileInfo FileInfo
	err = json.Unmarshal(decFileInfo, &fileInfo)
	if err != nil {
		return nil, err
	}
	// verify no tampering with MAC
	fileInfoMAC, err := FileMAC(fileInfo, fileInfo.Salt, accessToken)
	if err != nil {
		return nil, err
	}
	// verify the MAC
	hmacCheck = userlib.HMACEqual(fileInfoMAC, fileInfo.MAC)
	if !(hmacCheck) {
		return nil, errors.New("Error verifying MAC of FileInfo Struct")
	}

	// grab the signature
	signature, exists := userlib.DatastoreGet(certStruct.SignatureUUID)
	if !exists {
		return nil, errors.New("Error finding Certificate Struct Signature in Datastore")
	}
	// verify the signature w/ public key
	verifyKey, exists := userlib.KeystoreGet(sender + " verifyKey")
	if !exists {
		return nil, errors.New("Error finding Verify Key in Keystore")
	}
	userlib.DSVerify(verifyKey, encSymKey, signature) // is there a realistic point to this if u can make it here?
	return decCertStructBytes, nil
}

// Checks every AppendBlock MAC in the chain to verify integrity
func traverseAppendBlock(firstAppendBlock uuid.UUID, blockKey []byte) error {
	currUUID := firstAppendBlock
	// read the filedata from start append, until last append, using next append field.
	for currUUID != uuid.Nil {
		encCurrAppend, exists := userlib.DatastoreGet(currUUID)
		if !exists {
			return errors.New("Error getting CurrAppend from Datastore")
		}

		decCurrAppend := userlib.SymDec(blockKey, encCurrAppend)

		var currAppendBlock AppendBlock
		err := json.Unmarshal(decCurrAppend, &currAppendBlock)
		if err != nil {
			return err
		}
		correctCurrHMAC, err := AppendMAC(currAppendBlock, currAppendBlock.Salt, blockKey)
		if err != nil {
			return err
		}
		currAppendBlockHMAC := currAppendBlock.MAC
		hmacCheckBlock := userlib.HMACEqual(currAppendBlockHMAC, correctCurrHMAC)
		if !(hmacCheckBlock) {
			return errors.New("Failed verification test on Append Block Struct")
		}
		currAppendDataBytes, exists := userlib.DatastoreGet(currAppendBlock.FileData)
		if !exists {
			return errors.New("Couldn't find AppendData in Datastore")
		}
		currAppendDataDec := userlib.SymDec(blockKey, currAppendDataBytes)
		var currAppendData AppendData
		err = json.Unmarshal(currAppendDataDec, &currAppendData)
		if err != nil {
			return err
		}
		currAppendDataMAC := currAppendData.MAC
		currAppendDataHMAC, err := AppendDataMAC(currAppendData, blockKey)
		if err != nil {
			return err
		}
		hmacCheckData := userlib.HMACEqual(currAppendDataMAC, currAppendDataHMAC)
		if !(hmacCheckData) {
			return errors.New("Failed verification test on Append Data Struct")
		}

		currUUID = currAppendBlock.NextAppend
	}
	return nil
}

// Using the username, grab the User Struct from the Datastore, grab its Cert Struct and update the AccessToken
func (userdata *User) updateToken(username string, parentname string, certificateUUID uuid.UUID, newAccessToken []byte) (err error) {
	// grab the user struct from the Datastore
	recipientHash := userlib.Hash([]byte(username))[:16]
	recipientUUID, err := uuid.FromBytes(recipientHash)
	if err != nil {
		return err
	}
	recipientEncPassword, exists := userlib.DatastoreGet(recipientUUID)
	if !exists {
		return errors.New("Error finding User Struct in Datastore")
	}
	// find the UUID of the encrypted user struct
	recipientPassHKDF, err := userlib.HashKDF(recipientEncPassword[:16], []byte("UUID"))
	if err != nil {
		return err
	}
	recipientPassUUID, err := uuid.FromBytes(recipientPassHKDF)
	if err != nil {
		return err
	}
	// grab the encrypted user struct
	recipientEncUser, exists := userlib.DatastoreGet(recipientPassUUID)
	if !exists {
		return errors.New("Error finding User Struct in Datastore")
	}
	// decrypt the user struct
	decUserStruct := userlib.SymDec(recipientEncPassword, recipientEncUser)
	var recipientUser User
	err = json.Unmarshal(decUserStruct, &recipientUser)
	if err != nil {
		return err
	}
	// grab the recipient's certificate struct
	recipientEncCertStruct, exists := userlib.DatastoreGet(certificateUUID)
	if !exists {
		return errors.New("Error finding Certificate Struct in Datastore")
	}

	// decrypt the certificate struct using private decKey
	decCertStruct, err := userdata.certificateDecryption(parentname, userdata.Username, "", certificateUUID)
	if err != nil {
		return err
	}

	var recipientCertStruct Certificates
	err = json.Unmarshal(decCertStruct, &recipientCertStruct)
	if err != nil {
		return err
	}

	// check the signature inside the cefrtificate struct
	// grab verification key
	verifyKey, exists := userlib.KeystoreGet(parentname + " verifyKey")
	if !exists {
		return errors.New("Error finding Verify Key in Keystore")
	}
	// grab the signature
	signature, exists := userlib.DatastoreGet(recipientCertStruct.SignatureUUID)
	if !exists {
		return errors.New("Error finding Certificate Struct Signature in Datastore")
	}
	// verify signature using verifyKey
	verified := userlib.DSVerify(verifyKey, recipientEncCertStruct, signature)
	if verified != nil {
		return errors.New("Error verifying digital signature")
	}
	for username, certificateUUID := range recipientCertStruct.Recipients {
		// grab the user struct for each username
		var recipientUser User
		err = recipientUser.updateToken(username, recipientUser.Username, certificateUUID, newAccessToken)
		if err != nil {
			return err
		}
	}

	//setup new struct to replace old one in DataStore
	recipientCertStruct.AccessToken = newAccessToken

	// reencrypt certificate struct and put in Datastore (does a new signature need to be created?)
	newCertStruct, err := recipientUser.certificateReencryption(parentname, username, "", certificateUUID, recipientCertStruct)
	if err != nil {
		return err
	}
	userlib.DatastoreSet(certificateUUID, newCertStruct)

	return nil
}

// Go from the Name to the FileInfo Struct ???
func (userdata *User) nameToFileInfo(filename string) (fileInfo *FileInfo, certificate *Certificates, err error) {
	userdata, err = usernameToUserStruct(userdata.Username)
	if err != nil {
		return nil, nil, err
	}
	certificateUUID, exists := userdata.Certificates[filename]
	if exists {
		// decrypt the certificate struct using private decKey
		sender, exists := userdata.Invites[filename]
		if !exists {
			return nil, nil, errors.New("Error finding sender in Invites")
		}
		decCertStruct, err := userdata.certificateDecryption(sender, userdata.Username, filename, certificateUUID) // 6
		if err != nil {
			return nil, nil, err
		}

		var certificateStruct Certificates
		err = json.Unmarshal(decCertStruct, &certificateStruct)
		if err != nil {
			return nil, nil, err
		}
		// possibly need to check MAC
		// grabbing corresponding FileInfo from decrypted Certificate struct
		fileInfoUUID := certificateStruct.FileInfo
		encFileInfo, exists := userlib.DatastoreGet(fileInfoUUID)
		if !exists {
			return nil, nil, errors.New("Error finding FileInfo Struct in Datastore")
		}

		// use Access Token to decrypt fileinfo struct
		accessToken := certificateStruct.AccessToken
		decFileInfo := userlib.SymDec(accessToken, encFileInfo)
		// need to error if the accessToken is wrong

		var fileInfo FileInfo
		err = json.Unmarshal(decFileInfo, &fileInfo)
		if err != nil {
			return nil, nil, err
		}

		currentFileMAC, err := FileMAC(fileInfo, fileInfo.Salt, accessToken)
		if err != nil {
			return nil, nil, err
		}

		hmacCheck := userlib.HMACEqual(fileInfo.MAC, currentFileMAC)
		if !(hmacCheck) {
			return nil, nil, errors.New("Failed verification test on FileInfo Struct")
		}
		return &fileInfo, &certificateStruct, nil
	} else {
		return nil, nil, nil
	}
}

// Go from the username to the UserStruct
func usernameToUserStruct(username string) (user *User, err error) {
	// grab the user struct from the Datastore
	userHash := userlib.Hash([]byte(username))[:16]
	userUUID, err := uuid.FromBytes(userHash)
	if err != nil {
		return nil, err
	}
	passHash, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return nil, errors.New("Error finding User Hashed Password in Datastore")
	}

	passHKDF, err := userlib.HashKDF(passHash, []byte("UUID"))
	if err != nil {
		return nil, err
	}

	passUUID, err := uuid.FromBytes(passHKDF[:16])
	if err != nil {
		return nil, err
	}

	// grab the encrypted user struct
	userEncUser, exists := userlib.DatastoreGet(passUUID)
	if !exists {
		return nil, errors.New("Error finding User Struct in Datastore")
	}

	// decrypt the user struct
	var userStruct User
	decUserStruct := userlib.SymDec(passHash, userEncUser)
	err = json.Unmarshal(decUserStruct, &userStruct)
	if err != nil {
		return nil, err
	}

	// HMAC verification
	currUserMAC := userStruct.MAC
	decUserMAC, err := UserMAC(userStruct, userStruct.Salt, passHash)
	if err != nil {
		return nil, err
	}

	hmacCheck := userlib.HMACEqual(currUserMAC, decUserMAC)
	if !(hmacCheck) {
		return nil, errors.New("Failed verification test on User Struct")
	}

	return &userStruct, nil
}

func (userdata *User) reencryptUser() (err error) {
	// get Hash(username)[:16]
	userHash := userlib.Hash([]byte(userdata.Username))[:16]
	// get the UUID over the computed Hash
	userUUID, err := uuid.FromBytes(userHash)
	if err != nil {
		return err
	}

	passHash, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return errors.New("Error getting the Enc Pass from User in Datastore")
	}

	passHKDF, err := userlib.HashKDF(passHash, []byte("UUID"))
	if err != nil {
		return err
	}

	passUUID, err := uuid.FromBytes(passHKDF[:16])
	if err != nil {
		return err
	}

	userdata.MAC, err = UserMAC(*userdata, userdata.Salt, passHash)
	if err != nil {
		return err
	}

	marshalledStruct, err := json.Marshal(userdata)
	if err != nil {
		return err
	}

	encUserStruct := userlib.SymEnc(passHash, userdata.Salt, marshalledStruct)

	userlib.DatastoreSet(passUUID, encUserStruct)
	return nil
}

// END OF HELPER FUNCTIONS

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username     string
	Password     string
	SignKey      userlib.DSSignKey    // sign privately and verify publicly
	DecryptKey   userlib.PKEDecKey    // encrypt publicly and decrypt privately
	Certificates map[string]uuid.UUID // filename : UUID of Certificate
	Invites      map[string]string    // file that was shared to user : username of person who shared it
	MAC          []byte               // MAC to verify struct integrity
	Salt         []byte               // IV for HMAC Verification and Enc/Dec
}

type AppendData struct {
	AppendData []byte
	MAC        []byte
	Salt       []byte
}

type AppendBlock struct {
	FileData   uuid.UUID // UUID of the Append Data
	NextAppend uuid.UUID
	MAC        []byte
	Salt       []byte
}

type FileInfo struct {
	StartAppend uuid.UUID
	EndAppend   uuid.UUID
	BlockKey    []byte // Key that encrypts blocks
	MAC         []byte
	Salt        []byte
}

type Certificates struct {
	ParentFilename string // filename given to same file by your  the
	FileInfo       uuid.UUID
	SignatureUUID  uuid.UUID            // UUID of the Certificate's signature
	Recipients     map[string]uuid.UUID // username : UUID of Certificate
	AccessToken    []byte               // might just keep AccessToken, no need for SEAToken ? // encrypts files
	Salt           []byte
	MAC            []byte
}

func InitUser(username string, password string) (userdataptr *User, err error) {
	// error if username is empty string
	if len(username) == 0 {
		return nil, errors.New("Username cannot be zero characters long")
	}
	// get Hash(username)[:16]
	userHash := userlib.Hash([]byte(username))
	// get the UUID over the computed Hash
	userUUID, err := uuid.FromBytes(userHash[:16])
	if err != nil {
		return nil, err // error getting the UUID from userHash
	}
	// error if username exists
	_, exists := userlib.DatastoreGet(userUUID)
	if exists {
		return nil, errors.New("This Username already exists")
	}
	encKey, decKey, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, err
	}
	signKey, verifyKey, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}

	var userdata User
	userdata.Username = username
	userdata.Password = password
	userdata.SignKey = signKey
	userdata.DecryptKey = decKey
	userdata.Certificates = make(map[string]uuid.UUID)
	userdata.Invites = make(map[string]string)
	userdata.Salt = userlib.RandomBytes(16)

	userlib.KeystoreSet(userdata.Username+" verifyKey", verifyKey)
	userlib.KeystoreSet(userdata.Username+" encKey", encKey)

	// Argon2Key(password, username)
	encryptedPass := userlib.Argon2Key([]byte(password), []byte(username), 16)
	// Encrypted Password --> Hash the Argon2Key
	hashedEncPass := userlib.Hash(encryptedPass)[:16]
	// setting the MAC of the new User Struct with the hashedEncPass
	userdata.MAC, err = UserMAC(userdata, userdata.Salt, hashedEncPass)
	if err != nil {
		return nil, err
	}
	// DatastoreSet(Hashed Username UUID, Encrypted Password)
	userlib.DatastoreSet(userUUID, hashedEncPass)

	// HKDF(Encrypted Password[:16], 'UUID')[:16]
	passHKDF, err := userlib.HashKDF(hashedEncPass, []byte("UUID"))
	if err != nil {
		return nil, err
	}

	passUUID, err := uuid.FromBytes(passHKDF[:16])
	if err != nil {
		return nil, err
	}
	// Encrypted User Struct --> SymEnc(Encrypted Password, UUID, Marshall(UserStruct))
	marshalledStruct, err := json.Marshal(userdata)
	if err != nil {
		return nil, err
	}

	encUserStruct := userlib.SymEnc(hashedEncPass, userdata.Salt, marshalledStruct)

	// DatastoreSet(Hashed and Encrypted Password, Ecrypted User Struct)
	userlib.DatastoreSet(passUUID, encUserStruct)

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata

	// use Hash(username)[:16] to grab the appropriate login entry UUID.
	userHash := userlib.Hash([]byte(username))[:16]
	userUUID, err := uuid.FromBytes(userHash)
	if err != nil {
		return nil, err // error getting the UUID from userHash
	}
	// check if userUUID exists, if yes,
	userHashedEncPass, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return nil, errors.New("Error finding Encrypted Password in Datastore")
	}
	// use the same process of HMAC(Argon2Key(username, password))
	encryptedPass := userlib.Argon2Key([]byte(password), []byte(username), 16)
	hashedEncPass := userlib.Hash(encryptedPass)[:16]
	// check that UUID - EncPass and hashedEncPass are the same, if not return err. // do we even need to do this ?
	hmacCheck := userlib.HMACEqual(userHashedEncPass, hashedEncPass)
	if !(hmacCheck) {
		return nil, errors.New("Failed verification test on User Struct")
	}

	// HashKDF(encrypted password, “UUID”)[:16]
	passHKDF, err := userlib.HashKDF(hashedEncPass, []byte("UUID"))
	if err != nil {
		return nil, err
	}
	passUUID, err := uuid.FromBytes(passHKDF[:16])
	if err != nil {
		return nil, err
	}
	// grab the appropriate UUID for the encrypted User struct
	userHashedEncStruct, exists := userlib.DatastoreGet(passUUID)
	// if exists decrypt encrypted user struct
	if !exists {
		return nil, errors.New("Error finding User Struct in Datastore")
	}
	decUserStruct := userlib.SymDec(userHashedEncPass, userHashedEncStruct)
	err = json.Unmarshal(decUserStruct, userdataptr)
	if err != nil {
		return nil, err
	}

	decUserStructMAC, err := UserMAC(*userdataptr, userdata.Salt, hashedEncPass)
	if err != nil {
		return nil, err
	}
	hmacCheck = userlib.HMACEqual(userdataptr.MAC, decUserStructMAC)
	if !(hmacCheck) {
		return nil, errors.New("Failed verification test on User Struct")
	}

	return userdataptr, nil
}

/*
Case 1: New file completely and need to set up file and first append block
Case 2: Wipe existing appends blocks except first and add new content
*/
func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	fileInfo, certificate, err := userdata.nameToFileInfo(filename)
	if err != nil {
		return err
	}
	if fileInfo != nil {
		// overwrite EXISTING file in Datastore
		accessToken := certificate.AccessToken
		fileInfoUUID := certificate.FileInfo
		firstAppendUUID := fileInfo.StartAppend
		blockKey := fileInfo.BlockKey

		// Check if anything has been tampered with
		err := traverseAppendBlock(firstAppendUUID, blockKey)
		if err != nil {
			return err
		}

		// create new AppendData to represent content of data in the append block
		var appendData AppendData
		appendDataUUID := uuid.New()
		appendData.AppendData = content
		appendData.Salt = userlib.RandomBytes(16)
		appendData.MAC, err = AppendDataMAC(appendData, blockKey)
		if err != nil {
			return err
		}

		// reencrypt and store new AppendBlock in Datastore
		marshalledNewAppendData, err := json.Marshal(appendData)
		if err != nil {
			return err
		}
		encNewAppendData := userlib.SymEnc(blockKey, appendData.Salt, marshalledNewAppendData)
		userlib.DatastoreSet(appendDataUUID, encNewAppendData)

		// create new AppendBlock to represent new file
		var appendBlock AppendBlock
		appendBlockUUID := uuid.New()
		appendBlock.FileData = appendBlockUUID
		appendBlock.NextAppend = uuid.Nil
		appendBlock.Salt = userlib.RandomBytes(16)
		appendBlock.MAC, err = AppendMAC(appendBlock, appendBlock.Salt, blockKey)
		if err != nil {
			return err
		}

		// reencrypt and store new AppendBlock in Datastore
		marshalledNewAppendBlock, err := json.Marshal(appendBlock)
		if err != nil {
			return err
		}
		encNewAppendBlock := userlib.SymEnc(blockKey, appendBlock.Salt, marshalledNewAppendBlock)
		userlib.DatastoreSet(appendBlockUUID, encNewAppendBlock)

		// "delete" previous AppendBlock and reset the "append chain"
		fileInfo.StartAppend = appendBlockUUID
		fileInfo.EndAppend = appendBlockUUID
		fileInfo.BlockKey = userlib.RandomBytes(16)
		fileInfo.MAC, err = FileMAC(*fileInfo, fileInfo.Salt, accessToken)
		if err != nil {
			return err
		}

		// reencrypt fileInfo and update it on Datastore
		marshalledFile, err := json.Marshal(fileInfo)
		if err != nil {
			return err
		}
		encFileInfo := userlib.SymEnc(accessToken, fileInfo.Salt, marshalledFile)
		userlib.DatastoreSet(fileInfoUUID, encFileInfo)
	} else {
		// overwrite EXISTING file in Datastore
		blockKey := userlib.RandomBytes(16) // create new blockKey
		var appendData AppendData
		appendDataUUID := uuid.New()
		appendData.AppendData = content
		appendData.Salt = userlib.RandomBytes(16)
		appendData.MAC, err = AppendDataMAC(appendData, blockKey)
		if err != nil {
			return err
		}
		// reencrypt and store new AppendBlock in Datastore
		marshalledNewAppendData, err := json.Marshal(appendData)
		if err != nil {
			return err
		}
		encNewAppendData := userlib.SymEnc(blockKey, appendData.Salt, marshalledNewAppendData)
		userlib.DatastoreSet(appendDataUUID, encNewAppendData)

		var appendBlock AppendBlock
		appendUUID := uuid.New()              // create a UUID for the corresponding AppendBlock
		appendBlock.FileData = appendDataUUID // set new File to have input content as filedata
		appendBlock.NextAppend = uuid.Nil     // has no nextappend since its first append in chain.
		appendBlock.Salt = userlib.RandomBytes(16)
		appendBlock.MAC, err = AppendMAC(appendBlock, appendBlock.Salt, blockKey)
		if err != nil {
			return err
		}

		// need to store in dataStore
		marshalledAppend, err := json.Marshal(appendBlock)
		if err != nil {
			return err
		}
		// encrypting the marshalled Struct
		encNewAppend := userlib.SymEnc(blockKey, appendBlock.Salt, marshalledAppend)
		if err != nil {
			return err
		}

		// store the computed encrypt-then-MAC, and store encrypted Struct in Datastore
		userlib.DatastoreSet(appendUUID, encNewAppend)

		var fileInfo FileInfo
		// Both Start and End must point to same AppendBlock
		accessToken := userlib.RandomBytes(16)
		FileUUID := uuid.New()
		fileInfo.StartAppend = appendUUID
		fileInfo.EndAppend = appendUUID
		fileInfo.BlockKey = blockKey
		fileInfo.Salt = userlib.RandomBytes(16)
		fileInfo.MAC, err = FileMAC(fileInfo, fileInfo.Salt, accessToken)
		if err != nil {
			return err
		}

		marshalledFile, err := json.Marshal(fileInfo)
		if err != nil {
			return err
		}
		// encrypting the marshalled Struct
		encNewFile := userlib.SymEnc(accessToken, fileInfo.Salt, marshalledFile)
		if err != nil {
			return err
		}

		// store the computed encrypt-then-MAC, and store encrypted Struct in Datastore
		userlib.DatastoreSet(FileUUID, encNewFile)

		// need to update the User struct with info on new file created:
		var certificate Certificates
		certificate.FileInfo = FileUUID
		certificate.Recipients = make(map[string]uuid.UUID)
		certificate.AccessToken = accessToken
		certificate.ParentFilename = filename
		certificate.Salt = userlib.RandomBytes(16)

		_, certificateUUID, err := userdata.certificateEncryption(userdata.Username, userdata.Username, filename, certificate)
		if err != nil {
			return err
		}

		userdata.Certificates[filename] = certificateUUID
		userdata.Invites[filename] = userdata.Username

		err = userdata.reencryptUser()
		if err != nil {
			return err
		}
	}

	return nil
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	// have to find endAppend (previous block in the AppendBlock chain)
	decFileInfo, decCertStruct, err := userdata.nameToFileInfo(filename)
	if err != nil {
		return err
	}
	if decFileInfo == nil {
		return errors.New("File Doesnt Exist in Users Namespace")
	}

	// set nextAppend in endAppend to newAppendBlock.
	endUUID := decFileInfo.EndAppend
	blockKey := decFileInfo.BlockKey
	accessToken := decCertStruct.AccessToken
	fileInfoUUID := decCertStruct.FileInfo
	encEndAppend, exists := userlib.DatastoreGet(endUUID)
	if !exists {
		return errors.New("Error finding End Append in Datastore")
	}

	var endAppend AppendBlock
	decEndAppend := userlib.SymDec(blockKey, encEndAppend)
	err = json.Unmarshal(decEndAppend, &endAppend) // how to decrypt ?
	if err != nil {
		return err
	}

	// creating AppendData
	var appendData AppendData
	appendDataUUID := uuid.New()
	appendData.AppendData = content
	appendData.Salt = userlib.RandomBytes(16)
	appendData.MAC, err = AppendDataMAC(appendData, blockKey)
	if err != nil {
		return err
	}

	// reencrypt and store new AppendBlock in Datastore
	marshalledNewAppendData, err := json.Marshal(appendData)
	if err != nil {
		return err
	}
	encNewAppendData := userlib.SymEnc(blockKey, appendData.Salt, marshalledNewAppendData)
	userlib.DatastoreSet(appendDataUUID, encNewAppendData)

	// creating AppendBlock with new data
	var appendBlock AppendBlock
	appendBlock.FileData = appendDataUUID
	appendBlock.NextAppend = uuid.Nil
	appendBlock.Salt = userlib.RandomBytes(16)
	appendBlock.MAC, err = AppendMAC(appendBlock, appendBlock.Salt, blockKey)
	if err != nil {
		return err
	}

	// need to encrypt and store new AppendBlock in Datastore
	currAppendUUID := uuid.New() // UUID for new AppendBlock
	marshalledAppend, err := json.Marshal(appendBlock)
	if err != nil {
		return err
	}
	// encrypt with block key
	encNewAppend := userlib.SymEnc(blockKey, appendBlock.Salt, marshalledAppend)
	userlib.DatastoreSet(currAppendUUID, encNewAppend)

	// update end append's next append to the curr append
	endAppend.NextAppend = currAppendUUID
	endAppend.MAC, err = AppendMAC(endAppend, endAppend.Salt, blockKey)
	if err != nil {
		return err
	}

	// marshal and reencrypt endAppendBlock
	marshalledEndAppend, err := json.Marshal(endAppend)
	if err != nil {
		return err
	}

	encEndAppend = userlib.SymEnc(blockKey, endAppend.Salt, marshalledEndAppend)

	// update the previous AppendBlock in datastore to have new nextAppend.
	userlib.DatastoreSet(endUUID, encEndAppend)

	// update and reencrypt FileInfo
	decFileInfo.EndAppend = currAppendUUID
	decFileInfo.MAC, err = FileMAC(*decFileInfo, decFileInfo.Salt, accessToken)
	if err != nil {
		return err
	}
	marshalledFileInfo, err := json.Marshal(decFileInfo)
	if err != nil {
		return err
	}

	encFileInfo := userlib.SymEnc(accessToken, decFileInfo.Salt, marshalledFileInfo)

	// update FileInfo in datastore to have new endAppend.
	userlib.DatastoreSet(fileInfoUUID, encFileInfo)

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	// find certificate and use keys to get access token to decrypt fileinfo struct
	decFileInfo, _, err := userdata.nameToFileInfo(filename)
	if err != nil {
		return nil, err
	}
	if decFileInfo == nil {
		return nil, errors.New("File Doesnt Exist in Users Namespace")
	}

	currUUID := decFileInfo.StartAppend
	blockKey := decFileInfo.BlockKey

	// check integrity of the AppendBlock
	err = traverseAppendBlock(currUUID, blockKey)
	if err != nil {
		return nil, err
	}

	// read the filedata from start append, until last append, using next append field.
	for currUUID != uuid.Nil {
		encCurrAppend, exists := userlib.DatastoreGet(currUUID)
		if !exists {
			return nil, errors.New("Error Finding Current Append Block Struct in Datastore") // value doesnt exist
		}

		var currAppend AppendBlock
		decCurrAppend := userlib.SymDec(blockKey, encCurrAppend)
		err = json.Unmarshal(decCurrAppend, &currAppend)
		if err != nil {
			return nil, err
		}

		currContentUUID := currAppend.FileData
		encAppendData, exists := userlib.DatastoreGet(currContentUUID)
		if !exists {
			return nil, errors.New("Error finding Append Data Struct in Datastore")
		}

		var appendData AppendData
		decCurrAppendData := userlib.SymDec(blockKey, encAppendData)
		err = json.Unmarshal(decCurrAppendData, &appendData)
		if err != nil {
			return nil, err
		}

		currContent := appendData.AppendData
		content = append(content, currContent...)
		currUUID = currAppend.NextAppend
	}

	return content, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	// check if recipientUsername exists
	recipientHash := userlib.Hash([]byte(recipientUsername))[:16]
	// get the UUID over the computed Hash of recipientUsername
	recipientUUID, err := uuid.FromBytes(recipientHash)
	if err != nil {
		return uuid.Nil, err
	}
	// checking that the recipient Exists ?
	_, exists := userlib.DatastoreGet(recipientUUID)
	if !exists {
		return uuid.Nil, errors.New("Recipient does not exist")
	}
	// Grabbing the User Struct
	ownerUser, err := usernameToUserStruct(userdata.Username)
	if err != nil {
		return uuid.Nil, err
	}
	// check if user has certificate of file then if user is owner of file
	certificateUUID, exists := ownerUser.Certificates[filename]
	if !exists {
		return uuid.Nil, errors.New("User does not have file")
	}
	// decrypt the certificate struct using private decKey
	decCert, err := userdata.certificateDecryption("", userdata.Username, filename, certificateUUID)
	if err != nil {
		return uuid.Nil, err
	}
	var ownerCert Certificates
	err = json.Unmarshal(decCert, &ownerCert)
	if err != nil {
		return uuid.Nil, err
	}

	var newCertificate Certificates
	newCertificate.FileInfo = ownerCert.FileInfo //gives UUID of the given file
	newCertificate.Recipients = make(map[string]uuid.UUID)
	newCertificate.AccessToken = ownerCert.AccessToken
	newCertificate.SignatureUUID = uuid.New() // careful of circular logic here
	newCertificate.ParentFilename = filename
	newCertificate.Salt = userlib.RandomBytes(16)

	_, encCertUUID, err := userdata.certificateEncryption(userdata.Username, recipientUsername, filename, newCertificate)
	if err != nil {
		return uuid.Nil, err
	}

	// return UUID to give the user access to the file
	return encCertUUID, nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	// check if senderUsername exists
	userHash := userlib.Hash([]byte(senderUsername))[:16]
	// get the UUID over the computed Hash of recipientUsername
	userUUID, err := uuid.FromBytes(userHash)
	if err != nil {
		return err // error getting the UUID from userHash
	}
	_, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return errors.New("Sender does not exist")
	}

	// check that a file with filename doesnt exist in users namespace
	_, exists = userdata.Certificates[filename]
	if exists {
		return errors.New("File with this name already exists")
	}

	decCertStructBytes, err := userdata.certificateDecryption(senderUsername, userdata.Username, filename, invitationPtr)
	if err != nil {
		return err
	}

	var certInfo Certificates
	err = json.Unmarshal(decCertStructBytes, &certInfo)
	if err != nil {
		return err
	}

	// update sender's recipient list and update signature
	ParentFilename := certInfo.ParentFilename
	senderInfo, err := usernameToUserStruct(senderUsername)
	if err != nil {
		return err
	}
	senderParent, exists := senderInfo.Invites[ParentFilename]
	if !exists {
		return errors.New("Sender's Parent not in Invites Map")
	}
	senderCertUUID, exists := senderInfo.Certificates[ParentFilename]
	if !exists {
		return errors.New("Sender's Parent not in Certificate Map")
	}
	decParentStructBytes, err := senderInfo.certificateDecryption(senderParent, senderUsername, "", senderCertUUID)
	if err != nil {
		return err
	}
	var senderCert Certificates
	err = json.Unmarshal(decParentStructBytes, &senderCert)
	if err != nil {
		return err
	}

	senderCert.Recipients[userdata.Username] = invitationPtr

	// update sender's certificate struct in DataStore
	_, err = senderInfo.certificateReencryption(senderParent, senderUsername, ParentFilename, senderCertUUID, senderCert)
	if err != nil {
		return err
	}

	// change name of the file to the given file and store in certificates
	userdata.Certificates[filename] = invitationPtr
	userdata.Invites[filename] = senderUsername

	err = userdata.reencryptUser()
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	// get the owners certificate struct for the given filename
	ownersCertUUID, exist := userdata.Certificates[filename]
	if !exist {
		return errors.New("File does not exist in users namespace")
	}

	// proper cert decryption
	ownersDecCertStruct, err := userdata.certificateDecryption("", userdata.Username, filename, ownersCertUUID)
	if err != nil {
		return err
	}

	var ownersCertStruct Certificates
	err = json.Unmarshal(ownersDecCertStruct, &ownersCertStruct)
	if err != nil {
		return err
	}

	// verify signature on CertStruct
	// access the recipient map within it
	recipientMap := ownersCertStruct.Recipients
	_, exists := recipientMap[recipientUsername]
	if !exists {
		// check that filename is shared with recipientUsername before deleting
		return errors.New("Person being revoked does not currently have access to the file.")
	}

	// recursively go through recipient's recipients
	// delete the certificate struct from the recipient
	// delete the datastore entry for the certificate struct
	// delete recipientUsername from recipient
	delete(recipientMap, recipientUsername)

	// create a new access token
	newAccessToken := userlib.RandomBytes(16)
	// update Access Token in users still with access
	for username, certificateStructUUID := range recipientMap {
		// grab the userstruct and call it on UpdateToken
		var recipientUser User
		err = recipientUser.updateToken(username, userdata.Username, certificateStructUUID, newAccessToken)
		if err != nil {
			return err
		}
	}

	// Change way FileInfo is encrypted --> using new access token
	fileUUID := ownersCertStruct.FileInfo
	encFileInfo, exists := userlib.DatastoreGet(fileUUID)
	if !exists {
		return errors.New("Error finding FileInfo Struct in Datastore")
	}
	decFile := userlib.SymDec(ownersCertStruct.AccessToken, encFileInfo)

	var newFile FileInfo
	err = json.Unmarshal(decFile, &newFile)
	if err != nil {
		return err
	}

	newFile.MAC, err = FileMAC(newFile, newFile.Salt, newAccessToken)
	if err != nil {
		return err
	}

	marshalledFile, err := json.Marshal(newFile)
	if err != nil {
		return err
	}
	newEncFile := userlib.SymEnc(newAccessToken, newFile.Salt, marshalledFile)

	userlib.DatastoreSet(fileUUID, newEncFile)

	// update the owners AccessToken
	ownersCertStruct.AccessToken = newAccessToken
	_, err = userdata.certificateReencryption(userdata.Username, userdata.Username, filename, ownersCertUUID, ownersCertStruct)
	if err != nil {
		return err
	}
	err = userdata.reencryptUser()
	if err != nil {
		return err
	}
	return nil
}

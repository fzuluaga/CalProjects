package client

///////////////////////////////////////////////////
//                                               //
// Everything in this file will NOT be graded!!! //
//                                               //
///////////////////////////////////////////////////

// In this unit tests file, you can write white-box unit tests on your implementation.
// These are different from the black-box integration tests in client_test.go,
// because in this unit tests file, you can use details specific to your implementation.

// For example, in this unit tests file, you can access struct fields and helper methods
// that you defined, but in the integration tests (client_test.go), you can only access
// the 8 functions (StoreFile, LoadFile, etc.) that are common to all implementations.

// In this unit tests file, you can write InitUser where you would write client.InitUser in the
// integration tests (client_test.go). In other words, the "client." in front is no longer needed.

import (
	"testing"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	_ "encoding/hex"
	"encoding/json"

	_ "errors"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"

	_ "strconv"

	_ "strings"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Unit Tests")
}

const defaultPassword = "LeButlerJordanJR"
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "

var _ = Describe("Client Unit Tests", func() {

	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	// charlesFile := "charlesFile.txt"

	BeforeEach(func() {
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Unit Tests", func() {
		Specify("Basic Test: Check that the Username field is set for a new user", func() {
			userlib.DebugMsg("Initializing user Alice.")
			// Note: In the integration tests (client_test.go) this would need to
			// be client.InitUser, but here (client_unittests.go) you can write InitUser.
			alice, err := InitUser("alice", "password")
			Expect(err).To(BeNil())

			// Note: You can access the Username field of the User struct here.
			// But in the integration tests (client_test.go), you cannot access
			// struct fields because not all implementations will have a username field.
			Expect(alice.Username).To(Equal("alice"))
		})
	})

	Describe("Malicious Activity Tests - User and File Functions", func() {
		Specify("Malicious Init User P1", func() {
			// init real user
			_, _ = InitUser("alice", defaultPassword)
			// get alice UUID
			aliceHash := userlib.Hash([]byte("alice"))
			aliceUUID, _ := uuid.FromBytes(aliceHash[:16])

			userlib.DebugMsg("Maliciously Changing Alice's Hashed Encrypted Password")
			// store garbage at alices UUID.
			userlib.DatastoreSet(aliceUUID, []byte("Lebron + Jimmy Butler >>> Michael Jordan."))
			_, err := GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Init User P2", func() {
			// init real user
			_, _ = InitUser("alice", defaultPassword)
			// get alice Pass UUID
			aliceHash := userlib.Hash([]byte("alice"))
			aliceUUID, _ := uuid.FromBytes(aliceHash[:16])
			aliceHashedEncPass, _ := userlib.DatastoreGet(aliceUUID)
			passHKDF, _ := userlib.HashKDF(aliceHashedEncPass, []byte("UUID"))
			passUUID, _ := uuid.FromBytes(passHKDF[:16])

			userlib.DebugMsg("Maliciously Changing Alice's User Struct.")
			// store garbage at alices User Struct.
			userlib.DatastoreSet(passUUID, []byte("Lebron + Jimmy Butler >>> Michael Jordan."))
			_, err := GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Store File", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			_, aliceCertStruct, _ := alice.nameToFileInfo(aliceFile)
			aliceFileUUID := aliceCertStruct.FileInfo

			userlib.DebugMsg("Maliciously Changing Alice's FileInfo Struct - Trying to Store File")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceFileUUID, []byte("very bad things were done here..."))
			// try to store again in the corrupted file.
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Load File", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			_, aliceCertStruct, _ := alice.nameToFileInfo(aliceFile)
			aliceFileUUID := aliceCertStruct.FileInfo

			userlib.DebugMsg("Maliciously Changing Alice's FileInfo Struct - Trying to Load File")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceFileUUID, []byte("very bad things were done here..."))
			// try to load the corrupted file.
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Append to File", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			_, aliceCertStruct, _ := alice.nameToFileInfo(aliceFile)
			aliceFileUUID := aliceCertStruct.FileInfo

			userlib.DebugMsg("Maliciously Changing Alice's FileInfo Struct - Trying to Append To File")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceFileUUID, []byte("very bad things were done here..."))
			// try to append to the corrupted file.
			err := alice.AppendToFile(aliceFile, []byte("I hope this file is not corrupted :P"))
			// try to load the corrupted file.
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Maliciously Changing Everything ?", func() {
		Specify("Malicious EncPass Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceHash := userlib.Hash([]byte("alice"))
			aliceUUID, _ := uuid.FromBytes(aliceHash[:16])

			userlib.DebugMsg("Maliciously Changing a Hashed and Encrypted Password.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious User Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceHash := userlib.Hash([]byte("alice"))
			aliceUUID, _ := uuid.FromBytes(aliceHash[:16])
			aliceHashedEncPass, _ := userlib.DatastoreGet(aliceUUID)
			passHKDF, _ := userlib.HashKDF(aliceHashedEncPass, []byte("UUID"))
			passUUID, _ := uuid.FromBytes(passHKDF[:16])

			userlib.DebugMsg("Maliciously Changing a User Struct.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(passUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Certificate Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceCertUUID := alice.Certificates[aliceFile]

			userlib.DebugMsg("Maliciously Changing a Certificate Struct.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceCertUUID, []byte("very bad things were done here..."))
			// try to store again in the corrupted file.
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious File Info Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			_, aliceCertStruct, _ := alice.nameToFileInfo(aliceFile)
			aliceFileUUID := aliceCertStruct.FileInfo

			userlib.DebugMsg("Maliciously Changing an FileInfo Struct.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceFileUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Append Block Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceFileInfoStruct, _, _ := alice.nameToFileInfo(aliceFile)
			aliceAppendUUID := aliceFileInfoStruct.StartAppend

			userlib.DebugMsg("Maliciously Changing an AppendBlock Struct.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceAppendUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Append Data Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceFileInfoStruct, _, _ := alice.nameToFileInfo(aliceFile)
			aliceAppendUUID := aliceFileInfoStruct.StartAppend
			encAppendBlock, _ := userlib.DatastoreGet(aliceAppendUUID)

			var appendBlock AppendBlock
			_ = json.Unmarshal(userlib.SymDec(aliceFileInfoStruct.BlockKey, encAppendBlock), &appendBlock)
			aliceAppendDataUUID := appendBlock.FileData

			userlib.DebugMsg("Maliciously Changing an AppendData Struct.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceAppendDataUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Certificate SymKey Tampering", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			aliceCertUUID := alice.Certificates[aliceFile]
			aliceKeyUUID, _ := getCertStructKeyUUID("alice", "alice", aliceCertUUID)

			userlib.DebugMsg("Maliciously Changing a Certificate Symmetrict Key.")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceKeyUUID, []byte("very bad things were done here..."))
			// try to store again after things have been corrupted
			err := alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Malicious Activity Tests - Invitation Functions", func() {
		Specify("File Maliciously Changed - CreateInvitation", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			// store real file
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			// get alice file UUID
			_, aliceCertStruct, _ := alice.nameToFileInfo(aliceFile)
			aliceFileUUID := aliceCertStruct.FileInfo

			userlib.DebugMsg("Maliciously Changing File Info - Trying To Create Invite")
			// store garbage at aliceFile UUID.
			userlib.DatastoreSet(aliceFileUUID, []byte("very bad things were done here..."))
			// try to store again in the corrupted file.
			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		Specify("Revoked User Trying to Gain Access", func() {
			// init real user
			alice, _ := InitUser("alice", defaultPassword)
			bob, _ := InitUser("bob", defaultPassword)
			// store real file
			alice.StoreFile(aliceFile, []byte(contentOne))
			// get invitationPtr for bob
			invitation, _ := alice.CreateInvitation(aliceFile, "bob")
			// bob accepts invitation and renames aliceFile to bobFile
			bob.AcceptInvitation("alice", invitation, bobFile)
			// bob loads old file and updates file
			bob.LoadFile(bobFile)
			bob.StoreFile(bobFile, []byte("Lakers in 4"))
			// bob is revoked Access from aliceFile
			alice.RevokeAccess("bob", aliceFile)

			userlib.DebugMsg("Revoked User Tries to Load File")
			// bob tries to load file that he no longer has access to.
			_, err := bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())
		})

		// Specify("Revoked User and Child Trying to Gain Access", func() {
		// 	// init real user
		// 	alice, _ := InitUser("alice", defaultPassword)
		// 	bob, _ := InitUser("bob", defaultPassword)
		// 	charles, _ := InitUser("charles", defaultPassword)
		// 	// store real file
		// 	alice.StoreFile(aliceFile, []byte(contentOne))
		// 	// get invitationPtr for bob
		// 	invitationBob, _ := alice.CreateInvitation(aliceFile, "bob")
		// 	// bob accepts invitation
		// 	bob.AcceptInvitation("alice", invitationBob, bobFile)
		// 	// bob operations on file
		// 	bob.LoadFile(bobFile)
		// 	bob.StoreFile(bobFile, []byte("Lakers in 4"))
		// 	// get invitationPtr for charles
		// 	invitationCharles, _ := bob.CreateInvitation(bobFile, "charles")
		// 	// charles accepts invitation
		// 	err := charles.AcceptInvitation("bob", invitationCharles, charlesFile)

		// 	userlib.DebugMsg("Charles My Only Opp") // ...
		// 	// charles loads
		// 	_, err = charles.LoadFile(charlesFile) // !!!
		// 	Expect(err).To(BeNil())

		// 	// charles stores the file
		// 	err = charles.StoreFile(charlesFile, []byte("Lakers in 4"))
		// 	Expect(err).To(BeNil())

		// 	// bob banned on file
		// 	alice.RevokeAccess("bob", aliceFile)
		// 	// bob tries to load file
		// 	userlib.DebugMsg("Child of Revoked User Tries to Load File")
		// 	_, err = charles.LoadFile(charlesFile)
		// 	Expect(err).ToNot(BeNil())
		// })
	})
})

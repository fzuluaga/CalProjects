package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const wrongPassword = "thisiswrong"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const maliciousContent = "HAHA Im so evil and malicious!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s", aliceFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite under name %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Extra Tests", func() {
		Specify("Edge Case Tests", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			// Init User Errors
			userlib.DebugMsg("Testing Init User Basic Error Cases")
			// user with the same name already exists
			_, err = client.InitUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())
			// username is empty string
			_, err = client.InitUser(emptyString, defaultPassword)
			Expect(err).ToNot(BeNil())

			// Get User Errors
			userlib.DebugMsg("Testing Get User Basic Error Cases")
			// there is no initialized user with that username
			_, err = client.GetUser("alex", defaultPassword)
			Expect(err).ToNot(BeNil())
			// wrong user credentials
			_, err = client.GetUser("alice", wrongPassword)
			Expect(err).ToNot(BeNil())

			// Store aliceFile:
			alice.StoreFile(aliceFile, []byte(contentOne))

			// Load File Errors
			userlib.DebugMsg("Testing Load File Basic Error Cases")
			// given file doesnt exist in personal namespace
			_, err = alice.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			// Append To File Errors
			userlib.DebugMsg("Testing Append To File Basic Error Cases")
			// file doesnt exist in personal namespace
			err = alice.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			// Create Invitiation Errors
			userlib.DebugMsg("Testing Create Invitation Basic Error Cases")
			// file doesnt exist in personal namespace of the caller
			_, err = alice.CreateInvitation(bobFile, "bob")
			Expect(err).ToNot(BeNil())
			// recipient with that username doesnt exist
			_, err = alice.CreateInvitation(aliceFile, "legoat")
			Expect(err).ToNot(BeNil())

			// Accept Invitation Errors
			userlib.DebugMsg("Testing Accept Invitation Basic Error Cases")
			bobInv, _ := alice.CreateInvitation(aliceFile, "bob")
			_ = bob.AcceptInvitation("alice", bobInv, aliceFile)
			// user already has access to this file
			bobInv, err = alice.CreateInvitation(aliceFile, "bob")
			err = bob.AcceptInvitation("alice", bobInv, aliceFile)
			Expect(err).ToNot(BeNil())
			// accepting from wrong username
			bobInv, err = alice.CreateInvitation(aliceFile, "bob")
			err = bob.AcceptInvitation("ey3", bobInv, aliceFile)
			Expect(err).ToNot(BeNil())
			// accepting invitation that wasnt sent to me
			err = charles.AcceptInvitation("alice", bobInv, aliceFile)
			Expect(err).ToNot(BeNil())
			// accepting invitation with inivitationPtr that doesnt exist
			err = bob.AcceptInvitation("alice", uuid.New(), aliceFile)
			Expect(err).ToNot(BeNil())
			// revoke aliceFile access from bob
			_ = alice.RevokeAccess(aliceFile, "bob")

			// Revoke Invitiation Errors
			userlib.DebugMsg("Testing Revoke Invitation Basic Error Cases")
			// filename doesnt exist in personal namespace
			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).ToNot(BeNil())
			// person being revoked doesnt have access to the file.
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
			// person being revoked doesnt exist
			err = alice.RevokeAccess(aliceFile, "jimmy")
			Expect(err).ToNot(BeNil())

			// Invitiation is no longer valid due to Revokation
			userlib.DebugMsg("Testing Invitation after Revokation Error Cases")
			// alice shares to bob
			bobInv, _ = alice.CreateInvitation(aliceFile, "bob")
			// bob accepts the invitaiton
			bob.AcceptInvitation("alice", bobInv, bobFile)
			// bob shares to charles
			charlesInv, _ := bob.CreateInvitation(bobFile, "charles")
			// alice revokes bob --> charles
			alice.RevokeAccess(aliceFile, "bob")
			// errors if charles tries to accept
			err = charles.AcceptInvitation("bob", charlesInv, charlesFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Custom Test: Testing Datastore Entry Deletion", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			userlib.DatastoreClear()

			userlib.DebugMsg("Trying to access users Alice, should error.")
			alice, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil())

		})

		Specify("Custom Test: Testing Keystore Entry Deletion", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
			userlib.KeystoreClear()

			userlib.DebugMsg("Trying to access key, should error while looking for key.")
			err = alice.StoreFile(aliceFile, []byte("content"))
			Expect(err).ToNot(BeNil())
		})

		Specify("Bandwidth Tests", func() {
			measureBandwidth := func(probe func()) (bandwidth int) {
				before := userlib.DatastoreGetBandwidth()
				probe()
				after := userlib.DatastoreGetBandwidth()
				return after - before
			}

			alice, err = client.InitUser("alice", defaultPassword)
			size10 := make([]byte, 10)
			err = alice.StoreFile("alice10.txt", size10)
			bw1 := measureBandwidth(func() {
				alice.AppendToFile("alice10.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size10))

			size50 := make([]byte, 50)
			err = alice.StoreFile("alice50.txt", size50)
			bw2 := measureBandwidth(func() {
				alice.AppendToFile("alice50.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size50))

			size200 := make([]byte, 200)
			err = alice.StoreFile("alice200.txt", size200)
			bw3 := measureBandwidth(func() {
				alice.AppendToFile("alice200.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size200))

			size1000 := make([]byte, 1000)
			err = alice.StoreFile("alice1000.txt", size1000)
			bw4 := measureBandwidth(func() {
				alice.AppendToFile("alice1000.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size1000))

			size2000 := make([]byte, 2000)
			err = alice.StoreFile("alice2000.txt", size2000)
			bw5 := measureBandwidth(func() {
				alice.AppendToFile("alice2000.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size2000))

			size10000 := make([]byte, 10000)
			err = alice.StoreFile("alice10000.txt", size10000)
			bw6 := measureBandwidth(func() {
				alice.AppendToFile("alice10000.txt", []byte("append"))
			})
			userlib.DebugMsg("Appending to a File of Size %d.", len(size10000))

			maxDiff := bw2 - bw1
			diff2 := bw4 - bw3
			diff3 := bw6 - bw5

			if diff2 > maxDiff {
				maxDiff = diff2
			}
			if diff3 > maxDiff {
				maxDiff = diff3
			}
			userlib.DebugMsg("Max Difference in Bandwidth between Appends is %d.", maxDiff)
			isGreaterThan500 := maxDiff > 500
			Expect(isGreaterThan500).To(BeFalse())
		})
	})

	Describe("Malicious Activity", func() {
		Specify("Malicious Activity Check - Get User", func() {
			_, _ = client.InitUser("alice", defaultPassword)
			userlib.DebugMsg("Maliciously Changing Data - Get User.")
			for uuid := range userlib.DatastoreGetMap() {
				userlib.DatastoreSet(uuid, []byte(maliciousContent))
				_, err := client.GetUser("alice", defaultPassword)
				Expect(err).ToNot(BeNil())
			}
		})

		Specify("Malicious Activty Check - Store File", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			beforeStore := len(userlib.DatastoreGetMap())
			_ = alice.StoreFile(aliceFile, []byte(contentOne))

			i := 0
			userlib.DebugMsg("Maliciously Changing Data - Store File.")
			for uuid := range userlib.DatastoreGetMap() {
				if i >= beforeStore {
					userlib.DatastoreSet(uuid, []byte(maliciousContent))
					err := alice.StoreFile(aliceFile, []byte(contentOne))
					Expect(err).ToNot(BeNil())
				}
				i++
			}
		})

		Specify("Malicious Activty Check - Load File", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			beforeStore := len(userlib.DatastoreGetMap())
			_ = alice.StoreFile(aliceFile, []byte(contentOne))

			i := 0
			userlib.DebugMsg("Maliciously Changing Data - Load File.")
			for uuid := range userlib.DatastoreGetMap() {
				if i >= beforeStore {
					userlib.DatastoreSet(uuid, []byte(maliciousContent))
				}
				i++
			}
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Activty Check - Append To File", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			beforeStore := len(userlib.DatastoreGetMap())
			_ = alice.StoreFile(aliceFile, []byte(contentOne))

			i := 0
			userlib.DebugMsg("Maliciously Changing Data - Append To File.")
			for uuid := range userlib.DatastoreGetMap() {
				if i >= beforeStore {
					userlib.DatastoreSet(uuid, []byte(maliciousContent))
				}
				i++
			}
			err := alice.AppendToFile(aliceFile, []byte("I hope this file is not corrupted :P"))
			Expect(err).ToNot(BeNil())
		})

		Specify("Malicious Activty Check - Create Invitation", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			bob, _ = client.InitUser("bob", defaultPassword)
			beforeStore := len(userlib.DatastoreGetMap())
			_ = alice.StoreFile(aliceFile, []byte(contentOne))

			i := 0
			userlib.DebugMsg("Maliciously Changing Data - Create Invitation.")
			for uuid := range userlib.DatastoreGetMap() {
				if i >= beforeStore {
					userlib.DatastoreSet(uuid, []byte(maliciousContent))
				}
				i++
			}
			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		// Accept Invitation doesnt check for malicious activity, we do error checking in edge case tests.

		Specify("Malicious Activty Check - Revoke Access", func() {
			alice, _ = client.InitUser("alice", defaultPassword)
			bob, _ = client.InitUser("bob", defaultPassword)
			_ = alice.StoreFile(aliceFile, []byte(contentOne))
			bobInv, _ := alice.CreateInvitation(aliceFile, "bob")
			beforeAccept := len(userlib.DatastoreGetMap())
			_ = bob.AcceptInvitation("alice", bobInv, aliceFile)

			i := 0
			userlib.DebugMsg("Maliciously Changing Data - Revoke Access.")
			for uuid := range userlib.DatastoreGetMap() {
				if i >= beforeAccept {
					userlib.DatastoreSet(uuid, []byte(maliciousContent))
					err := alice.RevokeAccess(aliceFile, "bob")
					Expect(err).ToNot(BeNil())
				}
				i++
			}
		})
	})
})

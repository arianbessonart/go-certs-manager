package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"
)

// You'll need a user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func main() {

	storageClient, err := NewStorageClient("directory-certificates-arian")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Starting")
	time.Sleep(120 * time.Second)
	fmt.Println("After 120 seconds")

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	fmt.Println(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: "arian@austinsoftware.com",
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Docker in a VM.
	config.CADirURL = lego.LEDirectoryStaging // "http://192.168.99.100:4000/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an http port of 5002 and an tls port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	if err != nil {
		log.Fatal(err)
	}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"gomesh.ml"},
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)

	err = storageClient.AddFileToBucket(certificates.Domain, certificates.Certificate)
	log.Fatalf("Trying upload certificate, %v", err)
}

// func main() {
// 	http.HandleFunc("/", helloServer)
// 	http.ListenAndServe(":5002", nil)
// }

// func helloServer(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
// }

// curl -X GET http://gomesh.ml

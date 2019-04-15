package cloudfuncs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/iterator"
)

type Machine struct {
	Name  string `json:"name,omitempty" firestore:"name,omitempty"`
	Power int    `json:"power,omitempty" firestore:"power,omitempty"`
}

// MachineHttp is an HTTP Cloud Function that CRUDs a Machine document in the firestore
func MachineHttp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := fmt.Fprintf(w, "%s method not currently supported", r.Method)
		if err != nil {
			http.Error(w, "Error printing error", http.StatusInternalServerError)
			log.Printf("fmt.Fprintf( method not currently supported): %v", err)
			return
		}
		return
	}
	client = getClient(w)
	if client == nil {
		http.Error(w, "Error connecting to firestore", http.StatusInternalServerError)
		log.Printf("Error connecting to firestore")
		return
	}
	iter := client.Collection("machines").Documents(r.Context())
	defer iter.Stop()
	for {
		var machine Machine
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error locating machine data", http.StatusInternalServerError)
			log.Printf("firestore.Collection(machines).Documents: %v", err)
			return
		}
		err2 := doc.DataTo(&machine)
		if err2 != nil {
			http.Error(w, "Error loading machine data", http.StatusInternalServerError)
			log.Printf("firestore.Document.DataTo: %v", err2)
			return
		}
		jsonData, _ := json.MarshalIndent(machine, "", "    ")
		_, err2 = fmt.Fprintln(w, string(jsonData))
		if err2 != nil {
			http.Error(w, "Error printing machine data", http.StatusInternalServerError)
			log.Printf("fmt.Fprintln(w, string(jsonData))): %v", err)
			return
		}
	}

}

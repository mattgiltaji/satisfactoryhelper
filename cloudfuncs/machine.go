package cloudfuncs

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/iterator"
)

// MachineHttp is an HTTP Cloud Function that CRUDs a Machine document in the firestore
func MachineHttp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err := fmt.Fprintf(w, "%s method not currently supported", r.Method)
		if err != nil {
			http.Error(w, "Error printing error", http.StatusInternalServerError)
			log.Printf(" fmt.Fprintf( method not currently supported): %v", err)
			return
		}
		return
	}
	client = getClient(w)
	iter := client.Collection("machines").Documents(r.Context())
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error locating machine data", http.StatusInternalServerError)
			log.Printf("firestore.Collection(machines).Documents: %v", err)
			return
		}
		_, err2 := fmt.Fprintln(w, doc.Data())
		if err2 != nil {
			http.Error(w, "Error printing machine data", http.StatusInternalServerError)
			log.Printf(" fmt.Fprintln(doc.Data()): %v", err)
			return
		}
	}

}

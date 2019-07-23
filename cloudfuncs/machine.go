package cloudfuncs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
	"google.golang.org/api/iterator"
)

const pathToMachinesInFirestore = "machines"

type Machine struct {
	Name  string `json:"name,omitempty" firestore:"name,omitempty"`
	Power int    `json:"power,omitempty" firestore:"power,omitempty"`
}

// MachineHttp is an HTTP Cloud Function that CRUDs a Machine document in the firestore
func MachineHttp(w http.ResponseWriter, r *http.Request) {
	//todo: parse url to figure out what this is being called with so we know what we should do
	// equivalent of api/machines/ -> get all, post all, delete all, batch update
	// maybe get all with query methods filters?
	// api/machines/$name -> get/add/update/delete specific machine

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

	machines, err := getAllMachines(r.Context(), client)
	if err != nil {
		http.Error(w, "Error locating machine data", http.StatusInternalServerError)
		log.Printf("getAllMachines(): %v", err)
		return
	}
	jsonData, _ := json.MarshalIndent(machines, "", "    ")
	_, err = fmt.Fprintln(w, string(jsonData))
	if err != nil {
		http.Error(w, "Error printing machine data", http.StatusInternalServerError)
		log.Printf("fmt.Fprintln(w, string(jsonData))): %v", err)
		return
	}

}

//getAllMachines returns all machines in the firestore
func getAllMachines(ctx context.Context, client *firestore.Client) (machines []Machine, err error) {
	iter := client.Collection(pathToMachinesInFirestore).Documents(ctx)
	defer iter.Stop()
	for {
		var machine Machine
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return machines, errors.Annotate(err, "iterator broke while loading all machine data from firestore.")
		}
		err = doc.DataTo(&machine)
		if err != nil {
			return machines, errors.Annotatef(err, "unable to parse specific machine's data: %v", doc.Data())
		}
		machines = append(machines, machine)
	}
	return
}

//getMachineByName returns a single machine with a matching name or an IsNotFound error if no such machine is in firestore
//we rely on the addMachine logic to ensure that there are no Machines with duplicate names
func getMachineByName(ctx context.Context, client *firestore.Client, name string) (machine Machine, err error) {
	q := client.Collection(pathToMachinesInFirestore).Where("name", "==", name)
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return machine, errors.Annotate(err, "iterator broke while loading specific machine data from firestore.")
		}
		err = doc.DataTo(&machine)
		if err != nil {
			return machine, errors.Annotatef(err, "unable to parse specific machine's data: %v", doc.Data())
		}
		//assume no multiple results from firestore so no need to loop again
		return machine, err
	}
	err = errors.NotFoundf("No machine with name %v", name)
	return
}

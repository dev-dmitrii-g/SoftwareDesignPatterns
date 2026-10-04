package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"hospitaladmission/internal/equipmentreserver"
	"hospitaladmission/internal/facades"
	"hospitaladmission/internal/insuranceverifier"
	"hospitaladmission/internal/pharmacydispatcher"
	"hospitaladmission/internal/store"
)

func run(args []string) error {
	flags := flag.NewFlagSet("hospital", flag.ContinueOnError)
	path := flags.String("db", "hospital.db", "SQLite database path")
	addr := flags.String("addr", "127.0.0.1:8080", "web server address")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}
	db, err := store.Open(*path)
	if err != nil {
		return err
	}
	defer db.Close()
	i := insuranceverifier.Verifier{Rule: insuranceverifier.Rule{}, Repository: insuranceverifier.Repository{DB: db}}
	e := equipmentreserver.Reserver{Rule: equipmentreserver.Rule{}, Repository: equipmentreserver.Repository{DB: db}}
	p := pharmacydispatcher.Dispatcher{Rule: pharmacydispatcher.Rule{}, Repository: pharmacydispatcher.Repository{DB: db}}
	f := facades.IntakeFacade{Gateway: facades.Gateway{DB: db}, Insurance: i, Equipment: e, Pharmacy: p}
	fmt.Printf("http://%s\n", *addr)
	return http.ListenAndServe(*addr, newServer(db, f, facades.AdminFacade{Intake: f}))
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

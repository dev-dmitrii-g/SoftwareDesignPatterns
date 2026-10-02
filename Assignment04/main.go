package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"hospitaladmission/internal/admission"
	"hospitaladmission/internal/equipment"
	"hospitaladmission/internal/insurance"
	"hospitaladmission/internal/pharmacy"
	"hospitaladmission/internal/store"
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: go run . -db hospital.db COMMAND [arguments]
Commands:
  stock equipment|pharmacy NAME QUANTITY
  start ID IIN POLICY EQUIPMENT MEDICINE   # stop after insurance
  resume ID                              # continue after restart
  admit ID IIN POLICY EQUIPMENT MEDICINE  # full chain
  inspect ID
  compensate ID

IIN is a 12-digit demo identifier. Set stock before admission.
Example: go run . stock equipment wheelchair 1`)
}
func run(args []string) error {
	flags := flag.NewFlagSet("hospital", flag.ContinueOnError)
	path := flags.String("db", "hospital.db", "SQLite database path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	a := flags.Args()
	if len(a) == 0 {
		usage()
		return nil
	}
	db, err := store.Open(*path)
	if err != nil {
		return err
	}
	defer db.Close()
	i := insurance.Verifier{Rule: insurance.Rule{}, Repository: insurance.Repository{DB: db}}
	e := equipment.Reserver{Rule: equipment.Rule{}, Repository: equipment.Repository{DB: db}}
	p := pharmacy.Dispatcher{Rule: pharmacy.Rule{}, Repository: pharmacy.Repository{DB: db}}
	f := admission.IntakeFacade{Gateway: admission.Gateway{DB: db}, Insurance: i, Equipment: e, Pharmacy: p}
	admin := admission.AdminFacade{Intake: f}
	switch a[0] {
	case "stock":
		if len(a) != 4 {
			usage()
			return fmt.Errorf("stock needs kind, name, quantity")
		}
		n, err := strconv.Atoi(a[3])
		if err != nil {
			return err
		}
		switch a[1] {
		case "equipment":
			err = e.Repository.SetStock(a[2], n)
		case "pharmacy":
			err = p.Repository.SetStock(a[2], n)
		default:
			return fmt.Errorf("unknown stock kind %q", a[1])
		}
		if err == nil {
			fmt.Printf("%s %s stock = %d\n", a[1], a[2], n)
		}
		return err
	case "start", "admit":
		if len(a) != 6 {
			usage()
			return fmt.Errorf("%s needs ID, IIN, policy, equipment, medicine", a[0])
		}
		request := admission.Admission{ID: a[1], IIN: a[2], Policy: a[3], Equipment: a[4], Medicine: a[5]}
		if a[0] == "start" {
			err = f.Start(request)
		} else {
			err = f.Admit(request)
		}
		if err != nil {
			_ = printView(admin, a[1])
			return err
		}
		return printView(admin, a[1])
	case "resume":
		if len(a) != 2 {
			usage()
			return fmt.Errorf("resume needs ID")
		}
		if err = f.Resume(a[1]); err != nil {
			_ = printView(admin, a[1])
			return err
		}
		return printView(admin, a[1])
	case "inspect":
		if len(a) != 2 {
			usage()
			return fmt.Errorf("inspect needs ID")
		}
		return printView(admin, a[1])
	case "compensate":
		if len(a) != 2 {
			usage()
			return fmt.Errorf("compensate needs ID")
		}
		if err = admin.Compensate(a[1]); err != nil {
			return err
		}
		return printView(admin, a[1])
	default:
		usage()
		return fmt.Errorf("unknown command %q", a[0])
	}
}
func printView(f admission.AdminFacade, id string) error {
	v, err := f.Inspect(id)
	if err != nil {
		return err
	}
	fmt.Printf("Admission %s: %s\nInsurance: %s | Equipment: %s | Pharmacy: %s\n", id, v.Admission.Status, v.Insurance, v.Equipment, v.Pharmacy)
	if v.Admission.Reason != "" {
		fmt.Printf("Reason: %s\n", v.Admission.Reason)
	}
	return nil
}
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

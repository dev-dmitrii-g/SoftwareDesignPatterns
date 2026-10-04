package facades_test

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"hospitaladmission/internal/equipmentreserver"
	"hospitaladmission/internal/facades"
	"hospitaladmission/internal/insuranceverifier"
	"hospitaladmission/internal/pharmacydispatcher"
	"hospitaladmission/internal/store"
)

type hospital struct {
	t      *testing.T
	path   string
	db     *sql.DB
	intake facades.IntakeFacade
	admin  facades.AdminFacade
}

func aHospital(t *testing.T) *hospital {
	t.Helper()
	h := &hospital{t: t, path: filepath.Join(t.TempDir(), "test.db")}
	h.reopen()
	return h
}
func (h *hospital) reopen() {
	h.t.Helper()
	if h.db != nil {
		if err := h.db.Close(); err != nil {
			h.t.Fatal(err)
		}
	}
	db, err := store.Open(h.path)
	if err != nil {
		h.t.Fatal(err)
	}
	h.db = db
	h.intake = facades.IntakeFacade{Gateway: facades.Gateway{DB: db},
		Insurance: insuranceverifier.Verifier{Rule: insuranceverifier.Rule{}, Repository: insuranceverifier.Repository{DB: db}},
		Equipment: equipmentreserver.Reserver{Rule: equipmentreserver.Rule{}, Repository: equipmentreserver.Repository{DB: db}},
		Pharmacy:  pharmacydispatcher.Dispatcher{Rule: pharmacydispatcher.Rule{}, Repository: pharmacydispatcher.Repository{DB: db}}}
	h.admin = facades.AdminFacade{Intake: h.intake}
}
func (h *hospital) stock(equipmentQty, medicineQty int) *hospital {
	h.t.Helper()
	if err := h.intake.Equipment.Repository.SetStock("wheelchair", equipmentQty); err != nil {
		h.t.Fatal(err)
	}
	if err := h.intake.Pharmacy.Repository.SetStock("paracetamol", medicineQty); err != nil {
		h.t.Fatal(err)
	}
	return h
}
func request(id string) facades.Admission {
	return facades.Admission{ID: id, IIN: "990101123456", Coverage: "OSMS", Equipment: "wheelchair", Medicine: "paracetamol"}
}

func (h *hospital) admits(id string) error { h.t.Helper(); return h.intake.Admit(request(id)) }
func (h *hospital) records(id, status, ins, equip, pharm string) {
	h.t.Helper()
	v, err := h.admin.Inspect(id)
	if err != nil {
		h.t.Fatal(err)
	}
	if v.Admission.Status != status || v.Insurance != ins || v.Equipment != equip || v.Pharmacy != pharm {
		h.t.Fatalf("persisted state = %s/%s/%s/%s, want %s/%s/%s/%s", v.Admission.Status, v.Insurance, v.Equipment, v.Pharmacy, status, ins, equip, pharm)
	}
}
func (h *hospital) quantity(kind string) int {
	h.t.Helper()
	var n int
	if err := h.db.QueryRow(`SELECT quantity FROM stock WHERE kind=?`, kind).Scan(&n); err != nil {
		h.t.Fatal(err)
	}
	return n
}

func TestHappyPathAndReload(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	if err := h.admits("ok"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	h.records("ok", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	if err := h.intake.Resume("ok"); err != nil {
		t.Fatal(err)
	}
	if h.quantity("equipment") != 0 || h.quantity("pharmacy") != 0 {
		t.Fatal("stock was not consumed exactly once")
	}
}
func TestStepwiseAdmissionAndRelease(t *testing.T) {
	h := aHospital(t).stock(2, 2)
	defer h.db.Close()
	if err := h.intake.Start(request("steps")); err != nil {
		t.Fatal(err)
	}
	h.records("steps", facades.PendingB, insuranceverifier.Verified, "NONE", "NONE")
	if err := h.intake.Dispatch("steps"); err == nil {
		t.Fatal("pharmacy dispatched before equipment")
	}
	if err := h.admin.Release("steps"); err == nil {
		t.Fatal("pending admission released")
	}
	if err := h.intake.Reserve("steps"); err != nil {
		t.Fatal(err)
	}
	h.records("steps", facades.PendingC, insuranceverifier.Verified, equipmentreserver.Reserved, "NONE")
	if h.quantity("equipment") != 1 || h.quantity("pharmacy") != 2 {
		t.Fatal("incorrect stock after reservation")
	}
	if err := h.intake.Reserve("steps"); err == nil {
		t.Fatal("equipment reserved twice")
	}
	if err := h.intake.Dispatch("steps"); err != nil {
		t.Fatal(err)
	}
	h.records("steps", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	if h.quantity("equipment") != 1 || h.quantity("pharmacy") != 1 {
		t.Fatal("incorrect stock after dispatch")
	}
	if err := h.admin.Release("steps"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	h.records("steps", facades.Compensated, insuranceverifier.Reversed, equipmentreserver.Released, pharmacydispatcher.Cancelled)
	if h.quantity("equipment") != 2 || h.quantity("pharmacy") != 2 {
		t.Fatal("release did not restore both stocks")
	}
	if err := h.admin.Release("steps"); err == nil {
		t.Fatal("admission released twice")
	}
}

func TestResumePartlyReleasedAdmission(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	if err := h.admits("partial-release"); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Gateway.Update("partial-release", facades.Compensating, "admission released"); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Pharmacy.Cancel("partial-release"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("partial-release"); err != nil {
		t.Fatal(err)
	}
	h.records("partial-release", facades.Compensated, insuranceverifier.Reversed, equipmentreserver.Released, pharmacydispatcher.Cancelled)
	if h.quantity("equipment") != 1 || h.quantity("pharmacy") != 1 {
		t.Fatal("partial release changed stock twice")
	}
}
func TestEquipmentFailureCompensatesInsurance(t *testing.T) {
	h := aHospital(t).stock(0, 1)
	defer h.db.Close()
	if err := h.admits("no-equipment"); err == nil || !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("got %v", err)
	}
	h.reopen()
	h.records("no-equipment", facades.Compensated, insuranceverifier.Reversed, "NONE", "NONE")
	if h.quantity("equipment") != 0 || h.quantity("pharmacy") != 1 {
		t.Fatal("stock changed on failed step")
	}
}
func TestPharmacyFailureReleasesEquipmentThenReversesInsurance(t *testing.T) {
	h := aHospital(t).stock(1, 0)
	defer h.db.Close()
	if err := h.admits("no-medicine"); err == nil || !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("got %v", err)
	}
	h.reopen()
	h.records("no-medicine", facades.Compensated, insuranceverifier.Reversed, equipmentreserver.Released, "NONE")
	if h.quantity("equipment") != 1 {
		t.Fatal("equipment was not returned")
	}
}
func TestRestartMidChain(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	if err := h.intake.Start(request("restart")); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	h.records("restart", facades.PendingB, insuranceverifier.Verified, "NONE", "NONE")
	if err := h.intake.Resume("restart"); err != nil {
		t.Fatal(err)
	}
	h.records("restart", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
}
func TestResumeReadsRowsWhenProgressMarkerIsBehind(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	if err := h.intake.Start(request("lag")); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Equipment.Reserve("lag", "wheelchair"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("lag"); err != nil {
		t.Fatal(err)
	}
	h.records("lag", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	if h.quantity("equipment") != 0 {
		t.Fatal("equipment reserved twice")
	}
}
func TestManualCompensationAndInvalidInput(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	bad := request("bad")
	bad.IIN = "abc"
	if err := h.intake.Admit(bad); err == nil {
		t.Fatal("invalid IIN accepted")
	}
	if _, err := h.intake.Gateway.Get("bad"); err != sql.ErrNoRows {
		t.Fatalf("invalid admission persisted: %v", err)
	}
	if err := h.intake.Start(request("manual")); err != nil {
		t.Fatal(err)
	}
	if err := h.admin.Compensate("manual"); err != nil {
		t.Fatal(err)
	}
	h.records("manual", facades.Compensated, insuranceverifier.Reversed, "NONE", "NONE")
	if err := h.admin.Compensate("manual"); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Resume("manual"); err == nil {
		t.Fatal("compensated admission resumed")
	}
}
func TestSubsystemsWorkAloneAndRejectIllegalTransitions(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	i, e, p := h.intake.Insurance, h.intake.Equipment, h.intake.Pharmacy
	if err := i.Verify("i", "990101123456", "OSMS"); err != nil {
		t.Fatal(err)
	}
	if err := i.Reverse("i"); err != nil {
		t.Fatal(err)
	}
	if err := i.Reverse("i"); err == nil {
		t.Fatal("insurance reversed twice")
	}
	if err := e.Reserve("e", "wheelchair"); err != nil {
		t.Fatal(err)
	}
	if err := e.Release("e"); err != nil {
		t.Fatal(err)
	}
	if err := e.Release("e"); err == nil {
		t.Fatal("equipment released twice")
	}
	if err := p.Dispatch("p", "paracetamol"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.Repository.Get("p"); err != nil {
		t.Fatal(err)
	}
	if err := i.Verify("invalid", "123", "OSMS"); err == nil {
		t.Fatal("invalid IIN accepted")
	}
	if err := e.Reserve("blank", ""); err == nil {
		t.Fatal("blank equipment accepted")
	}
	if err := p.Dispatch("blank", ""); err == nil {
		t.Fatal("blank medicine accepted")
	}
}

func TestValidationAndAdministrativeGuards(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	for _, bad := range []facades.Admission{
		{IIN: "990101123456", Coverage: "OSMS", Equipment: "wheelchair", Medicine: "paracetamol"},
		{ID: "bad-coverage", IIN: "990101123456", Equipment: "wheelchair", Medicine: "paracetamol"},
		{ID: "no-coverage", IIN: "990101123456", Coverage: insuranceverifier.NoCoverage, Equipment: "wheelchair", Medicine: "paracetamol"},
		{ID: "unknown-coverage", IIN: "990101123456", Coverage: "OTHER", Equipment: "wheelchair", Medicine: "paracetamol"},
		{ID: "bad-equipment", IIN: "990101123456", Coverage: "OSMS", Medicine: "paracetamol"},
		{ID: "bad-medicine", IIN: "990101123456", Coverage: "OSMS", Equipment: "wheelchair"},
		{ID: "bad-digit", IIN: "99010112345x", Coverage: "OSMS", Equipment: "wheelchair", Medicine: "paracetamol"},
	} {
		if err := h.intake.Start(bad); err == nil {
			t.Fatalf("accepted %+v", bad)
		}
	}
	if err := h.intake.Equipment.Repository.SetStock("wheelchair", -1); err == nil {
		t.Fatal("negative equipment stock accepted")
	}
	if err := h.intake.Pharmacy.Repository.SetStock("paracetamol", -1); err == nil {
		t.Fatal("negative medicine stock accepted")
	}
	if _, err := h.admin.Inspect("missing"); err != sql.ErrNoRows {
		t.Fatalf("inspect missing: %v", err)
	}
	if err := h.intake.Resume("missing"); err != sql.ErrNoRows {
		t.Fatalf("resume missing: %v", err)
	}
	if err := h.admin.Compensate("missing"); err != sql.ErrNoRows {
		t.Fatalf("compensate missing: %v", err)
	}
	if err := h.admits("ok"); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Start(request("ok")); err == nil {
		t.Fatal("duplicate admission accepted")
	}
	if err := h.admin.Compensate("ok"); err == nil {
		t.Fatal("completed admission compensated")
	}
	p := pharmacydispatcher.Dispatch{Status: pharmacydispatcher.Dispatched}
	if err := p.Transition(pharmacydispatcher.Cancelled); err != nil {
		t.Fatal(err)
	}
	if err := p.Transition(pharmacydispatcher.Dispatched); err == nil {
		t.Fatal("illegal pharmacy transition accepted")
	}
}

func TestAcceptedCoverageChoices(t *testing.T) {
	h := aHospital(t).stock(3, 3)
	defer h.db.Close()
	for _, coverage := range []string{insuranceverifier.GOBMP, insuranceverifier.OSMS, insuranceverifier.DMS} {
		request := request(coverage)
		request.Coverage = coverage
		if err := h.intake.Admit(request); err != nil {
			t.Fatalf("%s: %v", coverage, err)
		}
		h.records(coverage, facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	}
}

func TestResumeFromPersistedPendingAndCompensating(t *testing.T) {
	h := aHospital(t).stock(1, 1)
	defer h.db.Close()
	if err := h.intake.Gateway.Create(request("pending")); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("pending"); err != nil {
		t.Fatal(err)
	}
	h.records("pending", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	if err := h.intake.Start(request("undo")); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Gateway.Update("undo", facades.Compensating, "operator stopped"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("undo"); err != nil {
		t.Fatal(err)
	}
	h.records("undo", facades.Compensated, insuranceverifier.Reversed, "NONE", "NONE")
}

func TestResumeAfterCommittedStepBeforeProgressUpdate(t *testing.T) {
	h := aHospital(t).stock(2, 2)
	defer h.db.Close()
	if err := h.intake.Gateway.Create(request("after-a")); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Insurance.Verify("after-a", "990101123456", "OSMS"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("after-a"); err != nil {
		t.Fatal(err)
	}
	h.records("after-a", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
	if err := h.intake.Start(request("after-c")); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Equipment.Reserve("after-c", "wheelchair"); err != nil {
		t.Fatal(err)
	}
	if err := h.intake.Pharmacy.Dispatch("after-c", "paracetamol"); err != nil {
		t.Fatal(err)
	}
	h.reopen()
	if err := h.intake.Resume("after-c"); err != nil {
		t.Fatal(err)
	}
	h.records("after-c", facades.Completed, insuranceverifier.Verified, equipmentreserver.Reserved, pharmacydispatcher.Dispatched)
}

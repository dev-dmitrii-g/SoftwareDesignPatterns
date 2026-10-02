const $ = (id) => document.getElementById(id);
let selectedId = "";

function notice(message, error = false) {
  const box = $("notice");
  box.className = `alert ${error ? "alert-danger" : "alert-success"}`;
  box.textContent = message;
  box.scrollIntoView({ behavior: "smooth", block: "nearest" });
}

async function api(path, options = {}) {
  const response = await fetch(path, {
    ...options,
    headers: options.body ? { "Content-Type": "application/json" } : {},
  });
  const data = await response.json();
  if (!response.ok) {
    const failure = new Error(data.error || "Request failed");
    failure.view = data.view;
    throw failure;
  }
  return data;
}

function stateBadge(status) {
  const badge = document.createElement("span");
  badge.className = "badge-state " + (status === "COMPLETED" ? "completed" : status === "COMPENSATED" ? "compensated" : "pending");
  badge.textContent = status.replaceAll("_", " ");
  return badge;
}

function cell(row, value) {
  const td = document.createElement("td");
  td.textContent = value;
  row.append(td);
  return td;
}

function renderStock(items) {
  const body = $("stock-rows");
  body.replaceChildren();
  if (!items.length) {
    const row = document.createElement("tr");
    cell(row, "No stock set yet").colSpan = 3;
    body.append(row);
    return;
  }
  for (const item of items) {
    const row = document.createElement("tr");
    cell(row, item.kind === "pharmacy" ? "Medicine" : "Equipment");
    cell(row, item.name);
    cell(row, String(item.quantity)).classList.add("text-end", item.quantity === 0 ? "text-danger" : "text-success");
    body.append(row);
  }
}

function renderAdmissions(items) {
  const body = $("admission-rows");
  body.replaceChildren();
  if (!items.length) {
    const row = document.createElement("tr");
    cell(row, "No admissions yet").colSpan = 5;
    body.append(row);
    return;
  }
  for (const view of items) {
    const row = document.createElement("tr");
    cell(row, view.Admission.ID).classList.add("fw-semibold");
    const status = document.createElement("td");
    status.append(stateBadge(view.Admission.Status));
    row.append(status);
    cell(row, view.Admission.Equipment);
    cell(row, view.Admission.Medicine);
    const action = document.createElement("td");
    action.className = "text-end";
    const button = document.createElement("button");
    button.type = "button";
    button.className = "btn btn-outline-secondary btn-sm";
    button.textContent = "Inspect";
    button.addEventListener("click", () => inspect(view.Admission.ID));
    action.append(button);
    row.append(action);
    body.append(row);
  }
}

function showView(view) {
  selectedId = view.Admission.ID;
  $("detail-empty").classList.add("d-none");
  $("detail").classList.remove("d-none");
  $("detail-id").textContent = selectedId;
  $("detail-status").replaceChildren(stateBadge(view.Admission.Status));
  $("detail-insurance").textContent = view.Insurance;
  $("detail-equipment").textContent = view.Equipment;
  $("detail-pharmacy").textContent = view.Pharmacy;
  $("detail-reason").textContent = view.Admission.Reason ? `Reason: ${view.Admission.Reason}` : "";
  $("detail-reason").classList.toggle("d-none", !view.Admission.Reason);
  $("resume").disabled = !["PENDING_A", "PENDING_B", "PENDING_C", "COMPENSATING"].includes(view.Admission.Status);
  $("compensate").disabled = ["COMPLETED", "COMPENSATED"].includes(view.Admission.Status) || view.Pharmacy === "DISPATCHED";
}

async function refresh() {
  const data = await api("/api/state");
  renderStock(data.stock);
  renderAdmissions(data.admissions);
  if (selectedId) {
    const current = data.admissions.find((v) => v.Admission.ID === selectedId);
    if (current) showView(current);
  }
}

async function inspect(id) {
  try { showView((await api(`/api/admissions/${encodeURIComponent(id)}`)).view); }
  catch (error) { notice(error.message, true); }
}

async function act(path, payload, successText) {
  try {
    const data = await api(path, { method: "POST", body: JSON.stringify(payload) });
    if (data.view) showView(data.view);
    await refresh();
    notice(successText);
  } catch (error) {
    if (error.view) showView(error.view);
    await refresh();
    notice(error.message, true);
  }
}

document.querySelectorAll("[data-admission-action]").forEach((button) => {
  button.addEventListener("click", async () => {
    const form = $("admission-form");
    if (!form.reportValidity()) return;
    const values = Object.fromEntries(new FormData(form));
    const action = button.dataset.admissionAction;
    await act("/api/admissions", { ...values, action }, action === "start" ? "Insurance verified. Use Resume chain to continue." : "Admission completed.");
  });
});

$("stock-kind").addEventListener("change", (event) => {
  $("stock-name").value = event.target.value === "pharmacy" ? $("medicine").value : $("equipment").value;
});
$("stock-form").addEventListener("submit", async (event) => {
  event.preventDefault();
  const values = Object.fromEntries(new FormData(event.target));
  await act("/api/stock", { ...values, quantity: Number(values.quantity) }, "Stock saved.");
});
$("prepare-b").addEventListener("click", async () => {
  await act("/api/stock", { kind: "equipment", name: $("equipment").value.trim(), quantity: 0 }, "B failure ready. Run full admission with a new ID.");
});
$("prepare-c").addEventListener("click", async () => {
  try {
    await api("/api/stock", { method: "POST", body: JSON.stringify({ kind: "equipment", name: $("equipment").value.trim(), quantity: 1 }) });
    await act("/api/stock", { kind: "pharmacy", name: $("medicine").value.trim(), quantity: 0 }, "C failure ready. Run full admission with a new ID.");
  } catch (error) { notice(error.message, true); }
});
$("resume").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/resume`, {}, "Chain resumed."));
$("compensate").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/compensate`, {}, "Admission compensated."));
$("refresh").addEventListener("click", () => refresh().catch((error) => notice(error.message, true)));
refresh().catch((error) => notice(error.message, true));

const $ = (id) => document.getElementById(id);
let selectedId = "";
let selectedStatus = "";

function showTab(name) {
  for (const tab of ["admissions", "stock"]) {
    const active = tab === name;
    $(`tab-${tab}`).classList.toggle("active", active);
    $(`tab-${tab}`).setAttribute("aria-selected", String(active));
    $(`panel-${tab}`).classList.toggle("d-none", !active);
  }
  $("notice").className = "alert d-none";
}

function notice(message, error = false) {
  const box = $("notice");
  box.className = `alert ${error ? "alert-danger" : "alert-success"}`;
  box.textContent = message;
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
  badge.className = "status " + (status === "COMPLETED" ? "completed" : status === "COMPENSATED" ? "compensated" : "pending");
  badge.textContent = status.replaceAll("_", " ");
  return badge;
}

function cell(row, value) {
  const td = document.createElement("td");
  td.textContent = value;
  row.append(td);
  return td;
}

function coverageLabel(value) {
  if (value === "NONE") return "No coverage";
  if (["GOBMP", "OSMS", "DMS"].includes(value)) return value;
  return value;
}

function renderAdmissions(items) {
  const body = $("admission-rows");
  body.replaceChildren();
  if (!items.length) {
    const row = document.createElement("tr");
    cell(row, "No admissions").colSpan = 7;
    body.append(row);
    return;
  }
  for (const view of items) {
    const admission = view.Admission;
    const row = document.createElement("tr");
    cell(row, admission.ID);
    cell(row, admission.IIN);
    cell(row, coverageLabel(admission.Coverage));
    cell(row, admission.Equipment);
    cell(row, admission.Medicine);
    const status = document.createElement("td");
    status.append(stateBadge(admission.Status));
    row.append(status);
    const action = document.createElement("td");
    const button = document.createElement("button");
    button.className = "btn btn-outline-secondary btn-sm";
    button.type = "button";
    button.textContent = "View";
    button.addEventListener("click", () => showView(view));
    action.append(button);
    row.append(action);
    body.append(row);
  }
}

function renderStockRows(bodyId, items) {
  const body = $(bodyId);
  body.replaceChildren();
  if (!items.length) {
    const row = document.createElement("tr");
    cell(row, "No items").colSpan = 3;
    body.append(row);
    return;
  }
  for (const item of items) {
    const row = document.createElement("tr");
    cell(row, item.name);
    const quantity = document.createElement("td");
    const input = document.createElement("input");
    input.className = "form-control form-control-sm quantity-input";
    input.type = "number";
    input.min = "0";
    input.value = String(item.quantity);
    input.setAttribute("aria-label", `${item.name} quantity`);
    quantity.append(input);
    row.append(quantity);
    const action = document.createElement("td");
    const button = document.createElement("button");
    button.className = "btn btn-outline-primary btn-sm";
    button.type = "button";
    button.textContent = "Save";
    button.addEventListener("click", async () => {
      if (!input.reportValidity()) return;
      await act("/api/stock", { kind: item.kind, name: item.name, quantity: Number(input.value) }, "Stock saved");
    });
    action.append(button);
    row.append(action);
    body.append(row);
  }
}

function updateChoices(selectId, items, placeholder) {
  const select = $(selectId);
  const previous = select.value;
  const options = [];
  const empty = document.createElement("option");
  empty.value = "";
  empty.textContent = placeholder;
  options.push(empty);
  for (const item of items) {
    const option = document.createElement("option");
    option.value = item.name;
    option.textContent = `${item.name} (${item.quantity} available)`;
    options.push(option);
  }
  select.replaceChildren(...options);
  select.value = items.some((item) => item.name === previous) ? previous : (items[0]?.name || "");
}

function renderStock(items) {
  const equipment = items.filter((item) => item.kind === "equipment");
  const medicine = items.filter((item) => item.kind === "pharmacy");
  renderStockRows("equipment-rows", equipment);
  renderStockRows("medicine-rows", medicine);
  updateChoices("equipment", equipment, "Choose equipment");
  updateChoices("medicine", medicine, "Choose medicine");
}

function showView(view) {
  selectedId = view.Admission.ID;
  selectedStatus = view.Admission.Status;
  $("detail").classList.remove("d-none");
  $("detail-id").textContent = selectedId;
  $("detail-status").replaceChildren(stateBadge(view.Admission.Status));
  $("detail-insurance").textContent = view.Insurance;
  $("detail-equipment").textContent = view.Equipment;
  $("detail-pharmacy").textContent = view.Pharmacy;
  $("detail-reason").textContent = view.Admission.Reason ? `Reason: ${view.Admission.Reason}` : "";
  $("detail-reason").classList.toggle("d-none", !view.Admission.Reason);
  $("reserve").disabled = view.Admission.Status !== "PENDING_B";
  $("dispatch").disabled = view.Admission.Status !== "PENDING_C";
  $("resume").disabled = !["PENDING_A", "PENDING_B", "PENDING_C", "COMPENSATING"].includes(view.Admission.Status);
  $("resume").textContent = view.Admission.Status === "COMPENSATING" ? "Finish reversal" : "Complete all";
  $("compensate").disabled = ["COMPLETED", "COMPENSATED"].includes(view.Admission.Status);
  $("release").disabled = view.Admission.Status !== "COMPLETED";
}

async function refresh() {
  const data = await api("/api/state");
  renderAdmissions(data.admissions);
  renderStock(data.stock);
  if (selectedId) {
    const current = data.admissions.find((view) => view.Admission.ID === selectedId);
    if (current) showView(current);
  }
}

async function act(path, payload, successMessage) {
  let view;
  let message = successMessage;
  let error = false;
  try {
    const data = await api(path, { method: "POST", body: JSON.stringify(payload) });
    view = data.view;
  } catch (failure) {
    view = failure.view;
    message = failure.message;
    error = true;
  }
  if (view) showView(view);
  try {
    await refresh();
  } catch (failure) {
    message = failure.message;
    error = true;
  }
  notice(message, error);
  return Boolean(view);
}

$("tab-admissions").addEventListener("click", () => showTab("admissions"));
$("tab-stock").addEventListener("click", () => showTab("stock"));

for (const button of document.querySelectorAll("[data-admission-action]")) {
  button.addEventListener("click", async () => {
    const form = $("admission-form");
    if (!form.reportValidity()) return;
    const values = Object.fromEntries(new FormData(form));
    const action = button.dataset.admissionAction;
    await act("/api/admissions", { ...values, action }, action === "start" ? "Insurance verified" : "Admission completed");
  });
}

for (const [formId, kind] of [["equipment-form", "equipment"], ["medicine-form", "pharmacy"]]) {
  $(formId).addEventListener("submit", async (event) => {
    event.preventDefault();
    const values = Object.fromEntries(new FormData(event.target));
    await act("/api/stock", { ...values, kind, quantity: Number(values.quantity) }, "Stock saved");
  });
}
$("resume").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/resume`, {}, selectedStatus === "COMPENSATING" ? "Compensation completed" : "Admission completed"));
$("reserve").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/reserve`, {}, "Equipment reserved"));
$("dispatch").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/dispatch`, {}, "Medicine dispatched"));
$("compensate").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/compensate`, {}, "Admission compensated"));
$("release").addEventListener("click", () => act(`/api/admissions/${encodeURIComponent(selectedId)}/release`, {}, "Admission released"));
$("refresh").addEventListener("click", () => refresh().catch((failure) => notice(failure.message, true)));
refresh().catch((failure) => notice(failure.message, true));

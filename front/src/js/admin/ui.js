import { getMeta } from "./size";

const BACKEND_URL = import.meta.env.PROD
  ? location.protocol + "//" + location.host + "/"
  : import.meta.env.VITE_BACKEND_URL;

let token = document.getElementById("token");
let h = document.getElementById("input-h");
let w = document.getElementById("input-w");
let s = document.getElementById("status-text");

getMeta(BACKEND_URL + "place.png").then((meta) => {
  h.value = meta.h;
  w.value = meta.w;
});

document.getElementById("ch-size").addEventListener("click", () => {
  let h = parseInt(document.getElementById("input-h").value);
  let w = parseInt(document.getElementById("input-w").value);
  fetch(BACKEND_URL + "admin/resize", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Internal-Request": token.value,
    },
    body: JSON.stringify({ height: h, width: w }),
  }).then((resp) => {
    if (!resp.ok) {
      alert("Error changing size");
    }
  });
});

document.getElementById("save-place").addEventListener("click", () => {
  fetch(BACKEND_URL + "admin/save", {
    method: "GET",
    headers: {
      "X-Internal-Request": token.value,
    },
  }).then((resp) => {
    if (!resp.ok) {
      alert("Error changing size");
    }
  });
});

document.getElementById("pause-place").addEventListener("click", () => {
  fetch(BACKEND_URL + "admin/pause", {
    method: "POST",
    headers: {
      "X-Internal-Request": token.value,
    },
  }).then((resp) => {
    if (!resp.ok) {
      alert("Error pausing place");
    }
    return resp.json()
  }).then((resp) => {
    if (resp.paused) {
      s.innerText = "paused";
    } else {
      s.innerText = "not paused";
    }
  });
});

fetch(BACKEND_URL + "admin/pause", {
  method: "GET",
}).then((resp) => {
  if (!resp.ok) {
    alert("Error getting place status");
  }
  return resp.json()
}).then((resp) => {
  if (resp.paused) {
    s.innerText = "paused";
  } else {
    s.innerText = "not paused";
  }
});
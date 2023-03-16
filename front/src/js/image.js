const BACKEND_URL = import.meta.env.PROD
  ? location.protocol + "//" + location.host + "/"
  : import.meta.env.VITE_BACKEND_URL;
console.log("BACKEND_URL", BACKEND_URL);

let loaded = false;
let loadingp = document.querySelector("#loading-p");
let uiwrapper = document.querySelector("#ui-wrapper");

async function downloadProgress(resp) {
  let len = resp.headers.get("Content-Length");
  let a = new Uint8Array(len);
  let pos = 0;
  let reader = resp.body.getReader();
  while (true) {
    let { done, value } = await reader.read();
    if (value) {
      a.set(value, pos);
      pos += value.length;
      loadingp.innerHTML =
        "downloading map " + Math.round((pos / len) * 100) + "%";
    }
    if (done) break;
  }
  return a;
}

async function setImage(data, glWindow) {
  let img = new Image();
  let blob = new Blob([data], { type: "image/png" });
  let blobUrl = URL.createObjectURL(blob);
  img.src = blobUrl;
  let promise = new Promise((resolve, reject) => {
    img.onload = () => {
      glWindow.setTexture(img);
      glWindow.draw();
      resolve();
    };
    img.onerror = reject;
  });
  await promise;
}

export function loadBaseImage(glWindow) {
  fetch(BACKEND_URL + "place.png").then(async (resp) => {
    if (!resp.ok) {
      console.error("Error downloading map.");
      return null;
    }

    let buf = await downloadProgress(resp);
    await setImage(buf, glWindow);

    loaded = true;
    loadingp.innerHTML = "";
    uiwrapper.setAttribute("hide", true);

    // Check out the GET parameters
    let params = new URLSearchParams(window.location.search);
    if (params.has("x") && params.has("y")) {
      let x = parseInt(params.get("x"));
      let y = parseInt(params.get("y"));
      glWindow.setPos(x, y);
    }
    if (params.has("z")) {
      let zoom = parseInt(params.get("z"));
      glWindow.setZoom(zoom);
    } else {
      // zoom in
      glWindow.setZoom(100);
    }
    glWindow.draw();
  });
}

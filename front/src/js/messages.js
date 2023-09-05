import { getUint32, getUint64 } from "./utils";
import { setTimeout } from "./timeout";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export function handleSocketSetPixel(glWindow, b) {
  if (b.byteLength != 11) return;
  let x = getUint32(b, 0);
  let y = getUint32(b, 4);
  let color = new Uint8Array(b.slice(8));
  glWindow.setPixelColor(x, y, color);
  glWindow.draw();
}

export function handleSocketSetTimeout(b) {
  if (b.byteLength != 8) return;
  let timeout = getUint64(b, 0);
  setTimeout(timeout);
}

export function handleSocketStatus(b) {
  if (b.byteLength != 4) return;
  let clients = getUint32(b, 0);
  let s = document.querySelector("#status-p");
  s.innerText = clients + " online";
  s.style.visibility = "visible";
}

export async function getUser() {
  let r = await fetch(BACKEND_URL + "auth/self", {credentials: "include"})
  let j = await r.json()
  
  return j
}
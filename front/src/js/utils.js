export function putUint32(b, offset, n) {
  let view = new DataView(b);
  view.setUint32(offset, n, false);
}

export function getUint32(b, offset) {
  let view = new DataView(b);
  return view.getUint32(offset, false);
}

export function getUint64(b, offset) {
  let view = new DataView(b);
  return view.getBigUint64(offset, false);
}

export function hexToRgb(hex) {
  let result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result
    ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16),
      }
    : null;
}

export function secondFormat(seconds) {
  let minutes = Math.floor(seconds / 60);
  let secs = (seconds % 60) + 1;
  return minutes + ":" + (secs < 10 ? "0" : "") + secs;
}

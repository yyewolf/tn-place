import { hexToRgb, secondFormat } from "./utils.js";
import { setTimeout as setTempTimeout, getTimeout } from "./timeout.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const GUI = (cvs, glWindow, gateway) => {
  let color = new Uint8Array([0, 0, 0]);
  let dragdown = false;
  let touchID = 0;
  let touchScaling = false;
  let lastMovePos = { x: 0, y: 0 };
  let lastWindowPos = { x: 0, y: 0 };
  let lastScalingDist = 0;
  let touchstartTime;
  let popup = document.querySelector("#popup");
  let tooltipDelay;

  setInterval(() => {
    if (getTimeout() > 0) {
      setTempTimeout(getTimeout() - 1);
      document.querySelector("#timer-p").innerHTML = secondFormat(getTimeout());
    } else {
      document.querySelector("#timer-p").innerHTML = "";
    }
  }, 1000);

  const colorWrapper = document.querySelector("#color-wrapper");
  const picker = document.querySelector("#color-picker");
  // prevent clicks on color wrapper from propagating to canvas

  picker.addEventListener("input", (e) => {
    let rgb = hexToRgb(picker.value);
    color[0] = rgb.r;
    color[1] = rgb.g;
    color[2] = rgb.b;
  });

  // ***************************************************
  // ***************************************************
  // Event Listeners
  //
  document.addEventListener("keydown", (ev) => {
    switch (ev.keyCode) {
      case 189:
      case 173:
        ev.preventDefault();
        zoomOut(1.2);
        break;
      case 187:
      case 61:
        ev.preventDefault();
        zoomIn(1.2);
        break;
    }
  });

  window.addEventListener("wheel", (ev) => {
    const zoom = glWindow.getZoom();
    const mousePos = { x: ev.clientX, y: ev.clientY };
    const canvasBounds = glWindow.getCanvas().getBoundingClientRect();
    const canvasCenter = {
      x: canvasBounds.left + canvasBounds.width / 2,
      y: canvasBounds.top + canvasBounds.height / 2,
    };
    const mouseOffset = {
      x: mousePos.x - canvasCenter.x,
      y: mousePos.y - canvasCenter.y,
    };
    const zoomFactor = ev.deltaY > 0 ? 1 / 1.05 : 1.05;
    const newZoom = zoom * zoomFactor;
    const cameraOffset = {
      x: -mouseOffset.x * (1 / zoom - 1 / newZoom),
      y: mouseOffset.y * (1 / zoom - 1 / newZoom),
    };
    const newCameraPos = {
      x: glWindow.getPos().x - cameraOffset.x,
      y: glWindow.getPos().y + cameraOffset.y,
    };

    glWindow.setZoom(newZoom);
    glWindow.setPos(newCameraPos.x, newCameraPos.y);
    glWindow.draw();

    const url = new URL(window.location.href);
    url.searchParams.set("x", newCameraPos.x);
    url.searchParams.set("y", newCameraPos.y);
    url.searchParams.set("z", newZoom);
    window.history.replaceState({}, "", url);
  });

  // window.addEventListener("wheel", (ev) => {
  //   let zoom = glWindow.getZoom();
  //   if (ev.deltaY > 0) {
  //     zoom /= 1.05;
  //   } else {
  //     zoom *= 1.05;
  //   }

  //   glWindow.setZoom(zoom);
  //   glWindow.draw();

  //   let url = new URL(window.location.href);
  //   url.searchParams.set("x", glWindow.getPos().x);
  //   url.searchParams.set("y", glWindow.getPos().y);
  //   url.searchParams.set("z", glWindow.getZoom());
  //   window.history.replaceState({}, "", url);
  // });

  document.querySelector("#zoom-in").addEventListener("click", () => {
    zoomIn(1.2);
  });

  document.querySelector("#zoom-out").addEventListener("click", () => {
    zoomOut(1.2);
  });

  document.querySelector("#place-color").addEventListener("click", (e) => {
    e.preventDefault();
    drawPixel({ x: window.innerWidth / 2, y: window.innerHeight / 2 }, color);
  });

  window.addEventListener("resize", (ev) => {
    glWindow.updateViewScale();
    glWindow.draw();
  });

  cvs.addEventListener("mousedown", (ev) => {
    switch (ev.button) {
      case 0:
        dragdown = true;
        lastMovePos = { x: ev.clientX, y: ev.clientY };
        lastWindowPos.x = glWindow.getPos().x;
        lastWindowPos.y = glWindow.getPos().y;
        break;
      case 1:
        pickColor({ x: ev.clientX, y: ev.clientY });
        break;
      case 2:
        if (ev.ctrlKey) {
          pickColor({ x: ev.clientX, y: ev.clientY });
        } else {
          const clickPos = { x: ev.clientX, y: ev.clientY };
          const pixel_pos = glWindow.click({ x: ev.clientX, y: ev.clientY });
          if (
            Math.floor(pixel_pos.x) == glWindow.outline.x &&
            Math.floor(pixel_pos.y) == glWindow.outline.y
          ) {
            return;
          }

          const movePos = {
            x:
              glWindow.getPos().x +
              (clickPos.x - window.innerWidth / 2) / glWindow.getZoom(),
            y:
              glWindow.getPos().y +
              (clickPos.y - window.innerHeight / 2) / glWindow.getZoom(),
          };
          glWindow.transitionToPos({ x: movePos.x, y: movePos.y }, 250);
          glWindow.draw();
          lastMovePos = movePos;
        }
    }
  });

  cvs.addEventListener("mouseup", (ev) => {
    dragdown = false;
    document.body.style.cursor = "auto";

    if (
      lastWindowPos.x == glWindow.getPos().x &&
      lastWindowPos.y == glWindow.getPos().y
    ) {
      if (ev.button === 0) {
        if (ev.ctrlKey) {
          pickColor({ x: ev.clientX, y: ev.clientY });
        } else {
          const clickPos = { x: ev.clientX, y: ev.clientY };
          const pixel_pos = glWindow.click({ x: ev.clientX, y: ev.clientY });
          if (
            Math.floor(pixel_pos.x) == glWindow.outline.x &&
            Math.floor(pixel_pos.y) == glWindow.outline.y
          ) {
            return;
          }

          const movePos = {
            x:
              glWindow.getPos().x +
              (clickPos.x - window.innerWidth / 2) / glWindow.getZoom(),
            y:
              glWindow.getPos().y +
              (clickPos.y - window.innerHeight / 2) / glWindow.getZoom(),
          };
          glWindow.transitionToPos({ x: movePos.x, y: movePos.y }, 250);
          glWindow.draw();
          lastMovePos = movePos;
        }
      }
    }

    let url = new URL(window.location.href);
    url.searchParams.set("x", glWindow.getPos().x);
    url.searchParams.set("y", glWindow.getPos().y);
    url.searchParams.set("z", glWindow.getZoom());
    window.history.replaceState({}, "", url);
  });

  document.addEventListener("mousemove", (ev) => {
    const movePos = { x: ev.clientX, y: ev.clientY };
    if (dragdown) {
      glWindow.move(movePos.x - lastMovePos.x, movePos.y - lastMovePos.y);
      document.body.style.cursor = "grab";
      glWindow.draw();
    }
    lastMovePos = movePos;

    // Hide tooltip
    popup.style.display = "none";

    // Handle outline if mouse is over canvas
    try {
      if (tooltipDelay) {
        clearTimeout(tooltipDelay);
      }
      tooltipDelay = setTimeout(() => {
        tooltipDelay = null;
        // Move popup to mouse position
        popup.style.left = ev.clientX + 10 + "px";
        popup.style.top = ev.clientY + 10 + "px";

        // Get popup text from server
        fetch(
          BACKEND_URL +
            "pixel/" +
            glWindow.outline.x +
            "/" +
            glWindow.outline.y +
            "/"
        )
          .then((res) => res.json())
          .then((data) => {
            popup.innerHTML = data.placer;
            popup.style.display = "block";
          });
      }, 500);
    } catch (e) {
      // ignore
    }
    glWindow.draw();
  });

  cvs.addEventListener("touchstart", (ev) => {
    let thisTouch = touchID;
    touchstartTime = new Date().getTime();
    lastMovePos = { x: ev.touches[0].clientX, y: ev.touches[0].clientY };
    if (ev.touches.length === 2) {
      touchScaling = true;
      lastScalingDist = null;
    }

    setTimeout(() => {
      if (thisTouch == touchID) {
        pickColor(lastMovePos);
        navigator.vibrate(200);
      }
    }, 350);

    let url = new URL(window.location.href);
    url.searchParams.set("x", glWindow.getPos().x);
    url.searchParams.set("y", glWindow.getPos().y);
    url.searchParams.set("z", glWindow.getZoom());
    window.history.replaceState({}, "", url);
  });

  cvs.addEventListener("touchend", (ev) => {
    touchID++;
    let elapsed = new Date().getTime() - touchstartTime;
    if (elapsed < 100) {
      if (drawPixel(lastMovePos, color)) {
        navigator.vibrate(10);
      }
    }
    if (ev.touches.length === 0) {
      touchScaling = false;
    }
  });

  document.addEventListener("touchmove", (ev) => {
    touchID++;
    if (touchScaling) {
      let dist = Math.hypot(
        ev.touches[0].pageX - ev.touches[1].pageX,
        ev.touches[0].pageY - ev.touches[1].pageY
      );
      if (lastScalingDist != null) {
        let delta = lastScalingDist - dist;
        if (delta < 0) {
          zoomIn(1 + Math.abs(delta) * 0.003);
        } else {
          zoomOut(1 + Math.abs(delta) * 0.003);
        }
      }
      lastScalingDist = dist;
    } else {
      let movePos = { x: ev.touches[0].clientX, y: ev.touches[0].clientY };
      glWindow.move(movePos.x - lastMovePos.x, movePos.y - lastMovePos.y);
      glWindow.draw();
      lastMovePos = movePos;
      // console.log("move");
      // Add x and y to GET parameters
      let url = new URL(window.location.href);
      url.searchParams.set("x", glWindow.getPos().x);
      url.searchParams.set("y", glWindow.getPos().y);
      url.searchParams.set("z", glWindow.getZoom());
      window.history.replaceState({}, "", url);
    }
  });

  cvs.addEventListener("contextmenu", () => {
    return false;
  });

  const pickColor = (pos) => {
    color = glWindow.getColor(glWindow.click(pos));
    let hex = "#";
    for (let i = 0; i < color.length; i++) {
      let d = color[i].toString(16);
      if (d.length == 1) d = "0" + d;
      hex += d;
    }
    picker.value = hex;
  };

  const drawPixel = (pos, color) => {
    if (getTimeout() > 0) {
      return false;
    }
    pos = glWindow.click(pos);
    if (pos) {
      const oldColor = glWindow.getColor(pos);
      for (let i = 0; i < oldColor.length; i++) {
        if (oldColor[i] != color[i]) {
          gateway.setPixel(pos.x, pos.y, color);
          return true;
        }
      }
    }
    return false;
  };

  const zoomIn = (factor) => {
    let zoom = glWindow.getZoom();
    glWindow.setZoom(zoom * factor);

    let url = new URL(window.location.href);
    url.searchParams.set("x", glWindow.getPos().x);
    url.searchParams.set("y", glWindow.getPos().y);
    url.searchParams.set("z", glWindow.getZoom());
    window.history.replaceState({}, "", url);

    glWindow.draw();
  };

  const zoomOut = (factor) => {
    let zoom = glWindow.getZoom();
    glWindow.setZoom(zoom / factor);

    let url = new URL(window.location.href);
    url.searchParams.set("x", glWindow.getPos().x);
    url.searchParams.set("y", glWindow.getPos().y);
    url.searchParams.set("z", glWindow.getZoom());
    window.history.replaceState({}, "", url);

    glWindow.draw();
  };
};

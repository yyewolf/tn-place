import { hexToRgb, secondFormat } from "./utils.js";
import { setTimeout as setTempTimeout, getTimeout } from "./timeout.js";

const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;

export const GUI = (cvs, glWindow, gateway) => {
	let color = new Uint8Array([242, 243, 244]);
	let dragdown = false;
	let touchID = 0;
	let touchScaling = false;
	let lastMovePos = { x: 0, y: 0 };
	let lastWindowPos = { x: 0, y: 0 };
	let lastScalingDist = 0;
	let touchstartTime;
	let outline = { x: 0, y: 0, originalColor: new Uint8Array([0, 0, 0]) }
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

	// // fill colors
	// let sixteenColorsPalette = ["#000000", "#0000FF", "#00FF00", "#00FFFF", "#FF0000", "#FF00FF", "#FFFF00", "#FFFFFF", "#808080", "#000080", "#008000", "#008080", "#800000", "#800080", "#808000", "#C0C0C0"];
	// // <input type="button" class="color-square">
	// // <div class="inside-square" style="background-color: #000000;"></div>
	// // </input>
	// for (let i = 0; i < 16; i++) {
	// 	let btn = document.createElement("div");
	// 	let inside = document.createElement("input");
	// 	inside.type = "button";
	// 	btn.classList.add("color-square");
	// 	inside.classList.add("inside-square");
	// 	if (i == 0) inside.classList.add("inside-square-selected");
	// 	inside.style.backgroundColor = sixteenColorsPalette[i];
	// 	btn.addEventListener("click", (e) => {
	// 		e.preventDefault();
	// 		let rgb = hexToRgb(sixteenColorsPalette[i]);
	// 		color[0] = rgb.r;
	// 		color[1] = rgb.g;
	// 		color[2] = rgb.b;
	// 		inside.style.backgroundColor = sixteenColorsPalette[i];
	// 		document.querySelector(".inside-square-selected").classList.remove("inside-square-selected");
	// 		inside.classList.add("inside-square-selected");
	// 	});
	// 	btn.appendChild(inside);
	// 	colorWrapper.appendChild(btn);
	// }

	// prevent clicks on color wrapper from propagating to canvas
	colorWrapper.addEventListener("click", (e) => {
		// trigger color picker
		picker.click();
	});

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
	document.addEventListener("keydown", ev => {
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

	window.addEventListener("wheel", ev => {
		let zoom = glWindow.getZoom();
		if (ev.deltaY > 0) {
			zoom /= 1.05;
		} else {
			zoom *= 1.05;
		}

		glWindow.setZoom(zoom);
		glWindow.draw();

		let url = new URL(window.location.href);
		url.searchParams.set("x", glWindow.getPos().x);
		url.searchParams.set("y", glWindow.getPos().y);
		url.searchParams.set("z", glWindow.getZoom());
		window.history.replaceState({}, "", url);
	});

	document.querySelector("#zoom-in").addEventListener("click", () => {
		zoomIn(1.2);
	});

	document.querySelector("#zoom-out").addEventListener("click", () => {
		zoomOut(1.2);
	});

	window.addEventListener("resize", ev => {
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
					drawPixel({ x: ev.clientX, y: ev.clientY }, color);
				}
		}
	});

	cvs.addEventListener("mouseup", (ev) => {
		dragdown = false;
		document.body.style.cursor = "auto";

		if (lastWindowPos.x == glWindow.getPos().x && lastWindowPos.y == glWindow.getPos().y) {
			if (ev.button === 0) {
				if (ev.ctrlKey) {
					pickColor({ x: ev.clientX, y: ev.clientY });
				} else {
					drawPixel({ x: ev.clientX, y: ev.clientY }, color);
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
			let color = glWindow.getColor(outline);
			glWindow.setPixelColor(outline.x + 0.5, outline.y + 0.5, color);

			let pos = glWindow.click({ x: ev.clientX, y: ev.clientY });
			outline = { x: pos.x, y: pos.y }
			outline.x = Math.floor(outline.x);
			outline.y = Math.floor(outline.y);
			color = glWindow.getColor(outline);
			glWindow.setPixelBorder(outline.x, outline.y, color);

			if (tooltipDelay) {
				clearTimeout(tooltipDelay);
			}
			tooltipDelay = setTimeout(() => {
				tooltipDelay = null;
				// Move popup to mouse position
				popup.style.left = ev.clientX + 10 + "px";
				popup.style.top = ev.clientY + 10 + "px";

				// Get popup text from server
				fetch(BACKEND_URL + "pixel/" + outline.x + "/" + outline.y + "/")
					.then(res => res.json())
					.then(data => {
						popup.innerHTML = data.placer;
						popup.style.display = "block";
					})
			}, 500);
		} catch (e) {
			// ignore
		}
		glWindow.draw();
	});

	cvs.addEventListener("touchstart", (ev) => {
		let thisTouch = touchID;
		touchstartTime = (new Date()).getTime();
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
		let elapsed = (new Date()).getTime() - touchstartTime;
		if (elapsed < 100) {
			if (drawPixel(lastMovePos, color)) {
				navigator.vibrate(10);
			};
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
				ev.touches[0].pageY - ev.touches[1].pageY);
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

	cvs.addEventListener("contextmenu", () => { return false; });

	const pickColor = (pos) => {
		color = glWindow.getColor(glWindow.click(pos));
		let hex = "#";
		for (let i = 0; i < color.length; i++) {
			let d = color[i].toString(16);
			if (d.length == 1) d = "0" + d;
			hex += d;
		}
		picker.value = hex;
	}

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
	}

	const zoomIn = (factor) => {
		let zoom = glWindow.getZoom();
		glWindow.setZoom(zoom * factor);

		let url = new URL(window.location.href);
		url.searchParams.set("x", glWindow.getPos().x);
		url.searchParams.set("y", glWindow.getPos().y);
		url.searchParams.set("z", glWindow.getZoom());
		window.history.replaceState({}, "", url);

		glWindow.draw();
	}

	const zoomOut = (factor) => {
		let zoom = glWindow.getZoom();
		glWindow.setZoom(zoom / factor);

		let url = new URL(window.location.href);
		url.searchParams.set("x", glWindow.getPos().x);
		url.searchParams.set("y", glWindow.getPos().y);
		url.searchParams.set("z", glWindow.getZoom());
		window.history.replaceState({}, "", url);

		glWindow.draw();
	}
}
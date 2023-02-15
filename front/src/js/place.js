import { Gateway } from "./gateway";
import { GLWindow } from "./glwindow";
import { loadBaseImage } from "./image";
import { GUI } from "./gui";
import { handleSocketSetPixel, handleSocketSetTimeout, handleSocketStatus } from "./messages";

// this is the listener
let listeners = [
    [
        "timeout",
        (b) => {
            handleSocketSetTimeout(b);
        }
    ],
    [
        "refresh",
        (b) => {
            if (b.byteLength == 32) {
                loadBaseImage(glWindow);
            }
        }
    ],
    [
        "status",
        (b) => {
            handleSocketStatus(b);
        }
    ]
]

let gateway = new Gateway(listeners);
gateway.initConnection();

let cvs = document.querySelector("#viewport-canvas");
let glWindow = new GLWindow(cvs);

if (!glWindow.ok()) {
    alert("WebGL not supported");
}

// Add pixel listener
gateway.addListener("pixel", (b) => {
    handleSocketSetPixel(glWindow, b);
});

loadBaseImage(glWindow);

GUI(cvs, glWindow, gateway);
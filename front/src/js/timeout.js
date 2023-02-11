let timeout = 0;

export function setTimeout(t) {
    timeout = parseInt(t);
}

export function getTimeout() {
    return parseInt(timeout);
}

"use strict";

import htmx from "htmx.org";

import { toggleNavigation } from "./handlers";

window.htmx = htmx;
window.toggleNavigation = toggleNavigation;

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("htmx:beforeSwap", (event) => {
        const status = event.detail.xhr.status;
        if (status >= 200 && status < 500) {
            // I want to process client errors as well
            event.detail.shouldSwap = true;
            event.detail.isError = false;
        }
    });
});

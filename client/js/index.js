"use strict";

import htmx from "htmx.org";

import { toggleNavigation, htmxBeforeSwap } from "./handlers";

window.htmx = htmx;
window.toggleNavigation = toggleNavigation;

document.addEventListener("DOMContentLoaded", () => {
    document.body.addEventListener("htmx:beforeSwap", htmxBeforeSwap);
});

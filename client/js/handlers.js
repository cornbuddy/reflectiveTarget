"use strict";

export function toggleNavigation() {
    const NAVBAR_WIDTH = "250px";
    const nav = document.querySelector("nav");
    const main = document.querySelector("main");
    const toggler = document.querySelector("#sidebar-toggler");
    const isOpened = nav.style.width == NAVBAR_WIDTH;
    if (isOpened) {
        nav.style.visibility = "hidden";
        nav.style.width = "0";
        main.style.marginLeft = "0";
        toggler.style.marginLeft = "0";
    } else {
        nav.style.visibility = "visible";
        nav.style.width = NAVBAR_WIDTH;
        main.style.marginLeft = NAVBAR_WIDTH;
        toggler.style.marginLeft = NAVBAR_WIDTH;
    }
}

export function htmxBeforeSwap(event) {
    // TODO: do not handle 3xx codes
    const status = event.detail.xhr.status;
    if (status >= 200 && status < 500) {
        // I want to process client errors as well
        event.detail.shouldSwap = true;
        event.detail.isError = false;
    }
};

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
    const status = event.detail.xhr.status;
    if (status >= 200 && status <= 299) {
        event.detail.shouldSwap = true;
        event.detail.isError = false;
    } else if (status >= 300 && status <= 399) {
        event.detail.shouldSwap = false;
        event.detail.isError = false;
    } else if (status >= 400 && status <= 499) {
        event.detail.shouldSwap = true;
        event.detail.isError = true;
    } else {
        event.detail.shouldSwap = false;
        event.detail.isError = true;
    }
};

"use strict";

import htmx from "htmx.org";

import "../css/index.css";
import { htmxBeforeSwap, toggleNavigation, addQuestion } from "./handlers.js";

window.htmx = htmx;
window.toggleNavigation = toggleNavigation;
window.addQuestion = addQuestion;

htmx.onLoad(() => {
    document.body.addEventListener("htmx:beforeSwap", htmxBeforeSwap);
});

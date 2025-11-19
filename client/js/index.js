"use strict";

import htmx from "htmx.org";

import { toggleNavigation } from "./handlers";

window.htmx = htmx;
window.toggleNavigation = toggleNavigation;

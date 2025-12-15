"use strict";

import { test, expect } from "@jest/globals";

import { htmxBeforeSwap } from "./handlers";

test.each([
    { status: 101, shouldSwap: false, isError: true },
    { status: 200, shouldSwap: true, isError: false },
    { status: 205, shouldSwap: true, isError: false },
    { status: 299, shouldSwap: true, isError: false },
    { status: 300, shouldSwap: false, isError: false },
    { status: 307, shouldSwap: false, isError: false },
    { status: 399, shouldSwap: false, isError: false },
    { status: 400, shouldSwap: true, isError: true },
    { status: 418, shouldSwap: true, isError: true },
    { status: 499, shouldSwap: true, isError: true },
    { status: 500, shouldSwap: false, isError: true },
    { status: 505, shouldSwap: false, isError: true },
    { status: 599, shouldSwap: false, isError: true },
    { status: 666, shouldSwap: false, isError: true },
])("htmxBeforeSwap, status code $status", ({ status, shouldSwap, isError }) => {
    const event = { detail: { xhr: { status } } };
    htmxBeforeSwap(event);
    expect(event.detail.shouldSwap).toEqual(shouldSwap);
    expect(event.detail.isError).toEqual(isError);
});

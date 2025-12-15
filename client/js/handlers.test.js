"use strict";

import { test, expect } from "@jest/globals";

import { htmxBeforeSwap } from "./handlers";

test.each([
    { status: 200, shouldSwap: true, isError: false },
])("htmxBeforeSwap, status code $status", ({ status, shouldSwap, isError }) => {
    const event = { status };
    htmxBeforeSwap(event);
    expect(event.detail.shouldSwap).toEqual(shouldSwap);
    expect(event.detail.isError).toEqual(isError);
});

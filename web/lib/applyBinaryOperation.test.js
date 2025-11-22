import { applyBinaryOperation } from "./ArrayDecodeProtocol.js";

function encodeJSON(json) {
    return new TextEncoder().encode(JSON.stringify(json));
}

/** -----------------------------------------
 *  Decode header fields
 *  -----------------------------------------
 *  bit 0-1 → op
 *  bit 2-3 → posSize   (0..3 bytes)
 *  bit 4-5 → dataSize  (0..3 bytes)
 *  bit 6   → partial
 *  bit 7   → bulk
 */
function parseHeader(header) {
    return {
        op: header & 0b11,
        posSize: (header >> 2) & 0b11,
        dataSize: (header >> 4) & 0b11,
        partial: (header >> 6) & 1,
        bulk: (header >> 7) & 1,
    };
}

/** ---------------------------------------------------------
 *  NEW buildBuffer — now header controls posSize/dataSize
 *  ---------------------------------------------------------
 */
function buildBuffer({ header, pos = 0, json = null, rangeEnd = null }) {
    const { posSize, dataSize, bulk } = parseHeader(header);
    const bufParts = [];

    // Header
    bufParts.push(Uint8Array.from([header]));

    // -----------------------------
    // Position or Range Begin
    // -----------------------------
    if (posSize > 0) {
        const arr = new Uint8Array(posSize);
        for (let i = 0; i < posSize; i++) {
            arr[posSize - 1 - i] = (pos >> (8 * i)) & 0xff;
        }
        bufParts.push(arr);
    }

    // -----------------------------
    // Range End (bulk = 1)
    // -----------------------------
    if (bulk && rangeEnd !== null) {
        const arr = new Uint8Array(posSize);
        for (let i = 0; i < posSize; i++) {
            arr[posSize - 1 - i] = (rangeEnd >> (8 * i)) & 0xff;
        }
        bufParts.push(arr);
    }

    // -----------------------------
    // JSON Payload
    // -----------------------------
    if (json !== null) {
        const jsonBytes = encodeJSON(json);
        let lenBuf = null;

        if (dataSize === 1) {
            lenBuf = Uint8Array.from([jsonBytes.length]);
        } else if (dataSize === 2) {
            lenBuf = new Uint8Array([jsonBytes.length >> 8, jsonBytes.length & 0xff]);
        } else if (dataSize === 3) {
            lenBuf = new Uint8Array([
                (jsonBytes.length >> 16) & 0xff,
                (jsonBytes.length >> 8) & 0xff,
                jsonBytes.length & 0xff,
            ]);
        }

        if (lenBuf) bufParts.push(lenBuf);
        bufParts.push(jsonBytes);
    }

    // Merge all parts
    const total = bufParts.reduce((n, b) => n + b.length, 0);
    const final = new Uint8Array(total);
    let offset = 0;

    for (const part of bufParts) {
        final.set(part, offset);
        offset += part.length;
    }

    return final;
}

/*
DELETE — header:
    op = 00
    posSize = 01
*/
test("DELETE single element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00000100, // op=00, posSize=01
        pos: 1,
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, 30]);
});

/*
UPDATE — header:
    op = 01
    posSize = 01
    dataSize = 01
*/
test("UPDATE element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00010101, // UPDATE, posSize=1, dataSize=1
        pos: 1,
        json: { a: 999 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, { a: 999 }, 30]);
});

/*
INSERT — header:
    op = 11
    posSize = 01
    dataSize = 01
*/
test("INSERT element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00010111, // CORRECTED header for insert!
        pos: 1,
        json: { x: 42 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, { x: 42 }, 20, 30]);
});

/*
BULK DELETE — header:
    bulk = 1
    op = 00
    posSize = 1
*/
test("BULK DELETE", () => {
    const target = [10, 20, 30, 40, 50];

    const header = 0b10000100;

    const buffer = buildBuffer({
        header,
        pos: 1,        // begin
        rangeEnd: 3,   // end
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, 50]);
});

/*
PARTIAL PATCH — header:
    partial = 1
    UPDATE ignored (partial overrides)
    posSize = 1
    dataSize = 1
*/
test("PARTIAL PATCH", () => {
    const target = [{ n: 1 }, { n: 2 }];

    const buffer = buildBuffer({
        header: 0b01010101, // partial=1, pos=1, len=1
        pos: 1,
        json: { n: 999 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([{ n: 1 }, { n: 999 }]);
});

import { applyBinaryOperation } from "./ArrayDecodeProtocol.js";

function encodeJSON(json) {
    const bytes = new TextEncoder().encode(JSON.stringify(json));
    return bytes;
}

function buildBuffer({ header, posSize = 0, pos = 0, dataSize = 0, json = null }) {
    const bufParts = [];

    // Header
    bufParts.push(Uint8Array.from([header]));

    // Position
    if (posSize > 0) {
        const arr = new Uint8Array(posSize);
        for (let i = 0; i < posSize; i++) {
            arr[posSize - 1 - i] = (pos >> (8 * i)) & 0xff;
        }
        bufParts.push(arr);
    }

    // JSON payload
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

        bufParts.push(lenBuf);
        bufParts.push(jsonBytes);
    }

    // Merge
    return new Uint8Array(bufParts.reduce((acc, b) => acc + b.length, 0))
        .map((_, i, __, flat = bufParts.flat()) => flat[i]);
}

/*
Header bits:

    op = 00
    posSize = 01 → 1 byte position
    no data
    no partial
    no bulk
*/
test("DELETE single element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00000100,
        posSize: 1,
        pos: 1,    // delete index 1
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, 30]);
});

/*
Header:
    op = 01
    posSize = 01
    dataSize = 01 (1-byte JSON length)
*/
test("UPDATE element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00010101,
        posSize: 1,
        pos: 1,
        dataSize: 1,
        json: { a: 999 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, { a: 999 }, 30]);
});

/*
Header:

    op = 11
    posSize = 1
    dataSize = 1
*/
test("INSERT element", () => {
    const target = [10, 20, 30];

    const buffer = buildBuffer({
        header: 0b00011101,
        posSize: 1,
        pos: 1,
        dataSize: 1,
        json: { x: 42 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, { x: 42 }, 20, 30]);
});

/*
Header:

    op = 00
    posSize = 1
    bulk = 1
*/
test("BULK DELETE", () => {
    const target = [10, 20, 30, 40, 50];

    const header = 0b10000100; // bulk + delete + posSize=1
    const buffer = new Uint8Array([header, 1, 3]); // delete from 1 to 3

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([10, 50]);
});

/*
Header:

    partial = 1 (bit 6)
    op = update (op ignored; partial makes PATCH mode)
    posSize = 1
    dataSize = 1
*/
test("PARTIAL PATCH", () => {
    const target = [{ n: 1 }, { n: 2 }];

    const buffer = buildBuffer({
        header: 0b01010101,
        posSize: 1,
        pos: 1,
        dataSize: 1,
        json: { n: 999 },
    });

    applyBinaryOperation(buffer, target);

    expect(target).toEqual([{ n: 1 }, { n: 999 }]);
});

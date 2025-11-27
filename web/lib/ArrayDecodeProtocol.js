//
//  Binary Decoder and Array/Object Transformer Protocol
//

export function applyBinaryOperation(buffer, target, debug) {
    const view = new DataView(buffer.buffer, buffer.byteOffset, buffer.length);
    debug = debug || false;
    let offset = 0;

    if (debug) {
        console.log("Buffer received at the start of applyBinaryOperation:", buffer);
    }

    // ---- 1. Read header ----
    const header = view.getUint8(offset++);
    if (debug) {
        console.log("Header:", header);
    }

    const op       = header & 0b11;           // bits 0-1
    const posSize  = (header >> 2) & 0b11;    // bits 2-3
    const dataSize = (header >> 4) & 0b11;    // bits 4-5
    const partial  = (header >> 6) & 1;       // bit 6
    const bulk     = (header >> 7) & 1;       // bit 7

    if (debug) {
        console.log("op", op);
        console.log("posSize", posSize);
        console.log("dataSize", dataSize);
        console.log("partial", partial);
        console.log("bulk", bulk);
        console.log("------------");
    }

    // ---- Reading helpers ----
    function readSizedInt(size) {
        switch (size) {
            case 0: return 0;
            case 1: return view.getUint8(offset++);
            case 2: {
                const v = view.getUint16(offset, false);
                offset += 2;
                return v;
            }
            case 3: {
                const v =
                    (view.getUint8(offset) << 16) |
                    (view.getUint8(offset + 1) << 8) |
                    view.getUint8(offset + 2);
                offset += 3;
                return v;
            }
        }
    }

    function readJSON(sizeIndicator) {
        const dataLen = readSizedInt(sizeIndicator);
        const bytes = new Uint8Array(buffer.buffer, buffer.byteOffset + offset, dataLen);
        offset += dataLen;
        return JSON.parse(new TextDecoder().decode(bytes));
    }

    //
    // =====================================================================
    //   SIMPLE (non-bulk) OPERATIONS
    // =====================================================================
    //
    if (!bulk) {
        const pos = readSizedInt(posSize);

        if (!partial) {
            switch (op) {
                case 0b00: // DELETE
                    if (pos < target.length) {
                        target.splice(pos, 1);
                    }
                    return;

                case 0b01: { // UPDATE
                    const value = readJSON(dataSize);
                    target[pos] = value;
                    return;
                }

                case 0b11: { // INSERT
                    const value = readJSON(dataSize);

                    // Ensure the array is long enough
                    if (pos > target.length) {
                        target.length = pos; // creates holes
                    }

                    target.splice(pos, 0, value);
                    return;
                }
            }
        }

        // ---- PARTIAL UPDATE (non-bulk) ----
        const patch = readJSON(dataSize);
        applyPartialPatch(target, pos, patch);
        return;
    }

    //
    // =====================================================================
    //   BULK OPERATIONS
    // =====================================================================
    //
    let start = readSizedInt(posSize);
    let end   = readSizedInt(posSize);
    let count = (end - start) + 1;

    // ----------------------------------------------------
    // Bulk Delete (no partial mode)
    // ----------------------------------------------------
    if (op === 0b00 && !partial) {
        target.splice(start, count);
        return;
    }

    // ----------------------------------------------------
    // Bulk FULL Update / Insert (dense update)
    // ----------------------------------------------------
    if (!partial) {
        switch (op) {
            case 0b01: { // FULL BULK UPDATE
                if (end >= target.length) target.length = end + 1;
                for (let i = start; i <= end; i++) {
                    target[i] = readJSON(dataSize);
                }
                return;
            }

            case 0b11: { // FULL BULK INSERT
                const values = new Array(count);
                for (let i = 0; i < count; i++) {
                    values[i] = readJSON(dataSize);
                }
                if (start > target.length) {
                    target.length = start;
                }
                target.splice(start, 0, ...values);
                return;
            }
        }
    }

    //
    // =====================================================================
    //   BULK + PARTIAL  (Sparse Partial Updates)
    // =====================================================================
    //
    while (offset < view.byteLength) {
        const innerPos = readSizedInt(posSize);
        const patch = readJSON(dataSize);
        applyPartialPatch(target, innerPos, patch);
    }
}


// ===================================================================
//  Utility: apply a partial patch to array or object
// ===================================================================
function applyPartialPatch(target, posOrKey, patchValue) {
    if (Array.isArray(target)) {
        target[posOrKey] = patchValue;
        return;
    }

    if (target && typeof target === "object") {
        target[posOrKey] = patchValue;
        return;
    }

    throw new Error("Partial patch applied on a non-container type");
}

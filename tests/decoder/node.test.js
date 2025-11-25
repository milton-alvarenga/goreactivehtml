import { applyBinaryOperation } from "./../../web/lib/ArrayDecodeProtocol.js";

const buffer = Buffer.from([23, 0, 3, 34, 65, 34]);  // Mimic the Go buffer
console.log("Received input:", buffer);

const target = [];
applyBinaryOperation(buffer, target, true);
console.log("Target after applying binary operation:", target);
process.stdout.write(JSON.stringify(target) + "\n");

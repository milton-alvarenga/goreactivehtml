import { applyBinaryOperationByte } from "./../../web/lib/ArrayDecodeProtocolByte.js";
import fs from 'fs';

let debug = false

if(debug){
    process.stdout.write("Starting node\n");
}

(async () => {
    try {
        // Expect JSON input with optional `initial` and `payload`
        const raw = fs.readFileSync(0); // Read from STDIN

        if (raw.length === 0) {
            throw new Error("Received empty input.");
        }

        if (raw.length < 4) throw new Error("Input too short for length prefix");

        // Read 4-byte length prefix
        const initialLen = raw.readUInt32BE(0);

        // Read initial JSON
        const initialJSON = raw.slice(4, 4 + initialLen).toString();
        const initial = JSON.parse(initialJSON);

        // Remaining bytes are payload
        const payloadBuffer = raw.slice(4 + initialLen);

        const target = Array.isArray(initial) ? initial : [];

        if(debug){
            console.log("Received input:", raw.toString()); // Ensure the data is as expected
            console.log("Received initialJSON:", initialJSON.toString()); // Ensure the data is as expected
            console.log("Received initial:", initial); // Ensure the data is as expected
            console.log("Initial target:", target);
            console.log("Payload buffer:", payloadBuffer);
        }

        applyBinaryOperationByte(payloadBuffer, target, debug); 
        
        process.stdout.write(JSON.stringify(target) + "\n");
    } catch (err) {

        console.error("Error processing input:", err);
        process.exit(1);
    }
})();

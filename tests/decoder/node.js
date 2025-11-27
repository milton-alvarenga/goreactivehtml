import { applyBinaryOperation } from "./../../web/lib/ArrayDecodeProtocol.js";
import fs from 'fs';

let debug = false

if(debug){
    process.stdout.write("Starting node\n");
}

(async () => {
    try {
        const buffer = fs.readFileSync(0); // Read from STDIN
        if(debug){
            console.log("Received input:", buffer.toString()); // Ensure the data is as expected
        }

        if (buffer.length === 0) {
            throw new Error("Received empty input.");
        }

        const target = [];
        applyBinaryOperation(buffer, target, debug);
        
        process.stdout.write(JSON.stringify(target) + "\n");
    } catch (err) {

        console.error("Error processing input:", err);
        process.exit(1);
    }
})();

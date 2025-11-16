class WebSocketEvents {
    ws;
    pendingRequests;
    RequestId = new Uint8Array(1);

    constructor(ws) {
        this.ws = ws;
        this.subscriptions = {};
        //zero is reserved for all subscriptions requests
        this.pendingRequests = new Array(256);
        this.RequestId[0] = 0;
        this.ws.onmessage = this.processBinaryResponse.bind(this)
        this.ws.onerror = this.processBinaryErrorResponse.bind(this)
        this.ws.onclose = () => {
            ws.close()
        }
    }

    getNextRequestId() {
        while (this.RequestId[0] > 0 && this.pendingRequests[this.RequestId[0]] !== null){
            this.RequestId[0]++; // wrap around after 255
        }
        return this.RequestId[0];
    }

    // Method to subscribe to a specific event (eventType could be MsgType, Destination, or other criteria)
    subscribe(eventType, callback) {
        if (!this.subscriptions[eventType]) {
            this.subscriptions[eventType] = [];
        }
        this.subscriptions[eventType].push(callback);
    }

    // Unsubscribe from a specific event type and callback
    unsubscribe(eventType, callback) {
        const eventCallbacks = this.subscriptions[eventType];
        
        if (!eventCallbacks) return; // No such eventType subscription
        
        // Find the index of the callback in the event's callback list
        const callbackIndex = eventCallbacks.indexOf(callback);
        
        if (callbackIndex !== -1) {
            // Remove the callback from the list
            eventCallbacks.splice(callbackIndex, 1);
        }

        // If no more callbacks are left for this event type, delete the event type from subscriptions
        if (eventCallbacks.length === 0) {
            delete this.subscriptions[eventType];
        }
    }

    processBinaryResponse(event) {
        let response = this.unmarshalBinaryResponse(event.data)
        if (response.ReqId == 0) {
            // Trigger the subscription callbacks if applicable
            this.triggerSubscriptions(response);
        } else if (this.pendingRequests[response.ReqId] ){
            if( response.MsgType === "S" ){
                this.pendingRequests[response.ReqId].resolve(response)
            } else {
                this.pendingRequests[response.ReqId].reject(response)
            }
            this.pendingRequests[response.ReqId] = null
        }
    }

    processBinaryErrorResponse(error){
        // If there was an error, reject the promise
        this.WebSocketEvents.pendingRequests[reqId].reject(error);
    }

    /*
    processResponse(event) {
        let response = this.unmarshalResponse(event.data)
        if (response.ReqId && this.pendingRequests[response.ReqId] ){
            if( response.MsgType === "S" ){
                this.pendingRequests[response.ReqId].resolve(response)
            } else {
                this.pendingRequests[response.ReqId].reject(response)
            }
            this.pendingRequests[response.ReqId] = null
        }
    }

    unmarshalResponse(wsStringResponse) {
        const parts = wsStringResponse.split(';');
        let top = parts[0];
        let ReqId = null;

        // E => Error
        // S => Success
        if (!['E','S'].includes(top[0])){
            ReqId = top[0].charCodeAt(0)
            top = top.substring(1)
        }

        return {
            ReqId: ReqId,
            MsgType: top[0],
            Destination: top.substring(1),
            Data: parts[1],
            Headers: null
        };
    }
    */
    // Method to trigger all relevant subscriptions based on the response
    triggerSubscriptions(response) {
        // Check if there are any subscriptions for the given response's eventType
        const eventType = `${response.Destination}`;
        
        if (this.subscriptions[eventType]) {
            // Loop through all the callbacks for the given eventType and invoke them
            for (const callback of this.subscriptions[eventType]) {
                callback(response);
            }
        }
    }

    /*
    Steps for Unmarshalling in JavaScript:
        Read ReqId (1 byte).
        Read MsgType (1 byte for length, followed by the ascii).
        Read Destination (2 bytes for length, followed by the string).
        Read Data (4 bytes for length, followed by the JSON string).
        Read Header (2 bytes for length, followed by the JSON string).
    */
    /*
    // Example usage:

    // Simulate the binary data received (as in Go Marshal output)
    let binaryData = new Uint8Array([
        1,             // ReqId
        83,            // MsgType 'S' (MsgTypeSuccessOutputMessage)
        0, 17,         // Length of destination ("exampleDestination")
        ...new TextEncoder().encode("exampleDestination"),
        0, 0, 0, 10,   // Length of Data (JSON string length)
        ...new TextEncoder().encode(JSON.stringify("Some data")),
        0, 0, 0, 15,   // Length of Header (JSON string length)
        ...new TextEncoder().encode(JSON.stringify({key: "value"}))
    ]);

    try {
        let clientOutput = unmarshal(binaryData);
        console.log(clientOutput);
    } catch (err) {
        console.error("Error unmarshalling:", err);
    }
    */
    unmarshalBinaryResponse(binaryData) {
        let offset = 0;

        // Helper function to read 1 byte
        function readByte() {
            return binaryData[offset++];
        }

        // Helper function to read N bytes
        function readBytes(length) {
            let bytes = binaryData.slice(offset, offset + length);
            offset += length;
            return bytes;
        }

        // ReqId (1 byte) - now required (must be non-zero)
        let reqId = readByte();

        // MsgType (1 byte for length + N bytes for content)
        let msgTypeLen = readByte(); // 1 byte for length of MsgType
        let msgType = new TextDecoder().decode(readBytes(msgTypeLen)); // Read the MsgType string

        // Destination (2 bytes for length + N bytes for content)
        let destLen = (readByte() << 8) | readByte(); // 2 bytes for length
        let destination = new TextDecoder().decode(readBytes(destLen));

        // Data (4 bytes for length + N bytes for JSON serialized content)
        let dataLen = (readByte() << 24) | (readByte() << 16) | (readByte() << 8) | readByte(); // 4 bytes for length
        let data = JSON.parse(new TextDecoder().decode(readBytes(dataLen)));

        // Header (2 bytes for length + N bytes for JSON serialized content)
        let headerLen = (readByte() << 8) | readByte(); // 2 bytes length
        let header = JSON.parse(new TextDecoder().decode(readBytes(headerLen)));

        return {
            ReqId: reqId,
            MsgType: msgType,
            Destination: destination,
            Data: data,
            Header: header
        };
    }
}

export default WebSocketEvents;
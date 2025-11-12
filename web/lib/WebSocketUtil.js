import WebSocketEvents from './WebSocketEvents'

class WebSocketUtil {
    constructor(url) {
        this.url = url
        this.ws = new WebSocket(url);
        this.connected = false
        this.ws.onopen = () => {
            this.connected = true
            this.onopen && this.onopen()
        }

        this.conn_type = {
            SUBSCRIBE:1,
            RPC:2,
            ENDPOINT:3
        }

        this.WebSocketEvents = new WebSocketEvents(this.ws)
        this.encoder = new TextEncoder();
    }

    //Permanent listen connection
    subscribe(destination,callback,data,header){
        this.WebSocketEvents.subscribe(destination, callback)

        const reqId = this.WebSocketEvents.getNextRequestId()

        send(binaryData, reqId)
    }

    unsubscribe(destination, callback){
        this.WebSocketEvents.unsubscribe(destination, callback)
    }

    //Like a client server request
    requestEndpoint(rest_method,endpoint,data,headers,callback_success,callback_fail){
        const reqId = this.WebSocketEvents.getNextRequestId()

        const binaryData = formatBinaryRequestEndpoint(reqId, method, endpoint, origin, data, header);

        send(binaryData,reqId,callback_success,callback_fail)
    }

    requestRPC(bff,method,params,headers,callback_success,callback_fail){
        const reqId = this.WebSocketEvents.getNextRequestId();

        const binaryData = formatRequestRPC(reqId, bff, method, params, header);

        send(binaryData,reqId,callback_success,callback_fail)
    }

    formatRequestEndpoint(reqId, endpoint, operation, origin, data){
        return this.formatTextRequest(reqId,endpoint, operation, origin, data)
    }

    formatRequestRPC(reqId, bff, method, params, header = {}) {
        const encoder = new TextEncoder();

        // Encode individual parts
        const bffBytes = encoder.encode(bff);
        const methodBytes = encoder.encode(method);
        const paramsBytes = encoder.encode(JSON.stringify(params));
        const headerBytes = encoder.encode(JSON.stringify(header));

        // Compute total length
        const totalLength =
            1 + // conn_type
            1 + // reqId
            2 + bffBytes.length + // 2 bytes length + bff
            1 + methodBytes.length + // 1 byte length + method
            4 + paramsBytes.length + // 4 bytes length + params
            2 + headerBytes.length;  // 2 bytes length + header

        const message = new Uint8Array(totalLength);
        let offset = 0;

        // --- conn_type (1 byte)
        message[offset++] = this.conn_type.RPC;

        // --- reqId (1 byte)
        message[offset++] = reqId;

        // --- bff length (2 bytes)
        message[offset++] = (bffBytes.length >> 8) & 0xff;
        message[offset++] = bffBytes.length & 0xff;

        // --- bff bytes
        message.set(bffBytes, offset);
        offset += bffBytes.length;

        // --- method length (1 byte)
        message[offset++] = methodBytes.length;

        // --- method bytes
        message.set(methodBytes, offset);
        offset += methodBytes.length;

        // --- params length (4 bytes)
        message[offset++] = (paramsBytes.length >> 24) & 0xff;
        message[offset++] = (paramsBytes.length >> 16) & 0xff;
        message[offset++] = (paramsBytes.length >> 8) & 0xff;
        message[offset++] = paramsBytes.length & 0xff;

        // --- params bytes
        message.set(paramsBytes, offset);
        offset += paramsBytes.length;

        // --- header length (2 bytes)
        message[offset++] = (headerBytes.length >> 8) & 0xff;
        message[offset++] = headerBytes.length & 0xff;

        // --- header bytes
        message.set(headerBytes, offset);
        offset += headerBytes.length;

        return message;
    }


    formatRequestSubscribe(topic, data, header = {}) {
        const encoder = new TextEncoder();

        // Encode data (assuming it's a JSON string)
        const dataBytes = encoder.encode(JSON.stringify(data));
        // Encode header (also assuming it's a JSON string)
        const headerBytes = encoder.encode(JSON.stringify(header));


        // Encode topic as bytes (2 bytes for length, then the UTF-8 encoded topic string)
        const topicBytes = encoder.encode(topic);
        const topicLength = topicBytes.length;
        if (topicLength > 65535) {
            throw new Error("Topic length exceeds maximum of 65535 bytes");
        }


        // Calculate total message size
        const totalLength =
            1 + // conn_type (1 byte)
            2 + topicLength + // topic length (2 bytes + topic length)
            4 + dataBytes.length + // data length (4 bytes + data bytes)
            2 + headerBytes.length; // header length (2 bytes + header bytes)

        const message = new Uint8Array(totalLength);
        let offset = 0;

        // --- Write conn_type (1 byte)
        message[offset++] = this.conn_type.SUBSCRIBE;

        // --- Write topic length (2 bytes)
        message[offset++] = (topicLength >> 8) & 0xff;
        message[offset++] = topicLength & 0xff;

        // --- Write topic bytes
        message.set(topicBytes, offset);
        offset += topicBytes.length;

        // --- Write data length (4 bytes)
        message[offset++] = (dataBytes.length >> 24) & 0xff;
        message[offset++] = (dataBytes.length >> 16) & 0xff;
        message[offset++] = (dataBytes.length >> 8) & 0xff;
        message[offset++] = dataBytes.length & 0xff;

        // --- Write data bytes
        message.set(dataBytes, offset);
        offset += dataBytes.length;

        // --- Write header length (2 bytes)
        message[offset++] = (headerBytes.length >> 8) & 0xff;
        message[offset++] = headerBytes.length & 0xff;

        // --- Write header bytes
        message.set(headerBytes, offset);
        offset += headerBytes.length;

        return message;
    }


    formatBinaryRequestEndpoint(reqId, method, endpoint, origin, data, header = {}) {
        const encoder = new TextEncoder();

        // Encode strings
        const methodBytes = encoder.encode(method);
        const endpointBytes = encoder.encode(endpoint);
        const payloadBytes = encoder.encode(JSON.stringify(data));
        const headerBytes = encoder.encode(JSON.stringify(header));

        // Calculate total message size
        const totalLength =
            1 + // conn_type
            1 + // reqId
            1 + methodBytes.length + // method length (1 byte) + data
            2 + endpointBytes.length + // endpoint length (2 bytes)
            4 + payloadBytes.length + // payload length (4 bytes)
            2 + headerBytes.length; // header length (2 bytes)

        const message = new Uint8Array(totalLength);
        let offset = 0;

        // --- Write conn_type (1 byte)
        message[offset++] = this.conn_type.ENDPOINT;

        // --- Write reqId (1 byte)
        message[offset++] = reqId;

        // --- Write method length (1 byte)
        message[offset++] = methodBytes.length;

        // --- Write method bytes
        message.set(methodBytes, offset);
        offset += methodBytes.length;

        // --- Write endpoint length (2 bytes)
        message[offset++] = (endpointBytes.length >> 8) & 0xff;
        message[offset++] = endpointBytes.length & 0xff;

        // --- Write endpoint bytes
        message.set(endpointBytes, offset);
        offset += endpointBytes.length;

        // --- Write payload length (4 bytes)
        message[offset++] = (payloadBytes.length >> 24) & 0xff;
        message[offset++] = (payloadBytes.length >> 16) & 0xff;
        message[offset++] = (payloadBytes.length >> 8) & 0xff;
        message[offset++] = payloadBytes.length & 0xff;

        // --- Write payload bytes
        message.set(payloadBytes, offset);
        offset += payloadBytes.length;

        // --- Write header length (2 bytes)
        message[offset++] = (headerBytes.length >> 8) & 0xff;
        message[offset++] = headerBytes.length & 0xff;

        // --- Write header bytes
        message.set(headerBytes, offset);
        offset += headerBytes.length;

        return message;
    }


    send(binaryData, reqId, callback_success,callback_fail){
        if (!this.connected){
            throw new Error("WebSocket not connected "+this.url);
        }

        if (this.ws.readyState !== WebSocket.OPEN ) {
            throw new Error("WebSocket is not open. Could not connect on "+this.url)
        }
        
        // Add the resolve/reject callbacks to pendingRequests using reqId
        this.WebSocketEvents.pendingRequests[reqId] = [callback_success,callback_fail]

        this.ws.send(binaryData)
    }

    normalizeEndpoint(endpoint) {
        if (!endpoint.startsWith('/')) {
            return '/' + endpoint;
        }
        return endpoint;
    }
}

export default WebSocketUtil;
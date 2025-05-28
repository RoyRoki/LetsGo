/**
 *  npx ts-node websocket-matching-test.ts
 */

import WebSocket from "ws";

const SERVER_URL = "ws://localhost:8080/ws/chat";

// Constants from your service
const StatusCodes = {
  SERVER_READY: "902",
  PARTNER_MSG: "904",
  SYSTEM_MSG: "905",
  NO_MATCH: "912",
  MATCH_FOUND: "913",
  PARTNER_DISCONNECTED: "920",
};

interface Message {
  status: string;
  content?: string;
  from?: number;
  timestamp?: number;
}

class TestClient {
  name: string;
  socket: WebSocket;
  isReady = false;
  partnerConnected = false;
  lastMessage: Message | null = null;

  private messageQueue: Message[] = [];
  private resolveNextMsg: ((msg: Message) => void) | null = null;

  constructor(name: string) {
    this.name = name;
    this.socket = new WebSocket(SERVER_URL);

    this.socket.on("open", () => {
      console.log(`[${this.name}] Connected to server`);
    });

    this.socket.on("message", (data) => {
      try {
        const msg: Message = JSON.parse(data.toString());
        this.lastMessage = msg;
        this.enqueueMessage(msg);
        this.handleMessage(msg);
      } catch (e) {
        console.error(`[${this.name}] Failed to parse message`, e);
      }
    });

    this.socket.on("close", () => {
      console.log(`[${this.name}] Connection closed`);
    });

    this.socket.on("error", (e) => {
      console.error(`[${this.name}] WebSocket error`, e);
    });
  }

  private enqueueMessage(msg: Message) {
    if (this.resolveNextMsg) {
      this.resolveNextMsg(msg);
      this.resolveNextMsg = null;
    } else {
      this.messageQueue.push(msg);
    }
  }

  waitForMessage(timeoutMs = 5000): Promise<Message> {
    if (this.messageQueue.length > 0) {
      return Promise.resolve(this.messageQueue.shift()!);
    }
    return new Promise((resolve, reject) => {
      this.resolveNextMsg = resolve;
      setTimeout(() => {
        this.resolveNextMsg = null;
        reject(new Error("Timeout waiting for message"));
      }, timeoutMs);
    });
  }

  private handleMessage(msg: Message) {
    if (msg.status === StatusCodes.SERVER_READY) {
      this.isReady = true;
    } else if (msg.status === StatusCodes.MATCH_FOUND) {
      this.partnerConnected = true;
    } else if (msg.status === StatusCodes.PARTNER_DISCONNECTED) {
      this.partnerConnected = false;
      console.log(`[${this.name}] Partner disconnected.`);
    }
  }

  sendTags(tags: string[]) {
    const payload = { tags };
    this.sendJson(payload);
    console.log(`[${this.name}] Sent tags: ${JSON.stringify(tags)}`);
  }

  sendMessage(content: string) {
    this.sendJson({ content });
    console.log(`[${this.name}] Sent message: ${content}`);
  }

  private sendJson(obj: object) {
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(obj));
    } else {
      throw new Error(`[${this.name}] Socket not open`);
    }
  }

  close() {
    this.socket.close();
  }
}

type TestResult = { name: string; passed: boolean; reason?: string };

async function runTest(
  testName: string,
  testFn: () => Promise<void>
): Promise<TestResult> {
  try {
    await testFn();
    console.log(`[PASS] ${testName}`);
    return { name: testName, passed: true };
  } catch (err: any) {
    console.error(`[FAIL] ${testName} - ${err.message || err}`);
    return { name: testName, passed: false, reason: err.message || "Failed" };
  }
}

async function waitForReady(client: TestClient) {
  const msg = await client.waitForMessage(5000);
  if (msg.status !== StatusCodes.SERVER_READY) {
    throw new Error(
      `Expected server ready status (902), got: ${msg.status} content: ${msg.content}`
    );
  }
}

async function waitForMatch(client: TestClient) {
  while (true) {
    const msg = await client.waitForMessage(10000);
    if (msg.status === StatusCodes.MATCH_FOUND) {
      return;
    }
    if (msg.status === StatusCodes.NO_MATCH) {
      // continue waiting
      continue;
    }
    throw new Error(
      `Expected match found (913) or no match (912), got ${msg.status} with content: ${msg.content}`
    );
  }
}

// All your test cases

async function testNoTagsConnect() {
  const client = new TestClient("Client-NoTags");
  await waitForReady(client);
  client.sendTags([]);
  // Expect no match found yet
  const msg = await client.waitForMessage(3000);
  if (msg.status !== StatusCodes.NO_MATCH)
    throw new Error("Expected no match found (912)");
  client.close();
}

async function testSingleTagDifferent() {
  const client1 = new TestClient("Client1-TagA");
  const client2 = new TestClient("Client2-TagB");

  await Promise.all([waitForReady(client1), waitForReady(client2)]);

  client1.sendTags(["gaming"]);
  client2.sendTags(["music"]);

  // Both should get no match found
  const msgs = await Promise.all([
    client1.waitForMessage(3000),
    client2.waitForMessage(3000),
  ]);

  if (
    msgs[0].status !== StatusCodes.NO_MATCH ||
    msgs[1].status !== StatusCodes.NO_MATCH
  ) {
    throw new Error("Expected no match for different tags");
  }

  client1.close();
  client2.close();
}

async function testSingleTagSame() {
  const client1 = new TestClient("Client1-TagX");
  const client2 = new TestClient("Client2-TagX");

  await Promise.all([waitForReady(client1), waitForReady(client2)]);

  client1.sendTags(["sports"]);
  client2.sendTags(["sports"]);

  // Both should get match found
  await Promise.all([waitForMatch(client1), waitForMatch(client2)]);

  // Wait a bit to ensure backend has paired them
  await new Promise((res) => setTimeout(res, 1000));

  // Exchange ping-pong messages
  client1.sendMessage("ping");
  const partnerMsg = await client2.waitForMessage(3000);
  if (partnerMsg.content !== "ping")
    throw new Error("Partner did not receive ping message");

  client2.sendMessage("pong");
  const partnerMsg2 = await client1.waitForMessage(3000);
  if (partnerMsg2.content !== "pong")
    throw new Error("Partner did not receive pong message");

  client1.close();
  client2.close();
}

async function testMultipleTagsPartialOverlap() {
  const client1 = new TestClient("Client1-MultiTag1");
  const client2 = new TestClient("Client2-MultiTag2");

  await Promise.all([waitForReady(client1), waitForReady(client2)]);

  client1.sendTags(["gaming", "music"]);
  client2.sendTags(["coding", "music"]);

  // Both should get match found (since "music" overlaps)
  await Promise.all([waitForMatch(client1), waitForMatch(client2)]);

  client1.close();
  client2.close();
}

async function testMultipleTagsNoOverlap() {
  const client1 = new TestClient("Client1-NoOverlap1");
  const client2 = new TestClient("Client2-NoOverlap2");

  await Promise.all([waitForReady(client1), waitForReady(client2)]);

  client1.sendTags(["reading", "traveling"]);
  client2.sendTags(["gaming", "cooking"]);

  // Both should get no match found
  const msgs = await Promise.all([
    client1.waitForMessage(3000),
    client2.waitForMessage(3000),
  ]);

  if (
    msgs[0].status !== StatusCodes.NO_MATCH ||
    msgs[1].status !== StatusCodes.NO_MATCH
  ) {
    throw new Error("Expected no match for no overlapping tags");
  }

  client1.close();
  client2.close();
}

async function testMultipleTagsSubsetMatch() {
  const client1 = new TestClient("Client1-SubsetTags");
  const client2 = new TestClient("Client2-SubsetTags");

  await Promise.all([waitForReady(client1), waitForReady(client2)]);

  client1.sendTags(["music", "traveling", "gaming"]);
  client2.sendTags(["music", "gaming"]);

  // Wait for match on client2 and client1 (at least one must

  await Promise.all([waitForMatch(client1), waitForMatch(client2)]);

  // Exchange test messages to verify communication
  client1.sendMessage("hello from client1");
  const msgFromClient2 = await client2.waitForMessage(3000);
  if (msgFromClient2.content !== "hello from client1") {
    throw new Error("Client2 did not receive message from Client1");
  }

  client2.sendMessage("hello from client2");
  const msgFromClient1 = await client1.waitForMessage(3000);
  if (msgFromClient1.content !== "hello from client2") {
    throw new Error("Client1 did not receive message from Client2");
  }

  client1.close();
  client2.close();
}

async function main() {
  const tests = [
    { name: "No tags connect", fn: testNoTagsConnect },
    { name: "Single tag different", fn: testSingleTagDifferent },
    { name: "Single tag same", fn: testSingleTagSame },
    { name: "Multiple tags partial overlap", fn: testMultipleTagsPartialOverlap },
    { name: "Multiple tags no overlap", fn: testMultipleTagsNoOverlap },
    { name: "Multiple tags subset match", fn: testMultipleTagsSubsetMatch },
  ];

  const results: TestResult[] = [];

  for (const test of tests) {
    const result = await runTest(test.name, test.fn);
    results.push(result);
  }

  // Print final summary
  console.log("\n=== TEST SUMMARY ===");
  const passed = results.filter((r) => r.passed).length;
  const failed = results.length - passed;

  for (const r of results) {
    if (r.passed) {
      console.log(`✅ ${r.name}`);
    } else {
      console.log(`❌ ${r.name} - Reason: ${r.reason}`);
    }
  }
  console.log(`\nTotal: ${results.length}  Passed: ${passed}  Failed: ${failed}`);
}

main().catch((e) => {
  console.error("Fatal error running tests:", e);
  process.exit(1);
});

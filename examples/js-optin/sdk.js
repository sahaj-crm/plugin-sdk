import { createInterface } from "node:readline";

/**
 * Run an extension over the NDJSON JSON-RPC protocol.
 *
 * @param {{
 *   info: () => object | Promise<object>,
 *   contribute: (req: object) => object | Promise<object>,
 *   handle: (event: object) => object | Promise<object>,
 * }} extension
 */
export function serve(extension) {
  const rl = createInterface({ input: process.stdin, crlfDelay: Infinity });
  let queue = Promise.resolve();

  rl.on("line", (line) => {
    queue = queue.then(() => handleLine(extension, line));
  });

  rl.on("close", () => {
    queue.finally(() => process.exit(0));
  });
}

async function handleLine(extension, line) {
  if (!line.trim()) {
    return;
  }

  let request;
  try {
    request = JSON.parse(line);
  } catch (err) {
    replyError(0, -32700, "invalid JSON: " + err.message);
    return;
  }

  try {
    let result;
    switch (request.method) {
      case "info":
        result = await extension.info();
        break;
      case "contribute":
        result = await extension.contribute(request.params || {});
        break;
      case "handle":
        result = await extension.handle(request.params || {});
        break;
      case "shutdown":
        reply(request.id, { ok: true });
        process.exit(0);
        return;
      default:
        replyError(request.id, -32601, "unknown method: " + request.method);
        return;
    }
    reply(request.id, result);
  } catch (err) {
    replyError(request.id, -32000, err && err.message ? err.message : String(err));
  }
}

function reply(id, result) {
  process.stdout.write(
    JSON.stringify({ jsonrpc: "2.0", id, result }) + "\n",
  );
}

function replyError(id, code, message) {
  process.stdout.write(
    JSON.stringify({ jsonrpc: "2.0", id, error: { code, message } }) + "\n",
  );
}

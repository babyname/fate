/**
 * Capture real Web UI screenshots for the README.
 *
 * Prerequisites: the fate server is running (default http://127.0.0.1:18080)
 * and Chrome/Edge is installed. Uses the Chromium DevTools Protocol directly,
 * so no test-framework dependency is required (Node >= 22 provides WebSocket).
 *
 *   node scripts/capture-screenshots.mjs
 *
 * Environment overrides: FATE_BASE, CHROME_PATH, CDP_PORT.
 */
import { spawn } from 'node:child_process';
import { existsSync } from 'node:fs';
import { mkdir, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import path from 'node:path';

const BASE = process.env.FATE_BASE || 'http://127.0.0.1:18080';
const PORT = Number(process.env.CDP_PORT || 9229);
const OUT = path.resolve('docs/screenshots');
const CHROME = process.env.CHROME_PATH || [
  'C:/Program Files/Google/Chrome/Application/chrome.exe',
  'C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe',
  'C:/Program Files/Microsoft/Edge/Application/msedge.exe',
].find((p) => existsSync(p));

if (!CHROME) throw new Error('Chrome/Edge not found; set CHROME_PATH');

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function waitFor(check, label, timeoutMs = 30000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      if (await check()) return;
    } catch { /* endpoint not ready yet */ }
    await sleep(250);
  }
  throw new Error(`timeout waiting for ${label}`);
}

class CDP {
  constructor(ws) {
    this.ws = ws;
    this.nextId = 0;
    this.pending = new Map();
    ws.addEventListener('message', (event) => {
      const msg = JSON.parse(event.data);
      if (!msg.id || !this.pending.has(msg.id)) return;
      const { resolve, reject } = this.pending.get(msg.id);
      this.pending.delete(msg.id);
      if (msg.error) reject(new Error(JSON.stringify(msg.error)));
      else resolve(msg.result);
    });
  }

  send(method, params = {}) {
    return new Promise((resolve, reject) => {
      const id = ++this.nextId;
      this.pending.set(id, { resolve, reject });
      this.ws.send(JSON.stringify({ id, method, params }));
    });
  }

  async evaluate(expression) {
    const res = await this.send('Runtime.evaluate', {
      expression,
      awaitPromise: true,
      returnByValue: true,
    });
    if (res.exceptionDetails) {
      throw new Error(`evaluate failed: ${JSON.stringify(res.exceptionDetails)}`);
    }
    return res.result.value;
  }

  async waitForExpression(expression, label, timeoutMs = 30000) {
    await waitFor(async () => Boolean(await this.evaluate(expression)), label, timeoutMs);
  }

  async navigate(url) {
    await this.send('Page.navigate', { url });
    await this.waitForExpression(
      "document.readyState === 'complete' && !!document.querySelector('#root') && document.querySelector('#root').children.length > 0",
      `render of ${url}`,
    );
    await sleep(400);
  }

  async screenshot(file) {
    const { data } = await this.send('Page.captureScreenshot', {
      format: 'png',
      captureBeyondViewport: false,
    });
    await writeFile(path.join(OUT, file), Buffer.from(data, 'base64'));
    console.log(`captured ${file}`);
  }
}

async function main() {
  await mkdir(OUT, { recursive: true });
  const profile = path.join(tmpdir(), 'fate-cdp-profile');
  const chrome = spawn(CHROME, [
    '--headless=new',
    '--disable-gpu',
    '--no-first-run',
    '--no-default-browser-check',
    `--remote-debugging-port=${PORT}`,
    '--remote-allow-origins=*',
    `--user-data-dir=${profile}`,
    '--window-size=1440,900',
    '--hide-scrollbars',
    '--force-device-scale-factor=1',
    'about:blank',
  ], { stdio: 'ignore' });

  try {
    await waitFor(
      async () => (await fetch(`http://127.0.0.1:${PORT}/json/version`)).ok,
      'Chrome DevTools endpoint',
      20000,
    );
    const targets = await (await fetch(`http://127.0.0.1:${PORT}/json/list`)).json();
    const page = targets.find((t) => t.type === 'page');
    if (!page) throw new Error('no page target available');

    const ws = new WebSocket(page.webSocketDebuggerUrl);
    await new Promise((resolve, reject) => {
      ws.addEventListener('open', resolve, { once: true });
      ws.addEventListener('error', reject, { once: true });
    });
    const cdp = new CDP(ws);
    await cdp.send('Page.enable');
    await cdp.send('Runtime.enable');
    await cdp.send('Emulation.setDeviceMetricsOverride', {
      width: 1440, height: 900, deviceScaleFactor: 1, mobile: false,
    });

    // 1. Home page, as a first-time visitor sees it.
    await cdp.navigate(BASE + '/');
    await cdp.screenshot('home.png');

    // 2. Create a real task through the public API (same flow as the UI).
    const taskId = await cdp.evaluate(`(async () => {
      const res = await fetch('/api/v1/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ surname: '张', born: '2024/06/15 10:30', sex: 'boy', poetry_mode: 0 }),
      });
      const body = await res.json();
      if (!res.ok) throw new Error('generate failed: ' + JSON.stringify(body));
      return body.task_id;
    })()`);
    if (!taskId) throw new Error('generate returned no task_id');
    console.log(`task_id=${taskId}`);

    // 3. Results page: top names.
    await cdp.navigate(`${BASE}/results/${taskId}`);
    await cdp.waitForExpression(
      "document.body.innerText.includes('精选好名') && document.querySelectorAll('table tbody tr').length > 0",
      'results table rendering',
      45000,
    );
    await cdp.screenshot('results.png');

    // 4. Candidate table with pagination.
    await cdp.evaluate('window.scrollTo(0, 760)');
    await sleep(350);
    await cdp.screenshot('table.png');

    // 5. Card view with five-grid luck and same-name lookup.
    await cdp.evaluate(`(() => {
      const button = [...document.querySelectorAll('button')]
        .find((b) => b.textContent.trim() === '卡片');
      if (!button) throw new Error('card view button not found');
      button.click();
      return true;
    })()`);
    await cdp.waitForExpression(
      "!document.querySelector('table') && document.querySelectorAll('dl[aria-label=\"五格吉凶\"]').length > 0 && document.body.innerText.includes('查同名')",
      'card view rendering',
      15000,
    );
    await sleep(350);
    await cdp.screenshot('cards.png');

    ws.close();
    console.log('screenshot capture complete');
  } finally {
    chrome.kill();
  }
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});

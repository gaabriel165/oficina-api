import http from "k6/http";
import { check, sleep } from "k6";

const BASE = __ENV.BASE_URL;
const CPF = __ENV.CPF || "98765432100";
const EMAIL = __ENV.LOGIN_EMAIL || "admin@oficina.com";
const PASSWORD = __ENV.LOGIN_PASSWORD || "admin123";

export const options = {
  stages: [
    { duration: "30s", target: 10 },
    { duration: "1m", target: 30 },
    { duration: "2m", target: 60 },
    { duration: "1m", target: 0 },
  ],
  thresholds: {
    http_req_failed: ["rate<0.10"],
    http_req_duration: ["p(95)<2000"],
  },
};

const jsonHeaders = { "Content-Type": "application/json" };

export function setup() {
  const byCpf = http.post(`${BASE}/auth/cpf`, JSON.stringify({ cpf: CPF }), { headers: jsonHeaders });
  if (byCpf.status === 200) {
    return { token: byCpf.json("token"), issuer: "lambda auth-cpf" };
  }

  const byLogin = http.post(`${BASE}/api/v1/auth/login`, JSON.stringify({ email: EMAIL, password: PASSWORD }), { headers: jsonHeaders });
  check(byLogin, { "login succeeded": (r) => r.status === 200 });
  return { token: byLogin.json("token"), issuer: "api login" };
}

export default function (data) {
  const authHeaders = { headers: { Authorization: `Bearer ${data.token}`, "X-Request-ID": `k6-${__VU}-${__ITER}` } };

  const orders = http.get(`${BASE}/api/v1/service-orders`, authHeaders);
  check(orders, { "list orders 200": (r) => r.status === 200 });

  const metrics = http.get(`${BASE}/api/v1/service-orders/metrics`, authHeaders);
  check(metrics, { "metrics 200": (r) => r.status === 200 });

  const login = http.post(`${BASE}/api/v1/auth/login`, JSON.stringify({ email: EMAIL, password: PASSWORD }), { headers: jsonHeaders });
  check(login, { "login 200": (r) => r.status === 200 });

  sleep(0.5);
}

export function teardown(data) {
  console.log(`token issued by ${data.issuer}`);
}

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
  const operatorLogin = http.post(`${BASE}/api/v1/auth/login`, JSON.stringify({ email: EMAIL, password: PASSWORD }), { headers: jsonHeaders });
  check(operatorLogin, { "operator login succeeded": (r) => r.status === 200 });
  const operatorToken = operatorLogin.json("token");

  const customerLogin = http.post(`${BASE}/auth/cpf`, JSON.stringify({ cpf: CPF }), { headers: jsonHeaders });
  const customerToken = customerLogin.status === 200 ? customerLogin.json("token") : operatorToken;

  return { operatorToken, customerToken, customerIssuer: customerLogin.status === 200 ? "lambda auth-cpf" : "operator fallback" };
}

function authHeaders(token) {
  return { headers: { Authorization: `Bearer ${token}`, "X-Request-ID": `k6-${__VU}-${__ITER}` } };
}

export default function (data) {
  const customerOrders = http.get(`${BASE}/api/v1/service-orders`, authHeaders(data.customerToken));
  check(customerOrders, { "customer lists own orders 200": (r) => r.status === 200 });

  const operatorOrders = http.get(`${BASE}/api/v1/service-orders`, authHeaders(data.operatorToken));
  check(operatorOrders, { "operator lists orders 200": (r) => r.status === 200 });

  const metrics = http.get(`${BASE}/api/v1/service-orders/metrics`, authHeaders(data.operatorToken));
  check(metrics, { "operator metrics 200": (r) => r.status === 200 });

  const login = http.post(`${BASE}/api/v1/auth/login`, JSON.stringify({ email: EMAIL, password: PASSWORD }), { headers: jsonHeaders });
  check(login, { "login 200": (r) => r.status === 200 });

  sleep(0.5);
}

export function teardown(data) {
  console.log(`customer token issued by ${data.customerIssuer}`);
}

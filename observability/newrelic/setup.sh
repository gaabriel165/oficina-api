#!/usr/bin/env bash
set -euo pipefail

: "${NEW_RELIC_API_KEY:?user API key (NRAK-...)}"
: "${NEW_RELIC_ACCOUNT_ID:?account id}"
: "${ALERT_EMAIL:?e-mail that receives alerts}"
HEALTH_URL="${HEALTH_URL:-}"
ONLY_SYNTHETICS="${ONLY_SYNTHETICS:-false}"

ENDPOINT="https://api.newrelic.com/graphql"
HERE="$(cd "$(dirname "$0")" && pwd)"

graphql() {
  local query="$1" variables="${2:-{\}}"
  curl -sS "$ENDPOINT" \
    -H "Content-Type: application/json" \
    -H "API-Key: $NEW_RELIC_API_KEY" \
    --data "$(jq -n --arg q "$query" --argjson v "$variables" '{query: $q, variables: $v}')"
}

fail_on_errors() {
  local response="$1" step="$2"
  if echo "$response" | jq -e '.errors? // (.. | .errors? | select(. != null and . != [])) ' >/dev/null 2>&1; then
    echo "[$step] error:" >&2
    echo "$response" | jq . >&2
    exit 1
  fi
}

create_synthetics() {
  echo "==> Synthetics ping monitor"
  RESPONSE="$(graphql 'mutation($accountId: Int!, $url: String!) {
    syntheticsCreateSimpleMonitor(accountId: $accountId, monitor: {
      name: "oficina-api-health", uri: $url, period: EVERY_5_MINUTES, status: ENABLED,
      locations: { public: ["US_EAST_1"] },
      advancedOptions: { shouldBypassHeadRequest: true, responseValidationText: "ok" }
    }) { monitor { guid } errors { description type } }
  }' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --arg u "$HEALTH_URL" '{accountId: $id, url: $u}')")"
  fail_on_errors "$RESPONSE" synthetics
  echo "monitor: $(echo "$RESPONSE" | jq -r '.data.syntheticsCreateSimpleMonitor.monitor.guid')"
}

if [ "$ONLY_SYNTHETICS" = "true" ]; then
  : "${HEALTH_URL:?HEALTH_URL is required with ONLY_SYNTHETICS=true}"
  create_synthetics
  exit 0
fi

echo "==> Dashboard"
DASHBOARD="$(sed "s/ACCOUNT_ID/$NEW_RELIC_ACCOUNT_ID/g" "$HERE/dashboard.json")"
RESPONSE="$(graphql 'mutation($accountId: Int!, $dashboard: DashboardInput!) {
  dashboardCreate(accountId: $accountId, dashboard: $dashboard) {
    entityResult { guid name }
    errors { description type }
  }
}' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --argjson d "$DASHBOARD" '{accountId: $id, dashboard: $d}')")"
fail_on_errors "$RESPONSE" dashboard
DASHBOARD_GUID="$(echo "$RESPONSE" | jq -r '.data.dashboardCreate.entityResult.guid')"
echo "dashboard guid: $DASHBOARD_GUID"

echo "==> Alert policy"
RESPONSE="$(graphql 'mutation($accountId: Int!) {
  alertsPolicyCreate(accountId: $accountId, policy: { name: "Oficina API", incidentPreference: PER_CONDITION }) { id name }
}' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" '{accountId: $id}')")"
fail_on_errors "$RESPONSE" policy
POLICY_ID="$(echo "$RESPONSE" | jq -r '.data.alertsPolicyCreate.id')"
echo "policy id: $POLICY_ID"

create_condition() {
  local name="$1" nrql="$2" threshold="$3" duration="$4"
  RESPONSE="$(graphql 'mutation($accountId: Int!, $policyId: ID!, $condition: AlertsNrqlConditionStaticInput!) {
    alertsNrqlConditionStaticCreate(accountId: $accountId, policyId: $policyId, condition: $condition) { id name }
  }' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --arg p "$POLICY_ID" --arg n "$name" --arg q "$nrql" --argjson t "$threshold" --argjson d "$duration" '{
      accountId: $id, policyId: $p,
      condition: {
        name: $n, enabled: true,
        nrql: { query: $q },
        signal: { aggregationWindow: 60, aggregationMethod: "EVENT_FLOW", aggregationDelay: 120 },
        terms: [{ threshold: $t, thresholdOccurrences: "AT_LEAST_ONCE", thresholdDuration: $d, operator: "ABOVE", priority: "CRITICAL" }],
        violationTimeLimitSeconds: 3600
      }
    }')")"
  fail_on_errors "$RESPONSE" "condition:$name"
  echo "condition: $(echo "$RESPONSE" | jq -r '.data.alertsNrqlConditionStaticCreate.name')"
}

create_condition "Falha no processamento de ordens de serviço" \
  "SELECT count(*) FROM Log WHERE event = 'service_order.failed'" 0 300
create_condition "Erros de integração (Resend / BrasilAPI)" \
  "SELECT count(*) FROM Log WHERE event = 'integration.error'" 2 300
create_condition "Erros 5xx na API" \
  "SELECT count(*) FROM Log WHERE event = 'http.request' AND numeric(status) >= 500" 0 300
create_condition "Latência p95 acima de 1s" \
  "SELECT percentile(numeric(latency_ms), 95) FROM Log WHERE event = 'http.request'" 1000 300

echo "==> Notification destination, channel and workflow"
RESPONSE="$(graphql 'mutation($accountId: Int!, $email: String!) {
  aiNotificationsCreateDestination(accountId: $accountId, destination: {
    name: "Oficina API — e-mail", type: EMAIL,
    properties: [{ key: "email", value: $email }]
  }) { destination { id } error { ... on AiNotificationsResponseError { description } } }
}' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --arg e "$ALERT_EMAIL" '{accountId: $id, email: $e}')")"
fail_on_errors "$RESPONSE" destination
DESTINATION_ID="$(echo "$RESPONSE" | jq -r '.data.aiNotificationsCreateDestination.destination.id')"

RESPONSE="$(graphql 'mutation($accountId: Int!, $destinationId: ID!) {
  aiNotificationsCreateChannel(accountId: $accountId, channel: {
    name: "Oficina API — canal e-mail", type: EMAIL, product: IINT, destinationId: $destinationId,
    properties: [{ key: "subject", value: "{{ issueTitle }}" }]
  }) { channel { id } error { ... on AiNotificationsResponseError { description } } }
}' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --arg d "$DESTINATION_ID" '{accountId: $id, destinationId: $d}')")"
fail_on_errors "$RESPONSE" channel
CHANNEL_ID="$(echo "$RESPONSE" | jq -r '.data.aiNotificationsCreateChannel.channel.id')"

RESPONSE="$(graphql 'mutation($accountId: Int!, $policyId: String!, $channelId: ID!) {
  aiWorkflowsCreateWorkflow(accountId: $accountId, createWorkflowData: {
    name: "Oficina API — alertas por e-mail", workflowEnabled: true, mutingRulesHandling: NOTIFY_ALL_ISSUES,
    destinationConfigurations: [{ channelId: $channelId }],
    issuesFilter: { name: "policy", type: FILTER, predicates: [{ attribute: "labels.policyIds", operator: EXACTLY_MATCHES, values: [$policyId] }] }
  }) { workflow { id } errors { description } }
}' "$(jq -n --argjson id "$NEW_RELIC_ACCOUNT_ID" --arg p "$POLICY_ID" --arg c "$CHANNEL_ID" '{accountId: $id, policyId: $p, channelId: $c}')")"
fail_on_errors "$RESPONSE" workflow
echo "workflow: $(echo "$RESPONSE" | jq -r '.data.aiWorkflowsCreateWorkflow.workflow.id')"

if [ -n "$HEALTH_URL" ]; then
  create_synthetics
else
  echo "==> HEALTH_URL not set: skipping Synthetics monitor (run later with ONLY_SYNTHETICS=true HEALTH_URL=...)"
fi

echo
echo "Done. Dashboard: https://one.newrelic.com/dashboards/detail/$DASHBOARD_GUID"

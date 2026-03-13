#!/usr/bin/env bash
# Run a smoke-test query against every knowbe4 table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN="\033[0;32m"
  RED="\033[0;31m"
  YELLOW="\033[0;33m"
  RESET="\033[0m"
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-50s" "$table"
  local output
  if output=$(steampipe query "$query" 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    if echo "$output" | grep -q "missing_required_scope\|not_authorized\|403"; then
      printf "${YELLOW}SKIP${RESET} (missing API scope)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}

echo ""
echo "KnowBe4 Steampipe plugin — table smoke tests"
echo "============================================="
echo ""

run_test "knowbe4_account"                          "select name, subscription_level, number_of_seats from knowbe4_account limit 1"
run_test "knowbe4_account_risk_score_history"       "select risk_score, date from knowbe4_account_risk_score_history limit 1"
run_test "knowbe4_user"                             "select id, email, phish_prone_percentage from knowbe4_user limit 1"
run_test "knowbe4_group"                            "select id, name, member_count from knowbe4_group limit 1"
run_test "knowbe4_phishing_campaign"                "select campaign_id, name, psts_count, status from knowbe4_phishing_campaign limit 1"
run_test "knowbe4_phishing_security_test"           "select pst_id, campaign_id, status from knowbe4_phishing_security_test limit 1"
run_test "knowbe4_training_store_purchase"          "select store_purchased_id, name, type from knowbe4_training_store_purchase limit 1"
run_test "knowbe4_training_policy"                  "select policy_id, name, status from knowbe4_training_policy limit 1"
run_test "knowbe4_training_campaign"                "select campaign_id, name, status from knowbe4_training_campaign limit 1"
run_test "knowbe4_training_enrollment"              "select enrollment_id, module_name, status from knowbe4_training_enrollment limit 1"

echo ""
echo "  (The following tables require a required qual — tested separately)"

# Tables requiring a parent ID qual
USER_ID=$(steampipe query --output json "select id from knowbe4_user limit 1" 2>/dev/null | jq -r '.rows[0].id // empty')
GROUP_ID=$(steampipe query --output json "select id from knowbe4_group limit 1" 2>/dev/null | jq -r '.rows[0].id // empty')
PST_ID=$(steampipe query --output json "select pst_id from knowbe4_phishing_security_test limit 1" 2>/dev/null | jq -r '.rows[0].pst_id // empty')

if [ -n "$USER_ID" ]; then
  run_test "knowbe4_user_risk_score_history"        "select user_id, risk_score, date from knowbe4_user_risk_score_history where user_id = $USER_ID limit 1"
fi
if [ -n "$GROUP_ID" ]; then
  run_test "knowbe4_group_risk_score_history"       "select group_id, risk_score, date from knowbe4_group_risk_score_history where group_id = $GROUP_ID limit 1"
fi
if [ -n "$PST_ID" ]; then
  run_test "knowbe4_phishing_recipient"             "select recipient_id, pst_id from knowbe4_phishing_recipient where pst_id = $PST_ID limit 1"
fi

echo ""
echo "---------------------------------------------"
printf "Results: ${GREEN}%d passed${RESET}  ${RED}%d failed${RESET}  ${YELLOW}%d skipped${RESET}\n" "$PASS" "$FAIL" "$SKIP"
echo ""

exit "$FAIL"

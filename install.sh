#!/usr/bin/env bash

set -Eeuo pipefail

APP_NAME="albumer"
APP_USER="albumer"
APP_GROUP="albumer"
INSTALL_DIR="/opt/albumer"
SERVICE_FILE="/etc/systemd/system/albumer.service"
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

info() {
  printf '\n%s\n' "$1"
}

fail() {
  printf 'Error: %s\n' "$1" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || fail "Required command is missing: $1"
}

prompt_value() {
  local variable_name="$1"
  local label="$2"
  local default_value="$3"
  local value

  read -r -p "$label [$default_value]: " value
  printf -v "$variable_name" '%s' "${value:-$default_value}"
}

prompt_secret() {
  local variable_name="$1"
  local label="$2"
  local value

  read -r -s -p "$label: " value
  printf '\n'
  printf -v "$variable_name" '%s' "$value"
}

port_is_in_use() {
  local port="$1"

  if command -v ss >/dev/null 2>&1; then
    ss -H -ltn | awk '{print $4}' | grep -Eq "(^|:)$port$"
    return
  fi

  if command -v lsof >/dev/null 2>&1; then
    lsof -nP -iTCP:"$port" -sTCP:LISTEN >/dev/null 2>&1
    return
  fi

  fail "Cannot check the port. Install ss (the iproute2 package) or lsof."
}

neo4j_is_present() {
  command -v neo4j >/dev/null 2>&1 ||
    command -v cypher-shell >/dev/null 2>&1 ||
    systemctl list-unit-files neo4j.service --no-legend 2>/dev/null | grep -q '^neo4j\.service'
}

write_env() {
  local target="$1"

  {
    printf 'ADDR=:%s\n' "$PORT"
    printf 'NEO4J_URI=%s\n' "$NEO4J_URI"
    printf 'NEO4J_USER=%s\n' "$NEO4J_USER"
    printf 'NEO4J_PASSWORD=%s\n' "$NEO4J_PASSWORD"
    printf 'APP_SECRET=%s\n' "$APP_SECRET"
    printf 'UPLOAD_DIR=%s\n' "$UPLOAD_DIR"
    printf 'ADMIN_PASSWORD_HASH=%s\n' "$ADMIN_PASSWORD_HASH"
  } >"$target"
  chmod 600 "$target"
}

if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
  SUDO=()
else
  require_command sudo
  SUDO=(sudo)
fi

require_command go
require_command systemctl
require_command install
require_command awk
require_command grep
require_command od
require_command tr
require_command getent
require_command groupadd
require_command useradd

info "Checking Neo4j..."
neo4j_is_present || fail "Neo4j was not found. Install Neo4j 5+ before running this installer."
printf 'Neo4j is available.\n'

info "Albumer configuration"
prompt_value PORT "HTTP port" "4000"
while ! [[ "$PORT" =~ ^[0-9]+$ ]] || ((PORT < 1 || PORT > 65535)); do
  printf 'Enter a number between 1 and 65535.\n'
  prompt_value PORT "HTTP port" "4000"
done
while port_is_in_use "$PORT"; do
  printf 'Port %s is already in use. Choose another port.\n' "$PORT"
  prompt_value PORT "HTTP port" "4000"
  while ! [[ "$PORT" =~ ^[0-9]+$ ]] || ((PORT < 1 || PORT > 65535)); do
    printf 'Enter a number between 1 and 65535.\n'
    prompt_value PORT "HTTP port" "4000"
  done
done

prompt_value NEO4J_URI "NEO4J_URI" "neo4j://localhost:7687"
prompt_value NEO4J_USER "NEO4J_USER" "neo4j"
prompt_secret NEO4J_PASSWORD "NEO4J_PASSWORD"
[[ -n "$NEO4J_PASSWORD" ]] || fail "NEO4J_PASSWORD cannot be empty."

DEFAULT_SECRET="$(od -An -N32 -tx1 /dev/urandom | tr -d ' \n')"
prompt_value APP_SECRET "APP_SECRET" "$DEFAULT_SECRET"
prompt_value UPLOAD_DIR "UPLOAD_DIR" "/var/lib/albumer/uploads"
[[ "$UPLOAD_DIR" = /* ]] || fail "UPLOAD_DIR must be an absolute path."

printf '\nThe administrator password is stored as a bcrypt hash.\n'
printf 'If you leave it empty, the application will use the default password: admin\n'
prompt_secret ADMIN_PASSWORD "Administrator password"
if [[ -n "$ADMIN_PASSWORD" ]]; then
  prompt_secret ADMIN_PASSWORD_CONFIRM "Repeat administrator password"
  [[ "$ADMIN_PASSWORD" == "$ADMIN_PASSWORD_CONFIRM" ]] || fail "Administrator passwords do not match."
  ADMIN_PASSWORD_HASH="$(cd "$SCRIPT_DIR" && printf '%s' "$ADMIN_PASSWORD" | go run ./cmd/hash-password)"
else
  ADMIN_PASSWORD_HASH=""
fi

write_env "$TMP_DIR/.env"

info "Building the application..."
(
  cd "$SCRIPT_DIR"
  go build -o "$TMP_DIR/albumer" ./cmd/albumer
)

info "Installing files and the systemd service..."
if ! getent group "$APP_GROUP" >/dev/null 2>&1; then
  "${SUDO[@]}" groupadd --system "$APP_GROUP"
fi
if ! id -u "$APP_USER" >/dev/null 2>&1; then
  "${SUDO[@]}" useradd --system --gid "$APP_GROUP" --home-dir "$INSTALL_DIR" --shell /usr/sbin/nologin "$APP_USER"
fi

"${SUDO[@]}" install -d -o "$APP_USER" -g "$APP_GROUP" -m 750 "$INSTALL_DIR"
"${SUDO[@]}" install -d -o "$APP_USER" -g "$APP_GROUP" -m 750 "$UPLOAD_DIR"
"${SUDO[@]}" install -o root -g root -m 755 "$TMP_DIR/albumer" "$INSTALL_DIR/albumer"
"${SUDO[@]}" install -o "$APP_USER" -g "$APP_GROUP" -m 600 "$TMP_DIR/.env" "$INSTALL_DIR/.env"
"${SUDO[@]}" install -o root -g root -m 644 "$SCRIPT_DIR/deploy/albumer.service" "$SERVICE_FILE"

"${SUDO[@]}" systemctl daemon-reload
"${SUDO[@]}" systemctl enable --now "$APP_NAME"

info "Albumer has been installed."
printf 'Address: http://localhost:%s\n' "$PORT"
printf 'Configuration: %s/.env\n' "$INSTALL_DIR"
printf 'Uploads: %s\n' "$UPLOAD_DIR"

cat <<'EOF'

Service commands:
  sudo systemctl status albumer
  sudo systemctl start albumer
  sudo systemctl stop albumer
  sudo systemctl restart albumer
  sudo journalctl -u albumer -f
EOF

#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_PORT=8080
FRONTEND_PORT=5173
WEB_DIR="${ROOT_DIR}/web"
LOCAL_CONFIG_DIR="${ROOT_DIR}/configs"
LOCAL_CONFIG_FILE="${LOCAL_CONFIG_DIR}/config.yaml"

backend_pid=""
frontend_pid=""

cd "$ROOT_DIR"

stop_port_processes() {
	local port="$1"
	local pids

	pids="$(lsof -ti tcp:${port} -sTCP:LISTEN || true)"
	if [ -n "$pids" ]; then
		echo "Stopping processes on port ${port}: ${pids}"
		kill ${pids}

		sleep 1

		local remaining_pids
		remaining_pids="$(lsof -ti tcp:${port} -sTCP:LISTEN || true)"
		if [ -n "$remaining_pids" ]; then
			echo "Force killing remaining processes on port ${port}: ${remaining_pids}"
			kill -9 ${remaining_pids}
		fi
	else
		echo "Port ${port} is free"
	fi
}

cleanup() {
	if [ -n "${frontend_pid}" ] && kill -0 "${frontend_pid}" 2>/dev/null; then
		echo "Stopping frontend process ${frontend_pid}"
		kill "${frontend_pid}" 2>/dev/null || true
	fi

	if [ -n "${backend_pid}" ] && kill -0 "${backend_pid}" 2>/dev/null; then
		echo "Stopping backend process ${backend_pid}"
		kill "${backend_pid}" 2>/dev/null || true
	fi
}

trap cleanup EXIT INT TERM

stop_port_processes "${BACKEND_PORT}"
stop_port_processes "${FRONTEND_PORT}"

if [ ! -d "${WEB_DIR}" ]; then
	echo "Frontend directory not found: ${WEB_DIR}"
	exit 1
fi

if [ ! -f "${WEB_DIR}/package.json" ]; then
	echo "Frontend package.json not found: ${WEB_DIR}/package.json"
	exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
	echo "npm is required to start the frontend"
	exit 1
fi

if [ ! -d "${WEB_DIR}/node_modules" ]; then
	echo "Installing frontend dependencies..."
	(
		cd "${WEB_DIR}"
		npm install
	)
fi

if [ ! -f "${LOCAL_CONFIG_FILE}" ]; then
	echo "Backend config file not found: ${LOCAL_CONFIG_FILE}"
	exit 1
fi

echo "Starting backend on port ${BACKEND_PORT}..."
# Force local module mode to avoid picking parent go.work, and pin local config file.
GOWORK=off CONFIG_FILE="${LOCAL_CONFIG_FILE}" go run ./cmd/main.go &
backend_pid=$!

echo "Starting frontend on port ${FRONTEND_PORT}..."
(
	cd "${WEB_DIR}"
	npm run dev -- --host 0.0.0.0 --port "${FRONTEND_PORT}"
) &
frontend_pid=$!

echo "Backend:  http://localhost:${BACKEND_PORT}"
echo "Frontend: http://localhost:${FRONTEND_PORT}"
echo "Press Ctrl+C to stop both processes."

wait "${backend_pid}" "${frontend_pid}"

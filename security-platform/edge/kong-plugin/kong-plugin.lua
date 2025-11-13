-- Security Platform Kong Plugin
-- Integrates with control plane for dynamic security policies

local http = require "resty.http"
local cjson = require "cjson"

local SecurityPlatformHandler = {}

SecurityPlatformHandler.PRIORITY = 1000
SecurityPlatformHandler.VERSION = "0.1.0"

local function get_config()
    return {
        control_plane_url = kong.configuration.custom_plugins.security_platform.control_plane_url or
                           os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL") or
                           "https://console.securityplatform.com",
        auth_key = kong.configuration.custom_plugins.security_platform.auth_key or
                  os.getenv("SECURITY_PLATFORM_AUTH_KEY") or "",
        policy_cache_ttl = 60,
    }
end

local function fetch_policy(config)
    local httpc = http.new()
    local res, err = httpc:request_uri(config.control_plane_url .. "/api/v1/policy", {
        method = "GET",
        headers = {
            ["Authorization"] = "Bearer " .. config.auth_key,
        },
        timeout = 2000,
    })

    if err or not res or res.status ~= 200 then
        kong.log.err("Failed to fetch policy: ", err or "unknown error")
        return {}
    end

    local policy = cjson.decode(res.body)
    return policy.rules or {}
end

function SecurityPlatformHandler:access(conf)
    local config = get_config()
    local rules = fetch_policy(config)

    -- Check if request should be blocked
    local request_path = kong.request.get_path()
    for pattern, action in pairs(rules) do
        if action == "block" and string.find(request_path, pattern) then
            kong.log.warn("Request blocked by security policy: ", pattern)
            return kong.response.exit(403, {
                error = "Request blocked by security policy",
                rule = pattern,
            }, {
                ["Content-Type"] = "application/json",
            })
        end
    end

    -- Add security headers
    kong.response.set_header("X-Security-Platform", "enabled")
end

function SecurityPlatformHandler:response(conf)
    -- Capture response metadata and send to control plane
    local trace_id = kong.request.get_header("X-Trace-ID") or kong.ctx.shared.trace_id
    local request_id = kong.request.get_header("X-Request-ID") or kong.ctx.shared.request_id

    -- Send telemetry (async)
    kong.timer.at(0, function()
        local httpc = http.new()
        httpc:request_uri(config.control_plane_url .. "/api/v1/telemetry", {
            method = "POST",
            headers = {
                ["Authorization"] = "Bearer " .. config.auth_key,
                ["Content-Type"] = "application/json",
            },
            body = cjson.encode({
                trace_id = trace_id,
                request_id = request_id,
                method = kong.request.get_method(),
                path = kong.request.get_path(),
                status = kong.response.get_status(),
                timestamp = ngx.time(),
            }),
        })
    end)
end

return SecurityPlatformHandler

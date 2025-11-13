-- Security Platform WAF Plugin for Nginx
-- Integrates with control plane for dynamic rule updates

local http = require "resty.http"
local cjson = require "cjson"

local control_plane_url = os.getenv("SECURITY_PLATFORM_CONTROL_PLANE_URL") or "https://console.securityplatform.com"
local auth_key = os.getenv("SECURITY_PLATFORM_AUTH_KEY") or ""

-- Cache for policy rules (TTL: 60 seconds)
local policy_cache = {}
local cache_ttl = 60
local last_fetch = 0

local function fetch_policy()
    local now = ngx.time()
    if now - last_fetch < cache_ttl and policy_cache.rules then
        return policy_cache.rules
    end

    local httpc = http.new()
    local res, err = httpc:request_uri(control_plane_url .. "/api/v1/policy", {
        method = "GET",
        headers = {
            ["Authorization"] = "Bearer " .. auth_key,
            ["Content-Type"] = "application/json",
        },
        timeout = 2000,
    })

    if err or not res or res.status ~= 200 then
        ngx.log(ngx.ERR, "Failed to fetch policy: ", err or "unknown error")
        return policy_cache.rules or {}
    end

    local policy = cjson.decode(res.body)
    policy_cache = {
        rules = policy.rules or {},
        fetched_at = now,
    }
    last_fetch = now

    return policy_cache.rules
end

local function check_block_rule(rules, request_uri, headers)
    -- Check if request matches any block rules
    for pattern, action in pairs(rules) do
        if action == "block" then
            -- Simple pattern matching (in production, use proper regex)
            if string.find(request_uri, pattern) or string.find(ngx.var.request_uri, pattern) then
                return true, pattern
            end
        end
    end
    return false, nil
end

local function send_incident(incident_data)
    -- Send incident to control plane (async, non-blocking)
    ngx.timer.at(0, function()
        local httpc = http.new()
        httpc:request_uri(control_plane_url .. "/api/v1/incidents", {
            method = "POST",
            headers = {
                ["Authorization"] = "Bearer " .. auth_key,
                ["Content-Type"] = "application/json",
            },
            body = cjson.encode(incident_data),
            timeout = 2000,
        })
    end)
end

-- Main access phase handler
local function access()
    local rules = fetch_policy()
    
    -- Check block rules
    local should_block, matched_rule = check_block_rule(rules, ngx.var.request_uri, ngx.req.get_headers())
    
    if should_block then
        -- Log incident
        send_incident({
            type = "waf_block",
            rule = matched_rule,
            request_uri = ngx.var.request_uri,
            remote_addr = ngx.var.remote_addr,
            user_agent = ngx.var.http_user_agent,
            timestamp = ngx.time(),
        })

        ngx.log(ngx.WARN, "Request blocked by WAF rule: ", matched_rule)
        ngx.status = 403
        ngx.header.content_type = "application/json"
        ngx.say(cjson.encode({ error = "Request blocked by security policy", rule = matched_rule }))
        ngx.exit(403)
    end
end

-- Run access handler
access()
